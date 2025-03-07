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
)

type challengeSubmissionRepository struct {
	db *sqlx.DB
}

func NewChallengeSubmissionRepository(conn *sqlx.DB) contract.IChallengeSubmissionRepository {
	return &challengeSubmissionRepository{
		db: conn,
	}
}

func (r *challengeSubmissionRepository) CreateSubmission(ctx context.Context, txWrapper database.ITransaction,
	submission *entity.ChallengeSubmission) error {
	tx := txWrapper.GetTx()

	query := `
		INSERT INTO challenge_submissions (
			id, challenge_id, student_id, url
		) VALUES (
			:id, :challenge_id, :student_id, :url
		)
	`

	_, err := tx.NamedExecContext(ctx, query, submission)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "challenge_submissions_challenge_id_student_id_key":
				return errors.New("student has already submitted for this challenge")
			}
		}
		return fmt.Errorf("failed to create challenge submission: %w", err)
	}

	updateQuery := `
		UPDATE challenges
		SET submission_count = submission_count + 1, updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, updateQuery, submission.ChallengeID)
	if err != nil {
		return fmt.Errorf("failed to update challenge submission count: %w", err)
	}

	return nil
}

func (r *challengeSubmissionRepository) GetSubmissionByID(ctx context.Context,
	id uuid.UUID) (*entity.ChallengeSubmission, error) {
	query := `
		SELECT id, challenge_id, student_id, url, created_at
		FROM challenge_submissions
		WHERE id = $1
	`

	submission := &entity.ChallengeSubmission{}
	err := r.db.GetContext(ctx, submission, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("submission not found")
		}
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}

	return submission, nil
}

func (r *challengeSubmissionRepository) GetSubmissionByStudent(ctx context.Context, challengeID,
	studentID uuid.UUID) (*entity.ChallengeSubmission, error) {
	query := `
		SELECT cs.id, cs.challenge_id, cs.student_id, cs.url, cs.created_at,
			u.id as "student.id", u.name as "student.name", u.has_avatar as "student.has_avatar"
		FROM challenge_submissions cs
		JOIN users u ON cs.student_id = u.id
		WHERE cs.challenge_id = $1 AND cs.student_id = $2
	`

	submission := &entity.ChallengeSubmission{
		Student: &entity.User{},
	}

	err := r.db.GetContext(ctx, submission, query, challengeID, studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("submission not found")
		}
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}

	// Fetch feedback if it exists
	feedbackQuery := `
		SELECT f.submission_id, f.mentor_id, f.score, f.feedback, f.created_at,
			u.id as "mentor.id", u.name as "mentor.name", u.has_avatar as "mentor.has_avatar"
		FROM challenge_submission_feedbacks f
		JOIN users u ON f.mentor_id = u.id
		WHERE f.submission_id = (
			SELECT id
			FROM challenge_submissions
			WHERE challenge_id = $1 AND student_id = $2
		)
	`

	feedback := &entity.ChallengeSubmissionFeedback{
		Mentor: &entity.User{},
	}

	err = r.db.GetContext(ctx, feedback, feedbackQuery, challengeID, studentID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get submission feedback: %w", err)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		submission.Feedback = feedback
	}

	return submission, nil
}

func (r *challengeSubmissionRepository) GetSubmissionsByChallenge(ctx context.Context, challengeID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*entity.ChallengeSubmission, dto.PaginationResponse, error) {
	baseQuery := `
		SELECT cs.id, cs.challenge_id, cs.student_id, cs.url, cs.created_at,
			u.id as "student.id", u.name as "student.name", u.has_avatar as "student.has_avatar"
		FROM challenge_submissions cs
		JOIN users u ON cs.student_id = u.id
		WHERE cs.challenge_id = $1
	`

	var sqlQuery string
	var args []interface{}

	args = append(args, challengeID)

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

		sqlQuery = baseQuery + fmt.Sprintf(" AND cs.id %s $2 ORDER BY cs.id %s LIMIT $3", operator, orderDirection)
		args = append(args, pageReq.Cursor, pageReq.Limit+1)
	} else {
		sqlQuery = baseQuery + " ORDER BY cs.id DESC LIMIT $2"
		args = append(args, pageReq.Limit+1)
	}

	rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("failed to get challenge submissions: %w", err)
	}
	defer rows.Close()

	var submissions []*entity.ChallengeSubmission
	for rows.Next() {
		submission := &entity.ChallengeSubmission{
			Student: &entity.User{},
		}
		if err := rows.StructScan(submission); err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan submission row: %w", err)
		}
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("error iterating over submission rows: %w", err)
	}

	hasMore := false
	if len(submissions) > pageReq.Limit {
		hasMore = true
		submissions = submissions[:pageReq.Limit]
	}

	if pageReq.Direction == "prev" && pageReq.Cursor != uuid.Nil {
		for i, j := 0, len(submissions)-1; i < j; i, j = i+1, j-1 {
			submissions[i], submissions[j] = submissions[j], submissions[i]
		}
	}

	if len(submissions) > 0 {
		feedbackQuery := `
			SELECT f.submission_id, f.mentor_id, f.score, f.feedback, f.created_at,
				u.id as "mentor.id", u.name as "mentor.name", u.has_avatar as "mentor.has_avatar"
			FROM challenge_submission_feedbacks f
			JOIN users u ON f.mentor_id = u.id
			WHERE f.submission_id IN (
				SELECT id FROM challenge_submissions
				WHERE challenge_id = $1 AND student_id = ANY($2)
			)
		`

		studentIDs := make([]uuid.UUID, len(submissions))
		submissionMap := make(map[string]*entity.ChallengeSubmission)

		for i, sub := range submissions {
			studentIDs[i] = sub.StudentID
			submissionMap[sub.ID.String()] = sub
		}

		feedbackRows, err := r.db.QueryxContext(ctx, feedbackQuery, challengeID, studentIDs)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to get submission feedbacks: %w", err)
		}

		if err == nil {
			defer feedbackRows.Close()

			for feedbackRows.Next() {
				feedback := &entity.ChallengeSubmissionFeedback{
					Mentor: &entity.User{},
				}
				if err := feedbackRows.StructScan(feedback); err != nil {
					return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan feedback row: %w", err)
				}

				// Find the corresponding submission and attach the feedback
				if sub, ok := submissionMap[feedback.SubmissionID.String()]; ok {
					sub.Feedback = feedback
				}
			}

			if err := feedbackRows.Err(); err != nil {
				return nil, dto.PaginationResponse{}, fmt.Errorf("error iterating over feedback rows: %w", err)
			}
		}
	}

	return submissions, dto.PaginationResponse{HasMore: hasMore}, nil
}

func (r *challengeSubmissionRepository) CreateFeedback(ctx context.Context, txWrapper database.ITransaction,
	feedback *entity.ChallengeSubmissionFeedback) error {
	tx := txWrapper.GetTx()

	query := `
		INSERT INTO challenge_submission_feedbacks (
			submission_id, mentor_id, score, feedback
		) VALUES (
			:submission_id, :mentor_id, :score, :feedback
		)
	`

	_, err := tx.NamedExecContext(ctx, query, feedback)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "challenge_submission_feedbacks_pkey":
				return errors.New("feedback already exists for this submission")
			case "challenge_submission_feedbacks_submission_id_fkey":
				return errors.New("submission not found")
			}
		}
		return fmt.Errorf("failed to create feedback: %w", err)
	}

	return nil
}
