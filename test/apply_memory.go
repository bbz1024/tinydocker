package main

import (
	"fmt"
	"time"
)

func main() {
	// 每秒申请 10MB 的内存
	memorySize := 10 * 1024 * 1024 // 10MB

	for {
		// 申请内存
		data := make([]byte, memorySize)
		_ = data // 避免编译器优化掉这块内存

		// 打印当前时间
		fmt.Println("申请了 10MB 内存")

		// 等待 1 秒
		time.Sleep(1 * time.Second)
	}
}
