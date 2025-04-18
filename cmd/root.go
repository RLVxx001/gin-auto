package cmd

import (
	"os"

	"github.com/n8sPxD/gin-auto/auto"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gin-auto",
	Short: "自动化gin框架代码生成工具",
	Long: `自动化gin框架代码生成工具

使用方式：
  gin-auto -i <输入目录> -o <输出目录> [-t <模板目录>] [-s <配置文件路径>]

参数说明：
  -i, --input string   指定API文件目录路径（必填）
  -o, --output string  指定输出目录路径（必填）
  -s, --setting string 指定配置文件路径（可选，默认为当前目录下的.setting）
  -t, --template string 指定模板目录路径（可选，默认为当前目录下的.templates）

示例：
  gin-auto -i ./input -o ./output
  gin-auto -i ./input -o ./output -t ./my-templates`,

	Run: func(cmd *cobra.Command, args []string) {
		inputDir, err := cmd.Flags().GetString("input")
		if err != nil {
			cmd.PrintErrf("获取输入目录失败: %v\n", err)
			return
		}

		outputDir, err := cmd.Flags().GetString("output")
		if err != nil {
			cmd.PrintErrf("获取输出目录失败: %v\n", err)
			return
		}

		templateDir, err := cmd.Flags().GetString("template")
		if err != nil {
			cmd.PrintErrf("获取模板目录失败: %v\n", err)
			return
		}

		settingFile, err := cmd.Flags().GetString("setting")
		if err != nil {
			cmd.PrintErrf("获取配置文件路径失败: %v\n", err)
			return
		}

		// 必填参数没有填，弹帮助
		if inputDir == "" || outputDir == "" {
			cmd.Help()
			return
		}

		// 如果未指定模板目录，使用默认值
		if templateDir == "" {
			templateDir = ".templates"
		}

		// 初始化auto包
		auto.Initialize(settingFile)

		// 初始化目标输出目录
		auto.InitWorkDir(outputDir)

		// 初始化模板目录
		auto.InitTemplateDir(templateDir)

		auto.GetApi()
		auto.A.InsertContext()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("input", "i", "", "输入目录（必填）")
	rootCmd.Flags().StringP("output", "o", "", "输出目录（必填）")
	rootCmd.Flags().StringP("template", "t", "", "模板目录（可选，默认为.templates）")
	rootCmd.Flags().StringP("setting", "s", "", "配置文件路径（可选，默认为当前目录下的.setting）")
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
