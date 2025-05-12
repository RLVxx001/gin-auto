/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"embed"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "项目初始化",
	Long: `创建初始化配置文件
使用方式：
  gin-auto init [-c] -o <输出目录>

参数说明：
  -c, --cover bool   如果存在原文件的话是否覆盖
  -o, --output string  指定输出目录路径（默认当前路径）

示例：
  gin-auto init -c -o ./auto-test`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
		cover, err := cmd.Flags().GetBool("cover")
		if err != nil {
			cmd.PrintErrf("参数错误：%v\n", err)
			return
		}
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			cmd.PrintErrf("参数错误：%v\n", err)
			return
		}
		err = createInitSettingsFile(filepath.Join(output, ".settings"), cover)
		if err != nil {
			cmd.PrintErrf("创建初始化配置文件失败：%v\n", err)
			return
		}
		err = createInitTemplateFile(filepath.Join(output, ".template"), cover)
		if err != nil {
			cmd.PrintErrf("创建初始化模板文件失败：%v\n", err)
			return
		}
		cmd.Printf("初始化配置文件成功")
	},
}

// createInitSettingsFile 创建初始化配置文件
func createInitSettingsFile(settingsUrl string, cover bool) error {
	_, err := os.ReadDir(settingsUrl)
	if err != nil {
		err = os.Mkdir(settingsUrl, os.ModePerm)
		if err != nil {
			return errors.New("创建目录失败:" + err.Error())
		}
	}
	for k, v := range defaultSettingDict {
		fileUrl := filepath.Join(settingsUrl, k)
		_, err = os.Stat(fileUrl)
		if err != nil || cover {
			err = os.WriteFile(fileUrl, v, os.ModePerm)
			if err != nil {
				return errors.New("创建文件失败:" + err.Error())
			}
		}
	}
	return nil
}

// createInitTemplateFile 创建初始化模板文件
func createInitTemplateFile(templateUrl string, cover bool) error {
	//如果覆盖操作
	_, err := os.ReadDir(templateUrl)
	if err != nil {
		err = os.Mkdir(templateUrl, os.ModePerm)
		if err != nil {
			return errors.New("创建目录失败:" + err.Error())
		}
	}
	for k, v := range defaultTemplateDict {
		fileUrl := filepath.Join(templateUrl, k)
		_, err = os.Stat(fileUrl)
		if err != nil || cover {
			err = os.WriteFile(fileUrl, v, os.ModePerm)
			if err != nil {
				return errors.New("创建文件失败:" + err.Error())
			}
		}
	}
	return nil
}
func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("cover", "c", false, "是否覆盖")
	initCmd.Flags().StringP("output", "o", "", "输出目录（默认当前路径）")

	initDefaultSettingAndTemplate()
}

var defaultSettingDict = map[string][]byte{}
var defaultTemplateDict = map[string][]byte{}

//go:embed .settings/*
var defaultSettingFS embed.FS

//go:embed .template/*
var defaultTemplateFS embed.FS

// 初始化默认配置和模板
func initDefaultSettingAndTemplate() {
	var err error
	defaultSettingDict, err = getDirToMap(&defaultSettingFS, ".settings")
	if err != nil {
		panic("获取默认配置失败：" + err.Error())
	}
	defaultTemplateDict, err = getDirToMap(&defaultTemplateFS, ".template")
	if err != nil {
		panic("获取默认模板失败：" + err.Error())
	}
}
func getDirToMap(fsys *embed.FS, name string) (map[string][]byte, error) {
	fileContents := map[string][]byte{}
	dir, err := fs.ReadDir(*fsys, name)
	if err != nil {
		return nil, err
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		data, _ := fsys.ReadFile(name + "/" + f.Name())
		fileContents[f.Name()] = data
	}
	return fileContents, nil
}
