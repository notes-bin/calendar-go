package uuid

import (
	"crypto/rand"
	"fmt"
	"io"
)

// 生成UUID v4（随机）
func Generate() (string, error) {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		return "", err
	}

	// 设置版本（v4）和变体
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // 版本4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // 变体RFC4122

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16]), nil
}
