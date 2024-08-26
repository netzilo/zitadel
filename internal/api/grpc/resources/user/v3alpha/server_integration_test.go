//go:build integration

package user_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	user "github.com/zitadel/zitadel/pkg/grpc/resources/user/v3alpha"
)

var (
	IAMOwnerCTX, SystemCTX context.Context
	Tester                 *integration.Tester
	Client                 user.ZITADELUsersClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, _, cancel := integration.Contexts(time.Hour)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		IAMOwnerCTX = Tester.WithAuthorization(ctx, integration.IAMOwner)
		SystemCTX = Tester.WithAuthorization(ctx, integration.SystemUser)
		Client = Tester.Client.UserV3Alpha
		return m.Run()
	}())
}

func ensureFeatureEnabled(t *testing.T, iamOwnerCTX context.Context) {
	f, err := Tester.Client.FeatureV2.GetInstanceFeatures(iamOwnerCTX, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)
	if f.UserSchema.GetEnabled() {
		return
	}
	_, err = Tester.Client.FeatureV2.SetInstanceFeatures(iamOwnerCTX, &feature.SetInstanceFeaturesRequest{
		UserSchema: gu.Ptr(true),
	})
	require.NoError(t, err)
	retryDuration := time.Minute
	if ctxDeadline, ok := iamOwnerCTX.Deadline(); ok {
		retryDuration = time.Until(ctxDeadline)
	}
	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			f, err := Tester.Client.FeatureV2.GetInstanceFeatures(iamOwnerCTX, &feature.GetInstanceFeaturesRequest{
				Inheritance: true,
			})
			require.NoError(ttt, err)
			if f.UserSchema.GetEnabled() {
				return
			}
		},
		retryDuration,
		100*time.Millisecond,
		"timed out waiting for ensuring instance feature")
}
