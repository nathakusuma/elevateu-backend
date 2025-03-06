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
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
	"github.com/nathakusuma/elevateu-backend/pkg/sqlutil"
)

type courseFeedbackRepository struct {
	db *sqlx.DB
}

func NewCourseFeedbackRepository(conn *sqlx.DB) contract.ICourseFeedbackRepository {
	return &courseFeedbackRepository{
		db: conn,
	}
}

func (r *courseFeedbackRepository) CreateFeedback(ctx context.Context, txWrapper database.ITransaction,
	feedback *entity.CourseFeedback) error {
	tx := txWrapper.GetTx()

	query := `
		INSERT INTO course_feedbacks (
			id, course_id, student_id, rating, comment, created_at, updated_at
		) VALUES (
			:id, :course_id, :student_id, :rating, :comment, NOW(), NOW()
		)
	`

	_, err := sqlx.NamedExecContext(ctx, tx, query, feedback)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "course_feedbacks_course_student_key" {
				return fmt.Errorf("student has already submitted feedback for this course: %w", err)
			}
			if pgErr.ConstraintName == "course_feedbacks_course_id_fkey" {
				return fmt.Errorf("course not found: %w", err)
			}
		}
		return fmt.Errorf("failed to create course feedback: %w", err)
	}

	return nil
}

func (r *courseFeedbackRepository) GetFeedbacksByCourseID(ctx context.Context, courseID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*entity.CourseFeedback, dto.PaginationResponse, error) {

	baseQuery := `
		SELECT
			cf.id, cf.course_id, cf.student_id, cf.rating, cf.comment, cf.created_at, cf.updated_at,
			u.id AS "user.id", u.name AS "user.name", u.has_avatar AS "user.has_avatar"
		FROM course_feedbacks cf
		JOIN students s ON cf.student_id = s.user_id
		JOIN users u ON s.user_id = u.id
		WHERE cf.course_id = $1
	`

	var sqlQuery string
	var args []interface{}

	args = append(args, courseID)

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

		sqlQuery = baseQuery + fmt.Sprintf(" AND cf.id %s $2 ORDER BY cf.id %s LIMIT $3",
			operator, orderDirection)
		args = append(args, pageReq.Cursor, pageReq.Limit+1)
	} else {
		// Initial query without cursor
		sqlQuery = baseQuery + " ORDER BY cf.id DESC LIMIT $2"
		args = append(args, pageReq.Limit+1)
	}

	rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var feedbacks []*entity.CourseFeedback
	for rows.Next() {
		var feedback entity.CourseFeedback
		feedback.User = &entity.User{}

		if err = rows.StructScan(&feedback); err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan row: %w", err)
		}
		feedbacks = append(feedbacks, &feedback)
	}

	if err = rows.Err(); err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("error iterating over rows: %w", err)
	}

	hasMore := false
	if len(feedbacks) > pageReq.Limit {
		hasMore = true
		feedbacks = feedbacks[:pageReq.Limit]
	}

	// Reverse results for "prev" direction
	if pageReq.Direction == "prev" && len(feedbacks) > 0 {
		for i, j := 0, len(feedbacks)-1; i < j; i, j = i+1, j-1 {
			feedbacks[i], feedbacks[j] = feedbacks[j], feedbacks[i]
		}
	}

	return feedbacks, dto.PaginationResponse{HasMore: hasMore}, nil
}

func (r *courseFeedbackRepository) GetFeedbackByID(ctx context.Context,
	feedbackID uuid.UUID) (*entity.CourseFeedback, error) {
	query := `
		SELECT
			cf.id, cf.course_id, cf.student_id, cf.rating, cf.comment, cf.created_at, cf.updated_at,
			u.id AS "user.id", u.name AS "user.name", u.has_avatar AS "user.has_avatar"
		FROM course_feedbacks cf
		JOIN students s ON cf.student_id = s.user_id
		JOIN users u ON s.user_id = u.id
		WHERE cf.id = $1
	`

	var feedback entity.CourseFeedback
	feedback.User = &entity.User{}

	err := r.db.QueryRowxContext(ctx, query, feedbackID).StructScan(&feedback)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("feedback not found")
		}
		return nil, fmt.Errorf("failed to get feedback: %w", err)
	}

	return &feedback, nil
}

func (r *courseFeedbackRepository) UpdateFeedback(ctx context.Context, txWrapper database.ITransaction,
	feedbackID uuid.UUID, updates dto.CourseFeedbackUpdate) error {
	tx := txWrapper.GetTx()

	builder := sqlutil.NewSQLUpdateBuilder("course_feedbacks").
		WithUpdatedAt().
		Where("id = ?", feedbackID)

	query, args, err := builder.BuildFromStruct(updates)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// No fields to update
	if query == "" {
		return nil
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update feedback: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("feedback not found")
	}

	return nil
}

func (r *courseFeedbackRepository) DeleteFeedback(ctx context.Context, txWrapper database.ITransaction,
	feedbackID uuid.UUID) error {
	tx := txWrapper.GetTx()

	query := "DELETE FROM course_feedbacks WHERE id = $1"
	result, err := tx.ExecContext(ctx, query, feedbackID)
	if err != nil {
		return fmt.Errorf("failed to delete feedback: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("feedback not found")
	}

	return nil
}

func (r *courseFeedbackRepository) UpdateCourseRating(ctx context.Context, txWrapper database.ITransaction,
	courseID uuid.UUID, count int64, rating, total float64) error {
	tx := txWrapper.GetTx()

	query := `
		UPDATE courses
		SET
			rating = $1,
			rating_count = $2,
			total_rating = $3
		WHERE id = $4
	`

	_, err := tx.ExecContext(ctx, query, rating, count, total, courseID)
	if err != nil {
		return fmt.Errorf("failed to update course rating: %w", err)
	}

	return nil
}
