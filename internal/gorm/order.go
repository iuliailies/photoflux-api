package gorm

type OrderType int

const (
	OrderASC OrderType = iota
	OrderDESC
)

/*
	These functions are defined to avoid code duplication
	in buildQueryRecursively.
*/

// ReverseOrder returns the reverse of a given order.
func ReverseOrder(o OrderType) OrderType {
	if o == OrderASC {
		return OrderDESC
	}
	return OrderASC
}

// KeepOrder returns the provided order.
func KeepOrder(o OrderType) OrderType {
	return o
}

// Get returns the order type.
func (o OrderType) Get() OrderType {
	return o
}

// Reverse returns the reverse of a given order.
func (o OrderType) Reverse() OrderType {
	if o == OrderASC {
		return OrderDESC
	}
	return OrderASC
}

// ToDBRepresentation returns the internal representation of the order in a db
// query.
func (o OrderType) ToDBRepresentation() string {
	if o == OrderASC {
		return "ASC"
	}
	return "DESC"
}

// ToStrictSymbol returns a symbol that is used to indicate the relation of the
// records that follow after some record based on the given ordering, e.g. if
// the order is ascending, for all items "e" that follow after "a", the relation
// e > a holds true, thus the symbol ">" is returned.
func (o OrderType) ToStrictSymbol() string {
	if o == OrderASC {
		return ">"
	}
	return "<"
}

// ToSymbol does the same as ToStrictSymbol, except that it returns the symbol
// of a non-strict inequality (i.e. it considers "a" to follow after "a").
func (o OrderType) ToSymbol() string {
	if o == OrderASC {
		return ">="
	}
	return "<="
}

// OrderedColumn holds both a column and its requested ordering.
type OrderedColumn struct {
	Column string
	Order  OrderType
}
