package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

const root = "." // 生成在当前目录

func Rearrange(svc string) error {
	// 0. 如果目标已存在，先清空  && TODO 这块逻辑有点问题，想想如果可以做成以后改了什么接口 还可以自动生成，不删目录
	if _, err := os.Stat(svc); err == nil {
		_ = os.RemoveAll(svc)
	}

	// 1. 创建公司目录
	if err := os.Rename("tmp-gen", svc); err != nil {
		return fmt.Errorf("move tmp-gen->%s: %w", svc, err)
	}
	dirs := []string{
		"cmd/app", "cmd/cli",
		"api/openapi/" + svc, "api/proto",
		"config", "dev/docker", "dev/k8s",
		"internal/dao/db", "internal/handler/" + svc,
		"internal/middleware/" + svc, "internal/server",
		"model", "pkg", "scripts",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(root, d), 0755); err != nil {
			return err
		}
	}

	// 2. 移动文件
	moves := []struct{ src, dst string }{
		{filepath.Join(svc, svc+".go"), filepath.Join("cmd/app", svc+".go")},
		{filepath.Join(svc, "etc", svc+".yaml"), filepath.Join("config", svc+".yaml")},
		{filepath.Join(svc, "desc", svc+".api"), filepath.Join("api/openapi", svc, svc+".api")},
		{filepath.Join(svc, "rpc", "pb", svc+".proto"), filepath.Join("api/proto", svc+".proto")},
		{filepath.Join(svc, "internal/handler"), filepath.Join("internal/handler", svc)},
		{filepath.Join(svc, "internal/logic"), filepath.Join("internal/logic")},
		{filepath.Join(svc, "internal/svc"), filepath.Join("internal/svc")},
		{filepath.Join(svc, "internal/types"), filepath.Join("internal/types")},
	}
	for _, m := range moves {
		if err := os.Rename(m.src, m.dst); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("move %s->%s: %w", m.src, m.dst, err)
		}
	}

	// 3. 清理空目录
	_ = os.RemoveAll(svc)
	return nil
}
