package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/iotames/miniutils"
	"github.com/joho/godotenv"
)

func LoadEnv() {
	setEnvFile()
	efiles := initEnvFile()
	err := godotenv.Load(efiles...)
	if err != nil {
		panic(fmt.Errorf("godotenv.Load(%v)err(%v)", efiles, err))
	}
	getConfByEnv()
}

func UpdateConf(mp map[string]string) error {
	for k, v := range mp {
		os.Setenv(k, v)
	}
	return godotenv.Write(mp, ENV_FILE)
}

func GetLogger() *miniutils.Logger {
	lgdir := filepath.Join(runtimeDir, "logs")
	if runtimeDir == "" {
		lgdir = filepath.Join(DEFAULT_RUNTIME_DIR, "logs")
	}
	return miniutils.GetLogger(lgdir)
}
