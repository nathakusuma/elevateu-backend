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
)

type mentoringRepository struct {
	db *sqlx.DB
}

func NewMentoringRepository(conn *sqlx.DB) contract.IMentoringRepository {
	return &mentoringRepository{
		db: conn,
	}
}

func (r *mentoringRepository) CreateChat(ctx context.Context, chat *entity.MentoringChat) error {
	return r.createChat(ctx, r.db, chat)
}

func (r *mentoringRepository) createChat(ctx context.Context, tx sqlx.ExtContext, chat *entity.MentoringChat) error {
	query := `
       INSERT INTO mentoring_chats (
          id, mentor_id, student_id, expires_at, is_trial
       ) VALUES (
          :id, :mentor_id, :student_id, :expires_at, :is_trial
       )
       ON CONFLICT (student_id, mentor_id)
       DO UPDATE SET expires_at = :expires_at
    `

	_, err := sqlx.NamedExecContext(ctx, tx, query, chat)
	if err != nil {
		return fmt.Errorf("failed to create or update chat: %w", err)
	}

	return nil
}

func (r *mentoringRepository) CreateTrialChat(ctx context.Context, chat *entity.MentoringChat) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	query1 := `
		INSERT INTO mentoring_trials (student_id) VALUES (:student_id)
	`
	_, err = tx.NamedExecContext(ctx, query1, chat)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("trial chat already exists: %w", err)
		}

		return fmt.Errorf("failed to create trial chat: %w", err)
	}

	if err = r.createChat(ctx, tx, chat); err != nil {
		return fmt.Errorf("failed to create chat: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *mentoringRepository) GetChatByID(ctx context.Context, chatID uuid.UUID) (*entity.MentoringChat, error) {
	query := `
       SELECT id, student_id, mentor_id, expires_at, is_trial
       FROM mentoring_chats
       WHERE id = $1
    `

	var chat entity.MentoringChat
	err := r.db.GetContext(ctx, &chat, query, chatID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("chat not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	return &chat, nil
}

func (r *mentoringRepository) GetChatsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.MentoringChat, error) {
	type chatJoin struct {
		ID        uuid.UUID `db:"id"`
		StudentID uuid.UUID `db:"student_id"`
		MentorID  uuid.UUID `db:"mentor_id"`
		ExpiresAt time.Time `db:"expires_at"`
		IsTrial   bool      `db:"is_trial"`

		LastMessageID        sql.NullString `db:"last_message.id"`
		LastMessageChatID    sql.NullString `db:"last_message.chat_id"`
		LastMessageSenderID  sql.NullString `db:"last_message.sender_id"`
		LastMessageText      sql.NullString `db:"last_message.message"`
		LastMessageCreatedAt sql.NullTime   `db:"last_message.created_at"`
	}

	query := `
        SELECT mc.id, mc.student_id, mc.mentor_id, mc.expires_at, mc.is_trial,
               mm.id as "last_message.id",
               mm.chat_id as "last_message.chat_id",
               mm.sender_id as "last_message.sender_id",
               mm.message as "last_message.message",
               mm.created_at as "last_message.created_at"
        FROM mentoring_chats mc
        LEFT JOIN LATERAL (
            SELECT *
            FROM mentoring_messages
            WHERE chat_id = mc.id
            ORDER BY created_at DESC
            LIMIT 1
        ) mm ON true
        WHERE mc.student_id = $1 OR mc.mentor_id = $1
        ORDER BY mm.created_at DESC NULLS LAST
    `

	rows, err := r.db.QueryxContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chats: %w", err)
	}
	defer rows.Close()

	var chats []*entity.MentoringChat
	for rows.Next() {
		var cj chatJoin
		if err := rows.StructScan(&cj); err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)
		}

		chat := &entity.MentoringChat{
			ID:        cj.ID,
			StudentID: cj.StudentID,
			MentorID:  cj.MentorID,
			ExpiresAt: cj.ExpiresAt,
			IsTrial:   cj.IsTrial,
		}

		if cj.LastMessageID.Valid {
			messageID, _ := uuid.Parse(cj.LastMessageID.String)
			chatID, _ := uuid.Parse(cj.LastMessageChatID.String)
			senderID, _ := uuid.Parse(cj.LastMessageSenderID.String)

			chat.LastMessage = &entity.MentoringMessage{
				ID:        messageID,
				ChatID:    chatID,
				SenderID:  senderID,
				Message:   cj.LastMessageText.String,
				CreatedAt: cj.LastMessageCreatedAt.Time,
			}
		}

		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return chats, nil
}

func (r *mentoringRepository) GetChatByMentorAndStudent(ctx context.Context, mentorID,
	studentID uuid.UUID) (*entity.MentoringChat, error) {
	query := `
	   SELECT id, student_id, mentor_id, expires_at, is_trial
	   FROM mentoring_chats
	   WHERE mentor_id = $1
	     AND student_id = $2
	`

	var chat entity.MentoringChat
	err := r.db.GetContext(ctx, &chat, query, mentorID, studentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("chat not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	return &chat, nil
}

func (r *mentoringRepository) SendMessage(ctx context.Context, message *entity.MentoringMessage) error {
	query := `
		INSERT INTO mentoring_messages (
			id, chat_id, sender_id, message
		) VALUES (
			:id, :chat_id, :sender_id, :message
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (r *mentoringRepository) GetMessages(ctx context.Context, chatID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*entity.MentoringMessage, dto.PaginationResponse, error) {

	baseQuery := `
		SELECT id, chat_id, sender_id, message, created_at
		FROM mentoring_messages
		WHERE chat_id = $1
	`

	var sqlQuery string
	var args []interface{}
	args = append(args, chatID)

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

		sqlQuery = baseQuery + fmt.Sprintf(" AND id %s $2 ORDER BY id %s LIMIT $3", operator, orderDirection)
		args = append(args, pageReq.Cursor, pageReq.Limit+1)
	} else {
		sqlQuery = baseQuery + " ORDER BY id DESC LIMIT $2"
		args = append(args, pageReq.Limit+1)
	}

	rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*entity.MentoringMessage
	for rows.Next() {
		message := &entity.MentoringMessage{}
		if err := rows.StructScan(message); err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan message row: %w", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("error iterating over message rows: %w", err)
	}

	hasMore := len(messages) > pageReq.Limit
	if hasMore {
		messages = messages[:pageReq.Limit]
	}

	if pageReq.Direction == "prev" && pageReq.Cursor != uuid.Nil {
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}
	}

	return messages, dto.PaginationResponse{HasMore: hasMore}, nil
}
