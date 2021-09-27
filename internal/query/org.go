package query

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/query/projection"
)

var (
	orgsTable = table{
		name: projection.OrgProjectionTable,
	}
	OrgColumnID = Column{
		name:  projection.OrgColumnID,
		table: orgsTable,
	}
	OrgColumnCreationDate = Column{
		name:  projection.OrgColumnCreationDate,
		table: orgsTable,
	}
	OrgColumnChangeDate = Column{
		name:  projection.OrgColumnChangeDate,
		table: orgsTable,
	}
	OrgColumnResourceOwner = Column{
		name:  projection.OrgColumnResourceOwner,
		table: orgsTable,
	}
	OrgColumnState = Column{
		name:  projection.OrgColumnState,
		table: orgsTable,
	}
	OrgColumnSequence = Column{
		name:  projection.OrgColumnSequence,
		table: orgsTable,
	}
	OrgColumnName = Column{
		name:  projection.OrgColumnName,
		table: orgsTable,
	}
	OrgColumnDomain = Column{
		name:  projection.OrgColumnDomain,
		table: orgsTable,
	}
	OrgsColumnCount = Column{
		name:  "COUNT(*) OVER ()",
		table: orgsTable,
	}
)

func prepareOrgQuery() (sq.SelectBuilder, func(*sql.Row) (*Org, error)) {
	return sq.Select(
			OrgColumnID.identifier(),
			OrgColumnCreationDate.identifier(),
			OrgColumnChangeDate.identifier(),
			OrgColumnResourceOwner.identifier(),
			OrgColumnState.identifier(),
			OrgColumnSequence.identifier(),
			OrgColumnName.identifier(),
			OrgColumnDomain.identifier(),
		).
			From(orgsTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Org, error) {
			o := new(Org)
			err := row.Scan(
				&o.ID,
				&o.CreationDate,
				&o.ChangeDate,
				&o.ResourceOwner,
				&o.State,
				&o.Sequence,
				&o.Name,
				&o.Domain,
			)
			if err != nil {
				if errs.Is(err, sql.ErrNoRows) {
					return nil, errors.ThrowNotFound(err, "QUERY-iTTGJ", "errors.orgs.not_found")
				}
				return nil, errors.ThrowInternal(err, "QUERY-pWS5H", "errors.internal")
			}
			return o, nil
		}
}

func (q *Queries) prepareOrgsQuery() (sq.SelectBuilder, func(*sql.Rows) (*Orgs, error)) {
	return sq.Select(
			OrgColumnID.identifier(),
			OrgColumnCreationDate.identifier(),
			OrgColumnChangeDate.identifier(),
			OrgColumnResourceOwner.identifier(),
			OrgColumnState.identifier(),
			OrgColumnSequence.identifier(),
			OrgColumnName.identifier(),
			OrgColumnDomain.identifier(),
			OrgsColumnCount.identifier()).
			From(projection.OrgProjectionTable).PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Orgs, error) {
			orgs := make([]*Org, 0)
			var count uint64
			for rows.Next() {
				org := new(Org)
				err := rows.Scan(
					&org.ID,
					&org.CreationDate,
					&org.ChangeDate,
					&org.ResourceOwner,
					&org.State,
					&org.Sequence,
					&org.Name,
					&org.Domain,
					&count,
				)
				if err != nil {
					return nil, err
				}
				orgs = append(orgs, org)
			}

			if err := rows.Close(); err != nil {
				return nil, errors.ThrowInternal(err, "QUERY-QMXJv", "unable to close rows")
			}

			return &Orgs{
				Orgs: orgs,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func (q *Queries) prepareOrgUniqueQuery() (sq.SelectBuilder, func(*sql.Row) (bool, error)) {
	return sq.Select("COUNT(*) = 0").
			From(orgsTable.identifier()).PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (isUnique bool, err error) {
			err = row.Scan(&isUnique)
			if err != nil {
				return false, errors.ThrowInternal(err, "QUERY-e6EiG", "errors.internal")
			}
			return isUnique, err
		}
}

func (q *Queries) OrgByID(ctx context.Context, id string) (*Org, error) {
	stmt, scan := prepareOrgQuery()
	query, args, err := stmt.Where(sq.Eq{
		OrgColumnID.identifier(): id,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-AWx52", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) OrgByDomainGlobal(ctx context.Context, domain string) (*Org, error) {
	stmt, scan := prepareOrgQuery()
	query, args, err := stmt.Where(sq.Eq{
		OrgColumnDomain.identifier(): domain,
	}).ToSql()
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-TYUCE", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, query, args...)
	return scan(row)
}

func (q *Queries) IsOrgUnique(ctx context.Context, name, domain string) (isUnique bool, err error) {
	query, scan := q.prepareOrgUniqueQuery()
	stmt, args, err := query.Where(sq.Eq{
		OrgColumnDomain.identifier(): domain,
	}).ToSql()
	if err != nil {
		return false, errors.ThrowInternal(err, "QUERY-TYUCE", "unable to create sql stmt")
	}

	row := q.client.QueryRowContext(ctx, stmt, args...)
	return scan(row)
}

func (q *Queries) ExistsOrg(ctx context.Context, id string) (err error) {
	_, err = q.OrgByID(ctx, id)
	return err
}

func (q *Queries) SearchOrgs(ctx context.Context, queries *OrgSearchQueries) (orgs *Orgs, err error) {
	query, scan := q.prepareOrgsQuery()
	stmt, args, err := queries.toQuery(query).ToSql()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "QUERY-wQ3by", "Errors.orgs.invalid.request")
	}

	rows, err := q.client.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, errors.ThrowInternal(err, "QUERY-M6mYN", "Errors.orgs.internal")
	}
	orgs, err = scan(rows)
	if err != nil {
		return nil, err
	}
	orgs.LatestSequence, err = q.latestSequence(ctx, projection.OrgProjectionTable)
	return orgs, err
}

type Orgs struct {
	SearchResponse
	Orgs []*Org
}

type Org struct {
	ID            string
	CreationDate  time.Time
	ChangeDate    time.Time
	ResourceOwner string
	State         domain.OrgState
	Sequence      uint64

	Name   string
	Domain string
}

type OrgSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func NewOrgDomainSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(OrgColumnDomain, value, method)
}

func NewOrgNameSearchQuery(method TextComparison, value string) (SearchQuery, error) {
	return NewTextQuery(OrgColumnName, value, method)
}

func (q *OrgSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.ToQuery(query)
	}
	return query
}
