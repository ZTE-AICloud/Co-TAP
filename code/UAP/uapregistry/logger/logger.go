package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/core/logs"
	"github.com/jinzhu/configor"

	"uapregistry/types"
)

var (
	log         *logs.BeeLogger
	LoggerLevel = map[string]int{"Emergency": 0, "Alert": 1, "Critical": 2, "Error": 3,
		"Warn": 4, "Notice": 5, "Info": 6, "Debug": 7}
	logger MyLoggerInterface
)

const (
	logFileName    string = "log.yml"
	ConfigFilePath string = "config"
)

type Logger MyLoggerInterface

type ConsoleCfg struct {
	Level int `json:"level"`
}

type FileCfg struct {
	Filename string `json:"filename"`
	Level    int    `json:"level"`
	MaxLines int    `json:"maxlines"`
	MaxSize  int    `json:"maxsize"`
	Daily    bool   `json:"daily"`
	MaxDays  int64  `json:"maxdays"`
	Rotate   bool   `json:"rotate"`
	Perm     string `json:"perm"`
}

func init() {
	initBeegoLogger()
}

var initBeegoLogger = func() {
	log = logs.NewLogger()
	log.EnableFuncCallDepth(true)
	log.SetLogFuncCallDepth(2)

	//use custom configor
	if initCustomLogger() {
		logger = newLogger()
		logger.Warn("mylogger is init success")
		return
	}

	//if custom failed,use default
	initDedaultLogger()

	logger = newLogger()
	logger.Warn("mylogger is init success")
}

var GetLogger = func() MyLoggerInterface {
	return logger
}

// NewLogger is the creator for logger object
func newLogger() MyLoggerInterface {
	log.SetLogFuncCallDepth(3)

	return &MyLogger{
		hostName:      getHostName(),
		module:        moduleName,
		transactionID: "null",
		instanceID:    getInstanceID(),
		logger:        log,
	}
}

var initCustomLogger = func() bool {
	loggerCfg := readLoggerCfg()

	if loggerCfg == nil {
		return false
	}

	//console
	if !setConsoleLogger(loggerCfg) {
		fmt.Printf("set console logger failed.")
		return false
	}

	//file
	if !setFileLogger(loggerCfg) {
		if err := log.DelLogger(logs.AdapterConsole); err != nil {
			fmt.Printf("delete file logger failed:%s", err.Error())
		}
		fmt.Printf("set file logger failed.")
		return false
	}

	log.Warn("---------custom logger conf file-------")
	printLoggerCfg(loggerCfg)

	return true
}

var readLoggerCfg = func() *types.Logger {
	loggerCfg := &types.Logger{}

	confdir := GetCfgFilePath()
	pthSep := string(os.PathSeparator)
	filPath := confdir + pthSep + logFileName

	err := configor.Load(loggerCfg, filPath)

	if err != nil {
		fmt.Printf("read config file failed:%s", err.Error())
		return nil
	}

	return loggerCfg
}

var initDedaultLogger = func() bool {
	loggerCfg := &types.Logger{}
	loggerCfg.Console.Level = "Warn"

	loggerCfg.File.Filename = "uapregistry.log"
	loggerCfg.File.Level = "Info"
	loggerCfg.File.Maxlines = 100000
	loggerCfg.File.Maxsize = 30
	loggerCfg.File.Daily = true
	loggerCfg.File.Maxdays = 10
	loggerCfg.File.Rotate = true
	loggerCfg.File.Perm = "0640"

	//////////////////////////////////////
	setConsoleLogger(loggerCfg)
	setFileLogger(loggerCfg)

	log.Warn("---------default logger conf file-------")
	printLoggerCfg(loggerCfg)

	return true
}

// set console
// Level    int  `json:"level"`
// Colorful bool `json:"color"`
var setConsoleLogger = func(lc *types.Logger) bool {
	//console
	consolecfg := &ConsoleCfg{}
	consolecfg.Level = LoggerLevel[lc.Console.Level]

	byteconfig, err := json.Marshal(consolecfg)
	if err != nil {
		fmt.Printf("set console logger,change to json failed:%s", err.Error())
		return false
	}

	seterr := log.SetLogger(logs.AdapterConsole, string(byteconfig))
	if seterr != nil {
		fmt.Printf("set console logger failed:%s", seterr.Error())
		return false
	}

	return true
}

