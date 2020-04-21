package path_utils

import (
	"github.com/chendeke/config/config/err_code"
	"os"
	"path/filepath"
)

var SupportedExts = []string{"json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "dotenv", "env"}

func searchInPath(in ,configName string) (filename string) {
	for _, ext := range SupportedExts {
		if b := IsExist(filepath.Join(in, configName+"."+ext)); b {
			return filepath.Join(in, configName+"."+ext)
		}
	}
	return ""
}

func FindConfigFile(configName string, configPaths ...string) (string, error) {
	for _, cp := range configPaths {
		file := searchInPath(cp, configName)
		if file != "" {
			return file, nil
		}
	}
	return "", err_code.NotFoundError
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}