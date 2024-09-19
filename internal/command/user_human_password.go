package command

import (
	"context"
	"errors"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/passwap"

	"github.com/zitadel/zitadel/v2/internal/crypto"
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/repository/user"
	"github.com/zitadel/zitadel/v2/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

var (
	ErrPasswordInvalid = func(err error) error {
		return zerrors.ThrowInvalidArgument(err, "COMMAND-3M0fs", "Errors.User.Password.Invalid")
	}
	ErrPasswordUnchanged = func(err error) error {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-Aesh5", "Errors.User.Password.NotChanged")
	}
)

func (c *Commands) SetPassword(ctx context.Context, orgID, userID, password string, oneTime bool) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.IDMissing")
	}
	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	return c.setPassword(
		ctx,
		wm,
		password,
		"", // current api implementations never provide an encoded password
		"",
		oneTime,
		c.setPasswordWithPermission(wm.AggregateID, wm.ResourceOwner),
	)
}

func (c *Commands) SetPasswordWithVerifyCode(ctx context.Context, orgID, userID, code, password, userAgentID string, changeRequired bool) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M9fs", "Errors.IDMissing")
	}
	if password == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-Mf0sd", "Errors.User.Password.Empty")
	}
	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	return c.setPassword(
		ctx,
		wm,
		password,
		"",
		userAgentID,
		changeRequired,
		c.setPasswordWithVerifyCode(wm.CodeCreationDate, wm.CodeExpiry, wm.Code, code),
	)
}

// ChangePassword change password of existing user
func (c *Commands) ChangePassword(ctx context.Context, orgID, userID, oldPassword, newPassword, userAgentID string, changeRequired bool) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.IDMissing")
	}
	if oldPassword == "" || newPassword == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Empty")
	}
	wm, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return nil, err
	}
	return c.setPassword(
		ctx,
		wm,
		newPassword,
		"",
		userAgentID,
		changeRequired,
		c.checkCurrentPassword(newPassword, "", oldPassword, wm.EncodedHash),
	)
}

type setPasswordVerification func(ctx context.Context) (newEncodedPassword string, err error)

// setPasswordWithPermission returns a permission check as [setPasswordVerification] implementation
func (c *Commands) setPasswordWithPermission(userID, orgID string) setPasswordVerification {
	return func(ctx context.Context) (_ string, err error) {
		return "", c.checkPermission(ctx, domain.PermissionUserWrite, orgID, userID)
	}
}

// setPasswordWithVerifyCode returns a password code check as [setPasswordVerification] implementation
func (c *Commands) setPasswordWithVerifyCode(
	passwordCodeCreationDate time.Time,
	passwordCodeExpiry time.Duration,
	passwordCode *crypto.CryptoValue,
	code string,
) setPasswordVerification {
	return func(ctx context.Context) (_ string, err error) {
		if passwordCode == nil {
			return "", zerrors.ThrowPreconditionFailed(nil, "COMMAND-2M9fs", "Errors.User.Code.NotFound")
		}
		_, spanCrypto := tracing.NewNamedSpan(ctx, "crypto.VerifyCode")
		defer func() {
			spanCrypto.EndWithError(err)
		}()
		return "", crypto.VerifyCode(passwordCodeCreationDate, passwordCodeExpiry, passwordCode, code, c.userEncryption)
	}
}

// checkCurrentPassword returns a password check as [setPasswordVerification] implementation
func (c *Commands) checkCurrentPassword(
	newPassword, newEncodedPassword, currentPassword, currentEncodePassword string,
) setPasswordVerification {
	// in case the new password is already encoded, we only need to verify the current
	if newEncodedPassword != "" {
		return func(ctx context.Context) (_ string, err error) {
			_, spanPasswap := tracing.NewNamedSpan(ctx, "passwap.Verify")
			_, err = c.userPasswordHasher.Verify(currentEncodePassword, currentPassword)
			spanPasswap.EndWithError(err)
			return "", convertPasswapErr(err)
		}
	}

	// otherwise let's directly verify and return the new generate hash, so we can reuse it in the event
	return func(ctx context.Context) (string, error) {
		return c.verifyAndUpdatePassword(ctx, currentEncodePassword, currentPassword, newPassword)
	}
}

