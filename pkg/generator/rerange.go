package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Rearrange 把 goctl 骨架摆成你要求的固定结构
//
//	func Rearrange(svc string) error {
//		// 0. 目录提升：tmp-gen/<svc> -> <svc>  （解决双层）TODO 这里逻辑还是有点问题，后续应该支持目录自定义
//		//src := filepath.Join("tmp-gen", svc)
//		if err := os.Rename(src, svc); err != nil {
//			return err
//		}
//
//		// 1. 只清本次会覆盖”的子树，保留用户其他文件
//		dirsToClean := []string{
//			filepath.Join(svc, "cmd"),
//			filepath.Join(svc, "config"),
//			filepath.Join(svc, "api"),
//			filepath.Join(svc, "internal"),
//		}
//		for _, d := range dirsToClean {
//			_ = os.RemoveAll(d)
//		}
//
//		// 2. 把 goctl 原目录再搬回来（保证完整骨架）
//		_ = os.MkdirAll(filepath.Join(svc, "cmd/app"), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "config"), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "api/openapi"), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "api/proto"), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "internal/handler", svc), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "internal/logic", svc), 0755) //新增logic
//		_ = os.MkdirAll(filepath.Join(svc, "internal/middleware", svc), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "internal/dao/db"), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "dev/docker"), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "model"), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "pkg"), 0755)
//		_ = os.MkdirAll(filepath.Join(svc, "scripts"), 0755)
//
//		// 3. 精确移动：源 -> 目标（空目录已建好，直接覆盖）
//		moves := []struct{ src, dst string }{
//			{filepath.Join("tmp-gen", svc+".go"), filepath.Join(svc, "cmd/app", svc+".go")},
//			{filepath.Join("tmp-gen", "etc", svc+".yaml"), filepath.Join(svc, "config", svc+".yaml")},
//			{filepath.Join("tmp-gen", "desc", svc+".api"), filepath.Join(svc, "api/openapi", svc, svc+".api")},
//			{filepath.Join("tmp-gen", "rpc/", svc+".proto"), filepath.Join(svc, "api/proto", svc+".proto")},
//			{filepath.Join("tmp-gen", "internal/handler"), filepath.Join(svc, "internal/handler", svc)},
//			{filepath.Join("tmp-gen", "internal/logic"), filepath.Join(svc, "internal/logic")},
//			{filepath.Join("tmp-gen", "internal/svc"), filepath.Join(svc, "internal/svc")},
//			{filepath.Join("tmp-gen", "internal/types"), filepath.Join(svc, "internal/types")},
//		}
//		for _, m := range moves {
//			if _, err := os.Stat(m.src); err == nil {
//				if err := os.Rename(m.src, m.dst); err != nil {
//					// 如果重命名失败，尝试复制
//					if copyErr := copyFile(m.src, m.dst); copyErr != nil {
//						fmt.Printf("警告: 无法移动 %s 到 %s: %v\n", m.src, m.dst, err)
//					}
//				}
//			}
//		}
//
//
//
//		// 4. 清掉 tmp-gen 骨架残留
//		return os.RemoveAll("tmp-gen")
//	}
//
// // 添加文件复制辅助函数
//
//	func copyFile(src, dst string) error {
//		input, err := os.ReadFile(src)
//		if err != nil {
//			return err
//		}
//		return os.WriteFile(dst, input, 0644)
//	}
func Rearrange(svc string) error {
	targetDir := svc

	// 1. 创建标准目录结构
	dirs := []string{
		filepath.Join(targetDir, "cmd", "app"),
		filepath.Join(targetDir, "cmd", "rpc"),
		filepath.Join(targetDir, "config"),
		filepath.Join(targetDir, "api", "openapi", svc),
		filepath.Join(targetDir, "api", "proto"),
		filepath.Join(targetDir, "internal", "handler", svc),
		filepath.Join(targetDir, "internal", "logic", svc),
		filepath.Join(targetDir, "internal", "middleware", svc),
		filepath.Join(targetDir, "internal", "dao", "db"),
		filepath.Join(targetDir, "internal", "svc"),
		filepath.Join(targetDir, "dev", "docker"),
		filepath.Join(targetDir, "model"),
		filepath.Join(targetDir, "pkg"),
		filepath.Join(targetDir, "scripts"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	// 2. 复制 API 相关文件
	apiCopyOps := []struct {
		src  string
		dest string
	}{
		// API 主文件
		{filepath.Join("tmp-gen", svc+".go"), filepath.Join(targetDir, "cmd", "app", svc+".go")},
		{filepath.Join("tmp-gen", "etc", svc+".yaml"), filepath.Join(targetDir, "config", svc+".yaml")},
		{filepath.Join("tmp-gen", "desc", svc+".api"), filepath.Join(targetDir, "api", "openapi", svc, svc+".api")},

		// API 内部文件
		{filepath.Join("tmp-gen", "internal", "handler"), filepath.Join(targetDir, "internal", "handler", svc)},
		{filepath.Join("tmp-gen", "internal", "logic"), filepath.Join(targetDir, "internal", "logic", svc)},
		{filepath.Join("tmp-gen", "internal", "svc"), filepath.Join(targetDir, "internal", "svc")},
		{filepath.Join("tmp-gen", "internal", "types"), filepath.Join(targetDir, "internal", "types")},
	}

	// 3. 复制 RPC 相关文件
	rpcCopyOps := []struct {
		src  string
		dest string
	}{
		// RPC 主文件
		{filepath.Join("tmp-gen", svc, svc+".go"), filepath.Join(targetDir, "cmd", "rpc", svc+".go")},
		{filepath.Join("tmp-gen", svc, "etc", svc+".yaml"), filepath.Join(targetDir, "config", svc+"-rpc.yaml")},
		{filepath.Join("tmp-gen", svc, svc+".proto"), filepath.Join(targetDir, "api", "proto", svc+".proto")},

		// RPC 生成的文件
		{filepath.Join("tmp-gen", svc, svc+".pb.go"), filepath.Join(targetDir, "api", "proto", svc+".pb.go")},
		{filepath.Join("tmp-gen", svc, svc+"_grpc.pb.go"), filepath.Join(targetDir, "api", "proto", svc+"_grpc.pb.go")},

		// RPC 内部文件
		{filepath.Join("tmp-gen", svc, "internal", "config"), filepath.Join(targetDir, "internal", "config")},
		{filepath.Join("tmp-gen", svc, "internal", "server"), filepath.Join(targetDir, "internal", "server")},
		{filepath.Join("tmp-gen", svc, "internal", "svc"), filepath.Join(targetDir, "internal", "svc", "rpc")},
		{filepath.Join("tmp-gen", svc, "internal", "logic"), filepath.Join(targetDir, "internal", "logic", "rpc")},
	}

	fmt.Println("复制 API 文件...")
	for _, op := range apiCopyOps {
		if err := copyIfExists(op.src, op.dest); err != nil {
			fmt.Printf("警告: 无法复制 %s 到 %s: %v\n", op.src, op.dest, err)
		}
	}

	fmt.Println("复制 RPC 文件...")
	for _, op := range rpcCopyOps {
		if err := copyIfExists(op.src, op.dest); err != nil {
			fmt.Printf("警告: 无法复制 %s 到 %s: %v\n", op.src, op.dest, err)
		}
	}
	// 4.删除模板
	fmt.Println("清理不需要的文件...")
	unwantedFiles := []string{
		filepath.Join(targetDir, "api.api.tmpl"),
		filepath.Join(targetDir, "rpc.proto.tmpl"),
	}

	for _, file := range unwantedFiles {
		if err := os.RemoveAll(file); err != nil && !os.IsNotExist(err) {
			fmt.Printf("无法删除 %s: %v\n", file, err)
		}
	}

	// 5. 清理临时目录
	if err := os.RemoveAll("tmp-gen"); err != nil {
		fmt.Printf("无法删除临时目录 tmp-gen: %v\n", err)
	}

	return nil
}

// copyIfExists 如果源文件存在则复制
func copyIfExists(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 源文件不存在，跳过
		}
		return err
	}

	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

// 复制目录
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return copyFile(path, targetPath)
	})
}

// 复制文件
func copyFile(src, dst string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}
