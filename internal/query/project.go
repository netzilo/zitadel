package query

import (
	"context"
	"database/sql"
	errs "errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/query/projection"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
)

var (
	projectsTable = table{
		name: projection.ProjectProjectionTable,
	}
	ProjectColumnID = Column{
		name:  projection.ProjectColumnID,
		table: projectsTable,
	}
	ProjectColumnName = Column{
		name:  projection.ProjectColumnName,
		table: projectsTable,
	}
	ProjectColumnProjectRoleAssertion = Column{
		name:  projection.ProjectColumnProjectRoleAssertion,
		table: projectsTable,
	}
	ProjectColumnProjectRoleCheck = Column{
		name:  projection.ProjectColumnProjectRoleCheck,
		table: projectsTable,
	}
	ProjectColumnHasProjectCheck = Column{
		name:  projection.ProjectColumnHasProjectCheck,
		table: projectsTable,
	}
	ProjectColumnPrivateLabelingSetting = Column{
		name:  projection.ProjectColumnPrivateLabelingSetting,
		table: projectsTable,
	}
	ProjectColumnCreationDate = Column{
		name:  projection.ProjectColumnCreationDate,
		table: projectsTable,
	}
	ProjectColumnChangeDate = Column{
		name:  projection.ProjectColumnChangeDate,
		table: projectsTable,
	}
	ProjectColumnResourceOwner = Column{
		name:  projection.ProjectColumnResourceOwner,
		table: projectsTable,
	}
	ProjectColumnCreator = Column{
		name:  projection.ProjectColumnCreator,
		table: projectsTable,
	}
	ProjectColumnSequence = Column{
		name:  projection.ProjectColumnSequence,
		table: projectsTable,
	}
	ProjectColumnState = Column{
		name:  projection.ProjectColumnState,
		table: projectsTable,
	}
)

func prepareProjectQuery() (sq.SelectBuilder, func(*sql.Row) (*Project, error)) {
	return sq.Select(
			ProjectColumnID.identifier(),
			ProjectColumnCreationDate.identifier(),
			ProjectColumnChangeDate.identifier(),
			ProjectColumnResourceOwner.identifier(),
			ProjectColumnState.identifier(),
			ProjectColumnSequence.identifier(),
			ProjectColumnName.identifier(),
			ProjectColumnProjectRoleAssertion.identifier(),
			ProjectColumnProjectRoleCheck.identifier(),
			ProjectColumnHasProjectCheck.identifier(),
			ProjectColumnPrivateLabelingSetting.identifier()).
			From(projection.ProjectProjectionTable).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Project, error) {
			p := new(Project)
			err := row.Scan(
				&p.ID,
				&p.CreationDate,
				&p.ChangeDate,
				&p.ResourceOwner,
				&p.State,
				&p.Sequence,
				&p.Name,
				&p.ProjectRoleAssertion,
				&p.ProjectRoleCheck,
				&p.HasProjectCheck,
				&p.PrivateLabelingSetting,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-fk2fs", "errors.projects.not_found")
				}
				fmt.Printf("error: %v", err.Error())
				return nil, errors.ThrowInternal(err, "QUERY-dj2FF", "errors.internal")
			}
			return p, nil
		}
}

func (q *Queries) prepareProjectsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Projects, error)) {
	return sq.Select(
			ProjectColumnID.identifier(),
			ProjectColumnCreationDate.identifier(),
			ProjectColumnChangeDate.identifier(),
			ProjectColumnResourceOwner.identifier(),
			ProjectColumnState.identifier(),
			ProjectColumnSequence.identifier(),
			ProjectColumnName.identifier(),
			ProjectColumnProjectRoleAssertion.identifier(),
			ProjectColumnProjectRoleCheck.identifier(),
			ProjectColumnHasProjectCheck.identifier(),
			ProjectColumnPrivateLabelingSetting.identifier(),
			"COUNT(*) OVER ()").
			From(projection.ProjectProjectionTable).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Projects, error) {
			projects := make([]*Project, 0)
			var count uint64
			for rows.Next() {
				project := new(Project)
				err := rows.Scan(
					&project.ID,
					&project.CreationDate,
					&project.ChangeDate,
					&project.ResourceOwner,
					&project.State,
					&project.Sequence,
					&project.Name,
					&project.ProjectRoleAssertion,
					&project.ProjectRoleCheck,
					&project.HasProjectCheck,
					&project.PrivateLabelingSetting,
					&count,
				)
				if err != nil {
					return nil, err
				}
				projects = append(projects, project)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-QMXJv", "unable to close rows")
			}

			return &Projects{
				Projects: projects,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func (q *Queries) prepareProjectUniqueQuery() (sq.SelectBuilder, func(*sql.Row) (bool, error)) {
	return sq.Select("COUNT(*) = 0").
			From(projection.ProjectProjectionTable).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (isUnique bool, err error) {
			err = row.Scan(&isUnique)
			if err != nil {
				return false, errors.ThrowInternal(err, "QUERY-2n99F", "errors.internal")
			}
			return isUnique, err
		}
}

func (q *Queries) ProjectByID(ctx context.Context, id string) (*Project, error) {
	stmt, scan := prepareProjectQuery()
	query, args, err := stmt.Where(sq.Eq{
		ProjectColumnID.identifier(): id,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-2m00Q", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) ExistsProject(ctx context.Context, id string) (err error) {
	_, err = q.ProjectByID(ctx, id)
	return err
}

func (q *Queries) SearchProjects(ctx context.Context, queries *ProjectSearchQueries) (projects *Projects, err error) {
	query, scan := q.prepareProjectsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-fn9ew", "Errors.projects.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-2j00f", "Errors.projects.internal")
	}
	projects, err = scan(rows)
	if err != nil {
		return nil, err
	}
	projects.LatestSequence, err = q.latestSequence(ctx, projectsTable)
	return projects, err
}

type Projects struct {
	SearchResponse
	Projects []*Project
}

type Project struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.ProjectState
	Sequence      uint64

	Name                   string
	ProjectRoleAssertion   bool
	ProjectRoleCheck       bool
	HasProjectCheck        bool
	PrivateLabelingSetting domain.PrivateLabelingSetting
}

type ProjectSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func NewProjectNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnName, value, method)
}

func NewProjectIDSearchQuery(values []string) (SearchQuery, error) {
	list := make([]interface{}, len(values))
	for i, value := range values {
		list[i] = value
	}
	return NewListQuery(ProjectColumnID, list, ListIn)
}

func NewProjectResourceOwnerSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(ProjectColumnResourceOwner, value, method)
}

func (q *ProjectSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}

func (r *ProjectSearchQueries) AppendMyResourceOwnerQuery(orgID string) error {
	query, err := NewProjectResourceOwnerSearchQuery(TextEquals, orgID)
	if err != nil {
		return err
	}
	r.Queries = append(r.Queries, query)
	return nil
}

func (r ProjectSearchQueries) AppendPermissionQueries(permissions []string) error {
	if !authz.HasGlobalPermission(permissions) {
		ids := authz.GetAllPermissionCtxIDs(permissions)
		query, err := NewProjectIDSearchQuery(ids)
		if err != nil {
			return err
		}
		r.Queries = append(r.Queries, query)
	}
	return nil
}