// setPassword directly pushes the intent of [setPasswordCommand] to the eventstore and returns the [domain.ObjectDetails]
func (c *Commands) setPassword(
	ctx context.Context,
	wm *HumanPasswordWriteModel,
	password, encodedPassword, userAgentID string,
	changeRequired bool,
	verificationCheck setPasswordVerification,
) (*domain.ObjectDetails, error) {
	agg := user.NewAggregate(wm.AggregateID, wm.ResourceOwner)
	command, err := c.setPasswordCommand(ctx, &agg.Aggregate, wm.UserState, password, encodedPassword, userAgentID, changeRequired, verificationCheck)
	if err != nil {
		return nil, err
	}
	err = c.pushAppendAndReduce(ctx, wm, command)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

// setPasswordCommand creates the command / intent for changing a user's password.
// It will check the user's [domain.UserState] to be existing and not initial,
// if the caller is allowed to change the password (permission, by code or by providing the current password),
// and it will ensure the new password (if provided as plain) corresponds to the password complexity policy.
// If not already encoded, the new password will be hashed.
func (c *Commands) setPasswordCommand(ctx context.Context, agg *eventstore.Aggregate, userState domain.UserState, password, encodedPassword, userAgentID string, changeRequired bool, verificationCheck setPasswordVerification) (_ eventstore.Command, err error) {
	if !isUserStateExists(userState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-G8dh3", "Errors.User.Password.NotFound")
	}
	if isUserStateInitial(userState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-M9dse", "Errors.User.NotInitialised")
	}
	if verificationCheck != nil {
		newEncodedPassword, err := verificationCheck(ctx)
		if err != nil {
			return nil, err
		}
		// use the new hash from the verification in case there is one (e.g. existing pw check)
		if newEncodedPassword != "" {
			encodedPassword = newEncodedPassword
		}
	}
	// If password is provided, let's check if is compliant with the policy.
	// If only a encodedPassword is passed, we can skip this.
	if password != "" {
		if err = c.checkPasswordComplexity(ctx, password, agg.ResourceOwner); err != nil {
			return nil, err
		}
	}

	// In case only a plain password was passed, we need to hash it.
	if encodedPassword == "" {
		_, span := tracing.NewNamedSpan(ctx, "passwap.Hash")
		encodedPassword, err = c.userPasswordHasher.Hash(password)
		span.EndWithError(err)
		if err = convertPasswapErr(err); err != nil {
			return nil, err
		}
	}
	return user.NewHumanPasswordChangedEvent(ctx, agg, encodedPassword, changeRequired, userAgentID), nil
}

// verifyAndUpdatePassword verify if the old password is correct with the encoded hash and
// returns the hash of the new password if so
func (c *Commands) verifyAndUpdatePassword(ctx context.Context, encodedHash, oldPassword, newPassword string) (string, error) {
	if encodedHash == "" {
		return "", zerrors.ThrowPreconditionFailed(nil, "COMMAND-Fds3s", "Errors.User.Password.NotSet")
	}

	_, spanPasswap := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := c.userPasswordHasher.VerifyAndUpdate(encodedHash, oldPassword, newPassword)
	spanPasswap.EndWithError(err)
	return updated, convertPasswapErr(err)
}

// checkPasswordComplexity checks uf the given password can be used to be the password of a user
func (c *Commands) checkPasswordComplexity(ctx context.Context, newPassword string, resourceOwner string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	policy, err := c.getOrgPasswordComplexityPolicy(ctx, resourceOwner)
	if err != nil {
		return err
	}

	if err := policy.Check(newPassword); err != nil {
		return err
	}
	return nil
}

// RequestSetPassword generate and send out new code to change password for a specific user
func (c *Commands) RequestSetPassword(ctx context.Context, userID, resourceOwner string, notifyType domain.NotificationType, authRequestID string) (objectDetails *domain.ObjectDetails, err error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-M00oL", "Errors.User.UserIDMissing")
	}

	existingHuman, err := c.userWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if !isUserStateExists(existingHuman.UserState) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Hj9ds", "Errors.User.NotFound")
	}
	if existingHuman.UserState == domain.UserStateInitial {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-2M9sd", "Errors.User.NotInitialised")
	}
	userAgg := UserAggregateFromWriteModel(&existingHuman.WriteModel)
	passwordCode, err := c.newEncryptedCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordResetCode, c.userEncryption) //nolint:staticcheck
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, user.NewHumanPasswordCodeAddedEvent(ctx, userAgg, passwordCode.Crypted, passwordCode.Expiry, notifyType, authRequestID))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingHuman, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingHuman.WriteModel), nil
}

// PasswordCodeSent notification send with code to change password
func (c *Commands) PasswordCodeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-meEfe", "Errors.User.UserIDMissing")
	}

	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanPasswordCodeSentEvent(ctx, userAgg))
	return err
}

