package conf

import (
	"fmt"
)

const DEFAULT_ENV_FILE = "env.default"
const DEFAULT_DB_DRIVER = DRIVER_SQLITE3
const DEFAULT_DB_HOST = "127.0.0.1"
const DEFAULT_DB_PORT = 3306
const DEFAULT_DB_NAME = "downutils"
const DEFAULT_DB_USERNAME = "root"
const DEFAULT_DB_PASSWORD = "root"
const DEFAULT_DB_NODE_ID = 1

const DEFAULT_HTTP_PROXY = ""
const DEFAULT_USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
const DEFAULT_TIMEOUT = 20
const DEFAULT_RETRY_TIMES = 3
const DEFAULT_DELAY = 1
const DEFAULT_RANDOM_DELAY = 1
const DEFAULT_PARALLELISM = 12
const DEFAULT_IMAGE_HEIGHT = 76
const DEFAULT_COL_WIDTH = 76
const DEFAULT_IMAGE_WEBPAGE = false
const DEFAULT_IMAGES_PATH = "runtime/images"
const DEFAULT_DOWNLOADS_PATH = "runtime/downloads"
const DEFAULT_RUNTIME_DIR = "runtime"
const DEFAULT_FYNE_FONT = "resource/fonts/OPPOSans-H.ttf"
const DEFAULT_FYNE_TITLE = "DownUtils"

const ENV_FILE_CONTENT = `# 此文件由系统自动创建，配置项为默认值。可修改本目录下的 .env 文件，以更新默认值。
# DB_DRIVER support: mysql,sqlite3,postgres
DB_DRIVER = "%s"
DB_HOST = "%s"
# DB_PORT like: 3306(mysql); 5432(postgres)
DB_PORT = %d
DB_NAME = "%s"
# DB_USERNAME like: root, postgres
DB_USERNAME = "%s"
DB_PASSWORD = "%s"
DB_NODE_ID = %d

USER_AGENT = "%s"
HTTP_PROXY = "%s"
COOKIE = "%s"
TIMEOUT = %d
RETRY_TIMES = %d
DELAY = %d
RANDOM_DELAY = %d
PARALLELISM = %d
IMAGE_HEIGHT = %d
COL_WIDTH = %d
IMAGE_WEBPAGE = %t

# 该目录存放程序运行时产生的文件
RUNTIME_DIR = "%s"
IMAGES_PATH = "%s"
DOWNLOADS_DIR = "%s"
FYNE_FONT = "%s"
FYNE_TITLE = "%s"
`

func getAllConfEnvStrDefault() string {
	return fmt.Sprintf(ENV_FILE_CONTENT, DEFAULT_DB_DRIVER, DEFAULT_DB_HOST, DEFAULT_DB_PORT,
		DEFAULT_DB_NAME, DEFAULT_DB_USERNAME, DEFAULT_DB_PASSWORD, DEFAULT_DB_NODE_ID,
		DEFAULT_USER_AGENT, "", "", DEFAULT_TIMEOUT, DEFAULT_RETRY_TIMES, DEFAULT_DELAY,
		DEFAULT_RANDOM_DELAY, DEFAULT_PARALLELISM, DEFAULT_IMAGE_HEIGHT, DEFAULT_COL_WIDTH, DEFAULT_IMAGE_WEBPAGE, DEFAULT_RUNTIME_DIR, DEFAULT_IMAGES_PATH,
		DEFAULT_DOWNLOADS_PATH, DEFAULT_FYNE_FONT, DEFAULT_FYNE_TITLE)
}
