package collection_core

import "time"

const DefaultRequestTimeout = 10 * time.Second

type Error string

func (e Error) Error() string {
	return string(e)
}

type Operator string

const (
	GreaterThan      Operator = ">"
	LessThan         Operator = "<"
	GreaterThanEqual Operator = ">="
	LessThanEqual    Operator = "<="
	EqualTo          Operator = "=="
	NotEqualTo       Operator = "!="
	NotIn            Operator = "not-in"
	In               Operator = "in"
	ArrayContains    Operator = "array-contains"
	ArrayContainsAny Operator = "array-contains-any"
)

func (o Operator) ToString() string {
	return string(o)
}

type Direction int32

const (
	ASC  Direction = 1
	DESC Direction = 2
)

func (o Direction) ToInt32() int32 {
	return int32(o)
}
