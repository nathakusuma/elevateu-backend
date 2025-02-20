package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type authRepository struct {
	db  *sqlx.DB
	rds *redis.Client
}

func NewAuthRepository(db *sqlx.DB, rds *redis.Client) contract.IAuthRepository {
	return &authRepository{
		db:  db,
		rds: rds,
	}
}

func (r *authRepository) SetRegisterOTP(ctx context.Context, email string, otp string) error {
	if err := r.rds.Set(ctx, "auth:"+email+":register_otp", otp, 10*time.Minute).Err(); err != nil {
		return fmt.Errorf("failed to set otp: %w", err)
	}

	return nil
}

func (r *authRepository) GetRegisterOTP(ctx context.Context, email string) (string, error) {
	result, err := r.rds.Get(ctx, "auth:"+email+":register_otp").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("otp not found: %w", err)
		}

		return "", fmt.Errorf("failed to get otp: %w", err)
	}

	return result, nil
}

func (r *authRepository) DeleteRegisterOTP(ctx context.Context, email string) error {
	if err := r.rds.Del(ctx, "auth:"+email+":register_otp").Err(); err != nil {
		return fmt.Errorf("failed to delete otp: %w", err)
	}

	return nil
}

func (r *authRepository) CreateAuthSession(ctx context.Context, session *entity.AuthSession) error {
	return r.createAuthSession(ctx, r.db, session)
}

func (r *authRepository) createAuthSession(ctx context.Context, tx sqlx.ExtContext,
	authSession *entity.AuthSession) error {
	query := `INSERT INTO auth_sessions (token, user_id, expires_at)
				VALUES (:token, :user_id, :expires_at)
				ON CONFLICT (user_id) DO UPDATE SET token = :token, expires_at = :expires_at`

	_, err := sqlx.NamedExecContext(ctx, tx, query, authSession)
	if err != nil {
		return fmt.Errorf("failed to create auth session: %w", err)
	}

	return nil
}

func (r *authRepository) GetAuthSessionByToken(ctx context.Context, token string) (*entity.AuthSession, error) {
	var authSession entity.AuthSession

	statement := `SELECT
    		token,
			user_id,
			expires_at
		FROM auth_sessions
		WHERE token = $1
		`

	err := r.db.GetContext(ctx, &authSession, statement, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("auth session not found: %w", err)
		}

		return nil, fmt.Errorf("failed to get auth session by token: %w", err)
	}

	return &authSession, nil
}

func (r *authRepository) deleteAuthSession(ctx context.Context, tx sqlx.ExtContext, userID uuid.UUID) error {
	query := `DELETE FROM auth_sessions WHERE user_id = $1`

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete auth session: %w", err)
	}

	return nil
}

func (r *authRepository) DeleteAuthSession(ctx context.Context, userID uuid.UUID) error {
	return r.deleteAuthSession(ctx, r.db, userID)
}

func (r *authRepository) SetPasswordResetOTP(ctx context.Context, email, otp string) error {
	if err := r.rds.Set(ctx, "auth:"+email+":reset_password_otp", otp, 10*time.Minute).Err(); err != nil {
		return fmt.Errorf("failed to set otp: %w", err)
	}

	return nil
}

func (r *authRepository) GetPasswordResetOTP(ctx context.Context, email string) (string, error) {
	result, err := r.rds.Get(ctx, "auth:"+email+":reset_password_otp").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("otp not found: %w", err)
		}

		return "", fmt.Errorf("failed to get otp: %w", err)
	}

	return result, nil
}

func (r *authRepository) DeletePasswordResetOTP(ctx context.Context, email string) error {
	if err := r.rds.Del(ctx, "auth:"+email+":reset_password_otp").Err(); err != nil {
		return fmt.Errorf("failed to delete otp: %w", err)
	}

	return nil
}
