package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/pkg/sqlutil"
)

type challengeRepository struct {
	db *sqlx.DB
}

func NewChallengeRepository(conn *sqlx.DB) contract.IChallengeRepository {
	return &challengeRepository{
		db: conn,
	}
}

func (r *challengeRepository) CreateChallenge(ctx context.Context, challenge *entity.Challenge) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO challenges (
			id, group_id, title, subtitle, description, difficulty, is_free
		) VALUES (
			:id, :group_id, :title, :subtitle, :description, :difficulty, :is_free
		)
	`

	_, err = tx.NamedExecContext(ctx, query, challenge)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "challenges_group_id_fkey" {
			return fmt.Errorf("challenge group not found: %w", err)
		}

		return fmt.Errorf("failed to create challenge: %w", err)
	}

	// Increment challenge_count in challenge_groups
	updateQuery := `
		UPDATE challenge_groups
		SET challenge_count = challenge_count + 1, updated_at = NOW()
		WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, updateQuery, challenge.GroupID)
	if err != nil {
		return fmt.Errorf("failed to update challenge count: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *challengeRepository) GetChallenges(ctx context.Context, groupID uuid.UUID, difficulty enum.ChallengeDifficulty,
	pageReq dto.PaginationRequest) ([]*entity.Challenge, dto.PaginationResponse, error) {

	baseQuery := `
		SELECT
			id, group_id, title, subtitle, description, difficulty, is_free,
			submission_count, created_at, updated_at
		FROM challenges
		WHERE group_id = $1 AND difficulty = $2
	`

	var sqlQuery string
	var args []interface{}

	args = append(args, groupID, difficulty)

	if pageReq.Cursor != uuid.Nil {
		var operator string
		var orderDirection string

		if pageReq.Direction == "next" {
			operator = "<"
			orderDirection = "DESC"
		} else {
			operator = ">"
			orderDirection = "ASC"
		}

		sqlQuery = baseQuery + fmt.Sprintf(" AND id %s $3 ORDER BY id %s LIMIT $4", operator, orderDirection)
		args = append(args, pageReq.Cursor, pageReq.Limit+1)
	} else {
		sqlQuery = baseQuery + " ORDER BY id DESC LIMIT $3"
		args = append(args, pageReq.Limit+1)
	}

	rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, dto.PaginationResponse{}, err
	}
	defer rows.Close()

	var challenges []*entity.Challenge
	for rows.Next() {
		var challenge entity.Challenge
		if err = rows.StructScan(&challenge); err != nil {
			return nil, dto.PaginationResponse{}, err
		}
		challenges = append(challenges, &challenge)
	}

	hasMore := false
	if len(challenges) > pageReq.Limit {
		hasMore = true
		challenges = challenges[:pageReq.Limit]
	}

	if pageReq.Direction == "prev" && pageReq.Cursor != uuid.Nil {
		for i, j := 0, len(challenges)-1; i < j; i, j = i+1, j-1 {
			challenges[i], challenges[j] = challenges[j], challenges[i]
		}
	}

	return challenges, dto.PaginationResponse{HasMore: hasMore}, nil
}

func (r *challengeRepository) GetChallengeByID(ctx context.Context, id uuid.UUID) (*entity.Challenge, error) {
	query := `
		SELECT
			id, group_id, title, subtitle, description, difficulty, is_free,
			submission_count, created_at, updated_at
		FROM challenges
		WHERE id = $1
	`

	var challenge entity.Challenge
	err := r.db.GetContext(ctx, &challenge, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("challenge not found: %w", err)
		}
		return nil, err
	}

	return &challenge, nil
}

func (r *challengeRepository) UpdateChallenge(ctx context.Context, id uuid.UUID, updates *dto.ChallengeUpdate) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current group ID
	var currentGroupID uuid.UUID
	getCurrentQuery := "SELECT group_id FROM challenges WHERE id = $1"
	err = tx.GetContext(ctx, &currentGroupID, getCurrentQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("challenge not found: %w", err)
		}
		return fmt.Errorf("failed to get current challenge: %w", err)
	}

	builder := sqlutil.NewSQLUpdateBuilder("challenges").
		WithUpdatedAt().
		Where("id = ?", id)

	query, args, err := builder.BuildFromStruct(updates)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// No fields to update (query is empty)
	if query == "" {
		return nil
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "challenges_group_id_fkey" {
			return fmt.Errorf("challenge group not found: %w", err)
		}
		return fmt.Errorf("failed to update challenge: %w", err)
	}

	// If group ID was changed, update challenge counts
	if updates.GroupID != nil && *updates.GroupID != currentGroupID {
		decrementQuery := `
			UPDATE challenge_groups
			SET challenge_count = challenge_count - 1, updated_at = NOW()
			WHERE id = $1
		`
		_, err = tx.ExecContext(ctx, decrementQuery, currentGroupID)
		if err != nil {
			return fmt.Errorf("failed to decrement old group challenge count: %w", err)
		}

		incrementQuery := `
			UPDATE challenge_groups
			SET challenge_count = challenge_count + 1, updated_at = NOW()
			WHERE id = $1
		`
		_, err = tx.ExecContext(ctx, incrementQuery, updates.GroupID)
		if err != nil {
			return fmt.Errorf("failed to increment new group challenge count: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *challengeRepository) DeleteChallenge(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var groupID uuid.UUID
	getGroupQuery := "SELECT group_id FROM challenges WHERE id = $1"
	err = tx.GetContext(ctx, &groupID, getGroupQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("challenge not found: %w", err)
		}
		return fmt.Errorf("failed to get challenge group: %w", err)
	}

	deleteQuery := "DELETE FROM challenges WHERE id = $1"
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete challenge: %w", err)
	}

	updateQuery := `
		UPDATE challenge_groups
		SET challenge_count = challenge_count - 1, updated_at = NOW()
		WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, updateQuery, groupID)
	if err != nil {
		return fmt.Errorf("failed to update challenge count: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
