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
	"github.com/nathakusuma/elevateu-backend/domain/entity"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(conn *sqlx.DB) contract.IUserRepository {
	return &userRepository{
		db: conn,
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return r.createUser(ctx, r.db, user)
}

func (r *userRepository) createUser(ctx context.Context, tx sqlx.ExtContext, user *entity.User) error {
	_, err := sqlx.NamedExecContext(
		ctx,
		tx,
		`INSERT INTO users (
                   id, name, email, password_hash, role
                   ) VALUES (:id, :name, :email, :password_hash, :role)`,
		user,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "users_email_key" {
			return fmt.Errorf("conflict email: %w", err)
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) getUserByCondition(ctx context.Context, whereClause string,
	args ...interface{}) (*entity.User, error) {
	var user entity.User

	baseQuery := `SELECT
        id,
        name,
        email,
        password_hash,
        role,
        bio,
        avatar_url,
        created_at,
        updated_at
    FROM users
    WHERE %s
    AND deleted_at IS NULL`

	statement := fmt.Sprintf(baseQuery, whereClause)

	err := r.db.GetContext(ctx, &user, statement, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByField(ctx context.Context, field, value string) (*entity.User, error) {
	whereClause := field + " = $1"
	return r.getUserByCondition(ctx, whereClause, value)
}

func (r *userRepository) updateUser(ctx context.Context, tx sqlx.ExtContext, user *entity.User) error {
	_, err := sqlx.NamedExecContext(
		ctx,
		tx,
		`UPDATE users
		SET name = :name,
			email = :email,
			password_hash = :password_hash,
			role = :role,
			bio = :bio,
			avatar_url = :avatar_url,
			updated_at = now()
		WHERE id = :id`,
		user,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.updateUser(ctx, r.db, user)
}

func (r *userRepository) deleteUser(ctx context.Context, tx sqlx.ExtContext, id uuid.UUID) error {
	res, err := tx.ExecContext(ctx,
		`UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.deleteUser(ctx, r.db, id)
}
