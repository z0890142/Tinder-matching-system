package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"tinderMatchingSystem/pkg/utils/logger"

	"github.com/argoproj/pkg/file"

	"github.com/spf13/viper"
)

// Op op
type Op struct {
	FileName      string
	VerifyElement []string
}

// WithFileName WithFileName
func WithFileName(fileName string) Option {
	return func(op *Op) {
		op.FileName = fileName
	}
}

// WithVerify WithVerify
func WithVerify(verifyKey, verifyValue string) Option {
	return func(op *Op) {
		op.VerifyElement = append(op.VerifyElement, verifyKey, verifyValue)
	}
}

// Option option
type Option func(op *Op)

func LoadConf(paths []string, globalConfig interface{}, opts ...Option) {
	var op Op
	for _, option := range opts {
		option(&op)
	}
	v := viper.GetViper()
	loadLocal(v, paths, op.FileName, globalConfig)
	loadEnv(v, globalConfig)
}

// loadLocal load local file
func loadLocal(vp *viper.Viper, dirPaths []string, filename string, globalConfig interface{}) {

	defaultLoggingKeys := map[string]interface{}{
		"log_key": "local_config",
		"service": viper.GetString("SERVICE.NAME"),
		"env":     viper.GetString("ENV"),
	}
	if filename == "" {
		filename = "config"
	}
	for _, dirPath := range dirPaths {
		vp.AddConfigPath(dirPath)
	}
	vp.SetConfigName(filename)
	err := vp.ReadInConfig()
	if err != nil {
		logger.GetLoggerWithKeys(defaultLoggingKeys).Panic(fmt.Sprintf("ReadInConfig err: %s ", err))
	}

	settingsLocalFilename := filename + "_local"

	for _, dirPath := range dirPaths {
		if file.Exists(path.Join(dirPath, settingsLocalFilename+".yaml")) || file.Exists(path.Join(dirPath, settingsLocalFilename+".json")) {
			vp.SetConfigName(settingsLocalFilename)
			err = vp.MergeInConfig()
			if err != nil {
				logger.GetLoggerWithKeys(defaultLoggingKeys).Panic(fmt.Sprintf("fatal error MergeInConfig: err: %s ", err))
			}
		}
	}

	err = vp.Unmarshal(globalConfig)
	if err != nil {
		logger.GetLoggerWithKeys(defaultLoggingKeys).Panic(fmt.Sprintf("fatal error config unmarshal:  err: %s ", err))
	}
}

func loadEnv(vp *viper.Viper, globalConfig interface{}) {
	defaultLoggingKeys := map[string]interface{}{
		"log_key": "local_config",
		"service": viper.GetString("SERVICE.NAME"),
		"env":     viper.GetString("ENV"),
	}
	vp.AutomaticEnv()

	for _, envstr := range os.Environ() {
		parts := strings.SplitN(envstr, "=", 2)
		key := parts[0]
		value := ""

		if len(parts) == 2 {
			value = parts[1]
		}
		vp.Set(key, value)
	}

	err := vp.Unmarshal(globalConfig)
	if err != nil {
		logger.GetLoggerWithKeys(defaultLoggingKeys).Panic(fmt.Sprintf("fatal error apollo unmarshall err: %s ", err))
	}
}
