package generator

import (
	"fmt"
	"os"
	"os/exec"
)

// RunGoctl 先让原生 goctl 生成最原始骨架
func RunGoctl(svc string) error {
	//TODO 这里目录还是绝对路径，硬编码，需要想想是否考虑匿名临时目录 rpc其他服务的新增都会报错
	// 0. 创建并进入临时目录
	_ = os.RemoveAll("tmp-gen")
	_ = os.MkdirAll("tmp-gen", 0755)

	// 1. api 骨架
	cmd := exec.Command("goctl", "api", "new", svc, "--style=goZero")
	cmd.Dir = "tmp-gen"
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("goctl api new: %w", err)
	}
	// 2. rpc 骨架
	//cmd = exec.Command("goctl", "rpc", "new", svc, "--style=goZero")
	//cmd.Dir = filepath.Join("tmp-gen", svc)
	//if err := cmd.Run(); err != nil {
	//	return fmt.Errorf("goctl rpc new: %w", err)
	//}
	return nil
}
