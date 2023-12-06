package service

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/iotames/downutils/conf"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/iotames/miniutils"
)

const LOCAL_IMAGE_FILE_EXT = ".jpg"

type ErrorRequest struct {
	Error         error
	Retry         int64
	Url           string
	CollyResponse *colly.Response
}

type SiteSpider struct {
	Name            string
	BaseUrl         string
	AllowedDomains  []string
	sconf           conf.SpiderConf
	Collector       *colly.Collector
	HttpRequest     *miniutils.HttpRequest
	OnRequest       func(r *http.Request)
	args            map[string]interface{}
	ErrorRequests   []ErrorRequest
	isStartBrowser  bool
	isStartChromedp bool
	// collyCtx        *colly.Context
}

func (s SiteSpider) GetSiteUrl() string {
	return s.BaseUrl
}

func (s SiteSpider) GetErrorRequest(errUrl string) (int64, ErrorRequest) {
	errReq := ErrorRequest{}
	var errIndex int64 = -1
	for i, eachReq := range s.ErrorRequests {
		if eachReq.Url == errUrl {
			errIndex = int64(i)
			errReq = eachReq
		}
	}
	return errIndex, errReq
}

// OnCollyRetryError. 错误自动重试。
// POST必须传递ctx参数request_body。因为r.Request.Body内容被Readed了，再次读取内容为空
func (s *SiteSpider) OnCollyRetryError(retryFail func(errReq ErrorRequest), maxTrytimes int) {
	c := s.GetCollector()
	logger := miniutils.GetLogger("")
	c.OnError(func(r *colly.Response, e error) {
		errUrl := r.Request.URL.String()
		retryAny := r.Ctx.GetAny("retry")
		retry := 0
		if retryAny != nil {
			retry = retryAny.(int)
		}
		if retry >= maxTrytimes {
			retryFail(ErrorRequest{Error: e, Retry: int64(retry), Url: errUrl, CollyResponse: r})
			return
		}
		retry++
		logger.Warn(fmt.Sprintf("---OnCollyRetryError---Error(%v)---Url(%s)----BeginRetry(%d)-------", e, errUrl, retry))
		r.Request.Ctx.Put("retry", retry)
		requestBody := r.Ctx.Get("request_body")
		if strings.ToUpper(r.Request.Method) == "POST" && requestBody == "" {
			bb, berr := ioutil.ReadAll(r.Request.Body)
			logger.Error(fmt.Sprintf("-------ErrForRetryPost---ctx.Get(request_body)Empty---Url(%s)---PostRAW(%s)ReadRequestBodyErr(%v)---RequestBody(%+v)--", errUrl, string(bb), berr, r.Request.Body))
		}
		if requestBody == "" {
			c.Request(r.Request.Method, errUrl, nil, r.Request.Ctx, nil)
		} else {
			c.Request(r.Request.Method, errUrl, strings.NewReader(requestBody), r.Request.Ctx, nil)
		}
	})
}

func (s *SiteSpider) GetAny(key string) interface{} {
	if s.args == nil {
		return nil
	}
	dt, ok := s.args[key]
	if ok {
		return dt
	}
	return nil
}
func (s *SiteSpider) SetAny(key string, value interface{}) {
	if s.args == nil {
		s.args = make(map[string]interface{})
	}
	s.args[key] = value
}

func (s *SiteSpider) SetAsyncAndLimit(c *colly.Collector, domainRegexp string) {
	c.Async = true
	sconf := s.sconf
	// 修改并发和随机延时设置
	err := c.Limit(&colly.LimitRule{
		DomainRegexp: domainRegexp,
		Delay:        time.Second * time.Duration(sconf.Delay),
		RandomDelay:  time.Second * time.Duration(sconf.RandomDelay),
		Parallelism:  int(sconf.Parallelism),
	})
	if err != nil {
		panic(err)
	}
}

