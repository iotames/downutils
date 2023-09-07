package conf

import (
	"os"

	"github.com/iotames/miniutils"
	"github.com/joho/godotenv"
)

const ENV_PROD = "prod"
const ENV_DEV = "dev"
const ENV_FILE = ".env"
const DEFAULT_ENV_FILE = "env.default"

const DRIVER_SQLITE3 = "sqlite3"
const DRIVER_MYSQL = "mysql"
const DRIVER_POSTGRES = "postgres"
const SQLITE_FILENAME = "sqlite3.db"

const DEFAULT_ENV_FILE_CONTENT = `# 此文件由系统自动创建，配置项为默认值。可修改本目录下的 .env 文件，以更新默认值。
# DB_DRIVER support: mysql,sqlite3,postgres
DB_DRIVER = "sqlite3"
DB_HOST = "127.0.0.1"
# DB_PORT like: 3306(mysql); 5432(postgres)
DB_PORT = 3306
DB_NAME = "DownUtils"
# DB_USERNAME like: root, postgres
DB_USERNAME = "root"
DB_PASSWORD = "root"
DB_NODE_ID = 1

USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
HTTP_PROXY = ""
COOKIE = ""
TIMEOUT = 20
RETRY_TIMES = 3
DELAY = 1
RANDOM_DELAY = 1
PARALLELISM = 12
`

var envFile string

func setEnvFile() {
	envFile = os.Getenv("DUTI_ENV_FILE")
	if envFile == "" {
		envFile = ENV_FILE
	}
}

func LoadEnv() {
	setEnvFile()
	initEnvFile()
	err := godotenv.Load(ENV_FILE, DEFAULT_ENV_FILE)
	if err != nil {
		panic("godotenv Error: " + err.Error())
	}
}

func initEnvFile() {
	if !miniutils.IsPathExists(ENV_FILE) {
		f, err := os.Create(ENV_FILE)
		if err != nil {
			panic("Create .env Error: " + err.Error())
		}
		f.Close()
	}
	if !miniutils.IsPathExists(DEFAULT_ENV_FILE) {
		f, err := os.Create(DEFAULT_ENV_FILE)
		if err != nil {
			panic("Create .env Error: " + err.Error())
		}
		f.WriteString(DEFAULT_ENV_FILE_CONTENT)
		f.Close()
	}
}

func UpdateConf(mp map[string]string) error {
	for k, v := range mp {
		os.Setenv(k, v)
	}
	return godotenv.Write(mp, ENV_FILE)
}
