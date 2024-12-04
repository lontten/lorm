package return_type

type ReturnType int

const (
	PrimaryKey ReturnType = iota // 返回主键
	None                         // 不返回
	ZeroField                    // 返回零值字段
	AllField                     // 返回所有字段
)
