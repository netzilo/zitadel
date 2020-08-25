package model

import (
	"testing"

	user_model "github.com/caos/zitadel/internal/user/model"
	"golang.org/x/text/language"
)

func TestProfileChanges(t *testing.T) {
	type args struct {
		existingProfile *Profile
		newProfile      *Profile
	}
	type res struct {
		changesLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "all attributes changed",
			args: args{
				existingProfile: &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GenderFemale)},
				newProfile:      &Profile{UserName: "UserNameChanged", FirstName: "FirstNameChanged", LastName: "LastNameChanged", NickName: "NickNameChanged", DisplayName: "DisplayNameChanged", PreferredLanguage: language.English, Gender: int32(user_model.GenderMale)},
			},
			res: res{
				changesLen: 6,
			},
		},
		{
			name: "no changes",
			args: args{
				existingProfile: &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GenderFemale)},
				newProfile:      &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GenderFemale)},
			},
			res: res{
				changesLen: 0,
			},
		},
		{
			name: "username changed",
			args: args{
				existingProfile: &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GenderFemale)},
				newProfile:      &Profile{UserName: "UserNameChanged", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GenderFemale)},
			},
			res: res{
				changesLen: 0,
			},
		},
		{
			name: "empty names",
			args: args{
				existingProfile: &Profile{UserName: "UserName", FirstName: "FirstName", LastName: "LastName", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GenderFemale)},
				newProfile:      &Profile{UserName: "UserName", FirstName: "", LastName: "", NickName: "NickName", DisplayName: "DisplayName", PreferredLanguage: language.German, Gender: int32(user_model.GenderFemale)},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existingProfile.Changes(tt.args.newProfile)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}
