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
	"github.com/nathakusuma/elevateu-backend/domain/enum"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
	"github.com/nathakusuma/elevateu-backend/pkg/sqlutil"
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
			user_id, address, specialization, current_job, company, gender
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)`,
		userID, mentor.Address, mentor.Specialization, mentor.CurrentJob, mentor.Company, mentor.Gender,
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
		u.id,
		u.name,
		u.email,
		u.password_hash,
		u.role,
		u.has_avatar,
		u.created_at,
		u.updated_at,
		s.instance,
		s.major,
		s.point,
		s.subscribed_boost_until,
		s.subscribed_challenge_until,
		m.address,
		m.specialization,
		m.current_job,
		m.company,
		m.bio,
		m.gender,
		m.rating,
		m.rating_count,
		m.rating_total,
		m.price,
		m.balance
	FROM users u
	LEFT JOIN students s ON u.id = s.user_id AND u.role = 'student'
	LEFT JOIN mentors m ON u.id = m.user_id AND u.role = 'mentor'
	WHERE %s`

	statement := fmt.Sprintf(baseQuery, whereClause)

	rows, err := r.db.QueryxContext(ctx, statement, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying user: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("user not found: %w", sql.ErrNoRows)
	}

	// Struct to scan all fields
	type UserJoin struct {
		ID                       uuid.UUID       `db:"id"`
		Name                     string          `db:"name"`
		Email                    string          `db:"email"`
		PasswordHash             string          `db:"password_hash"`
		Role                     enum.UserRole   `db:"role"`
		HasAvatar                bool            `db:"has_avatar"`
		CreatedAt                time.Time       `db:"created_at"`
		UpdatedAt                time.Time       `db:"updated_at"`
		Instance                 sql.NullString  `db:"instance"`
		Major                    sql.NullString  `db:"major"`
		Point                    sql.NullInt64   `db:"point"`
		SubscribedBoostUntil     sql.NullTime    `db:"subscribed_boost_until"`
		SubscribedChallengeUntil sql.NullTime    `db:"subscribed_challenge_until"`
		Address                  sql.NullString  `db:"address"`
		Specialization           sql.NullString  `db:"specialization"`
		CurrentJob               sql.NullString  `db:"current_job"`
		Company                  sql.NullString  `db:"company"`
		Bio                      sql.NullString  `db:"bio"`
		Gender                   sql.NullString  `db:"gender"`
		Rating                   sql.NullFloat64 `db:"rating"`
		RatingCount              sql.NullInt64   `db:"rating_count"`
		RatingTotal              sql.NullFloat64 `db:"rating_total"`
		Price                    sql.NullInt64   `db:"price"`
		Balance                  sql.NullInt64   `db:"balance"`
	}

	var userJoin UserJoin
	if err := rows.StructScan(&userJoin); err != nil {
		return nil, fmt.Errorf("error scanning user: %w", err)
	}

	user = entity.User{
		ID:           userJoin.ID,
		Name:         userJoin.Name,
		Email:        userJoin.Email,
		PasswordHash: userJoin.PasswordHash,
		Role:         userJoin.Role,
		HasAvatar:    userJoin.HasAvatar,
		CreatedAt:    userJoin.CreatedAt,
		UpdatedAt:    userJoin.UpdatedAt,
	}

	if user.Role == enum.UserRoleStudent && userJoin.Instance.Valid {
		user.Student = &entity.Student{
			Instance:                 userJoin.Instance.String,
			Major:                    userJoin.Major.String,
			Point:                    int(userJoin.Point.Int64),
			SubscribedBoostUntil:     userJoin.SubscribedBoostUntil.Time,
			SubscribedChallengeUntil: userJoin.SubscribedChallengeUntil.Time,
		}
	}

	if user.Role == enum.UserRoleMentor && userJoin.Specialization.Valid {
		var bio *string
		if userJoin.Bio.Valid {
			bio = &userJoin.Bio.String
		}

		user.Mentor = &entity.Mentor{
			Address:        userJoin.Address.String,
			Specialization: userJoin.Specialization.String,
			CurrentJob:     userJoin.CurrentJob.String,
			Company:        userJoin.Company.String,
			Bio:            bio,
			Gender:         userJoin.Gender.String,
			Rating:         userJoin.Rating.Float64,
			RatingCount:    int(userJoin.RatingCount.Int64),
			RatingTotal:    userJoin.RatingTotal.Float64,
			Price:          int(userJoin.Price.Int64),
			Balance:        int(userJoin.Balance.Int64),
		}
	}

	return &user, nil
}

