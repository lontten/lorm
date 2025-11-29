package return_type

type ReturnType int

const (
	Auto      ReturnType = iota // 返回 自动生成字段
	None                        // 不返回
	ZeroField                   // 返回零值字段
	AllField                    // 返回所有字段
)
