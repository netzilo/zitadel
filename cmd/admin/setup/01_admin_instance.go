package setup

import (
	"context"

	command "github.com/caos/zitadel/internal/command/v2"
)

type DefaultInstance struct {
	cmd           *command.Command
	InstanceSetup command.InstanceSetup
}

func (mig *DefaultInstance) Execute(ctx context.Context) error {
	_, err := mig.cmd.SetUpTenant(ctx, &mig.InstanceSetup)

	return err
}

func (mig *DefaultInstance) String() string {
	return "01_default_tenant"
}