func (s *SiteSpider) GetNewCollector() *colly.Collector {
	c := colly.NewCollector()

	c.UserAgent = s.sconf.UserAgent
	// "Sec-Ch-Ua-Platform": "Windows",
	c.SetDebugger(&debug.LogDebugger{})
	c.SetRequestTimeout(time.Duration(s.sconf.Timeout) * time.Second)
	// httpTransport := &http.Transport{
	// 	Proxy: ProxyFromEnvironment,
	// 	DialContext: (&net.Dialer{
	// 		Timeout:   30 * time.Second,
	// 		KeepAlive: 30 * time.Second,
	// 	}).DialContext,
	// 	ForceAttemptHTTP2:     true,
	// 	MaxIdleConns:          100,
	// 	IdleConnTimeout:       90 * time.Second,
	// 	TLSHandshakeTimeout:   10 * time.Second,
	// 	ExpectContinueTimeout: 1 * time.Second,
	// }
	// c.WithTransport(httpTransport)
	hproxy := s.sconf.HttpProxy
	if hproxy != "" {
		// "ss://rc4-md5:123456@ss.server.com:1080"
		c.SetProxy(hproxy)
		// proxy := func(_ *http.Request) (*url.URL, error) {
		// 	return url.Parse(config.HttpProxy)
		// }
		// c.WithTransport(&http.Transport{Proxy: proxy})
		log.Printf("SUCCESS CONFIG:colly.Collector.SetProxy(%v)------\n\n", hproxy)
	}

	// // 默认并发和随机延时设置
	// c.Limit(&colly.LimitRule{
	// 	RandomDelay: time.Second * 2, // 随机延时2秒
	// 	Parallelism: 3,               // 并发数为3
	// })

	return c
}

// GetCollector TODO set request header connection:close and "Sec-Ch-Ua-Platform": "Windows",
func (s *SiteSpider) GetCollector() *colly.Collector {
	if s.Collector == nil {
		s.Collector = s.GetNewCollector()
	}
	return s.Collector
}

func (s SiteSpider) NewCollyCtx() *colly.Context {
	return colly.NewContext()
}

// GetHttpRequest Set SetRequestHeader like User-Agent After SetRequest
func (s *SiteSpider) GetHttpRequest(url string) *miniutils.HttpRequest {
	if s.HttpRequest != nil {
		s.HttpRequest.Url = url
		return s.HttpRequest
	}

	s.HttpRequest = miniutils.NewHttpRequest(url)

	// 请求失败重试次数
	s.HttpRequest.RetryTimes = uint8(s.sconf.RetryTimes)

	// 请求超时设置
	if s.sconf.Timeout != 0 {
		s.HttpRequest.SetTimeout(uint8(s.sconf.Timeout))
	}

	// HTTP代理设置
	if s.sconf.HttpProxy != "" {
		s.HttpRequest.SetProxy(s.sconf.HttpProxy)
	}

	// 发起HTTP请求前的回调函数设置。回调函数套娃
	s.HttpRequest.OnRequest = func(r *http.Request) {
		// 套娃外层: miniutils.HttpRequest
		r.Header.Set("Sec-Ch-Ua-Platform", "Windows")
		r.Header.Set("User-Agent", s.sconf.UserAgent)
		r.Header.Set("Connection", "Close")
		r.Header.Set("referer", s.BaseUrl+"/")
		if s.OnRequest != nil {
			// 套娃内层: service.SiteSpider
			s.OnRequest(r)
		}
	}

	return s.HttpRequest
}

func GetLocalImageFileName(url string) string {
	return miniutils.Md5(url) + LOCAL_IMAGE_FILE_EXT
}

// GetLocalFile. fileExt, Like: .jpg .png. .zip
func (s *SiteSpider) GetLocalFile(fileUrl, dirname, fileExt string) string {
	fileUrl = strings.TrimSpace(fileUrl)
	filename := miniutils.Md5(fileUrl) + fileExt
	imageExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	basePath := s.sconf.DownloadsPath
	imgIndex := miniutils.GetIndexOf[string](fileExt, imageExts)
	if imgIndex > -1 {
		basePath = s.sconf.ImagesPath
	}
	localFilePath := fmt.Sprintf("%s/%s/%s", basePath, dirname, filename)
	// modTime, err := GetFileModTime(localFilePath)
	return localFilePath
}

