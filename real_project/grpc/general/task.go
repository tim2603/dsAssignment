package general

type TaskState int

const (
	BeforeMapping TaskState = iota
	Mapping
	AfterMapping
	Intermediate
	BeforeReducing
	Reducing
	AfterReducing
)

type MapReduceTask struct {
	N_mappers  int
	N_reducers int
	State      TaskState
}
