package types

type Change[T any] struct {
	Value T
	Changed bool 
}
