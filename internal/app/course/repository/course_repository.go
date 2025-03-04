package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"github.com/nathakusuma/elevateu-backend/domain/contract"
	"github.com/nathakusuma/elevateu-backend/domain/dto"
	"github.com/nathakusuma/elevateu-backend/domain/entity"
	"github.com/nathakusuma/elevateu-backend/pkg/sqlutil"
)

type courseRepository struct {
	db *sqlx.DB
}

func NewCourseRepository(conn *sqlx.DB) contract.ICourseRepository {
	return &courseRepository{
		db: conn,
	}
}

func (r *courseRepository) BeginTx() (*sqlx.Tx, error) {
	return r.db.Beginx()
}

func (r *courseRepository) CreateCourse(ctx context.Context, course *entity.Course) error {
	query := `
		INSERT INTO courses (
			id, category_id, title, description, teacher_name
		) VALUES (
			:id, :category_id, :title, :description, :teacher_name
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, course)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "courses_category_id_fkey" {
			return fmt.Errorf("category not found: %w", err)
		}

		return err
	}

	return nil
}

func (r *courseRepository) GetCourseByID(ctx context.Context, id uuid.UUID) (*entity.Course, error) {
	query := `
		SELECT
			c.id, c.category_id, c.title, c.description, c.teacher_name,
			c.rating, c.rating_count, c.total_rating, c.enrollment_count,
			c.content_count, c.total_duration, c.created_at, c.updated_at,
			cat.id AS "category.id", cat.name AS "category.name"
		FROM courses c
		LEFT JOIN categories cat ON c.category_id = cat.id
		WHERE c.id = $1
	`

	var course entity.Course
	course.Category = &entity.Category{}

	err := r.db.QueryRowxContext(ctx, query, id).StructScan(&course)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("course not found")
		}

		return nil, err
	}

	return &course, nil
}

func (r *courseRepository) GetCourses(ctx context.Context, query dto.GetCoursesQuery,
	paginationReq dto.PaginationRequest) ([]*entity.Course, dto.PaginationResponse, error) {

	baseQuery := `
       SELECT
          c.id, c.category_id, c.title, c.description, c.teacher_name,
          c.rating, c.rating_count, c.total_rating, c.enrollment_count,
          c.content_count, c.total_duration, c.created_at, c.updated_at,
          cat.id AS "category.id", cat.name AS "category.name"
       FROM courses c
       LEFT JOIN categories cat ON c.category_id = cat.id
    `

	// WHERE clause based on query parameters
	var whereConditions []string
	var args []interface{}
	argIndex := 1

	if query.CategoryID != uuid.Nil {
		whereConditions = append(whereConditions, fmt.Sprintf("c.category_id = $%d", argIndex))
		args = append(args, query.CategoryID)
		argIndex++
	}

	if query.Title != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("c.title ILIKE $%d", argIndex))
		args = append(args, "%"+query.Title+"%")
		argIndex++
	}

	// cursor-based pagination
	if paginationReq.Cursor != uuid.Nil {
		var operator string
		var orderDirection string

		if paginationReq.Direction == "next" {
			operator = "<"
			orderDirection = "DESC"
		} else {
			operator = ">"
			orderDirection = "ASC"
		}

		var cursorRating float64
		err := r.db.GetContext(ctx, &cursorRating,
			"SELECT total_rating FROM courses WHERE id = $1", paginationReq.Cursor)
		if err != nil {
			return nil, dto.PaginationResponse{}, err
		}

		whereConditions = append(whereConditions,
			fmt.Sprintf("(c.total_rating %s $%d OR (c.total_rating = $%d AND c.id %s $%d))",
				operator, argIndex, argIndex, operator, argIndex+1))
		args = append(args, cursorRating, paginationReq.Cursor)
		argIndex += 2

		sqlQuery := baseQuery
		if len(whereConditions) > 0 {
			sqlQuery += " WHERE " + strings.Join(whereConditions, " AND ")
		}
		sqlQuery += fmt.Sprintf(" ORDER BY c.total_rating %s, c.id %s LIMIT $%d",
			orderDirection, orderDirection, argIndex)
		args = append(args, paginationReq.Limit+1)

		// Execute
		rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
		if err != nil {
			return nil, dto.PaginationResponse{}, err
		}
		defer rows.Close()

		// Process results
		var courses []*entity.Course
		for rows.Next() {
			var course entity.Course
			course.Category = &entity.Category{}

			if err := rows.StructScan(&course); err != nil {
				return nil, dto.PaginationResponse{}, err
			}
			courses = append(courses, &course)
		}

		// hasMore
		hasMore := false
		if len(courses) > paginationReq.Limit {
			hasMore = true
			courses = courses[:paginationReq.Limit]
		}

		// Reverse when "prev"
		if paginationReq.Direction == "prev" {
			for i, j := 0, len(courses)-1; i < j; i, j = i+1, j-1 {
				courses[i], courses[j] = courses[j], courses[i]
			}
		}

		return courses, dto.PaginationResponse{HasMore: hasMore}, nil
	} else {
		sqlQuery := baseQuery
		if len(whereConditions) > 0 {
			sqlQuery += " WHERE " + strings.Join(whereConditions, " AND ")
		}
		sqlQuery += fmt.Sprintf(" ORDER BY c.total_rating DESC, c.id DESC LIMIT $%d", argIndex)
		args = append(args, paginationReq.Limit+1)

		rows, err := r.db.QueryxContext(ctx, sqlQuery, args...)
		if err != nil {
			return nil, dto.PaginationResponse{}, err
		}
		defer rows.Close()

		var courses []*entity.Course
		for rows.Next() {
			var course entity.Course
			course.Category = &entity.Category{}

			if err := rows.StructScan(&course); err != nil {
				return nil, dto.PaginationResponse{}, err
			}
			courses = append(courses, &course)
		}

		hasMore := false
		if len(courses) > paginationReq.Limit {
			hasMore = true
			courses = courses[:paginationReq.Limit]
		}

		return courses, dto.PaginationResponse{HasMore: hasMore}, nil
	}
}

func (r *courseRepository) UpdateCourse(ctx context.Context, tx sqlx.ExtContext, updates *dto.CourseUpdate) error {
	if tx == nil {
		tx = r.db
	}

	builder := sqlutil.NewSQLUpdateBuilder("courses").
		WithUpdatedAt().
		Where("id = ?", updates.ID)

	query, args, err := builder.BuildFromStruct(updates)
	if err != nil {
		return fmt.Errorf("failed to build update query: %w", err)
	}

	// No fields to update (query is empty)
	if query == "" {
		return nil
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "courses_category_id_fkey" {
			return fmt.Errorf("category not found: %w", err)
		}
		return fmt.Errorf("failed to update course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("course not found")
	}

	return nil
}

func (r *courseRepository) DeleteCourse(ctx context.Context, tx sqlx.ExtContext, id uuid.UUID) error {
	if tx == nil {
		tx = r.db
	}

	query := "DELETE FROM courses WHERE id = $1"

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("course not found")
	}

	return nil
}

func (r *courseRepository) CreateEnrollment(ctx context.Context, courseID, studentID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO course_enrollments (course_id, student_id)
				VALUES ($1, $2)`

	_, err = tx.ExecContext(ctx, query, courseID, studentID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "course_enrollments_course_id_fkey":
				return errors.New("course not found")
			case "course_enrollments_pkey":
				return errors.New("student already enrolled in course")
			}
		}

		return err
	}

	// update enrollment count
	query = `UPDATE courses SET enrollment_count = enrollment_count + 1 WHERE id = $1`
	_, err = tx.ExecContext(ctx, query, courseID)
	if err != nil {
		return fmt.Errorf("failed to update enrollment count: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *courseRepository) GetEnrolledCourses(ctx context.Context, studentID uuid.UUID,
	pageReq dto.PaginationRequest) ([]*entity.Course, dto.PaginationResponse, error) {
	baseQuery := `
       SELECT
          c.id, c.category_id, c.title, c.description, c.teacher_name,
          c.rating, c.rating_count, c.total_rating, c.enrollment_count,
          c.content_count, c.total_duration, c.created_at, c.updated_at,
          cat.id AS "category.id", cat.name AS "category.name"
       FROM courses c
       LEFT JOIN categories cat ON c.category_id = cat.id
       JOIN course_enrollments ce ON c.id = ce.course_id
       WHERE ce.student_id = $1
    `

	// cursor-based pagination
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

		var cursorTimestamp time.Time
		err := r.db.GetContext(ctx, &cursorTimestamp,
			"SELECT last_accessed_at FROM course_enrollments WHERE course_id = $1 AND student_id = $2",
			pageReq.Cursor, studentID)
		if err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to get cursor timestamp: %w", err)
		}

		sqlQuery := baseQuery + fmt.Sprintf(
			" AND (ce.last_accessed_at %s $2 OR (ce.last_accessed_at = $2 AND c.id %s $3)) ORDER BY ce.last_accessed_at %s, c.id %s LIMIT $4",
			operator, operator, orderDirection, orderDirection)

		rows, err := r.db.QueryxContext(ctx, sqlQuery, studentID, cursorTimestamp, pageReq.Cursor, pageReq.Limit+1)
		if err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to execute query: %w", err)
		}
		defer rows.Close()

		var courses []*entity.Course
		for rows.Next() {
			var course entity.Course
			course.Category = &entity.Category{}

			if err := rows.StructScan(&course); err != nil {
				return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan row: %w", err)
			}
			courses = append(courses, &course)
		}

		hasMore := false
		if len(courses) > pageReq.Limit {
			hasMore = true
			courses = courses[:pageReq.Limit]
		}

		if pageReq.Direction == "prev" {
			for i, j := 0, len(courses)-1; i < j; i, j = i+1, j-1 {
				courses[i], courses[j] = courses[j], courses[i]
			}
		}

		return courses, dto.PaginationResponse{HasMore: hasMore}, nil
	} else {
		sqlQuery := baseQuery + " ORDER BY ce.last_accessed_at DESC, c.id DESC LIMIT $2"

		rows, err := r.db.QueryxContext(ctx, sqlQuery, studentID, pageReq.Limit+1)
		if err != nil {
			return nil, dto.PaginationResponse{}, fmt.Errorf("failed to execute query: %w", err)
		}
		defer rows.Close()

		var courses []*entity.Course
		for rows.Next() {
			var course entity.Course
			course.Category = &entity.Category{}

			if err := rows.StructScan(&course); err != nil {
				return nil, dto.PaginationResponse{}, fmt.Errorf("failed to scan row: %w", err)
			}
			courses = append(courses, &course)
		}

		hasMore := false
		if len(courses) > pageReq.Limit {
			hasMore = true
			courses = courses[:pageReq.Limit]
		}

		return courses, dto.PaginationResponse{HasMore: hasMore}, nil
	}
}
