package conf

import (
	"fmt"
	"os"

	"github.com/iotames/miniutils"
)

const ENV_PROD = "prod"
const ENV_DEV = "dev"
const ENV_FILE = ".env"

var WorkEnvFile string

func setEnvFile() {
	WorkEnvFile = os.Getenv("DUTI_ENV_FILE")
	if WorkEnvFile == "" {
		WorkEnvFile = ENV_FILE
	}
}

func initEnvFile() []string {
	var err error
	var files []string
	if !miniutils.IsPathExists(WorkEnvFile) {
		err = createEnvFile(WorkEnvFile)
		if err != nil {
			panic(err)
		}
		files = append(files, WorkEnvFile)
		fmt.Printf("Create file %s SUCCESS\n", WorkEnvFile)
	}
	files = append(files, WorkEnvFile)
	if !miniutils.IsPathExists(DEFAULT_ENV_FILE) {
		err = createEnvFile(DEFAULT_ENV_FILE)
		if err != nil {
			fmt.Printf("--------initEnvFile(%s)err(%v)\n", DEFAULT_ENV_FILE, err)
			return files
		}
		fmt.Printf("Create file %s SUCCESS\n", DEFAULT_ENV_FILE)
	}
	files = append(files, DEFAULT_ENV_FILE)
	return files
}

func createEnvFile(fpath string) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("create env file(%s)err(%v)", fpath, err)
	}
	_, err = f.WriteString(getAllConfEnvStrDefault())
	if err != nil {
		return fmt.Errorf("write env file(%s)err(%v)", fpath, err)
	}
	return f.Close()
}
