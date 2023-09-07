package conf

import (
	"os"
	"strconv"
)

const DEFALUT_USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"
const DEFALUT_TIMEOUT = 20
const DEFALUT_RETRY_TIMES = 3
const DEFALUT_DELAY = 1
const DEFALUT_RANDOM_DELAY = 1
const DEFALUT_PARALLELISM = 12
const DEFALUT_IMAGES_PATH = "runtime/images"
const DEFAULT_DOWNLOADS_PATH = "runtime/downloads"

type SpiderConf struct {
	UserAgent  string
	HttpProxy  string
	Cookie     string
	ImagesPath string
	Timeout    int
	RetryTimes int
	// Async, Limit bool
	// Delay is the duration to wait before creating a new request to the matching domains
	// RandomDelay is the extra randomized duration to wait added to Delay before creating a new request
	Delay, RandomDelay, Parallelism int
}

func GetSpiderConf() SpiderConf {
	imgsPath := os.Getenv("IMAGES_PATH")
	if imgsPath == "" {
		imgsPath = DEFALUT_IMAGES_PATH
	}
	userAgent := os.Getenv("USER_AGENT")
	if userAgent == "" {
		userAgent = DEFALUT_USER_AGENT
	}

	tout, _ := strconv.Atoi(os.Getenv("TIMEOUT"))
	if tout == 0 {
		tout = DEFALUT_TIMEOUT
	}

	retry, _ := strconv.Atoi(os.Getenv("RETRY_TIMES"))
	if retry == 0 {
		retry = DEFALUT_RETRY_TIMES
	}

	delay, _ := strconv.Atoi(os.Getenv("DELAY"))
	if delay == 0 {
		delay = DEFALUT_DELAY
	}

	randomDelay, _ := strconv.Atoi(os.Getenv("RANDOM_DELAY"))
	if randomDelay == 0 {
		randomDelay = DEFALUT_RANDOM_DELAY
	}

	parallelism, _ := strconv.Atoi(os.Getenv("PARALLELISM"))
	if parallelism == 0 {
		parallelism = DEFALUT_PARALLELISM
	}

	return SpiderConf{
		UserAgent:   userAgent,
		HttpProxy:   os.Getenv("HTTP_PROXY"),
		Cookie:      os.Getenv("COOKIE"),
		ImagesPath:  imgsPath,
		Timeout:     tout,
		RetryTimes:  retry,
		Delay:       delay,
		RandomDelay: randomDelay,
		Parallelism: parallelism,
	}
}
