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
			c.content_count, c.created_at, c.updated_at,
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
			c.content_count, c.created_at, c.updated_at,
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
			operator = ">"
			orderDirection = "ASC"
		} else {
			operator = "<"
			orderDirection = "DESC"
		}

		whereConditions = append(whereConditions, fmt.Sprintf("c.id %s $%d", operator, argIndex))
		args = append(args, paginationReq.Cursor)
		argIndex++

		// Build final query (with WHERE clause and pagination)
		sqlQuery := baseQuery
		if len(whereConditions) > 0 {
			sqlQuery += " WHERE " + strings.Join(whereConditions, " AND ")
		}
		sqlQuery += fmt.Sprintf(" ORDER BY c.id %s LIMIT $%d", orderDirection, argIndex)
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
		// When no cursor is provided, use only LIMIT
		sqlQuery := baseQuery
		if len(whereConditions) > 0 {
			sqlQuery += " WHERE " + strings.Join(whereConditions, " AND ")
		}
		sqlQuery += fmt.Sprintf(" ORDER BY c.id ASC LIMIT $%d", argIndex)
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

		return courses, dto.PaginationResponse{HasMore: hasMore}, nil
	}
}

func (r *courseRepository) UpdateCourse(ctx context.Context, tx sqlx.ExtContext, updates *dto.CourseUpdate) error {
	if tx == nil {
		tx = r.db
	}

	// Dynamic update query based on which fields are provided
	query := "UPDATE courses SET updated_at = NOW()"
	var args []interface{}
	argIndex := 1

	if updates.CategoryID != nil {
		query += fmt.Sprintf(", category_id = $%d", argIndex)
		args = append(args, *updates.CategoryID)
		argIndex++
	}

	if updates.Title != nil {
		query += fmt.Sprintf(", title = $%d", argIndex)
		args = append(args, *updates.Title)
		argIndex++
	}

	if updates.Description != nil {
		query += fmt.Sprintf(", description = $%d", argIndex)
		args = append(args, *updates.Description)
		argIndex++
	}

	if updates.TeacherName != nil {
		query += fmt.Sprintf(", teacher_name = $%d", argIndex)
		args = append(args, *updates.TeacherName)
		argIndex++
	}

	// Add WHERE id
	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, updates.ID)

	// Execute
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "courses_category_id_fkey" {
			return fmt.Errorf("category not found: %w", err)
		}

		return err
	}

	// Check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
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
