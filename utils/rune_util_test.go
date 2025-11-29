package utils

import (
	"fmt"
	"testing"
)

func TestNewRuneBuffer(t *testing.T) {
	buf := NewRuneBuffer(100)

	// 写入混合文本
	buf.WriteString("Hello, 世界! 👋")
	fmt.Println("初始内容:", buf.String()) // "Hello, 世界! 👋"
	fmt.Println("字符数:", buf.Len())     // 13 (H,e,l,l,o,,, ,世,界,!, ,👋)

	// 取回最后2个字符
	retrieved := buf.RetrieveLastChars(2)
	fmt.Println("取回内容:", retrieved)    // " 👋" (空格和表情符号)
	fmt.Println("剩余内容:", buf.String()) // "Hello, 世界!"
	fmt.Println("字符数:", buf.Len())     // 11

	// 取回最后3个字符
	retrieved = buf.RetrieveLastChars(3)
	fmt.Println("取回内容:", retrieved)    // "界!"
	fmt.Println("剩余内容:", buf.String()) // "Hello, 世"
	fmt.Println("字符数:", buf.Len())     // 8

	// 添加新内容
	buf.WriteString("欢迎!")
	fmt.Println("新内容:", buf.String()) // "Hello, 世欢迎!"
	fmt.Println("字符数:", buf.Len())    // 11

	// 高效查看最后2个字符
	lastTwo := buf.LastChars(2)
	fmt.Println("最后2字符:", lastTwo) // "迎!"
}