func (r *userRepository) GetUserByField(ctx context.Context, field string, value interface{}) (*entity.User, error) {
	whereClause := field + " = $1"
	user, err := r.getUserByCondition(ctx, whereClause, value)
	if err != nil {
		return nil, err
	}

	return user, nil
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
	builder := sqlutil.NewSQLUpdateBuilder("users").
		WithUpdatedAt().
		Where("id = ?", req.ID)

	query, args, err := builder.BuildFromStruct(req)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// No fields to update (query is empty)
	if query == "" {
		return nil
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *userRepository) updateStudentDynamic(ctx context.Context, tx sqlx.ExtContext, userID uuid.UUID,
	req dto.StudentUpdate) error {
	builder := sqlutil.NewSQLUpdateBuilder("students").
		Where("user_id = ?", userID)

	query, args, err := builder.BuildFromStruct(&req)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// No fields to update (query is empty)
	if query == "" {
		return nil
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}

	return nil
}

func (r *userRepository) updateMentorDynamic(ctx context.Context, tx sqlx.ExtContext, userID uuid.UUID,
	req dto.MentorUpdate) error {
	builder := sqlutil.NewSQLUpdateBuilder("mentors").
		Where("user_id = ?", userID)

	query, args, err := builder.BuildFromStruct(&req)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// No fields to update (query is empty)
	if query == "" {
		return nil
	}

	_, err = tx.ExecContext(ctx, query, args...)
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

func (r *userRepository) AddPoint(ctx context.Context, txWrapper database.ITransaction, userID uuid.UUID,
	point int) error {
	tx := txWrapper.GetTx()

	_, err := tx.ExecContext(ctx, `UPDATE students SET point = point + $1 WHERE user_id = $2`, point, userID)
	if err != nil {
		return fmt.Errorf("failed to add point: %w", err)
	}

	return nil
}

func (r *userRepository) GetTopPoints(ctx context.Context, limit int) ([]*entity.User, error) {
	query := `
		SELECT
			u.id,
			u.name,
			u.email,
			u.role,
			s.instance AS "student.instance",
			s.major AS "student.major",
			s.point AS "student.point"
		FROM students s
		JOIN users u ON s.user_id = u.id
		ORDER BY s.point DESC
		LIMIT $1
	`

	rows, err := r.db.QueryxContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top students: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err = rows.StructScan(&user); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) GetMentors(ctx context.Context,
	pageReq dto.PaginationRequest) ([]*entity.User, dto.PaginationResponse, error) {
	baseQuery := `
		SELECT
			u.id,
			u.name,
			u.email,
			u.password_hash,
			u.role,
			u.has_avatar,
			u.created_at,
			u.updated_at,
			m.address AS "mentor.address",
			m.specialization AS "mentor.specialization",
			m.current_job AS "mentor.current_job",
			m.company AS "mentor.company",
			m.bio AS "mentor.bio",
			m.gender AS "mentor.gender",
			m.rating AS "mentor.rating",
			m.rating_count AS "mentor.rating_count",
			m.rating_total AS "mentor.rating_total",
			m.price AS "mentor.price",
			m.balance AS "mentor.balance"
		FROM users u
		JOIN mentors m ON u.id = m.user_id
		WHERE u.role = 'mentor'
	`

	var sqlQuery string
	var args []interface{}

	if pageReq.Cursor != uuid.Nil {
		var cursorRatingTotal float64
		err := r.db.GetContext(ctx, &cursorRatingTotal,
			`SELECT m.rating_total FROM mentors m JOIN users u ON m.user_id = u.id WHERE u.id = $1`,
			pageReq.Cursor)
		if err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to get cursor rating total: %w", err)
		}

		var operator string
		var orderDirection string

		if pageReq.Direction == "next" {
			operator = "<"
			orderDirection = "DESC"
		} else {
			operator = ">"
			orderDirection = "ASC"
		}

		sqlQuery = baseQuery + fmt.Sprintf(
			` AND (m.rating_total %s $1 OR (m.rating_total = $1 AND u.id %s $2))
			 ORDER BY m.rating_total %s, u.id %s LIMIT $3`,
			operator, operator, orderDirection, orderDirection)
		args = append(args, cursorRatingTotal, pageReq.Cursor, pageReq.Limit+1)
	} else {
		sqlQuery = baseQuery + " ORDER BY m.rating_total DESC, u.id DESC LIMIT $1"
		args = append(args, pageReq.Limit+1)
	}

	rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, dto.PaginationResponse{}, fmt.Errorf("failed to query mentors: %w", err)
	}
	defer rows.Close()

	var mentors []*entity.User
	for rows.Next() {
		var user entity.User
		user.Mentor = &entity.Mentor{}

		if err = rows.StructScan(&user); err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan mentor: %w", err)
		}

		mentors = append(mentors, &user)
	}

	hasMore := false
	if len(mentors) > pageReq.Limit {
		hasMore = true
		mentors = mentors[:pageReq.Limit]
	}

	if pageReq.Direction == "prev" && pageReq.Cursor != uuid.Nil {
		// Reverse the results for "prev" direction
		for i, j := 0, len(mentors)-1; i < j; i, j = i+1, j-1 {
			mentors[i], mentors[j] = mentors[j], mentors[i]
		}
	}

	return mentors, dto.PaginationResponse{HasMore: hasMore}, nil
}
