package repository

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"

	admin_model "github.com/caos/zitadel/internal/admin/model"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	SetUpOrg(context.Context, *admin_model.SetupOrg) (*admin_model.SetupOrg, error)
	IsOrgUnique(ctx context.Context, name, domain string) (bool, error)
	OrgByID(ctx context.Context, id string) (*org_model.Org, error)
	SearchOrgs(ctx context.Context, query *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error)

	GetOrgIamPolicyByID(ctx context.Context, id string) (*iam_model.OrgIAMPolicy, error)
	CreateOrgIamPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error)
	ChangeOrgIamPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error)
	RemoveOrgIamPolicy(ctx context.Context, id string) error
}