func (s *SiteSpider) DownloadFile(url string, filepath string, referer string) {
	startAt := time.Now()
	if url == "" || !strings.Contains(url, "http") {
		log.Println("Skip: DownloadUrl is unavailable----------------")
		return
	}
	var lenfile int = 0
	isExist := miniutils.IsPathExists(filepath)
	if isExist {
		fbts, _ := ioutil.ReadFile(filepath)
		lenfile = len(fbts)
		// 文件大于1KB OR 大于500B才是有效文件。
		if lenfile > 500 {
			modTime, err := GetFileModTime(filepath)
			if err == nil {
				if time.Since(modTime).Hours() > 24*30 {
					os.Remove(filepath)
					log.Println("ReDownload---File--modTime > 24*30--", filepath)
					s.download(url, filepath, referer)
					return
				}
			}
			log.Println(time.Since(startAt).Milliseconds(), "ms. Skip. filepath is exist--SiteSpider DownloadFile---", filepath)
			return
		}
		delErr := os.Remove(filepath)
		if delErr != nil {
			log.Println("remove file Fail: ", delErr)
		}
		s.download(url, filepath, referer)
	}
	if !isExist {
		s.download(url, filepath, referer)
	}
}

func (s SiteSpider) download(url, filepath, referer string) {
	log.Println("------------dowloadFile--------------", url, filepath, referer)
	httpRequest := s.GetHttpRequest(url)
	s.OnRequest = func(req *http.Request) {
		if referer == "" {
			req.Header.Set("referer", s.BaseUrl+"/")
		} else {
			req.Header.Set("referer", referer)
		}
	}
	httpRequest.Download(filepath)
}

func (s SiteSpider) DownloadFileByColly(urls []string, getFilePath func(fileurl string) string, afterDownload func(fileurl, savepath string)) int {
	c := s.GetCollector()
	c.Async = true
	err := c.Limit(&colly.LimitRule{
		DomainRegexp: ".",
		RandomDelay:  time.Second * 1, // 随机延时1秒
		Parallelism:  12,              // 并发数为12
	})
	if err != nil {
		panic(err)
	}
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("referer", s.BaseUrl+"/")
	})
	totalresp := 0
	c.OnResponse(func(r *colly.Response) {
		totalresp++
		savePath := r.Ctx.Get("savePath")
		log.Println("Download To savePath:", savePath)
		f, err := os.Create(savePath)
		if err != nil {
			log.Printf("\n---Error Happened. os.Create:%v-------\n", err)
			return
		}
		_, err = io.Copy(f, bytes.NewReader(r.Body))
		if err != nil {
			log.Printf("---Error Happened. io.Copy:%v-------\n", err)
			return
		}
		if afterDownload != nil {
			afterDownload(r.Request.URL.String(), savePath)
		}
	})

	for _, furl := range urls {
		savePath := s.GetLocalFile(furl, s.Name, LOCAL_IMAGE_FILE_EXT)
		if getFilePath != nil {
			savePath = getFilePath(furl)
		}
		if miniutils.IsPathExists(savePath) {
			log.Printf("----SkipDownloadFile(%s)---SavePath(%s)---\n", furl, savePath)
			continue
		}
		ctx := colly.NewContext()
		ctx.Put("savePath", savePath)
		c.Request("GET", furl, nil, ctx, nil)
	}

	c.Wait()
	return totalresp
}

func CheckNetwork() error {
	var netWorkErr error
	ss := &SiteSpider{}
	c := ss.GetCollector()
	c.OnError(func(r *colly.Response, err error) {
		netWorkErr = err
	})
	c.Request("GET", "https://httpbin.org/get", nil, nil, nil)
	return netWorkErr
}

func NewSpider(siteName string) *SiteSpider {
	sconf := conf.GetSpiderConf()
	miniutils.Mkdir(sconf.ImagesPath + "/" + siteName)
	return &SiteSpider{
		Name:  siteName,
		sconf: sconf,
	}
}
