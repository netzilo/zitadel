//go:build integration

package action_test

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
	action "github.com/zitadel/zitadel/pkg/grpc/resources/action/v3alpha"
)

var (
	CTX      context.Context
	Instance *integration.Instance
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		Instance = integration.GetInstance(ctx)

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		return m.Run()
	}())
}

func ensureFeatureEnabled(t *testing.T, instance *integration.Instance) {
	ctx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	f, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
		Inheritance: true,
	})
	require.NoError(t, err)
	if f.Actions.GetEnabled() {
		return
	}
	_, err = instance.Client.FeatureV2.SetInstanceFeatures(ctx, &feature.SetInstanceFeaturesRequest{
		Actions: gu.Ptr(true),
	})
	require.NoError(t, err)
	retryDuration := time.Minute
	if ctxDeadline, ok := ctx.Deadline(); ok {
		retryDuration = time.Until(ctxDeadline)
	}
	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			f, err := instance.Client.FeatureV2.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
				Inheritance: true,
			})
			assert.NoError(ttt, err)
			assert.True(ttt, f.Actions.GetEnabled())
		},
		retryDuration,
		100*time.Millisecond,
		"timed out waiting for ensuring instance feature")

	require.EventuallyWithT(t,
		func(ttt *assert.CollectT) {
			_, err := instance.Client.ActionV3.ListExecutionMethods(ctx, &action.ListExecutionMethodsRequest{})
			assert.NoError(ttt, err)
		},
		retryDuration,
		100*time.Millisecond,
		"timed out waiting for ensuring instance feature call")
}
