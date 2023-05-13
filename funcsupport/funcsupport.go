package funcsupport

// Supplier supplier function
type Supplier[MType any] func() MType

// Predicate predicate function
type Predicate[MType any] func(in MType) bool
