package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/api/auth"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/project/model"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestProjectByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es      *ProjectEventstore
		project *model.Project
	}
	type res struct {
		project *model.Project
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "project from events, ok",
			args: args{
				es:      GetMockProjectByIDOK(ctrl),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
			},
		},
		{
			name: "project from events, no events",
			args: args{
				es:      GetMockProjectByIDNoEvents(ctrl),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsNotFound(err)
				},
			},
		},
		{
			name: "project from events, no id",
			args: args{
				es:      GetMockProjectByIDNoEvents(ctrl),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsPreconditionFailed(err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ProjectByID(nil, tt.args.project)

			if !tt.res.wantErr && result.ID != tt.res.project.ID {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.project.ID, result.ID)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestCreateProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es      *ProjectEventstore
		ctx     context.Context
		project *model.Project
	}
	type res struct {
		project *model.Project
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "project from events, ok",
			args: args{
				es:      GetMockCreateProject(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name"},
			},
			res: res{
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}, Name: "Name"},
			},
		},
		{
			name: "create project no name",
			args: args{
				es:      GetMockCreateProject(ctrl),
				ctx:     auth.NewMockContext("orgID", "userID"),
				project: &model.Project{ObjectRoot: es_models.ObjectRoot{ID: "ID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: func(err error) bool {
					return caos_errs.IsPreconditionFailed(err)
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.CreateProject(tt.args.ctx, tt.args.project)

			if !tt.res.wantErr && result.ID == "" {
				t.Errorf("result has no id")
			}
			if !tt.res.wantErr && result.Name != tt.res.project.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.project.Name, result.Name)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
