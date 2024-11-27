package gos

import (
	"github.com/name5566/leaf/log"
	"os"
)

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateDir(dir string) {
	exist, err := PathExists(dir)
	if err != nil {
		log.Error("Get PathExists dir error![%v]", err)
		return
	}
	if !exist {
		// 创建文件夹
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			log.Fatal("mkdir %v failed! %v", dir, err.Error())
		} else {
			log.Debug("mkdir %v success!", dir)
		}
	}
}
