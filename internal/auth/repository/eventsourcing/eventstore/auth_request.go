package eventstore

import (
	"context"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/auth_request/model"
	cache "github.com/caos/zitadel/internal/auth_request/repository"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/id"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type AuthRequestRepo struct {
	UserEvents   *user_event.UserEventstore
	AuthRequests cache.AuthRequestCache
	View         *view.View

	UserSessionViewProvider userSessionViewProvider
	UserViewProvider        userViewProvider
	UserEventProvider       userEventProvider

	IdGenerator id.Generator

	PasswordCheckLifeTime    time.Duration
	MfaInitSkippedLifeTime   time.Duration
	MfaSoftwareCheckLifeTime time.Duration
	MfaHardwareCheckLifeTime time.Duration
}

type userSessionViewProvider interface {
	UserSessionByIDs(string, string) (*view_model.UserSessionView, error)
	UserSessionsByAgentID(string) ([]*view_model.UserSessionView, error)
}
type userViewProvider interface {
	UserByID(string) (*view_model.UserView, error)
}

type userEventProvider interface {
	UserEventsByID(ctx context.Context, id string, sequence uint64) ([]*es_models.Event, error)
}

func (repo *AuthRequestRepo) Health(ctx context.Context) error {
	if err := repo.UserEvents.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequests.Health(ctx)
}

func (repo *AuthRequestRepo) CreateAuthRequest(ctx context.Context, request *model.AuthRequest) (*model.AuthRequest, error) {
	reqID, err := repo.IdGenerator.Next()
	if err != nil {
		return nil, err
	}
	request.ID = reqID
	ids, err := repo.View.AppIDsFromProjectByClientID(ctx, request.ApplicationID)
	if err != nil {
		return nil, err
	}
	request.Audience = ids
	err = repo.AuthRequests.SaveAuthRequest(ctx, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (repo *AuthRequestRepo) AuthRequestByID(ctx context.Context, id string) (*model.AuthRequest, error) {
	return repo.getAuthRequest(ctx, id, false)
}

func (repo *AuthRequestRepo) AuthRequestByIDCheckLoggedIn(ctx context.Context, id string) (*model.AuthRequest, error) {
	return repo.getAuthRequest(ctx, id, true)
}

func (repo *AuthRequestRepo) SaveAuthCode(ctx context.Context, id, code string) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return err
	}
	request.Code = code
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) AuthRequestByCode(ctx context.Context, code string) (*model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	steps, err := repo.nextSteps(ctx, request, true)
	if err != nil {
		return nil, err
	}
	request.PossibleSteps = steps
	return request, nil
}

func (repo *AuthRequestRepo) DeleteAuthRequest(ctx context.Context, id string) error {
	return repo.AuthRequests.DeleteAuthRequest(ctx, id)
}

func (repo *AuthRequestRepo) CheckUsername(ctx context.Context, id, username string) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return err
	}
	user, err := repo.View.UserByLoginName(username)
	if err != nil {
		return err
	}
	request.SetUserInfo(user.ID, user.UserName, user.ResourceOwner)
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) SelectUser(ctx context.Context, id, userID string) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return err
	}
	user, err := repo.View.UserByID(userID)
	if err != nil {
		return err
	}
	request.SetUserInfo(user.ID, user.UserName, user.ResourceOwner)
	return repo.AuthRequests.UpdateAuthRequest(ctx, request)
}

func (repo *AuthRequestRepo) VerifyPassword(ctx context.Context, id, userID, password string, info *model.BrowserInfo) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return err
	}
	if request.UserID != userID {
		return errors.ThrowPreconditionFailed(nil, "EVENT-ds35D", "user id does not match request id")
	}
	return repo.UserEvents.CheckPassword(ctx, userID, password, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) VerifyMfaOTP(ctx context.Context, authRequestID, userID string, code string, info *model.BrowserInfo) error {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, authRequestID)
	if err != nil {
		return err
	}
	if request.UserID != userID {
		return errors.ThrowPreconditionFailed(nil, "EVENT-ADJ26", "user id does not match request id")
	}
	return repo.UserEvents.CheckMfaOTP(ctx, userID, code, request.WithCurrentInfo(info))
}

func (repo *AuthRequestRepo) getAuthRequest(ctx context.Context, id string, checkLoggedIn bool) (*model.AuthRequest, error) {
	request, err := repo.AuthRequests.GetAuthRequestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	steps, err := repo.nextSteps(ctx, request, checkLoggedIn)
	if err != nil {
		return nil, err
	}
	request.PossibleSteps = steps
	return request, nil
}

func (repo *AuthRequestRepo) nextSteps(ctx context.Context, request *model.AuthRequest, checkLoggedIn bool) ([]model.NextStep, error) {
	if request == nil {
		return nil, errors.ThrowInvalidArgument(nil, "EVENT-ds27a", "request must not be nil")
	}
	steps := make([]model.NextStep, 0)
	if !checkLoggedIn && request.Prompt == model.PromptNone {
		return append(steps, &model.RedirectToCallbackStep{}), nil
	}
	if request.UserID == "" {
		steps = append(steps, &model.LoginStep{})
		if request.Prompt == model.PromptSelectAccount {
			users, err := repo.usersForUserSelection(request)
			if err != nil {
				return nil, err
			}
			steps = append(steps, &model.SelectUserStep{Users: users})
		}
		return steps, nil
	}
	user, err := userByID(ctx, repo.UserViewProvider, repo.UserEventProvider, request.UserID)
	if err != nil {
		return nil, err
	}
	userSession, err := userSessionByIDs(ctx, repo.UserSessionViewProvider, repo.UserEventProvider, request.AgentID, user)
	if err != nil {
		return nil, err
	}

	if user.InitRequired {
		return append(steps, &model.InitUserStep{PasswordSet: user.PasswordSet}), nil
	}
	if !user.PasswordSet {
		return append(steps, &model.InitPasswordStep{}), nil
	}

	if !checkVerificationTime(userSession.PasswordVerification, repo.PasswordCheckLifeTime) {
		return append(steps, &model.PasswordStep{}), nil
	}
	request.PasswordVerified = true
	request.AuthTime = userSession.PasswordVerification

	if step, ok := repo.mfaChecked(userSession, request, user); !ok {
		return append(steps, step), nil
	}

	if user.PasswordChangeRequired {
		steps = append(steps, &model.ChangePasswordStep{})
	}
	if !user.IsEmailVerified {
		steps = append(steps, &model.VerifyEMailStep{})
	}

	if user.PasswordChangeRequired || !user.IsEmailVerified {
		return steps, nil
	}

	//PLANNED: consent step
	return append(steps, &model.RedirectToCallbackStep{}), nil
}

