package helper

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
)

func IsEmpty(params interface{}) bool {
	//初始化变量
	var (
		flag         bool = true
		defaultValue reflect.Value
	)
	r := reflect.ValueOf(params)
	defaultValue = reflect.Zero(r.Type())
	//由于params 接口类型 所以default_value也要获取对应接口类型的值 如果获取不为接口类型 一直为返回false
	if !reflect.DeepEqual(r.Interface(), defaultValue.Interface()) {
		flag = false
	}
	return flag
}

func GetAppDir() string {
	appDir, err := os.Getwd()
	if err != nil {
		file, _ := exec.LookPath(os.Args[0])
		applicationPath, _ := filepath.Abs(file)
		appDir, _ = filepath.Split(applicationPath)
	}
	return appDir
}