package sql_type

type SqlType int

const (
	Insert SqlType = iota
	Delete
	Update
	Select
)