// set file
// config need to be correct JSON as string: {"interval":360}.
// It writes messages by lines limit, file size limit, or time frequency.
// (1-999) filename.date.num.log
// Filename   string `json:"filename"`
// MaxLines         int `json:"maxlines"` 1000000
// MaxSize        int `json:"maxsize"` 1 << 28 256M  length in bytes
// Daily         bool  `json:"daily"`
// MaxDays       int64 `json:"maxdays"` 7
// Rotate bool `json:"rotate"`
// Level int `json:"level"`
// Perm string `json:"perm"`
var setFileLogger = func(lc *types.Logger) bool {
	filecfg := &FileCfg{}
	filecfg.Filename = lc.File.Filename
	checkAndCreateLogDir(lc.File.Filename)
	filecfg.Level = LoggerLevel[lc.File.Level]
	filecfg.MaxLines = lc.File.Maxlines
	filecfg.MaxSize = lc.File.Maxsize * 1024 * 1024
	filecfg.Daily = lc.File.Daily
	filecfg.MaxDays = lc.File.Maxdays
	filecfg.Rotate = lc.File.Rotate
	filecfg.Perm = lc.File.Perm

	/////////////////////////////////////
	byteconfig, err := json.Marshal(filecfg)
	if err != nil {
		fmt.Printf("set file logger,change to json failed:%s", err.Error())
		return false
	}

	seterr := log.SetLogger(logs.AdapterFile, string(byteconfig))
	if seterr != nil {
		fmt.Printf("set file logger failed:%s", seterr.Error())
		return false
	}
	return true
}

func checkAndCreateLogDir(fileName string) {
	if fileName == "" {
		return
	}

	// no,/, ./ ,../
	var index int
	if index = strings.LastIndex(fileName, "/"); index <= 2 {
		return
	}

	perm, _ := strconv.ParseInt("0660", 8, 64)
	if mkerr := os.MkdirAll(fileName[0:index], os.FileMode(perm)); mkerr != nil {
		fmt.Printf("make dir failed,mkerr:%s", mkerr)
	}
}

func printLoggerCfg(lc *types.Logger) {
	log.Warn("---------console-------")
	log.Warn("level:%s", lc.Console.Level)

	log.Warn("---------file-------")
	log.Warn("filename:%s", lc.File.Filename)
	log.Warn("level:%s", lc.File.Level)
	log.Warn("maxlines:%d", lc.File.Maxlines)
	log.Warn("maxsize:%d", lc.File.Maxsize)
	log.Warn("daily:%t", lc.File.Daily)
	log.Warn("maxdays:%d", lc.File.Maxdays)
	log.Warn("rotate:%t", lc.File.Rotate)
	log.Warn("perm:%s", lc.File.Perm)
}

func GetCfgFilePath() string {
	var err error
	AppPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", AppPath)
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", workPath)
	appConfigPath := filepath.Join(workPath, ConfigFilePath)
	if !fileExists(appConfigPath) {
		appConfigPath = filepath.Join(AppPath, ConfigFilePath)
		if !fileExists(appConfigPath) {
			goPath := getGoPath()
			for _, val := range goPath {
				appConfigPath = filepath.Join(val, "src", "discover", ConfigFilePath)
				fmt.Println(appConfigPath)
				if fileExists(appConfigPath) {
					return appConfigPath
				}
			}
			appConfigPath = "/"
		}
	}

	return appConfigPath
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getGoPath() []string {
	goPath := os.Getenv("GOPATH")
	fmt.Println(goPath)
	if strings.Contains(goPath, ";") { //windows
		return strings.Split(goPath, ";")
	} else if strings.Contains(goPath, ":") { //linux
		return strings.Split(goPath, ":")
	} else { //only one
		path := make([]string, 1, 1)
		path[0] = goPath
		return path
	}
}
