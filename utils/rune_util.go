package utils

// RuneBuffer 高效字符缓冲区
type RuneBuffer struct {
	runes []rune // 直接存储字符
}

// 创建缓冲区
func NewRuneBuffer(capacity int) *RuneBuffer {
	return &RuneBuffer{
		runes: make([]rune, 0, capacity),
	}
}

// 写入字符串
func (b *RuneBuffer) WriteString(s string) {
	// 将字符串转换为 rune 切片并追加
	b.runes = append(b.runes, []rune(s)...)
}

// 获取最后 n 个字符
func (b *RuneBuffer) LastChars(n int) string {
	if n <= 0 {
		return ""
	}
	if len(b.runes) < n {
		return string(b.runes)
	}
	return string(b.runes[len(b.runes)-n:])
}

// 回退最后 n 个字符（仅移除，不返回）
func (b *RuneBuffer) Rewind(n int) {
	if n <= 0 {
		return
	}

	total := len(b.runes)

	// 如果请求回退的数量超过缓冲区长度，则清空缓冲区
	if total <= n {
		b.runes = b.runes[:0] // 清空缓冲区
		return
	}

	// 直接截断切片，移除最后 n 个字符
	b.runes = b.runes[:total-n]
}

// 取回最后 n 个字符（移除并返回）
func (b *RuneBuffer) RetrieveLastChars(n int) string {
	if n <= 0 {
		return ""
	}
	total := len(b.runes)

	if total < n {
		result := string(b.runes)
		b.runes = b.runes[:0] // 清空缓冲区
		return result
	}

	// 获取最后 n 个字符
	result := string(b.runes[total-n:])

	// 截断切片
	b.runes = b.runes[:total-n]

	return result
}

// 获取字符串表示
func (b *RuneBuffer) String() string {
	return string(b.runes)
}

// 获取 rune 切片
func (b *RuneBuffer) Runes() []rune {
	return b.runes
}

// 缓冲区长度（字符数）
func (b *RuneBuffer) Len() int {
	return len(b.runes)
}

// 重置缓冲区
func (b *RuneBuffer) Reset() {
	b.runes = b.runes[:0] // 重用底层数组
}
