package insert_type

type InsertType int

const (
	Err     InsertType = iota // 报错
	Ignore                    // 忽略
	Update                    // 更新，根据主键，唯一索引更新
	Replace                   // 替换，删除后插入
)
