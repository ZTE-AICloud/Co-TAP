package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Unknwon/goconfig"
)

var (
	cfg               *goconfig.ConfigFile
	consulIPGlobal    string
	consulHTTPPort    string
	listenPortGlobal  string
	healthCheckEnable = strings.ToLower(os.Getenv("HEALTH_CHECK_ENABLE"))
)

const (
	ConsulIP   = "127.0.0.1"
	ListenPort = "8080"
)

func InitConfig() {
	var (
		err     error
		cfgpath string
	)

	if _, set := os.LookupEnv("UT_CURRENT_PATH_MODE"); set {
		cfgpath = "config.ini"
	} else if _, set := os.LookupEnv("UT_MODE"); set {
		_, b, _, _ := runtime.Caller(0)
		cfgpath = filepath.Join(filepath.Dir(b), "config.ini")
	} else {
		cfgpath = "config/config.ini"
	}

	log.Println("path:", cfgpath)

	cfg, err = goconfig.LoadConfigFile(cfgpath)
	if err != nil {
		log.Printf("read config file(%s) failed", cfgpath)
	}
}

// get CONSUL_IP from conf
var GetConsulIPFromConf = func() string {
	consulIP, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "CONSUL_IP")

	return consulIP
}

func GetConsulHTTPPortFromConf() string {
	consulIP, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "CONSUL_HTTP_PORT")

	return consulIP
}

// get HTTP_LISTEN_PORT from conf
var getHTTPListenPortFromConf = func() string {
	listenPort, _ := cfg.GetValue(goconfig.DEFAULT_SECTION, "UAPREGISTRY_PORT")
	return listenPort
}

var GetHTTPListenPort = func() string {
	if listenPortGlobal == "" {
		listenPortGlobal = os.Getenv("UAPREGISTRY_PORT")
	}
	if listenPortGlobal == "" {
		listenPortGlobal = getHTTPListenPortFromConf()
	}

	if listenPortGlobal == "" {
		listenPortGlobal = ListenPort
	}

	return listenPortGlobal
}

func GetHTTPBindIP() string {
	return os.Getenv("UAPREGISTRY_BIND_IP")
}

func GetHTTPIPs() []string {
	return []string{""}
}

// initConsulIP
var initConsulIP = func() string {
	var consulIP string

	consulIP = os.Getenv("CONSUL_IP")
	if consulIP == "" {
		consulIP = GetConsulIPFromConf()
	}
	if consulIP == "" {
		consulIP = ConsulIP
	}

	return consulIP
}

// get CONSUL_IP from conf and env
var GetConsulIP = func() string {
	if consulIPGlobal == "" {
		consulIPGlobal = initConsulIP()
	}

	return consulIPGlobal
}

// get CONSUL_HTTP_PORT from env
func GetConsulHTTPPort() string {
	if consulHTTPPort == "" {
		consulHTTPPort = os.Getenv("CONSUL_HTTP_PORT")
		if consulHTTPPort == "" {
			consulHTTPPort = GetConsulHTTPPortFromConf()
		}
		if consulHTTPPort == "" {
			consulHTTPPort = "8500"
		}
	}

	return consulHTTPPort
}

func GetInsecureSkipVerify() bool {
	v, _ := cfg.Bool(goconfig.DEFAULT_SECTION, "INSECURE_SKIP_VERIFY")
	return v
}

func GetHealthCheckEnable() string {
	return healthCheckEnable
}
