package repository

type Condition interface {
	Args() []any
	String() string
	GetScenarioUserID() int
}
