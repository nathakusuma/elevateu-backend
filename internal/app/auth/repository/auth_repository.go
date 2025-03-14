package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/domain/enum"
)

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) contract.IAuthRepository {
	return &authRepository{
		db: db,
	}
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
	baseQuery := `SELECT
        s.token,
        s.user_id,
        s.created_at,
        s.expires_at,
        u.id,
        u.name,
        u.email,
        u.password_hash,
        u.role,
        u.has_avatar,
        u.created_at as user_created_at,
        u.updated_at as user_updated_at,
        st.instance,
        st.major,
        st.point,
        st.subscribed_boost_until,
        st.subscribed_challenge_until,
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
    FROM auth_sessions s
    JOIN users u ON u.id = s.user_id
    LEFT JOIN students st ON u.id = st.user_id AND u.role = 'student'
    LEFT JOIN mentors m ON u.id = m.user_id AND u.role = 'mentor'
    WHERE s.token = $1`

	rows, err := r.db.QueryxContext(ctx, baseQuery, token)
	if err != nil {
		return nil, fmt.Errorf("error querying auth session: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("auth session not found: %w", sql.ErrNoRows)
	}

	// Temporary struct to handle the flat join result
	type SessionJoin struct {
		// Auth Session fields
		Token     string    `db:"token"`
		UserID    uuid.UUID `db:"user_id"`
		CreatedAt time.Time `db:"created_at"`
		ExpiresAt time.Time `db:"expires_at"`

		// User fields
		ID            uuid.UUID     `db:"id"`
		Name          string        `db:"name"`
		Email         string        `db:"email"`
		PasswordHash  string        `db:"password_hash"`
		Role          enum.UserRole `db:"role"`
		HasAvatar     bool          `db:"has_avatar"`
		UserCreatedAt time.Time     `db:"user_created_at"`
		UserUpdatedAt time.Time     `db:"user_updated_at"`

		// Student fields
		Instance                 sql.NullString `db:"instance"`
		Major                    sql.NullString `db:"major"`
		Point                    sql.NullInt64  `db:"point"`
		SubscribedBoostUntil     sql.NullTime   `db:"subscribed_boost_until"`
		SubscribedChallengeUntil sql.NullTime   `db:"subscribed_challenge_until"`

		// Mentor fields
		Address        sql.NullString  `db:"address"`
		Specialization sql.NullString  `db:"specialization"`
		CurrentJob     sql.NullString  `db:"current_job"`
		Company        sql.NullString  `db:"company"`
		Bio            sql.NullString  `db:"bio"`
		Gender         sql.NullString  `db:"gender"`
		Rating         sql.NullFloat64 `db:"rating"`
		RatingCount    sql.NullInt64   `db:"rating_count"`
		RatingTotal    sql.NullFloat64 `db:"rating_total"`
		Price          sql.NullInt64   `db:"price"`
		Balance        sql.NullInt64   `db:"balance"`
	}

	var join SessionJoin
	if err = rows.StructScan(&join); err != nil {
		return nil, fmt.Errorf("error scanning auth session: %w", err)
	}

	user := entity.User{
		ID:           join.ID,
		Name:         join.Name,
		Email:        join.Email,
		PasswordHash: join.PasswordHash,
		Role:         join.Role,
		HasAvatar:    join.HasAvatar,
		CreatedAt:    join.UserCreatedAt,
		UpdatedAt:    join.UserUpdatedAt,
	}

	if join.Role == enum.UserRoleStudent && join.Instance.Valid {
		user.Student = &entity.Student{
			Instance: join.Instance.String,
			Major:    join.Major.String,
		}
	}

	if join.Role == enum.UserRoleMentor && join.Specialization.Valid {
		var bio *string
		if join.Bio.Valid {
			bio = &join.Bio.String
		}

		user.Mentor = &entity.Mentor{
			Address:        join.Address.String,
			Specialization: join.Specialization.String,
			CurrentJob:     join.CurrentJob.String,
			Company:        join.Company.String,
			Bio:            bio,
			Gender:         join.Gender.String,
			Rating:         join.Rating.Float64,
			RatingCount:    int(join.RatingCount.Int64),
			RatingTotal:    join.RatingTotal.Float64,
			Price:          int(join.Price.Int64),
			Balance:        int(join.Balance.Int64),
		}
	}

	authSession := entity.AuthSession{
		Token:     join.Token,
		UserID:    join.UserID,
		CreatedAt: join.CreatedAt,
		ExpiresAt: join.ExpiresAt,
		User:      user,
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
