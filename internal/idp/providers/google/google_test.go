package google

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/idp"
	oidc2 "github.com/zitadel/zitadel/internal/idp/providers/oidc"
)

func TestProvider_BeginAuth(t *testing.T) {
	type fields struct {
		clientID     string
		clientSecret string
		redirectURI  string
	}
	tests := []struct {
		name   string
		fields fields
		want   idp.Session
	}{
		{
			name: "successful auth",
			fields: fields{
				clientID:     "clientID",
				clientSecret: "clientSecret",
				redirectURI:  "redirectURI",
			},
			want: &oidc2.Session{
				AuthURL: "https://accounts.google.com/o/oauth2/v2/auth?client_id=clientID&redirect_uri=redirectURI&response_type=code&scope=openid&state=testState",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			provider, err := New(tt.fields.clientID, tt.fields.clientSecret, tt.fields.redirectURI)
			a.NoError(err)

			session, err := provider.BeginAuth("testState")
			a.NoError(err)

			//authUrl, err := url.Parse(session.GetAuthURL())
			//a.NoError(err)
			//
			a.Equal(tt.want.GetAuthURL(), session.GetAuthURL())
			//a.Equal("/authorize", authUrl.Path)
			//a.Equal("clientID", authUrl.Query().Get("client_id"))
			//a.Equal("testState", authUrl.Query().Get("state"))
			//a.Equal("redirectURI", authUrl.Query().Get("redirect_uri"))
			//a.Equal("openid", authUrl.Query().Get("scope"))
			//
			//if !tt.wantErr(t, err, fmt.Sprintf("BeginAuth(%v)", tt.fields.state)) {
			//	return
			//}
			//assert.Equalf(t, tt.want, got, "BeginAuth(%v)", tt.args.state)
		})
	}
}
