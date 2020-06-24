package grpc

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/golang/protobuf/ptypes/empty"
)

func (s *Server) GetMyOrg(ctx context.Context, _ *empty.Empty) (*OrgView, error) {
	org, err := s.org.OrgByID(ctx, auth.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return orgViewFromModel(org), nil
}

func (s *Server) GetOrgByDomainGlobal(ctx context.Context, in *Domain) (*OrgView, error) {
	org, err := s.org.OrgByDomainGlobal(ctx, in.Domain)
	if err != nil {
		return nil, err
	}
	return orgViewFromModel(org), nil
}

func (s *Server) DeactivateMyOrg(ctx context.Context, _ *empty.Empty) (*Org, error) {
	org, err := s.org.DeactivateOrg(ctx, auth.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) ReactivateMyOrg(ctx context.Context, _ *empty.Empty) (*Org, error) {
	org, err := s.org.ReactivateOrg(ctx, auth.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return orgFromModel(org), nil
}

func (s *Server) SearchMyOrgDomains(ctx context.Context, in *OrgDomainSearchRequest) (*OrgDomainSearchResponse, error) {
	domains, err := s.org.SearchMyOrgDomains(ctx, orgDomainSearchRequestToModel(in))
	if err != nil {
		return nil, err
	}
	return orgDomainSearchResponseFromModel(domains), nil
}

func (s *Server) AddMyOrgDomain(ctx context.Context, in *AddOrgDomainRequest) (*OrgDomain, error) {
	domain, err := s.org.AddMyOrgDomain(ctx, addOrgDomainToModel(in))
	if err != nil {
		return nil, err
	}
	return orgDomainFromModel(domain), nil
}

func (s *Server) RemoveMyOrgDomain(ctx context.Context, in *RemoveOrgDomainRequest) (*empty.Empty, error) {
	err := s.org.RemoveMyOrgDomain(ctx, in.Domain)
	return &empty.Empty{}, err
}

func (s *Server) OrgChanges(ctx context.Context, changesRequest *ChangeRequest) (*Changes, error) {
	response, err := s.org.OrgChanges(ctx, changesRequest.Id, changesRequest.SequenceOffset, changesRequest.Limit, changesRequest.Asc)
	if err != nil {
		return nil, err
	}
	return orgChangesToResponse(response, changesRequest.GetSequenceOffset(), changesRequest.GetLimit()), nil
}

func (s *Server) GetMyOrgIamPolicy(ctx context.Context, _ *empty.Empty) (_ *OrgIamPolicy, err error) {
	policy, err := s.org.GetMyOrgIamPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return orgIamPolicyFromModel(policy), err
}
