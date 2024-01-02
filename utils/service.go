package utils

import (
	"funoj-backend/consts"
	"os"
)

/**
 * services公用方法
 */

const (
	AcmCCodeFilePath    = "./resources/acmTemplate/c"
	AcmGoCodeFilePath   = "./resources/acmTemplate/go"
	AcmJavaCodeFilePath = "./resources/acmTemplate/java"
)

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
