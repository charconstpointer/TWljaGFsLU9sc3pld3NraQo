package fw

type Result struct {
	value string
}

func NewResult(value string) Result {
	return Result{
		value: value,
	}
}
