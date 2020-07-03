package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// 初始化配置文件
func Initialize(fileName string) error {
	splits := strings.Split(filepath.Base(fileName), ".")
	viper.SetConfigName(filepath.Base(splits[0]))
	viper.AddConfigPath(filepath.Dir(fileName))
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// 检查必须设置的参数，是否有设置
	err = checkMustSetArgs()
	if err != nil {
		return err
	}

	//viper.Debug()

	return nil
}

// 检查必须设置的参数，是否有设置
func checkMustSetArgs() error {
	return nil
}

func checkKey(key string) {
	if !viper.IsSet(key) {
		fmt.Printf("配置文件参数%q未设置，现在退出程序\n", key)
		os.Exit(1)
	}
}

func MustGetString(key string) string {
	checkKey(key)
	return viper.GetString(key)
}

func MustGetInt(key string) int {
	checkKey(key)
	return viper.GetInt(key)
}

func MustGetBool(key string) bool {
	checkKey(key)
	return viper.GetBool(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func GetUint64(key string) uint64 {
	return viper.GetUint64(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}
