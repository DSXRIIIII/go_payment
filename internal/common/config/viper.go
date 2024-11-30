package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var once sync.Once

func ViperInit() {
	if err := NewViperConfig(); err != nil {
		panic(err)
	}
}

func NewViperConfig() (err error) {
	once.Do(func() {
		err = newViperConfig()
	})
	return
}

func newViperConfig() error {
	relPath, err := getRelativePathFromCaller()
	if err != nil {
		return err
	}
	viper.SetConfigName("global")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(relPath)
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}

// getRelativePathFromCaller函数用于获取从调用者所在目录到当前函数所在目录的相对路径
// 它返回相对路径字符串以及可能出现的错误
func getRelativePathFromCaller() (relPath string, err error) {
	// 获取调用者所在的当前工作目录
	callerPwd, err := os.Getwd()
	if err != nil {
		// 如果获取当前工作目录出错，直接返回，此时relPath为空字符串，err为获取目录时的错误
		return
	}

	// 获取当前函数（即getRelativePathFromCaller函数）的调用栈信息
	// runtime.Caller(0)返回四个值，这里我们只关心第二个值，即当前函数所在文件的绝对路径
	// 第一个值是调用者的程序计数器（PC）值，暂未使用所以用下划线忽略
	// 第三个值是当前函数所在的函数名，暂未使用所以用下划线忽略
	// 第四个值是当前函数所在文件的行号，暂未使用所以用下划线忽略
	_, here, _, _ := runtime.Caller(0)

	// 通过filepath.Rel函数计算从调用者所在的工作目录（callerPwd）到当前函数所在目录（filepath.Dir(here)）的相对路径
	// 计算结果赋值给relPath，同时可能返回的错误也赋值给err
	relPath, err = filepath.Rel(callerPwd, filepath.Dir(here))

	// 打印输出调用者所在目录、当前函数所在文件的绝对路径以及计算得到的相对路径信息
	fmt.Printf("caller from: %s, here: %s, relpath: %s", callerPwd, here, relPath)

	// 返回计算得到的相对路径和可能出现的错误
	return
}
