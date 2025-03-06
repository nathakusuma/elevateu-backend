package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/pkg/sqlutil"
)

type challengeGroupRepository struct {
	db *sqlx.DB
}

func NewChallengeGroupRepository(conn *sqlx.DB) contract.IChallengeGroupRepository {
	return &challengeGroupRepository{
		db: conn,
	}
}

func (r *challengeGroupRepository) CreateGroup(ctx context.Context, group *entity.ChallengeGroup) error {
	query := `
		INSERT INTO challenge_groups (
			id, category_id, title, description
		) VALUES (
			:id, :category_id, :title, :description
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, group)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "challenge_groups_category_id_fkey" {
			return fmt.Errorf("category not found: %w", err)
		}

		return err
	}

	return nil
}

func (r *challengeGroupRepository) GetGroups(ctx context.Context, query dto.GetChallengeGroupQuery,
	paginationReq dto.PaginationRequest) ([]*entity.ChallengeGroup, dto.PaginationResponse, error) {

	baseQuery := `
		SELECT
			id, category_id, title, description, challenge_count,
			created_at, updated_at
		FROM challenge_groups
	`

	var whereConditions []string
	var args []interface{}
	argIndex := 1

	if query.CategoryID != nil && *query.CategoryID != uuid.Nil {
		whereConditions = append(whereConditions, fmt.Sprintf("category_id = $%d", argIndex))
		args = append(args, *query.CategoryID)
		argIndex++
	}

	if query.Title != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("title ILIKE $%d", argIndex))
		args = append(args, "%"+query.Title+"%")
		argIndex++
	}

	if paginationReq.Cursor != uuid.Nil {
		var operator string
		var orderDirection string

		if paginationReq.Direction == "next" {
			operator = "<"
			orderDirection = "DESC"
		} else {
			operator = ">"
			orderDirection = "ASC"
		}

		whereConditions = append(whereConditions, fmt.Sprintf("id %s $%d", operator, argIndex))
		args = append(args, paginationReq.Cursor)
		argIndex++

		sqlQuery := baseQuery
		if len(whereConditions) > 0 {
			sqlQuery += " WHERE " + strings.Join(whereConditions, " AND ")
		}
		sqlQuery += fmt.Sprintf(" ORDER BY id %s LIMIT $%d", orderDirection, argIndex)
		args = append(args, paginationReq.Limit+1)

		rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
		if err != nil {
			return nil, dto.PaginationResponse{}, err
		}
		defer rows.Close()

		var groups []*entity.ChallengeGroup
		for rows.Next() {
			var group entity.ChallengeGroup
			group.Category = &entity.Category{}

			if err := rows.StructScan(&group); err != nil {
				return nil, dto.PaginationResponse{}, err
			}
			groups = append(groups, &group)
		}

		hasMore := false
		if len(groups) > paginationReq.Limit {
			hasMore = true
			groups = groups[:paginationReq.Limit]
		}

		if paginationReq.Direction == "prev" {
			for i, j := 0, len(groups)-1; i < j; i, j = i+1, j-1 {
				groups[i], groups[j] = groups[j], groups[i]
			}
		}

		return groups, dto.PaginationResponse{HasMore: hasMore}, nil
	} else {
		sqlQuery := baseQuery
		if len(whereConditions) > 0 {
			sqlQuery += " WHERE " + strings.Join(whereConditions, " AND ")
		}
		sqlQuery += fmt.Sprintf(" ORDER BY id DESC LIMIT $%d", argIndex)
		args = append(args, paginationReq.Limit+1)

		rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
		if err != nil {
			return nil, dto.PaginationResponse{}, err
		}
		defer rows.Close()

		var groups []*entity.ChallengeGroup
		for rows.Next() {
			var group entity.ChallengeGroup
			group.Category = &entity.Category{}

			if err := rows.StructScan(&group); err != nil {
				return nil, dto.PaginationResponse{}, err
			}
			groups = append(groups, &group)
		}

		hasMore := false
		if len(groups) > paginationReq.Limit {
			hasMore = true
			groups = groups[:paginationReq.Limit]
		}

		return groups, dto.PaginationResponse{HasMore: hasMore}, nil
	}
}

func (r *challengeGroupRepository) UpdateGroup(ctx context.Context, groupID uuid.UUID,
	updates dto.ChallengeGroupUpdate) error {
	builder := sqlutil.NewSQLUpdateBuilder("challenge_groups").
		WithUpdatedAt().
		Where("id = ?", groupID)

	query, args, err := builder.BuildFromStruct(updates)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// No fields to update (query is empty)
	if query == "" {
		return nil
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "challenge_groups_category_id_fkey" {
			return fmt.Errorf("category not found: %w", err)
		}
		return fmt.Errorf("failed to update challenge group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("challenge group not found")
	}

	return nil
}

func (r *challengeGroupRepository) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	query := "DELETE FROM challenge_groups WHERE id = $1"

	result, err := r.db.ExecContext(ctx, query, groupID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("challenge group not found")
	}

	return nil
}
