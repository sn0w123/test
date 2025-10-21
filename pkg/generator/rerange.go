package generator

import (
	"os"
	"path/filepath"
)

// Rearrange 把 goctl 骨架摆成你要求的固定结构
func Rearrange(svc string) error {
	// 0. 目录提升：tmp-gen/<svc> -> <svc>  （解决双层）TODO 这里逻辑还是有点问题，后续应该支持目录自定义
	src := filepath.Join("tmp-gen", svc)
	if err := os.Rename(src, svc); err != nil {
		return err
	}

	// 1. 只清本次会覆盖”的子树，保留用户其他文件
	dirsToClean := []string{
		filepath.Join(svc, "cmd"),
		filepath.Join(svc, "config"),
		filepath.Join(svc, "api"),
		filepath.Join(svc, "internal"),
	}
	for _, d := range dirsToClean {
		_ = os.RemoveAll(d)
	}

	// 2. 把 goctl 原目录再搬回来（保证完整骨架）
	_ = os.MkdirAll(filepath.Join(svc, "cmd/app"), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "config"), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "api/openapi", svc), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "api/proto"), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "internal/handler", svc), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "internal/middleware", svc), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "internal/dao/db"), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "dev/docker"), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "model"), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "pkg"), 0755)
	_ = os.MkdirAll(filepath.Join(svc, "scripts"), 0755)

	// 3. 精确移动：源 -> 目标（空目录已建好，直接覆盖）
	moves := []struct{ src, dst string }{
		{filepath.Join("tmp-gen", svc+".go"), filepath.Join(svc, "cmd/app", svc+".go")},
		{filepath.Join("tmp-gen", "etc", svc+".yaml"), filepath.Join(svc, "config", svc+".yaml")},
		{filepath.Join("tmp-gen", "desc", svc+".api"), filepath.Join(svc, "api/openapi", svc, svc+".api")},
		{filepath.Join("tmp-gen", "rpc/pb", svc+".proto"), filepath.Join(svc, "api/proto", svc+".proto")},
		{filepath.Join("tmp-gen", "internal/handler"), filepath.Join(svc, "internal/handler", svc)},
		{filepath.Join("tmp-gen", "internal/logic"), filepath.Join(svc, "internal/logic")},
		{filepath.Join("tmp-gen", "internal/svc"), filepath.Join(svc, "internal/svc")},
		{filepath.Join("tmp-gen", "internal/types"), filepath.Join(svc, "internal/types")},
	}
	for _, m := range moves {
		if err := os.Rename(m.src, m.dst); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	// 4. 清掉 tmp-gen 骨架残留
	return os.RemoveAll("tmp-gen")
}
