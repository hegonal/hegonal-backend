package queries

type NoRowsAffectedError struct {
	Message string
}

func (e *NoRowsAffectedError) Error() string {
	return e.Message
}