package environment

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	AppName    string
	AppVersion string
	LogLevel   string

	RestListenAddr string

	BypassAuth bool

	ESRestAddr    string
	ESGoodsIndex  string
	ESBoardsIndex string

	JDCookie string
	TBCookie string
}

func Get() (Env, error) {
	var err error

	if err = load(); err != nil {
		return Env{}, err
	}

	var appName string
	if os.Getenv("APP_NAME") == "" {
		return Env{}, fmt.Errorf("app name is required")
	} else {
		appName = os.Getenv("APP_NAME")
	}

	var appVersion string
	if os.Getenv("APP_VERSION") == "" {
		return Env{}, fmt.Errorf("app version is required")
	} else {
		appVersion = os.Getenv("APP_VERSION")
	}

	var logLevel string
	if os.Getenv("LOG_LEVEL") == "" {
		logLevel = "INFO"
	} else {
		logLevel = os.Getenv("LOG_LEVEL")
	}

	var restListenAddr string
	if os.Getenv("REST_LISTEN_ADDR") == "" {
		restListenAddr = "localhost:8090"
	} else {
		restListenAddr = os.Getenv("REST_LISTEN_ADDR")
	}

	var bypassAuth bool
	if os.Getenv("BYPASS_AUTH") == "" {
		bypassAuth = false
	} else {
		bypassAuth, err = strconv.ParseBool(os.Getenv("BYPASS_AUTH"))
		if err != nil {
			return Env{}, err
		}
	}

	var esRestAddr string
	if os.Getenv("ES_REST_ADDR") == "" {
		esRestAddr = "http://localhost:9200"
	} else {
		esRestAddr = os.Getenv("ES_REST_ADDR")
	}

	var esGoodsIndex string
	if os.Getenv("ES_GOODS_INDEX") == "" {
		esGoodsIndex = "goods"
	} else {
		esGoodsIndex = os.Getenv("ES_GOODS_INDEX")
	}

	var esBoardsIndex string
	if os.Getenv("ES_BOARDS_INDEX") == "" {
		esBoardsIndex = "boards"
	} else {
		esBoardsIndex = os.Getenv("ES_BOARDS_INDEX")
	}

	var jdCookie string
	if os.Getenv("JD_COOKIE") == "" {
		return Env{}, fmt.Errorf("jd cookie is required")
	} else {
		jdCookie = os.Getenv("JD_COOKIE")
	}

	var tbCookie string
	if os.Getenv("TB_COOKIE") == "" {
		return Env{}, fmt.Errorf("taobao cookie is required")
	} else {
		tbCookie = os.Getenv("TB_COOKIE")
	}

	return Env{
		AppName:        appName,
		AppVersion:     appVersion,
		LogLevel:       logLevel,
		RestListenAddr: restListenAddr,
		BypassAuth:     bypassAuth,
		ESRestAddr:     esRestAddr,
		ESGoodsIndex:   esGoodsIndex,
		ESBoardsIndex:  esBoardsIndex,
		JDCookie:       jdCookie,
		TBCookie:       tbCookie,
	}, nil
}

var (
	envLoaded = false
)

func load() error {
	if envLoaded {
		return nil
	}

	// 定位到根目录的.env文件
	_, f, _, _ := runtime.Caller(0) // 当前执行的文件，即此文件environment/env.go
	basepath := filepath.Dir(f)
	envFile := path.Join(basepath, "../.env")
	if err := godotenv.Load(envFile); err != nil {
		return err
	}
	envLoaded = true
	return nil
}
