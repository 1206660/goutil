package configPath

import (
	"github.com/name5566/leaf/log"
	"os"
	"os/exec"
	"path/filepath"
)

func GetExecDir() string {
	execPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Fatal("err: %v", err)
	}
	//    Is Symlink
	fi, err := os.Lstat(execPath)
	if err != nil {
		log.Fatal("err: %v", err)
	}
	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		execPath, err = os.Readlink(execPath)
		if err != nil {
			log.Fatal("err: %v", err)
		}
	}
	execDir := filepath.Dir(execPath)
	if execDir == "." {
		execDir, err = os.Getwd()
		if err != nil {
			log.Fatal("err: %v", err)
		}
	}
	return execDir
}

func GetConfDir() string {
	execDir := GetExecDir()
	return execDir + "/conf/"
}
