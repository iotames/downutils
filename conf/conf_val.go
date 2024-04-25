package conf

import (
	"fmt"
	"os"
	"strings"

	"strconv"

	"github.com/iotames/miniutils"
)

var userAgent, imgsPath, cookie, httpProxy, runtimeDir, downloadsPath, FyneFont, FyneTitle string
var timeout, retry, delay, randomDelay, parallelism int
var dbDriver, dbName, dbUsername, dbPassword, dbHost string
var dbPort, dbNodeId int
var ImageHeight, ColWidth int
var ImageWebpage bool

func getConfByEnv() {
	FyneTitle = getEnvDefaultStr("FYNE_TITLE", DEFAULT_FYNE_TITLE)
	FyneFont = getEnvDefaultStr("FYNE_FONT", DEFAULT_FYNE_FONT)
	runtimeDir = getEnvDefaultStr("RUNTIME_DIR", DEFAULT_RUNTIME_DIR)
	downloadsPath = getEnvDefaultStr("DOWNLOADS_DIR", DEFAULT_DOWNLOADS_PATH)
	imgsPath = getEnvDefaultStr("IMAGES_PATH", DEFAULT_IMAGES_PATH)
	httpProxy = getEnvDefaultStr("HTTP_PROXY", DEFAULT_HTTP_PROXY)
	userAgent = getEnvDefaultStr("USER_AGENT", DEFAULT_USER_AGENT)
	timeout = getEnvDefaultInt("TIMEOUT", DEFAULT_TIMEOUT)
	retry = getEnvDefaultInt("RETRY_TIMES", DEFAULT_RETRY_TIMES)
	delay = getEnvDefaultInt("DELAY", DEFAULT_DELAY)
	randomDelay = getEnvDefaultInt("RANDOM_DELAY", DEFAULT_RANDOM_DELAY)
	parallelism = getEnvDefaultInt("PARALLELISM", DEFAULT_PARALLELISM)
	ImageHeight = getEnvDefaultInt("IMAGE_HEIGHT", DEFAULT_IMAGE_HEIGHT)
	ColWidth = getEnvDefaultInt("COL_WIDTH", DEFAULT_COL_WIDTH)
	ImageWebpage = getEnvDefaultBool("IMAGE_WEBPAGE", DEFAULT_IMAGE_WEBPAGE)

	dbDriver = getEnvDefaultStr("DB_DRIVER", DEFAULT_DB_DRIVER)
	dbName = getEnvDefaultStr("DB_NAME", DEFAULT_DB_NAME)
	dbUsername = getEnvDefaultStr("DB_USERNAME", DEFAULT_DB_USERNAME)
	dbPassword = getEnvDefaultStr("DB_PASSWORD", DEFAULT_DB_PASSWORD)
	dbHost = getEnvDefaultStr("DB_HOST", DEFAULT_DB_HOST)
	dbPort = getEnvDefaultInt("DB_PORT", DEFAULT_DB_PORT)
	dbNodeId = getEnvDefaultInt("DB_NODE_ID", DEFAULT_DB_NODE_ID)
	checkRuntimeDir()
}

func checkRuntimeDir() {
	logger := GetLogger()
	hitmsg := fmt.Sprintf("请检查配置文件，路径:(%s)\n", WorkEnvFile)

	if runtimeDir == "" {
		hitmsg += "RuntimeDir 配置项不能为空"
		logger.Error(hitmsg)
		panic("RuntimeDir can not be empty")
	}

	var err error
	if !miniutils.IsPathExists(runtimeDir) {
		fmt.Printf("------创建runtime目录(%s)--\n", runtimeDir)
		err = os.Mkdir(runtimeDir, 0755)
		if err != nil {
			fmt.Printf("----runtime目录(%s)创建失败(%v)---\n", runtimeDir, err)
			panic(err)
		}
	}
}

func getEnvDefaultStr(key, defval string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defval
	}
	return v
}

func getEnvDefaultInt(key string, defval int) int {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defval
	}
	vv, _ := strconv.Atoi(v)
	return vv
}

func getEnvDefaultBool(key string, defval bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defval
	}
	val := false
	if strings.EqualFold(v, "true") {
		val = true
	}
	return val
}
