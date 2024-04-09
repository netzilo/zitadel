//go:build integration

package execution_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/integration"
	execution "github.com/zitadel/zitadel/pkg/grpc/execution/v3alpha"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

func executionTargetsSingleTarget(id string) []*execution.ExecutionTargetType {
	return []*execution.ExecutionTargetType{{Type: &execution.ExecutionTargetType_Target{Target: id}}}
}

func executionTargetsSingleInclude(include string) []*execution.ExecutionTargetType {
	return []*execution.ExecutionTargetType{{Type: &execution.ExecutionTargetType_Include{Include: include}}}
}

func TestServer_SetExecution_Request(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *execution.SetExecutionRequest
		want    *execution.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.NotExistingService/List",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_All{
								All: true,
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_SetExecution_Request_Include(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	executionCond := "request"
	Tester.SetExecution(CTX, t,
		&execution.Condition{
			ConditionType: &execution.Condition_Request{
				Request: &execution.RequestExecution{
					Condition: &execution.RequestExecution_All{
						All: true,
					},
				},
			},
		},
		executionTargetsSingleTarget(targetResp.GetId()),
	)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *execution.SetExecutionRequest
		want    *execution.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "method, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Targets: executionTargetsSingleInclude(executionCond),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Targets: executionTargetsSingleInclude(executionCond),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_All{
								All: true,
							},
						},
					},
				},
				Targets: executionTargetsSingleInclude(executionCond),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_DeleteExecution_Request(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, request *execution.DeleteExecutionRequest) error
		req     *execution.DeleteExecutionRequest
		want    *execution.DeleteExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/NotExisting",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/GetSession",
							},
						},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_Service{
								Service: "zitadel.user.v2beta.UserService",
							},
						},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Request{
						Request: &execution.RequestExecution{
							Condition: &execution.RequestExecution_All{
								All: true,
							},
						},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}

			got, err := Client.DeleteExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_SetExecution_Response(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *execution.SetExecutionRequest
		want    *execution.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.NotExistingService/List",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_All{
								All: true,
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_DeleteExecution_Response(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, request *execution.DeleteExecutionRequest) error
		req     *execution.DeleteExecutionRequest
		want    *execution.DeleteExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/NotExisting",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/GetSession",
							},
						},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_Service{
								Service: "zitadel.user.v2beta.UserService",
							},
						},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_All{
								All: true,
							},
						},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}

			got, err := Client.DeleteExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_SetExecution_Event(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *execution.SetExecutionRequest
		want    *execution.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		/*
			//TODO event existing check

			{
				name: "event, not existing",
				ctx:  CTX,
				req: &execution.SetExecutionRequest{
					Condition: &execution.Condition{
						ConditionType: &execution.Condition_Event{
							Event: &execution.EventExecution{
								Condition: &execution.EventExecution_Event{
									Event: "xxx",
								},
							},
						},
					},
					Targets: []string{targetResp.GetId()},
				},
				wantErr: true,
			},
		*/
		{
			name: "event, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_Event{
								Event: "xxx",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		/*
			// TODO:

			{
				name: "group, not existing",
				ctx:  CTX,
				req: &execution.SetExecutionRequest{
					Condition: &execution.Condition{
						ConditionType: &execution.Condition_Event{
							Event: &execution.EventExecution{
								Condition: &execution.EventExecution_Group{
									Group: "xxx",
								},
							},
						},
					},
					Targets: []string{targetResp.GetId()},
				},
				wantErr: true,
			},
		*/
		{
			name: "group, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_Group{
								Group: "xxx",
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_All{
								All: true,
							},
						},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_DeleteExecution_Event(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, request *execution.DeleteExecutionRequest) error
		req     *execution.DeleteExecutionRequest
		want    *execution.DeleteExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{},
					},
				},
			},
			wantErr: true,
		},
		/*
			//TODO: add when check is implemented
			{
				name: "event, not existing",
				ctx:  CTX,
				req: &execution.DeleteExecutionRequest{
					Condition: &execution.Condition{
						ConditionType: &execution.Condition_Event{
							Event: &execution.EventExecution{
								Condition: &execution.EventExecution_Event{
									Event: "xxx",
								},
							},
						},
					},
				},
				wantErr: true,
			},
		*/
		{
			name: "event, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_Event{
								Event: "xxx",
							},
						},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "group, not existing",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_Group{
								Group: "xxx",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "group, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_Group{
								Group: "xxx",
							},
						},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, not existing",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "all, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Event{
						Event: &execution.EventExecution{
							Condition: &execution.EventExecution_All{
								All: true,
							},
						},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}

			got, err := Client.DeleteExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_SetExecution_Function(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *execution.SetExecutionRequest
		want    *execution.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		{
			name: "function, not existing",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Function{
						Function: &execution.FunctionExecution{Name: "xxx"},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			wantErr: true,
		},
		{
			name: "function, ok",
			ctx:  CTX,
			req: &execution.SetExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Function{
						Function: &execution.FunctionExecution{Name: "Action.Flow.Type.ExternalAuthentication.Action.TriggerType.PostAuthentication"},
					},
				},
				Targets: executionTargetsSingleTarget(targetResp.GetId()),
			},
			want: &execution.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)

			// cleanup to not impact other requests
			Tester.DeleteExecution(tt.ctx, t, tt.req.GetCondition())
		})
	}
}

func TestServer_DeleteExecution_Function(t *testing.T) {
	targetResp := Tester.CreateTarget(CTX, t, "", "https://notexisting", domain.TargetTypeWebhook, false)

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, request *execution.DeleteExecutionRequest) error
		req     *execution.DeleteExecutionRequest
		want    *execution.DeleteExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{
							Condition: &execution.ResponseExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Response{
						Response: &execution.ResponseExecution{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "function, not existing",
			ctx:  CTX,
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Function{
						Function: &execution.FunctionExecution{Name: "xxx"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "function, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *execution.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), executionTargetsSingleTarget(targetResp.GetId()))
				return nil
			},
			req: &execution.DeleteExecutionRequest{
				Condition: &execution.Condition{
					ConditionType: &execution.Condition_Function{
						Function: &execution.FunctionExecution{Name: "Action.Flow.Type.ExternalAuthentication.Action.TriggerType.PostAuthentication"},
					},
				},
			},
			want: &execution.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}

			got, err := Client.DeleteExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}
