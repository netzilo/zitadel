package user

import (
	"context"

	"github.com/muhlemmer/gu"

	resource_object "github.com/zitadel/zitadel/internal/api/grpc/resources/object/v3alpha"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v3alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

func (s *Server) SetContactEmail(ctx context.Context, req *user.SetContactEmailRequest) (_ *user.SetContactEmailResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	schemauser := setContactEmailRequestToChangeSchemaUserEmail(ctx, req)
	if err := s.command.ChangeSchemaUserEmail(ctx, schemauser, s.userCodeAlg); err != nil {
		return nil, err
	}
	return &user.SetContactEmailResponse{
		Details:          resource_object.DomainToDetailsPb(schemauser.Details, object.OwnerType_OWNER_TYPE_ORG, schemauser.Details.ResourceOwner),
		VerificationCode: gu.Ptr(schemauser.ReturnCode),
	}, nil
}

func setContactEmailRequestToChangeSchemaUserEmail(ctx context.Context, req *user.SetContactEmailRequest) *command.ChangeSchemaUserEmail {
	return &command.ChangeSchemaUserEmail{
		ResourceOwner: organizationToCreateResourceOwner(ctx, req.Organization),
		ID:            req.GetId(),
		Email:         setEmailToEmail(req.Email),
	}
}

func setEmailToEmail(setEmail *user.SetEmail) *command.Email {
	if setEmail == nil {
		return nil
	}
	email := &command.Email{
		Address: domain.EmailAddress(setEmail.Address),
	}
	if setEmail.GetIsVerified() {
		email.Verified = true
	}
	if setEmail.GetReturnCode() != nil {
		email.ReturnCode = true
	}
	if setEmail.GetSendCode() != nil {
		email.URLTemplate = setEmail.GetSendCode().GetUrlTemplate()
	}
	return email
}

func (s *Server) VerifyContactEmail(ctx context.Context, req *user.VerifyContactEmailRequest) (_ *user.VerifyContactEmailResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	details, err := s.command.VerifySchemaUserEmail(ctx, organizationToUpdateResourceOwner(req.Organization), req.GetId(), req.GetVerificationCode(), s.userCodeAlg)
	if err != nil {
		return nil, err
	}
	return &user.VerifyContactEmailResponse{
		Details: resource_object.DomainToDetailsPb(details, object.OwnerType_OWNER_TYPE_ORG, details.ResourceOwner),
	}, nil
}

func (s *Server) ResendContactEmailCode(ctx context.Context, req *user.ResendContactEmailCodeRequest) (_ *user.ResendContactEmailCodeResponse, err error) {
	if err := checkUserSchemaEnabled(ctx); err != nil {
		return nil, err
	}
	schemauser := resendContactEmailCodeRequestToResendSchemaUserEmailCode(ctx, req)
	if err = s.command.ResendSchemaUserEmailCode(ctx, schemauser, s.userCodeAlg); err != nil {
		return nil, err
	}
	return &user.ResendContactEmailCodeResponse{
		Details: resource_object.DomainToDetailsPb(schemauser.Details, object.OwnerType_OWNER_TYPE_ORG, schemauser.Details.ResourceOwner),
	}, nil
}

func resendContactEmailCodeRequestToResendSchemaUserEmailCode(ctx context.Context, req *user.ResendContactEmailCodeRequest) *command.ResendSchemaUserEmailCode {
	return &command.ResendSchemaUserEmailCode{
		ResourceOwner: organizationToCreateResourceOwner(ctx, req.Organization),
		ID:            req.GetId(),
		URLTemplate:   req.GetSendCode().GetUrlTemplate(),
		ReturnCode:    req.GetReturnCode() != nil,
	}
}
