package gorm

import (
	"gorm.io/gorm"
)

// The Paginate functions doesn't apply a limit if it receives a limit
// parameter smaller or equal to 0, so setting this to -1 will ensure
// that all elements are queried even when the limit is increased by
// 1 when calling the Paginate function.
const NO_LIMIT = -1

/*
	Implementation based on the following profile:
	https://jsonapi.org/profiles/ethanresnick/cursor-pagination/

	See the profile for more details.
*/

// equalCursors checks whether two cursors are equal.
func equalCursors(s1, s2 []any) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func ListMultiColumn[T any](db *gorm.DB,
	relation string,
	limit int,
	columns []OrderedColumn,
	retrieve_cursor func(*T) []any,
	before []any,
	after []any,
	order OrderType) (elems []T, previous, next []any, err error) {

	// If only before is specified, we need to do a reverse request in order to
	// get the correct elements for the previous batch.
	reversed := false

	if len(after) == 0 && len(before) != 0 {
		after, before = before, after
		order = ReverseOrder(order)
		reversed = true
	}

	// The library will treat any negative limit as if no limit was provided,
	// but note that the backend must respond with an "invalid query parameter
	// error" to any client that request a negative page size.
	if limit <= 0 {
		// The Paginate functions doesn't apply a limit if it receives a limit
		// parameter smaller or equal to 0, so setting this to -1 will ensure
		// that the call below would still receive all elements.
		limit = NO_LIMIT
	}

	extendedSize := limit + 1
	elems = make([]T, 0, extendedSize)
	// An extra element is requested to determine if the end of the listing had
	// been reached. This allows returning a null value for next, if the result
	// set has the same size as the provided limit, but also includes the last
	// record from the database.
	db, err = PaginateMultiColumn(
		db, relation, columns, extendedSize, before, after, order)

	if err != nil {
		return nil, nil, nil, err
	}

	err = db.Scan(&elems).Error
	if err != nil {
		return nil, nil, nil, err
	}

	// Fill in prev and next, if necessary.
	if len(elems) > 0 {

		// Only provide previous if we didn't start at the beginning. Setting
		// after to null is guaranteed to retrieve the first element.
		if len(after) > 0 {
			previous = retrieve_cursor(&elems[0])
		}

		// Only provide next if last element was not retrieved. Note that if
		// before is specified, it is guaranteed that all elements are not
		// retrieved.
		if (len(elems) == extendedSize) || equalCursors(before, retrieve_cursor(&elems[len(elems)-1])) {
			next = retrieve_cursor(&elems[len(elems)-1])
			elems = elems[:len(elems)-1]
		}
	}

	// If the request was reversed, we need to undo the reversion for return
	// values.
	if reversed {
		l := len(elems)
		{
			tmp := make([]T, 0, l)
			for i := len(elems) - 1; i >= 0; i -= 1 {
				tmp = append(tmp, elems[i])
			}
			elems = tmp
		}
		previous, next = next, previous

	}

	return
}
