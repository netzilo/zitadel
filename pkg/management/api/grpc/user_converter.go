package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"
)

func userFromModel(user *usr_model.User) *User {
	creationDate, err := ptypes.TimestampProto(user.CreationDate)
	logging.Log("GRPC-8duwe").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(user.ChangeDate)
	logging.Log("GRPC-ckoe3d").OnError(err).Debug("unable to parse timestamp")

	converted := &User{
		Id:                user.AggregateID,
		State:             userStateFromModel(user.State),
		CreationDate:      creationDate,
		ChangeDate:        changeDate,
		Sequence:          user.Sequence,
		UserName:          user.UserName,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		DisplayName:       user.DisplayName,
		NickName:          user.NickName,
		PreferredLanguage: user.PreferredLanguage.String(),
		Gender:            genderFromModel(user.Gender),
	}
	if user.Email != nil {
		converted.Email = user.EmailAddress
		converted.IsEmailVerified = user.IsEmailVerified
	}
	if user.Phone != nil {
		converted.Phone = user.PhoneNumber
		converted.IsPhoneVerified = user.IsPhoneVerified
	}
	if user.Address != nil {
		converted.Country = user.Country
		converted.Locality = user.Locality
		converted.PostalCode = user.PostalCode
		converted.Region = user.Region
		converted.StreetAddress = user.StreetAddress
	}
	return converted
}

func userCreateToModel(u *CreateUserRequest) *usr_model.User {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-cK5k2").OnError(err).Debug("language malformed")

	user := &usr_model.User{
		Profile: &usr_model.Profile{
			UserName:          u.UserName,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
			NickName:          u.NickName,
			DisplayName:       u.DisplayName,
			PreferredLanguage: preferredLanguage,
			Gender:            genderToModel(u.Gender),
		},
		Email: &usr_model.Email{
			EmailAddress:    u.Email,
			IsEmailVerified: u.IsEmailVerified,
		},
		Address: &usr_model.Address{
			Country:       u.Country,
			Locality:      u.Locality,
			PostalCode:    u.PostalCode,
			Region:        u.Region,
			StreetAddress: u.StreetAddress,
		},
	}
	if u.Password != "" {
		user.Password = &usr_model.Password{SecretString: u.Password}
	}
	if u.Phone != "" {
		user.Phone = &usr_model.Phone{PhoneNumber: u.Phone, IsPhoneVerified: u.IsPhoneVerified}
	}
	return user
}

func passwordRequestToModel(r *PasswordRequest) *usr_model.Password {
	return &usr_model.Password{
		ObjectRoot:   models.ObjectRoot{AggregateID: r.Id},
		SecretString: r.Password,
	}
}

func profileFromModel(profile *usr_model.Profile) *UserProfile {
	creationDate, err := ptypes.TimestampProto(profile.CreationDate)
	logging.Log("GRPC-dkso3").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(profile.ChangeDate)
	logging.Log("GRPC-ski8d").OnError(err).Debug("unable to parse timestamp")

	converted := &UserProfile{
		Id:                profile.AggregateID,
		CreationDate:      creationDate,
		ChangeDate:        changeDate,
		Sequence:          profile.Sequence,
		UserName:          profile.UserName,
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		DisplayName:       profile.DisplayName,
		NickName:          profile.NickName,
		PreferredLanguage: profile.PreferredLanguage.String(),
		Gender:            genderFromModel(profile.Gender),
	}
	return converted
}

func updateProfileToModel(u *UpdateUserProfileRequest) *usr_model.Profile {
	preferredLanguage, err := language.Parse(u.PreferredLanguage)
	logging.Log("GRPC-d8k2s").OnError(err).Debug("language malformed")

	return &usr_model.Profile{
		ObjectRoot:        models.ObjectRoot{AggregateID: u.Id},
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		NickName:          u.NickName,
		DisplayName:       u.DisplayName,
		PreferredLanguage: preferredLanguage,
		Gender:            genderToModel(u.Gender),
	}
}

func notifyTypeToModel(state NotificationType) usr_model.NotificationType {
	switch state {
	case NotificationType_NOTIFICATIONTYPE_EMAIL:
		return usr_model.NOTIFICATIONTYPE_EMAIL
	case NotificationType_NOTIFICATIONTYPE_SMS:
		return usr_model.NOTIFICATIONTYPE_SMS
	default:
		return usr_model.NOTIFICATIONTYPE_EMAIL
	}
}

func userStateFromModel(state usr_model.UserState) UserState {
	switch state {
	case usr_model.USERSTATE_ACTIVE:
		return UserState_USERSTATE_ACTIVE
	case usr_model.USERSTATE_INACTIVE:
		return UserState_USERSTATE_INACTIVE
	case usr_model.USERSTATE_LOCKED:
		return UserState_USERSTATE_LOCKED
	default:
		return UserState_USERSTATE_UNSPECIFIED
	}
}

func genderFromModel(gender usr_model.Gender) Gender {
	switch gender {
	case usr_model.GENDER_FEMALE:
		return Gender_GENDER_FEMALE
	case usr_model.GENDER_MALE:
		return Gender_GENDER_MALE
	case usr_model.GENDER_DIVERSE:
		return Gender_GENDER_DIVERSE
	default:
		return Gender_GENDER_UNSPECIFIED
	}
}

func genderToModel(gender Gender) usr_model.Gender {
	switch gender {
	case Gender_GENDER_FEMALE:
		return usr_model.GENDER_FEMALE
	case Gender_GENDER_MALE:
		return usr_model.GENDER_MALE
	case Gender_GENDER_DIVERSE:
		return usr_model.GENDER_DIVERSE
	default:
		return usr_model.GENDER_UNDEFINED
	}
}
