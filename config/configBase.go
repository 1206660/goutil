package config

import (
	"fmt"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type EConfigContextType int

const (
	_ EConfigContextType = iota
	EConfigContextType_Int
	EConfigContextType_String
)

type ConfigCell interface {
	GetAsInt() int
	GetAsFloat() float64
	GetAsString() string
	ToString() string
	ReadByString(string, string) bool
}
type ConfigCellString struct {
	context string
}

type ConfigCellInt struct {
	context int
}

type ConfigCellFloat struct {
	context float64
}

func (self *ConfigCellInt) GetAsInt() int {
	return self.context
}
func (self *ConfigCellInt) GetAsFloat() float64 {
	return float64(self.context)
}
func (self *ConfigCellInt) GetAsString() string {
	//panic("can not get int field as string")
	return self.ToString()
}
func (self *ConfigCellInt) ToString() string {
	return strconv.Itoa(self.context)
}
func (self *ConfigCellInt) ReadByString(val string, fieldName string) bool {
	num, err := strconv.Atoi(val)
	if err != nil {
		return false
	}
	self.context = num
	return true
}

func (self *ConfigCellString) GetAsFloat() float64 {
	panic("can not get int field as float32")
	f, err := strconv.ParseFloat(self.ToString(), 64)
	if err != nil {
		return 0
	}
	return f
}
func (self *ConfigCellString) GetAsInt() int {
	panic("can not get string field as int")
	return 0
}
func (self *ConfigCellString) GetAsString() string {
	return self.context
}
func (self *ConfigCellString) ToString() string {
	return self.context
}
func (self *ConfigCellString) ReadByString(val string, fieldName string) bool {
	self.context = val
	return true
}

func (self *ConfigCellFloat) GetAsFloat() float64 {
	return self.context
}
func (self *ConfigCellFloat) GetAsInt() int {
	panic("can not get string field as int")
	return 0
}
func (self *ConfigCellFloat) GetAsString() string {
	panic("can not get int field as string")
	return self.ToString()
}
func (self *ConfigCellFloat) ToString() string {
	return ""
}
func (self *ConfigCellFloat) ReadByString(val string, fieldName string) bool {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return false
	}
	self.context = f
	return true
}

type Callback func()

type ConfigMap struct {
	infos           map[int](map[string]ConfigCell)
	MgrInitCallback Callback // 重载时候的回调函数
	filePath        string
}

func (self *ConfigMap) GetAllInfo() map[int](map[string]ConfigCell) {
	return self.infos
}

func (self *ConfigMap) GetInfo(id int) map[string]ConfigCell {
	return self.infos[id]
}
func (self *ConfigMap) LoadConfigByFilePath(filePath string) bool {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("config string read error" + err.Error())
		return false
	}
	self.filePath = filePath
	return self.LoadConfig(string(file))
}

func (self *ConfigMap) Reload() {
	self.LoadConfigByFilePath(self.filePath)
	if self.MgrInitCallback != nil {
		self.MgrInitCallback()
	}
	log.Debug("ConfigFile Reload success %v", self.filePath)
}

func (self *ConfigMap) LoadConfig(context string) bool {
	infos := make(map[int]map[string](ConfigCell))
	lines := strings.Split(context, "\r\n")
	count := 0
	var fields []string
	for lineCount, line := range lines {
		if line == "" { // 空行
			continue
		}
		count++
		vals := strings.Split(line, "\t")
		if count == 1 {
			// 字段行
			fields = vals
			continue
		}
		info := make(map[string](ConfigCell))
		id, err := strconv.Atoi(vals[0])
		if err != nil {
			fmt.Println("error: id[" + vals[0] + "] is not a valid number")
			return false
		}
		for index, fieldName := range fields {
			tp := fieldName[0:1]
			var cell ConfigCell
			if tp == "i" {
				cell = new(ConfigCellInt)
			} else if tp == "f" {
				cell = new(ConfigCellFloat)
			} else {
				cell = new(ConfigCellString)
			}
			if !cell.ReadByString(vals[index], fieldName) {
				panic("error: config string read error. line[" + strconv.Itoa(lineCount) + "] is not a valid number fieldName =" + fieldName)
				return false
			}
			info[fieldName] = cell
		}
		infos[id] = info
	}
	self.infos = infos
	return true
}

var ostype = runtime.GOOS

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
func getParentDirectory(dirctory string) string {
	comma := "\\"
	if ostype != "windows" {
		comma = "/"
	}
	return substr(dirctory, 0, strings.LastIndex(dirctory, comma))
}
func GetCurrentConfigDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	comma := "\\"
	if ostype != "windows" {
		comma = "/"
	}
	parentDir := substr(dir, strings.LastIndex(dir, comma), 4)
	fmt.Println(dir + "   " + parentDir)
	if parentDir == comma+"bin" {
		return getParentDirectory(dir) + "/config"
	}
	return getParentDirectory(getParentDirectory(dir)) + "/config" //strings.Replace(dir, "\\", "/", -1)
}
