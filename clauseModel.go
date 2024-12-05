package lorm

// 句子类型，用于whereBuilder
type clauseType int

const (
	Eq clauseType = iota
	Neq
	Less
	LessEq
	Greater
	GreaterEq
	Like
	NotLike
	In
	NotIn
	Between
	NotBetween
	IsNull
	IsNotNull
	IsFalse

	// Contains 包含
	// pg 独有
	// [1] @< [1,2]
	Contains
)
