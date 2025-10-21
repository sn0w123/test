package generator

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed tpl/*
var tplFS embed.FS

func Overlay(svc string) error {
	return fs.WalkDir(tplFS, "tpl", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		// 计算目标路径
		rel, _ := filepath.Rel("tpl", path)
		dst := filepath.Join(".", rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		// 解析模板并写出
		t, err := template.ParseFS(tplFS, path)
		if err != nil {
			return err
		}
		f, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer f.Close()
		return t.Execute(f, map[string]string{"svc": svc})
	})
}
