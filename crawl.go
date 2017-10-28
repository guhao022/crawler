package main

import (
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/num5/logger"
	"math/rand"
)

var userAgent = [...]string{
	"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
	"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
	"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
	"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
	"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
	"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
	"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
	"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
	"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	}

const (
	//DefaultUserAgent         string        = `Mozilla/5.0 (Windows NT 6.1; rv:15.0) awcrawl/0.4 Gecko/20120716 Firefox/15.0a2`
	DefaultUserAgent string = `Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36`
	DefaultRobotUserAgent    string        = `google spider (awcrawl v0.8)`
	//DefaultRobotUserAgent    string        = ``
	DefaultEnqueueChanBuffer int           = 10000
	DefaultCrawlDelay        time.Duration = 100 * time.Millisecond
	DefaultIdleTTL           time.Duration = 1 * time.Second
)

var crawler *gocrawl.Crawler

type Options struct {

	// 爬取延时
	CrawlDelay time.Duration

	// WorkerIdleTTL是一个工作者在被清除（其goroutine终止）之前允许的空闲生存时间。 爬网延迟不是空闲时间的一部分，这特别是工作人员可用时的时间，但没有要处理的URL。
	WorkerIdleTTL time.Duration

	// 最大爬取数
	MaxVisits int

	// 用于向主机发出请求的用户代理值
	UserAgent string

	// RobotUserAgent是robot的用户代理值，用于在主机的robots.txt文件中查找匹配策略。 它不用于制作robots.txt请求，仅用于匹配策略。 应始终将其设置为搜寻器应用程序的名称，以便网站所有者可以相应地配置robots.txt。
	RobotUserAgent string

	// EnqueueChanBuffer是入列通道的缓冲区大小
	EnqueueChanBuffer int

	// SameHostOnly将URL限制为仅排入与来自种子URL的那些目标相同的主机的URL。
	SameHostOnly bool

	// HeadBeforeGet要求爬行器在发出最终GET请求之前发出HEAD请求。 如果设置为true，则在HEAD之后调用扩展器方法RequestGet以控制是否应发出GET。
	HeadBeforeGet bool
}

func GetUserAgent() string {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return userAgent[r.Intn(len(userAgent))]
}

func NewOptions() *Options {

	return &Options{
		DefaultCrawlDelay,
		DefaultIdleTTL,
		0,
		GetUserAgent(),
		DefaultRobotUserAgent,
		DefaultEnqueueChanBuffer,
		false,
		false,
	}
}

func Stop() {
	crawler.Stop()
}

func Run(op *Options) {
	err := NewCrawler(op)

	if err != nil {
		log.Info("爬虫停止工作：" + err.Error())
		os.Exit(1)
	}

}

func NewCrawler(op *Options) error {

	log.Info("皮皮虾，我们走......")

	seeds := FetchSeeds(true)

	ext := &Crawler{&gocrawl.DefaultExtender{}}
	opts := gocrawl.NewOptions(ext)

	opts.CrawlDelay = op.CrawlDelay
	opts.WorkerIdleTTL = op.WorkerIdleTTL
	opts.MaxVisits = op.MaxVisits
	opts.UserAgent = op.UserAgent
	opts.RobotUserAgent = op.RobotUserAgent
	opts.EnqueueChanBuffer = op.EnqueueChanBuffer
	opts.SameHostOnly = op.SameHostOnly
	opts.HeadBeforeGet = op.HeadBeforeGet

	opts.LogFlags = gocrawl.LogError

	crawler = gocrawl.NewCrawlerWithOptions(opts)

	//err := crawler.Run(seeds)
	err := crawler.Run(seeds)

	return err
}

func FetchSeeds(whole bool) []string {
	hosts := os.Getenv("CRAWL_SEEDS")
	seeds := strings.Split(hosts, ",")

	if whole {
		var wseeds = make([]string, len(seeds))

		for _, url := range seeds {
			wseeds = append(wseeds, url)
		}

		return wseeds
	}

	return seeds
}

var log *logger.Log

func init() {

	// 初始化
	log = logger.NewLog(1000)

	// 设置log级别
	log.SetLevel("Debug")

	// 设置输出引擎
	log.SetEngine("file", `{"level":4, "spilt":"size", "filename":".storage/logs/spider.log", "maxsize":10}`)

	//log.DelEngine("console")

	// 设置是否输出行号
	log.SetFuncCall(true)
}