func (repo *AuthRequestRepo) usersForUserSelection(request *model.AuthRequest) ([]model.UserSelection, error) {
	userSessions, err := userSessionsByUserAgentID(repo.UserSessionViewProvider, request.AgentID)
	if err != nil {
		return nil, err
	}
	users := make([]model.UserSelection, len(userSessions))
	for i, session := range userSessions {
		users[i] = model.UserSelection{
			UserID:           session.UserID,
			UserName:         session.UserName,
			UserSessionState: session.State,
		}
	}
	return users, nil
}

func (repo *AuthRequestRepo) mfaChecked(userSession *user_model.UserSessionView, request *model.AuthRequest, user *user_model.UserView) (model.NextStep, bool) {
	mfaLevel := request.MfaLevel()
	promptRequired := user.MfaMaxSetUp < mfaLevel
	if promptRequired || !repo.mfaSkippedOrSetUp(user) {
		return &model.MfaPromptStep{
			Required:     promptRequired,
			MfaProviders: user.MfaTypesSetupPossible(mfaLevel),
		}, false
	}
	switch mfaLevel {
	default:
		fallthrough
	case model.MfaLevelNotSetUp:
		if user.MfaMaxSetUp == model.MfaLevelNotSetUp {
			return nil, true
		}
		fallthrough
	case model.MfaLevelSoftware:
		if checkVerificationTime(userSession.MfaSoftwareVerification, repo.MfaSoftwareCheckLifeTime) {
			request.MfasVerified = append(request.MfasVerified, userSession.MfaSoftwareVerificationType)
			request.AuthTime = userSession.MfaSoftwareVerification
			return nil, true
		}
		fallthrough
	case model.MfaLevelHardware:
		if checkVerificationTime(userSession.MfaHardwareVerification, repo.MfaHardwareCheckLifeTime) {
			request.MfasVerified = append(request.MfasVerified, userSession.MfaHardwareVerificationType)
			request.AuthTime = userSession.MfaHardwareVerification
			return nil, true
		}
	}
	return &model.MfaVerificationStep{
		MfaProviders: user.MfaTypesAllowed(mfaLevel),
	}, false
}

func (repo *AuthRequestRepo) mfaSkippedOrSetUp(user *user_model.UserView) bool {
	if user.MfaMaxSetUp > model.MfaLevelNotSetUp {
		return true
	}
	return checkVerificationTime(user.MfaInitSkipped, repo.MfaInitSkippedLifeTime)
}

func checkVerificationTime(verificationTime time.Time, lifetime time.Duration) bool {
	return verificationTime.Add(lifetime).After(time.Now().UTC())
}

func userSessionsByUserAgentID(provider userSessionViewProvider, agentID string) ([]*user_model.UserSessionView, error) {
	session, err := provider.UserSessionsByAgentID(agentID)
	if err != nil {
		return nil, err
	}
	return view_model.UserSessionsToModel(session), nil
}

func userSessionByIDs(ctx context.Context, provider userSessionViewProvider, eventProvider userEventProvider, agentID string, user *user_model.UserView) (*user_model.UserSessionView, error) {
	session, err := provider.UserSessionByIDs(agentID, user.ID)
	if err != nil {
		if !errors.IsNotFound(err) {
			return nil, err
		}
		session = &view_model.UserSessionView{}
	}
	events, err := eventProvider.UserEventsByID(ctx, user.ID, session.Sequence)
	if err != nil {
		logging.Log("EVENT-Hse6s").WithError(err).Debug("error retrieving new events")
		return view_model.UserSessionToModel(session), nil
	}
	sessionCopy := *session
	for _, event := range events {
		switch event.Type {
		case es_model.UserPasswordCheckSucceeded,
			es_model.UserPasswordCheckFailed,
			es_model.MfaOtpCheckSucceeded,
			es_model.MfaOtpCheckFailed:
			eventData, err := view_model.UserSessionFromEvent(event)
			if err != nil {
				logging.Log("EVENT-sdgT3").WithError(err).Debug("error getting event data")
				return view_model.UserSessionToModel(session), nil
			}
			if eventData.UserAgentID != agentID {
				continue
			}
		}
		sessionCopy.AppendEvent(event)
	}
	return view_model.UserSessionToModel(&sessionCopy), nil
}

func userByID(ctx context.Context, viewProvider userViewProvider, eventProvider userEventProvider, userID string) (*user_model.UserView, error) {
	user, err := viewProvider.UserByID(userID)
	if err != nil {
		return nil, err
	}
	events, err := eventProvider.UserEventsByID(ctx, userID, user.Sequence)
	if err != nil {
		logging.Log("EVENT-dfg42").WithError(err).Debug("error retrieving new events")
		return view_model.UserToModel(user), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return view_model.UserToModel(user), nil
		}
	}
	return view_model.UserToModel(&userCopy), nil
}
