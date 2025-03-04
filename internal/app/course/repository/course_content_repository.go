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
	"github.com/nathakusuma/elevateu-backend/pkg/sqlutil"
)

type courseContentRepository struct {
	db *sqlx.DB
}

func NewCourseContentRepository(conn *sqlx.DB) contract.ICourseContentRepository {
	return &courseContentRepository{
		db: conn,
	}
}

func (r *courseContentRepository) CreateVideo(ctx context.Context, video *entity.CourseVideo) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO course_videos (
			id, course_id, title, description, duration, is_free, "order"
		) VALUES (
			:id, :course_id, :title, :description, :duration, :is_free, :order
		)
	`

	_, err = tx.NamedExecContext(ctx, query, video)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "course_videos_course_id_fkey" {
			return fmt.Errorf("course not found: %w", err)
		}

		return fmt.Errorf("failed to create video: %w", err)
	}

	// After creating the video, update the course's content_count and total_duration
	updateQuery := `
		UPDATE courses
		SET content_count = content_count + 1,
		    total_duration = total_duration + $1,
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, updateQuery, video.Duration, video.CourseID)
	if err != nil {
		return fmt.Errorf("failed to update course stats: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *courseContentRepository) UpdateVideo(ctx context.Context, id uuid.UUID, updates dto.CourseVideoUpdate) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the current video to calculate duration change if needed
	var currentVideo entity.CourseVideo
	getQuery := `SELECT course_id, duration FROM course_videos WHERE id = $1`
	err = tx.GetContext(ctx, &currentVideo, getQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("video not found")
		}
		return err
	}

	builder := sqlutil.NewSQLUpdateBuilder("course_videos").
		WithUpdatedAt().
		Where("id = ?", id)

	query, args, err := builder.BuildFromStruct(updates)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// No fields to update (query is empty)
	if query == "" {
		return tx.Commit()
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}

	// If duration was updated, update the course's total_duration
	if updates.Duration != nil {
		durationDiff := *updates.Duration - currentVideo.Duration
		if durationDiff != 0 {
			updateCourseQuery := `
				UPDATE courses
				SET total_duration = total_duration + $1,
				    updated_at = NOW()
				WHERE id = $2
			`
			_, err = tx.ExecContext(ctx, updateCourseQuery, durationDiff, currentVideo.CourseID)
			if err != nil {
				return fmt.Errorf("failed to update course duration: %w", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *courseContentRepository) DeleteVideo(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the video to update course stats after deletion
	var video entity.CourseVideo
	getQuery := `SELECT course_id, duration FROM course_videos WHERE id = $1`
	err = tx.GetContext(ctx, &video, getQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("video not found")
		}
		return fmt.Errorf("failed to get video: %w", err)
	}

	deleteQuery := "DELETE FROM course_videos WHERE id = $1"
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}

	// Update the course's content_count and total_duration
	updateQuery := `
		UPDATE courses
		SET content_count = content_count - 1,
		    total_duration = total_duration - $1,
		    updated_at = NOW()
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, updateQuery, video.Duration, video.CourseID)
	if err != nil {
		return fmt.Errorf("failed to update course stats: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *courseContentRepository) GetVideoByID(ctx context.Context, id uuid.UUID) (*entity.CourseVideo, error) {
	var video entity.CourseVideo
	query := `
		SELECT id, course_id, title, description, duration, is_free, "order", created_at, updated_at
		FROM course_videos
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &video, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("video not found")
		}
		return nil, fmt.Errorf("failed to get video: %w", err)
	}

	return &video, nil
}

func (r *courseContentRepository) CreateMaterial(ctx context.Context, material *entity.CourseMaterial) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO course_materials (
			id, course_id, title, subtitle, is_free, "order"
		) VALUES (
			:id, :course_id, :title, :subtitle, :is_free, :order
		)
	`

	_, err = tx.NamedExecContext(ctx, query, material)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "course_materials_course_id_fkey" {
			return fmt.Errorf("course not found: %w", err)
		}

		return err
	}

	// After creating the material, update the course's content_count
	updateQuery := `
		UPDATE courses
		SET content_count = content_count + 1,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, updateQuery, material.CourseID)
	if err != nil {
		return fmt.Errorf("failed to update course stats: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *courseContentRepository) UpdateMaterial(ctx context.Context, id uuid.UUID,
	updates dto.CourseMaterialUpdate) error {
	builder := sqlutil.NewSQLUpdateBuilder("course_materials").
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

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update material: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("material not found")
	}

	return nil
}

func (r *courseContentRepository) DeleteMaterial(ctx context.Context, id uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the material to update course stats after deletion
	var material entity.CourseMaterial
	getQuery := `SELECT course_id FROM course_materials WHERE id = $1`
	err = tx.GetContext(ctx, &material, getQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("material not found")
		}
		return fmt.Errorf("failed to get material: %w", err)
	}

	deleteQuery := "DELETE FROM course_materials WHERE id = $1"
	_, err = tx.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return fmt.Errorf("failed to delete material: %w", err)
	}

	// Update the course's content_count
	updateQuery := `
		UPDATE courses
		SET content_count = content_count - 1,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, updateQuery, material.CourseID)
	if err != nil {
		return fmt.Errorf("failed to update course stats: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *courseContentRepository) GetMaterialByID(ctx context.Context, id uuid.UUID) (*entity.CourseMaterial, error) {
	var material entity.CourseMaterial
	query := `
		SELECT id, course_id, title, subtitle, is_free, "order", created_at, updated_at
		FROM course_materials
		WHERE id = $1
	`
	err := r.db.GetContext(ctx, &material, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("material not found")
		}
		return nil, fmt.Errorf("failed to get material: %w", err)
	}

	return &material, nil
}

func (r *courseContentRepository) GetCourseContents(ctx context.Context,
	courseID uuid.UUID) ([]*entity.CourseVideo, []*entity.CourseMaterial, error) {
	// Check if course exists
	courseExistsQuery := `SELECT EXISTS(SELECT 1 FROM courses WHERE id = $1)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, courseExistsQuery, courseID)
	if err != nil {
		return nil, nil, err
	}

	if !exists {
		return nil, nil, errors.New("course not found")
	}

	// Get videos
	videosQuery := `
		SELECT id, course_id, title, description, duration, is_free, "order", created_at, updated_at
		FROM course_videos
		WHERE course_id = $1
		ORDER BY "order" ASC
	`
	var videos []*entity.CourseVideo
	err = r.db.SelectContext(ctx, &videos, videosQuery, courseID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get course videos: %w", err)
	}

	// Get materials
	materialsQuery := `
		SELECT id, course_id, title, subtitle, is_free, "order", created_at, updated_at
		FROM course_materials
		WHERE course_id = $1
		ORDER BY "order" ASC
	`
	var materials []*entity.CourseMaterial
	err = r.db.SelectContext(ctx, &materials, materialsQuery, courseID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get course materials: %w", err)
	}

	return videos, materials, nil
}
