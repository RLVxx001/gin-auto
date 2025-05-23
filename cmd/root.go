package cmd

import (
	"os"

	"github.com/RLVxx001/gin-auto/auto"

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
  -s, --settings string 指定配置文件路径（可选，默认为当前目录下的.settings）
  -t, --templates string 指定模板目录路径（可选，默认为当前目录下的.templates）

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

		templateDir, err := cmd.Flags().GetString("templates")
		if err != nil {
			cmd.PrintErrf("获取模板目录失败: %v\n", err)
			return
		}

		settingFile, err := cmd.Flags().GetString("settings")
		if err != nil {
			cmd.PrintErrf("获取配置文件路径失败: %v\n", err)
			return
		}

		// 必填参数没有填，弹帮助
		if inputDir == "" || outputDir == "" {
			cmd.Help()
			return
		}
		// 如果未指定配置文件路径，使用默认值
		if settingFile == "" {
			settingFile = ".settings"
		}
		// 如果未指定模板目录，使用默认值
		if templateDir == "" {
			templateDir = ".templates"
		}
		{
			//检查是否存在配置文件以及模板文件
			_, err := os.Stat(settingFile)
			if err != nil {
				err = createInitSettingsFile(settingFile, false)
				if err != nil {
					cmd.PrintErrf("创建配置文件失败: %v\n", err)
					return
				}
			}
			//  检查模板目录
			_, err = os.Stat(templateDir)
			if err != nil {
				err = createInitTemplateFile(templateDir, false)
				if err != nil {
					cmd.PrintErrf("创建模板目录失败: %v\n", err)
					return
				}
			}
		}
		// 1.初始化auto包
		auto.Initialize(settingFile, templateDir)
		// 2.setting_front（设置开头）
		auto.RunSettingFront()
		// 3.api_setting（api中用户端设置）
		auto.GetApi(inputDir)
		// 4.setting_back（设置结尾）
		auto.RunSettingEnd()
		// 5.api_list（api文件中各个接口信息）
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
	rootCmd.Flags().StringP("templates", "t", "", "模板目录（可选，默认为.templates）")
	rootCmd.Flags().StringP("settings", "s", "", "配置文件路径（可选，默认为当前目录下的.settings）")
	rootCmd.MarkFlagRequired("input")
	rootCmd.MarkFlagRequired("output")
}
