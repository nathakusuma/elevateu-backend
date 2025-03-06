package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/internal/infra/database"
)

type courseProgressRepository struct {
	db *sqlx.DB
}

func NewCourseProgressRepository(db *sqlx.DB) contract.ICourseProgressRepository {
	return &courseProgressRepository{
		db: db,
	}
}

func (r *courseProgressRepository) UpdateVideoProgress(ctx context.Context, txWrapper database.ITransaction,
	progress entity.CourseVideoProgress) (bool, error) {
	tx := txWrapper.GetTx()

	var existingProgress entity.CourseVideoProgress
	query := `
		SELECT student_id, video_id, last_position, is_completed
		FROM course_video_progresses
		WHERE student_id = $1 AND video_id = $2
	`
	err := tx.QueryRowxContext(ctx, query, progress.StudentID, progress.VideoID).StructScan(&existingProgress)

	if errors.Is(err, sql.ErrNoRows) {
		insertQuery := `
			INSERT INTO course_video_progresses (student_id, video_id, last_position, is_completed)
			VALUES ($1, $2, $3, $4)
		`
		_, err = tx.ExecContext(ctx, insertQuery, progress.StudentID, progress.VideoID, progress.LastPosition,
			progress.IsCompleted)
		if err != nil {
			return false, fmt.Errorf("failed to insert video progress: %w", err)
		}

		// Return true if the video is now completed, indicating that the content completion count should be incremented
		return progress.IsCompleted, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check existing video progress: %w", err)
	}

	newlyCompleted := progress.IsCompleted && !existingProgress.IsCompleted

	updateQuery := `
		UPDATE course_video_progresses
		SET last_position = $3
	`
	args := []interface{}{progress.StudentID, progress.VideoID, progress.LastPosition}

	// Once completed, a video should remain completed
	if newlyCompleted {
		updateQuery += `, is_completed = $4`
		args = append(args, progress.IsCompleted)
	}

	updateQuery += ` WHERE student_id = $1 AND video_id = $2`

	_, err = tx.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return false, fmt.Errorf("failed to update video progress: %w", err)
	}

	// Return true if the video is newly completed, indicating that the content completion count should be incremented
	return newlyCompleted, nil
}

func (r *courseProgressRepository) UpdateMaterialProgress(ctx context.Context, txWrapper database.ITransaction,
	progress entity.CourseMaterialProgress) (bool, error) {
	tx := txWrapper.GetTx()

	query := `
		INSERT INTO course_material_progresses (student_id, material_id)
		VALUES ($1, $2)
		ON CONFLICT (student_id, material_id) DO NOTHING
	`
	result, err := tx.ExecContext(ctx, query, progress.StudentID, progress.MaterialID)
	if err != nil {
		return false, fmt.Errorf("failed to update material progress: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	// Return true if the material was newly completed, indicating that the content completion count should be incremented
	return rowsAffected > 0, nil
}

func (r *courseProgressRepository) IncrementCourseProgress(ctx context.Context, txWrapper database.ITransaction,
	courseID, studentID uuid.UUID) (bool, error) {
	tx := txWrapper.GetTx()

	updateQuery := `
		UPDATE course_enrollments
		SET content_completed = content_completed + 1,
			last_accessed_at = NOW()
		WHERE course_id = $1 AND student_id = $2
		RETURNING is_completed
	`
	var wasCompleted bool
	err := tx.QueryRowContext(ctx, updateQuery, courseID, studentID).Scan(&wasCompleted)
	if err != nil {
		return false, fmt.Errorf("failed to update course enrollment: %w", err)
	}

	// Check if all content items for the course are now completed
	if !wasCompleted {
		checkQuery := `
			WITH course_stats AS (
				SELECT
					e.content_completed,
					c.content_count
				FROM course_enrollments e
				JOIN courses c ON e.course_id = c.id
				WHERE e.course_id = $1 AND e.student_id = $2
			)
			SELECT content_completed >= content_count
			FROM course_stats
		`

		var shouldComplete bool
		err = tx.QueryRowContext(ctx, checkQuery, courseID, studentID).Scan(&shouldComplete)
		if err != nil {
			return false, fmt.Errorf("failed to check course completion status: %w", err)
		}

		if shouldComplete {
			markCompletedQuery := `
				UPDATE course_enrollments
				SET is_completed = TRUE
				WHERE course_id = $1 AND student_id = $2
			`
			_, err = tx.ExecContext(ctx, markCompletedQuery, courseID, studentID)
			if err != nil {
				return false, fmt.Errorf("failed to mark course as completed: %w", err)
			}

			// Return true indicating that the course was just completed
			return true, nil
		}
	}

	// Return false if the course wasn't just completed
	return false, nil
}

func (r *courseProgressRepository) GetContentCourseID(ctx context.Context, contentID uuid.UUID,
	contentType string) (uuid.UUID, error) {
	var query string
	if contentType == "video" {
		query = `SELECT course_id FROM course_videos WHERE id = $1`
	} else if contentType == "material" {
		query = `SELECT course_id FROM course_materials WHERE id = $1`
	} else {
		return uuid.Nil, fmt.Errorf("invalid content type: %s", contentType)
	}

	var courseID uuid.UUID
	err := r.db.GetContext(ctx, &courseID, query, contentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("course content not found: %w", err)
		}
		return uuid.Nil, fmt.Errorf("failed to get course ID for content: %w", err)
	}

	return courseID, nil
}

func (r *courseProgressRepository) BatchDecrementCourseProgress(ctx context.Context, txWrapper database.ITransaction,
	courseID uuid.UUID, contentID uuid.UUID, contentType string) error {
	tx := txWrapper.GetTx()

	var query string
	if contentType == "video" {
		query = `
			UPDATE course_enrollments ce
			SET content_completed = content_completed - 1
			FROM course_video_progresses cvp
			WHERE ce.student_id = cvp.student_id
			  AND ce.course_id = $1
			  AND cvp.video_id = $2
			  AND cvp.is_completed = TRUE
		`
	} else if contentType == "material" {
		query = `
			UPDATE course_enrollments ce
			SET content_completed = content_completed - 1
			FROM course_material_progresses cmp
			WHERE ce.student_id = cmp.student_id
			  AND ce.course_id = $1
			  AND cmp.material_id = $2
		`
	} else {
		return fmt.Errorf("invalid content type: %s", contentType)
	}

	_, err := tx.ExecContext(ctx, query, courseID, contentID)
	if err != nil {
		return fmt.Errorf("failed to batch update course progress: %w", err)
	}

	return nil
}
