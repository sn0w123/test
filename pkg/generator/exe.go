package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
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
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("goctl api new: %w, output: %s", err, string(output))
	}

	// 2. 生成 RPC 骨架
	fmt.Println("生成 RPC 骨架...")
	cmd = exec.Command("goctl", "rpc", "new", svc, "--style=goZero")
	cmd.Dir = "tmp-gen"
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("goctl rpc new: %w, output: %s", err, string(output))
	}

	// 3. 确保 desc 目录存在
	descDir := filepath.Join("tmp-gen", "desc")
	if err := os.MkdirAll(descDir, 0755); err != nil {
		return fmt.Errorf("create desc directory: %w", err)
	}

	// 4. 用模板覆盖 API 文件
	apiFile := filepath.Join("tmp-gen", "desc", svc+".api")
	if err := writeTemplateFile(apiFile, "tpl/api.api.tmpl", map[string]string{"svc": svc}); err != nil {
		return err
	}

	// 5. 查找并覆盖 RPC 文件
	rpcFile, err := findRPCFile(svc)
	if err != nil {
		return err
	}
	fmt.Printf("找到 RPC 文件: %s\n", rpcFile)

	if err := writeTemplateFile(rpcFile, "tpl/rpc.proto.tmpl", map[string]string{"svc": svc}); err != nil {
		return err
	}

	// 6. 重新生成 API 代码
	fmt.Println("重新生成 API 代码...")
	apiRelativePath := filepath.Join("desc", svc+".api")
	cmd = exec.Command("goctl", "api", "go", "-api", apiRelativePath, "-dir", ".")
	cmd.Dir = "tmp-gen"
	fmt.Printf("Running: goctl api go -api %s -dir .\n", apiRelativePath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("goctl api go: %w, output: %s", err, string(output))
	}

	// 7. 重新生成 RPC 代码
	fmt.Println("重新生成 RPC 代码...")

	// 使用正确的 RPC 生成命令
	rpcDir := filepath.Join("tmp-gen", svc)
	cmd = exec.Command("goctl", "rpc", "protoc", svc+".proto", "--go_out=.", "--go-grpc_out=.", "--zrpc_out=.", "-m")
	cmd.Dir = rpcDir
	fmt.Printf("Running: goctl rpc protoc %s.proto --go_out=. --go-grpc_out=. --zrpc_out=. -m\n", svc)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("goctl rpc protoc: %w, output: %s", err, string(output))
	}

	return nil
}

// findRPCFile 查找 RPC 文件的实际位置
func findRPCFile(svc string) (string, error) {
	possiblePaths := []string{
		filepath.Join("tmp-gen", svc, svc+".proto"),
		filepath.Join("tmp-gen", "rpc", svc+".proto"),
		filepath.Join("tmp-gen", svc, "pb", svc+".proto"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// 如果没找到，打印目录结构帮助调试
	fmt.Println("目录结构:")
	filepath.Walk("tmp-gen", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasSuffix(path, ".proto") || strings.HasSuffix(path, ".api")) {
			fmt.Printf("  %s\n", path)
		}
		return nil
	})

	return "", fmt.Errorf("未找到 RPC 文件")
}

// 辅助函数：用模板生成文件
func writeTemplateFile(dst, tmplPath string, data map[string]string) error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// 解析模板并写出
	tmpl, err := template.ParseFS(tplFS, tmplPath)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}
	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}
