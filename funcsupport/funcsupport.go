package funcsupport

// Supplier supplier function
type Supplier[MType any] func() MType

// Predicate predicate function
type Predicate[MType any] func(in MType) bool

type Function[In, Out any] func(in In) Out
type BiFunction[In, In2, Out any] func(in In, in2 In2) Out
type BiFunctionDual[In, In2, Out, Out2 any] func(in In, in2 In2) (Out, Out2)

type Consumer[MType any] func(in MType)
type BiConsumer[AType, BType any] func(a AType, b BType)
