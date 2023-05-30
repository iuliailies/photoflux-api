package gorm

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// PaginateMultiColumn returns a partial query that takes care of all
// pagination operations on the data. See the README for more details about how
// pagination is implemented.
//
// This function assumes that the tuple of values each record has for the
// requested columns is unique accross the entire database. Failure to satisfy
// this constraint may lead to unexpected behaviour.
func PaginateMultiColumn(db *gorm.DB,
	// TODO: since the column is placed in the query, ensure that it has a
	// certain value to avoid sql injection.
	relation string,
	columns []OrderedColumn,
	limit int,
	before []any,
	after []any,
	order OrderType) (*gorm.DB, error) {

	if len(columns) == 0 {
		return nil, errors.New("at least one column required for pagination")
	}

	// Both before and after were provided values
	if len(before) != 0 && len(after) != 0 {
		if len(before) != len(columns) {
			return nil, fmt.Errorf("invalid cursor length: provided \"before\" cursor with length %d for %d columns",
				len(before), len(columns))
		}
		if len(after) != len(columns) {
			return nil, fmt.Errorf("invalid cursor length: provided \"after\" cursor with length %d for %d columns",
				len(after), len(columns))
		}
		if len(before) != len(after) {
			return nil, errors.New("non-nil \"before\" and \"after\" cursors do not match length")
		}
	}

	if order == OrderDESC {
		tmp := after
		after = before
		before = tmp

		for _, col := range columns {
			db = db.Order(fmt.Sprintf("%s %s", col.Column, col.Order.Reverse().ToDBRepresentation()))
		}

	} else {
		for _, col := range columns {
			db = db.Order(fmt.Sprintf("%s %s", col.Column, col.Order.ToDBRepresentation()))
		}
	}

	// definition required in order for the recursive call
	var buildQueryRecursively func(i int, cursor []any, values *[]any,
		change func(o OrderType) OrderType) string

	// buildQueryRecursively creates the raw query used to paginate the results.
	// See the README for a sample query.
	buildQueryRecursively = func(i int, cursor []any, values *[]any,
		change func(o OrderType) OrderType) string {

		var query string

		// A strict symbol is required for all queries that have an OR clause
		// at the same level, e.g. 'age<10 OR (age=10 AND ...)'. This is
		// necessary for all columns except the last one.
		if i == len(columns)-1 {
			query = fmt.Sprintf("%s %s ?", columns[i].Column, change(columns[i].Order).ToSymbol())
		} else {
			query = fmt.Sprintf("%s %s ?", columns[i].Column, change(columns[i].Order).ToStrictSymbol())
		}

		// This helps avoid sql injection. The query is built with question
		// marks, and this function creates a value list that is substituted
		// for the question marks by GORM.
		//
		// See https://gorm.io/docs/security.html for more details.
		*values = append(*values, cursor[i])
		if i < len(columns)-1 {
			*values = append(*values, cursor[i])
			query += fmt.Sprintf(" OR (%s = ? AND (%s))", columns[i].Column, buildQueryRecursively(i+1, cursor, values, change))
		}
		return query
	}

	// Note that below, len(nil) == 0.

	if len(after) != 0 {
		emptySession := db.Session(&gorm.Session{NewDB: true, DryRun: true})
		values := make([]any, 0)
		query := buildQueryRecursively(0, after, &values, KeepOrder)
		stmt := emptySession.Raw(fmt.Sprintf("%s", query)).Statement.SQL.String()
		db = db.Where(stmt, values...)
	}

	if len(before) != 0 {
		emptySession := db.Session(&gorm.Session{NewDB: true, DryRun: true})
		values := make([]any, 0)
		query := buildQueryRecursively(0, before, &values, ReverseOrder)
		stmt := emptySession.Raw(fmt.Sprintf("%s", query)).Statement.SQL.String()
		db = db.Where(stmt, values...)
	}

	if limit > 0 {
		db = db.Limit(limit)
	}

	return db, nil
}
