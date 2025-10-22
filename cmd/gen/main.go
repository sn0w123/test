package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"cmcc-scaffold/pkg/generator"

	"github.com/spf13/cobra"
)

const (
	appName = "scaffold"
)

var rootCmd = &cobra.Command{
	Use:   appName,
	Short: "公司级 go-zero 脚手架插件",
}

var newCmd = &cobra.Command{
	Use:   "new [service]",
	Short: "生成符合公司规范的 go-zero 微服务",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// 0. goctl 是否存在
		if _, err := exec.LookPath("goctl"); err != nil {
			return errors.New("❌ 未找到 goctl，请先安装: go install github.com/zeromicro/go-zero/tools/goctl@latest")
		}
		// 1. 参数合法性
		svc := args[0]
		if strings.TrimSpace(svc) == "" {
			return errors.New("❌ 服务名不能为空")
		}
		// 2. 不能以 - 开头
		if strings.HasPrefix(svc, "-") {
			return errors.New("❌ 服务名不能以 '-' 开头")
		}
		// 3. 只能含字母、数字、下划线、中划线
		valid := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString
		if !valid(svc) {
			return errors.New("❌ 服务名只能包含字母、数字、下划线、中划线")
		}
		// 4. 重名检测
		if _, err := os.Stat(svc); err == nil {
			return fmt.Errorf("❌ 目录/文件 %s 已存在，请先清理或换个名字", svc)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := args[0]
		if err := generator.Run(svc); err != nil {
			return fmt.Errorf("❌ 生成失败: %w", err)
		}
		fmt.Printf("✅ 服务 %s 创建完成！\n", svc)
		fmt.Printf("📂 快去看看: cd %s\n", svc)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
