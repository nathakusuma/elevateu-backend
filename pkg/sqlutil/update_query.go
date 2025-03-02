package sqlutil

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SQLUpdateBuilder struct {
	table        string
	whereParts   []string
	queryParts   []string
	args         []any
	argIndex     int
	setUpdatedAt bool
}

func NewSQLUpdateBuilder(table string) *SQLUpdateBuilder {
	return &SQLUpdateBuilder{
		table:    table,
		argIndex: 1,
	}
}

func (b *SQLUpdateBuilder) WithUpdatedAt() *SQLUpdateBuilder {
	b.setUpdatedAt = true
	return b
}

func (b *SQLUpdateBuilder) Where(condition string, args ...any) *SQLUpdateBuilder {
	placeholderIndex := b.argIndex
	for strings.Contains(condition, "?") {
		condition = strings.Replace(condition, "?", fmt.Sprintf("$%d", placeholderIndex), 1)
		placeholderIndex++
	}

	b.whereParts = append(b.whereParts, condition)
	b.args = append(b.args, args...)
	b.argIndex += len(args)
	return b
}

func (b *SQLUpdateBuilder) BuildFromStruct(dto any) (string, []any, error) {
	v := reflect.ValueOf(dto)
	if v.Kind() != reflect.Ptr && v.Kind() != reflect.Struct {
		return "", nil, errors.New("dto must be a struct or pointer to struct")
	}

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", nil, errors.New("dto cannot be nil")
		}
		v = v.Elem()
		if v.Kind() != reflect.Struct {
			return "", nil, errors.New("dto must be a pointer to struct")
		}
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		// Skip unexported fields
		if !fieldValue.CanInterface() {
			continue
		}

		dbTag := fieldType.Tag.Get("db")
		if dbTag == "-" {
			continue
		}
		if dbTag == "" {
			dbTag = strings.ToLower(fieldType.Name)
		}

		if dbTag == "id" || strings.HasPrefix(dbTag, "_") {
			continue
		}

		switch fieldValue.Kind() {
		case reflect.Ptr:
			if fieldValue.IsNil() {
				continue
			}
			b.addField(dbTag, fieldValue.Elem().Interface())
		case reflect.Slice, reflect.Map, reflect.Interface, reflect.Chan:
			if fieldValue.IsNil() {
				continue
			}
			b.addField(dbTag, fieldValue.Interface())
		default:
			switch fieldType.Type {
			case reflect.TypeOf(uuid.UUID{}):
				if fieldValue.Interface() == uuid.Nil {
					continue
				}
			case reflect.TypeOf(time.Time{}):
				if fieldValue.Interface().(time.Time).IsZero() {
					continue
				}
			}
			b.addField(dbTag, fieldValue.Interface())
		}
	}

	if b.setUpdatedAt {
		b.queryParts = append(b.queryParts, "updated_at = now()")
	}

	// If no fields to update, return a special case
	if len(b.queryParts) == 0 {
		return "", nil, nil
	}

	query := fmt.Sprintf("UPDATE %s SET %s", b.table, strings.Join(b.queryParts, ", "))

	if len(b.whereParts) > 0 {
		query += " WHERE " + strings.Join(b.whereParts, " AND ")
	}

	return query, b.args, nil
}

// addField adds a field to the update query
func (b *SQLUpdateBuilder) addField(fieldName string, value interface{}) {
	b.queryParts = append(b.queryParts, fmt.Sprintf("%s = $%d", fieldName, b.argIndex))
	b.args = append(b.args, value)
	b.argIndex++
}
