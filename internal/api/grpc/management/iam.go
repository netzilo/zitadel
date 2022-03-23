package management

import (
	"context"

	mgmt_pb "github.com/caos/zitadel/pkg/grpc/management"
)

func (s *Server) GetIAM(ctx context.Context, _ *mgmt_pb.GetIAMRequest) (*mgmt_pb.GetIAMResponse, error) {
	iam, err := s.query.IAM(ctx)
	if err != nil {
		return nil, err
	}
	return &mgmt_pb.GetIAMResponse{
		GlobalOrgId:  iam.GlobalOrgID,
		IamProjectId: iam.IAMProjectID,
	}, nil
}