// PasswordChangeSent notification sent that user changed password
func (c *Commands) PasswordChangeSent(ctx context.Context, orgID, userID string) (err error) {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-pqlm2n", "Errors.User.UserIDMissing")
	}

	existingPassword, err := c.passwordWriteModel(ctx, userID, orgID)
	if err != nil {
		return err
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return zerrors.ThrowPreconditionFailed(nil, "COMMAND-x902b2v", "Errors.User.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingPassword.WriteModel)
	_, err = c.eventstore.Push(ctx, user.NewHumanPasswordChangeSentEvent(ctx, userAgg))
	return err
}

// HumanCheckPassword check password for user with additional information from authRequest
func (c *Commands) HumanCheckPassword(ctx context.Context, orgID, userID, password string, authRequest *domain.AuthRequest) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-4Mfsf", "Errors.User.UserIDMissing")
	}
	if password == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-3n8fs", "Errors.User.Password.Empty")
	}

	loginPolicy, err := c.getOrgLoginPolicy(ctx, orgID)
	if err != nil {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-Edf3g", "Errors.Org.LoginPolicy.NotFound")
	}
	if !loginPolicy.AllowUsernamePassword {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-Dft32", "Errors.Org.LoginPolicy.UsernamePasswordNotAllowed")
	}
	commands, err := checkPassword(ctx, userID, password, c.eventstore, c.userPasswordHasher, authRequestDomainToAuthRequestInfo(authRequest))
	if len(commands) == 0 {
		return err
	}
	_, pushErr := c.eventstore.Push(ctx, commands...)
	logging.OnError(pushErr).Error("error create password check failed event")
	return err
}

func checkPassword(ctx context.Context, userID, password string, es *eventstore.Eventstore, hasher *crypto.Hasher, optionalAuthRequestInfo *user.AuthRequestInfo) ([]eventstore.Command, error) {
	if userID == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-Sfw3f", "Errors.User.UserIDMissing")
	}
	wm := NewHumanPasswordWriteModel(userID, "")
	err := es.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	if !wm.UserState.Exists() {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-3n77z", "Errors.User.NotFound")
	}
	if wm.UserState == domain.UserStateLocked {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-JLK35", "Errors.User.Locked")
	}
	if wm.EncodedHash == "" {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-3nJ4t", "Errors.User.Password.NotSet")
	}

	userAgg := UserAggregateFromWriteModel(&wm.WriteModel)
	ctx, spanPasswordComparison := tracing.NewNamedSpan(ctx, "passwap.Verify")
	updated, err := hasher.Verify(wm.EncodedHash, password)
	spanPasswordComparison.EndWithError(err)
	err = convertPasswapErr(err)
	commands := make([]eventstore.Command, 0, 2)

	// recheck for additional events (failed password checks or locks)
	recheckErr := es.FilterToQueryReducer(ctx, wm)
	if recheckErr != nil {
		return nil, recheckErr
	}
	if wm.UserState == domain.UserStateLocked {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMAND-SFA3t", "Errors.User.Locked")
	}

	if err == nil {
		commands = append(commands, user.NewHumanPasswordCheckSucceededEvent(ctx, userAgg, optionalAuthRequestInfo))
		if updated != "" {
			commands = append(commands, user.NewHumanPasswordHashUpdatedEvent(ctx, userAgg, updated))
		}
		return commands, nil
	}

	commands = append(commands, user.NewHumanPasswordCheckFailedEvent(ctx, userAgg, optionalAuthRequestInfo))

	lockoutPolicy, lockoutErr := getLockoutPolicy(ctx, wm.ResourceOwner, es.FilterToQueryReducer)
	logging.OnError(lockoutErr).Error("unable to get lockout policy")
	if lockoutPolicy != nil && lockoutPolicy.MaxPasswordAttempts > 0 && wm.PasswordCheckFailedCount+1 >= lockoutPolicy.MaxPasswordAttempts {
		commands = append(commands, user.NewUserLockedEvent(ctx, userAgg))
	}
	return commands, err
}

func (c *Commands) passwordWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *HumanPasswordWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanPasswordWriteModel(userID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func convertPasswapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, passwap.ErrPasswordMismatch) {
		return ErrPasswordInvalid(err)
	}
	if errors.Is(err, passwap.ErrPasswordNoChange) {
		return ErrPasswordUnchanged(err)
	}
	return zerrors.ThrowInternal(err, "COMMAND-CahN2", "Errors.Internal")
}
