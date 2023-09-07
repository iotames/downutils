package conf

type SpiderConf struct {
	UserAgent                 string
	HttpProxy                 string
	Cookie                    string
	ImagesPath, DownloadsPath string
	Timeout                   int
	RetryTimes                int
	// Async, Limit bool
	// Delay is the duration to wait before creating a new request to the matching domains
	// RandomDelay is the extra randomized duration to wait added to Delay before creating a new request
	Delay, RandomDelay, Parallelism int
}

func GetSpiderConf() SpiderConf {
	return SpiderConf{
		UserAgent:     userAgent,
		HttpProxy:     httpProxy,
		Cookie:        cookie,
		ImagesPath:    imgsPath,
		DownloadsPath: downloadsPath,
		Timeout:       timeout,
		RetryTimes:    retry,
		Delay:         delay,
		RandomDelay:   randomDelay,
		Parallelism:   parallelism,
	}
}
