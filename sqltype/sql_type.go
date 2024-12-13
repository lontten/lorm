package sqltype

type SqlType int

const (
	Insert SqlType = iota
	Delete
	Update
	Select
)
