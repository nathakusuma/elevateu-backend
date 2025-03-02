package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
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
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err = r.createUser(ctx, tx, user); err != nil {
		return err
	}

	// Handle student or mentor creation if role is specified
	switch user.Role {
	case enum.UserRoleStudent:
		if user.Student != nil {
			if err = r.createStudent(ctx, tx, user.ID, user.Student); err != nil {
				return err
			}
		}
	case enum.UserRoleMentor:
		if user.Mentor != nil {
			if err = r.createMentor(ctx, tx, user.ID, user.Mentor); err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *userRepository) createUser(ctx context.Context, tx sqlx.ExtContext, user *entity.User) error {
	_, err := sqlx.NamedExecContext(
		ctx,
		tx,
		`INSERT INTO users (
			id, name, email, password_hash, role
		) VALUES (
			:id, :name, :email, :password_hash, :role
		)`,
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

func (r *userRepository) createStudent(ctx context.Context, tx sqlx.ExtContext, userID uuid.UUID,
	student *entity.Student) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO students (
			user_id, instance, major
		) VALUES (
			$1, $2, $3
		)`,
		userID, student.Instance, student.Major,
	)
	if err != nil {
		return fmt.Errorf("failed to create student: %w", err)
	}

	return nil
}

func (r *userRepository) createMentor(ctx context.Context, tx sqlx.ExtContext, userID uuid.UUID,
	mentor *entity.Mentor) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO mentors (
			user_id, specialization, experience, price
		) VALUES (
			$1, $2, $3, $4
		)`,
		userID, mentor.Specialization, mentor.Experience, mentor.Price,
	)
	if err != nil {
		return fmt.Errorf("failed to create mentor: %w", err)
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
		has_avatar,
		created_at,
		updated_at
	FROM users
	WHERE %s`

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

func (r *userRepository) GetUserByField(ctx context.Context, field string, value interface{}) (*entity.User, error) {
	whereClause := field + " = $1"
	user, err := r.getUserByCondition(ctx, whereClause, value)
	if err != nil {
		return nil, err
	}

	// Fetch additional role-specific data
	if user.Role == enum.UserRoleStudent {
		student, err := r.getStudentByUserID(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		user.Student = student
	} else if user.Role == enum.UserRoleMentor {
		mentor, err := r.getMentorByUserID(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		user.Mentor = mentor
	}

	return user, nil
}

func (r *userRepository) getStudentByUserID(ctx context.Context, userID uuid.UUID) (*entity.Student, error) {
	var student entity.Student

	err := r.db.GetContext(ctx, &student,
		`SELECT instance, major
		FROM students
		WHERE user_id = $1`,
		userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("student data not found: %w", err)
		}
		return nil, err
	}

	return &student, nil
}

func (r *userRepository) getMentorByUserID(ctx context.Context, userID uuid.UUID) (*entity.Mentor, error) {
	var mentor entity.Mentor

	err := r.db.GetContext(ctx, &mentor,
		`SELECT
			specialization, experience, rating, rating_count,
			rating_total, price, balance
		FROM mentors
		WHERE user_id = $1`,
		userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("mentor data not found: %w", err)
		}
		return nil, err
	}

	return &mentor, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, req *dto.UserUpdate) error {
	if req.ID == uuid.Nil {
		return errors.New("cannot update user with empty ID")
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update base user information
	err = r.updateUserDynamic(ctx, tx, req)
	if err != nil {
		return err
	}

	// Update student data if provided
	if req.Student != nil {
		err = r.updateStudentDynamic(ctx, tx, req.ID, *req.Student)
		if err != nil {
			return err
		}
	}

	// Update mentor data if provided
	if req.Mentor != nil {
		err = r.updateMentorDynamic(ctx, tx, req.ID, *req.Mentor)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *userRepository) updateUserDynamic(ctx context.Context, tx sqlx.ExtContext, req *dto.UserUpdate) error {
	// Build dynamic SQL for non-nil fields
	var updates []string
	var args []interface{}
	argIndex := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}

	if req.HasAvatar != nil {
		updates = append(updates, fmt.Sprintf("has_avatar = $%d", argIndex))
		args = append(args, *req.HasAvatar)
		argIndex++
	}

	// Add updated_at timestamp to always update this field
	updates = append(updates, "updated_at = now()")

	// If no fields to update, just return (no need to execute a query)
	if len(args) == 0 {
		return nil
	}

	// Build and execute the query
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d",
		strings.Join(updates, ", "),
		argIndex)
	args = append(args, req.ID)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) updateStudentDynamic(ctx context.Context, tx sqlx.ExtContext, userID uuid.UUID,
	req dto.StudentUpdate) error {
	// Build dynamic SQL for non-nil fields
	var updates []string
	var args []interface{}
	argIndex := 1

	if req.Instance != nil {
		updates = append(updates, fmt.Sprintf("instance = $%d", argIndex))
		args = append(args, *req.Instance)
		argIndex++
	}

	if req.Major != nil {
		updates = append(updates, fmt.Sprintf("major = $%d", argIndex))
		args = append(args, *req.Major)
		argIndex++
	}

	// If no fields to update, just return (no need to execute a query)
	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE students SET %s WHERE user_id = $%d",
		strings.Join(updates, ", "),
		argIndex)
	args = append(args, userID)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}

	return nil
}

func (r *userRepository) updateMentorDynamic(ctx context.Context, tx sqlx.ExtContext, userID uuid.UUID,
	req dto.MentorUpdate) error {
	// Build dynamic SQL for non-nil fields
	var updates []string
	var args []interface{}
	argIndex := 1

	if req.Specialization != nil {
		updates = append(updates, fmt.Sprintf("specialization = $%d", argIndex))
		args = append(args, *req.Specialization)
		argIndex++
	}

	if req.Experience != nil {
		updates = append(updates, fmt.Sprintf("experience = $%d", argIndex))
		args = append(args, *req.Experience)
		argIndex++
	}

	if req.Price != nil {
		updates = append(updates, fmt.Sprintf("price = $%d", argIndex))
		args = append(args, *req.Price)
		argIndex++
	}

	// If no fields to update, just return (no need to execute a query)
	if len(updates) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE mentors SET %s WHERE user_id = $%d",
		strings.Join(updates, ", "),
		argIndex)
	args = append(args, userID)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update mentor: %w", err)
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// The cascade delete will take care of the role-specific tables
	err = r.deleteUser(ctx, tx, id)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *userRepository) deleteUser(ctx context.Context, tx sqlx.ExtContext, id uuid.UUID) error {
	// Tables with foreign keys will be deleted via CASCADE
	res, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
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
