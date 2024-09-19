//go:build integration

package oidc_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/profile"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	oidc_api "github.com/zitadel/zitadel/v2/internal/api/oidc"
	"github.com/zitadel/zitadel/v2/internal/domain"
)

func TestServer_JWTProfile(t *testing.T) {
	t.Parallel()

	user, name, keyData, err := Instance.CreateOIDCJWTProfileClient(CTX)
	require.NoError(t, err)

	type claims struct {
		name                       string
		username                   string
		updated                    time.Time
		resourceOwnerID            any
		resourceOwnerName          any
		resourceOwnerPrimaryDomain any
		orgDomain                  any
	}
	tests := []struct {
		name       string
		keyData    []byte
		scope      []string
		wantClaims claims
		wantErr    bool
	}{
		{
			name:    "success",
			keyData: keyData,
			scope:   []string{oidc.ScopeOpenID},
		},
		{
			name:    "openid, profile, email",
			keyData: keyData,
			scope:   []string{oidc.ScopeOpenID, oidc.ScopeProfile, oidc.ScopeEmail},
			wantClaims: claims{
				name:     name,
				username: name,
				updated:  user.GetDetails().GetChangeDate().AsTime(),
			},
		},
		{
			name:    "org id and domain scope",
			keyData: keyData,
			scope: []string{
				oidc.ScopeOpenID,
				domain.OrgIDScope + Instance.DefaultOrg.Id,
				domain.OrgDomainPrimaryScope + Instance.DefaultOrg.PrimaryDomain,
			},
			wantClaims: claims{
				resourceOwnerID:            Instance.DefaultOrg.Id,
				resourceOwnerName:          Instance.DefaultOrg.Name,
				resourceOwnerPrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
				orgDomain:                  Instance.DefaultOrg.PrimaryDomain,
			},
		},
		{
			name:    "invalid org domain filtered",
			keyData: keyData,
			scope: []string{
				oidc.ScopeOpenID,
				domain.OrgDomainPrimaryScope + Instance.DefaultOrg.PrimaryDomain,
				domain.OrgDomainPrimaryScope + "foo"},
			wantClaims: claims{
				orgDomain: Instance.DefaultOrg.PrimaryDomain,
			},
		},
		{
			name:    "invalid org id filtered",
			keyData: keyData,
			scope: []string{oidc.ScopeOpenID,
				domain.OrgIDScope + Instance.DefaultOrg.Id,
				domain.OrgIDScope + "foo",
			},
			wantClaims: claims{
				resourceOwnerID:            Instance.DefaultOrg.Id,
				resourceOwnerName:          Instance.DefaultOrg.Name,
				resourceOwnerPrimaryDomain: Instance.DefaultOrg.PrimaryDomain,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenSource, err := profile.NewJWTProfileTokenSourceFromKeyFileData(CTX, Instance.OIDCIssuer(), tt.keyData, tt.scope)
			require.NoError(t, err)

			tokens, err := tokenSource.TokenCtx(CTX)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, tokens)

			provider, err := rp.NewRelyingPartyOIDC(CTX, Instance.OIDCIssuer(), "", "", redirectURI, tt.scope)
			require.NoError(t, err)
			userinfo, err := rp.Userinfo[*oidc.UserInfo](CTX, tokens.AccessToken, oidc.BearerToken, user.GetUserId(), provider)
			require.NoError(t, err)
			assert.Equal(t, tt.wantClaims.resourceOwnerID, userinfo.Claims[oidc_api.ClaimResourceOwnerID])
			assert.Equal(t, tt.wantClaims.resourceOwnerName, userinfo.Claims[oidc_api.ClaimResourceOwnerName])
			assert.Equal(t, tt.wantClaims.resourceOwnerPrimaryDomain, userinfo.Claims[oidc_api.ClaimResourceOwnerPrimaryDomain])
			assert.Equal(t, tt.wantClaims.orgDomain, userinfo.Claims[domain.OrgDomainPrimaryClaim])
			assert.Equal(t, tt.wantClaims.name, userinfo.Name)
			assert.Equal(t, tt.wantClaims.username, userinfo.PreferredUsername)
			assertOIDCTime(t, userinfo.UpdatedAt, tt.wantClaims.updated)
			assert.Empty(t, userinfo.UserInfoProfile.FamilyName)
			assert.Empty(t, userinfo.UserInfoProfile.GivenName)
			assert.Empty(t, userinfo.UserInfoEmail)
			assert.Empty(t, userinfo.UserInfoPhone)
			assert.Empty(t, userinfo.Address)
		})
	}
}
