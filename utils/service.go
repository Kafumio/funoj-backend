package utils

import (
	"funoj-backend/config"
	"funoj-backend/consts"
	"os"
	"path"
)

/**
 * services公用方法
 */

const (
	AcmCCodeFilePath    = "./resources/acmTemplate/c"
	AcmGoCodeFilePath   = "./resources/acmTemplate/go"
	AcmJavaCodeFilePath = "./resources/acmTemplate/java"
)

// GetExecutePath 给用户的此次运行生成一个临时目录
func GetExecutePath(config *config.AppConfig) string {
	uuid := GetUUID()
	executePath := path.Join(config.FilePathConfig.TempDir, uuid)
	return executePath
}

// GetTempDir 获取一个随机的临时文件夹
func GetTempDir(config *config.AppConfig) string {
	uuid := GetUUID()
	executePath := config.FilePathConfig.TempDir + "/" + uuid
	return executePath
}

func GetAcmCodeTemplate(language string) (string, error) {
	var filePath string
	switch language {
	case consts.ProgramC:
		filePath = AcmCCodeFilePath
	case consts.ProgramGo:
		filePath = AcmGoCodeFilePath
	case consts.ProgramJava:
		filePath = AcmJavaCodeFilePath
	}
	code, err := os.ReadFile(filePath)
	return string(code), err
}
