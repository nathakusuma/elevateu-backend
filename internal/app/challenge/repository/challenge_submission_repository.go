package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

type submissionWithFeedback struct {
	ID          uuid.UUID `db:"id"`
	ChallengeID uuid.UUID `db:"challenge_id"`
	StudentID   uuid.UUID `db:"student_id"`
	URL         string    `db:"url"`
	CreatedAt   time.Time `db:"created_at"`

	StudentName      string `db:"student.name"`
	StudentHasAvatar bool   `db:"student.has_avatar"`

	FeedbackSubmissionID sql.NullString `db:"feedback.submission_id"`
	FeedbackMentorID     sql.NullString `db:"feedback.mentor_id"`
	FeedbackScore        sql.NullInt64  `db:"feedback.score"`
	FeedbackText         sql.NullString `db:"feedback.feedback"`
	FeedbackCreatedAt    sql.NullTime   `db:"feedback.created_at"`

	MentorID        sql.NullString `db:"feedback.mentor.id"`
	MentorName      sql.NullString `db:"feedback.mentor.name"`
	MentorHasAvatar sql.NullBool   `db:"feedback.mentor.has_avatar"`
}

func mapToSubmission(sr *submissionWithFeedback) *entity.ChallengeSubmission {
	submission := &entity.ChallengeSubmission{
		ID:          sr.ID,
		ChallengeID: sr.ChallengeID,
		StudentID:   sr.StudentID,
		URL:         sr.URL,
		CreatedAt:   sr.CreatedAt,
		Student: &entity.User{
			ID:        sr.StudentID,
			Name:      sr.StudentName,
			HasAvatar: sr.StudentHasAvatar,
		},
	}

	if sr.FeedbackSubmissionID.Valid {
		feedbackID, _ := uuid.Parse(sr.FeedbackSubmissionID.String)
		mentorID, _ := uuid.Parse(sr.MentorID.String)

		submission.Feedback = &entity.ChallengeSubmissionFeedback{
			SubmissionID: feedbackID,
			MentorID:     mentorID,
			Score:        int(sr.FeedbackScore.Int64),
			Feedback:     sr.FeedbackText.String,
			CreatedAt:    sr.FeedbackCreatedAt.Time,
			Mentor: &entity.User{
				ID:        mentorID,
				Name:      sr.MentorName.String,
				HasAvatar: sr.MentorHasAvatar.Bool,
			},
		}
	}

	return submission
}

func (r *challengeSubmissionRepository) GetSubmissionByStudent(ctx context.Context, challengeID,
	studentID uuid.UUID) (*entity.ChallengeSubmission, error) {

	query := `
        SELECT cs.id, cs.challenge_id, cs.student_id, cs.url, cs.created_at,
            u.name as "student.name", u.has_avatar as "student.has_avatar",
            f.submission_id as "feedback.submission_id", f.mentor_id as "feedback.mentor_id",
            f.score as "feedback.score", f.feedback as "feedback.feedback", f.created_at as "feedback.created_at",
            m.id as "feedback.mentor.id", m.name as "feedback.mentor.name", m.has_avatar as "feedback.mentor.has_avatar"
        FROM challenge_submissions cs
        JOIN users u ON cs.student_id = u.id
        LEFT JOIN challenge_submission_feedbacks f ON cs.id = f.submission_id
        LEFT JOIN users m ON f.mentor_id = m.id
        WHERE cs.challenge_id = $1 AND cs.student_id = $2
    `

	var result submissionWithFeedback
	err := r.db.GetContext(ctx, &result, query, challengeID, studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("submission not found")
		}
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}

	return mapToSubmission(&result), nil
}

func (r *challengeSubmissionRepository) GetSubmissionsByChallenge(ctx context.Context, challengeID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*entity.ChallengeSubmission, dto.PaginationResponse, error) {

	baseQuery := `
        SELECT cs.id, cs.challenge_id, cs.student_id, cs.url, cs.created_at,
            u.name as "student.name", u.has_avatar as "student.has_avatar",
            f.submission_id as "feedback.submission_id", f.mentor_id as "feedback.mentor_id",
            f.score as "feedback.score", f.feedback as "feedback.feedback", f.created_at as "feedback.created_at",
            m.id as "feedback.mentor.id", m.name as "feedback.mentor.name", m.has_avatar as "feedback.mentor.has_avatar"
        FROM challenge_submissions cs
        JOIN users u ON cs.student_id = u.id
        LEFT JOIN challenge_submission_feedbacks f ON cs.id = f.submission_id
        LEFT JOIN users m ON f.mentor_id = m.id
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
		var result submissionWithFeedback
		if err := rows.StructScan(&result); err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan submission row: %w", err)
		}

		submissions = append(submissions, mapToSubmission(&result))
	}

	if err := rows.Err(); err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("error iterating over submission rows: %w", err)
	}

	hasMore := len(submissions) > pageReq.Limit
	if hasMore {
		submissions = submissions[:pageReq.Limit]
	}

	if pageReq.Direction == "prev" && pageReq.Cursor != uuid.Nil {
		for i, j := 0, len(submissions)-1; i < j; i, j = i+1, j-1 {
			submissions[i], submissions[j] = submissions[j], submissions[i]
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
