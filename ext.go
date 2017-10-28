package main

import (
	"crawler/analyze"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"os"
	"strings"
)

// 实现爬虫扩展器
type Crawler struct {
	*gocrawl.DefaultExtender
}

//
func (e *Crawler) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {

	log.Tracf("爬取地址：%s", ctx.URL().String())

	p := analyze.NewAnalyze(doc)

	ana := os.Getenv("CRAWL_ANALYZE")

	anas := strings.Split(ana, ",")

	for _, a := range anas {

		p = p.SetHandler(a, "")
	}

	pro, usetime, err := p.Process()

	if err != nil {
		log.Warnf("解析失败：%s", err.Error())
		return nil, true
	}

	if pro {
		log.Tracf("爬取成功，用时： %f 秒, %v", usetime, pro)
	}

	return nil, true
}

func (e *Crawler) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}

	strf := os.Getenv("CRAWL_SEEDS")
	slf := strings.Split(strf, ",")
	for _, f := range slf {

		f = strings.TrimLeft(f, "http://|https://")

		if ctx.URL().Host == f {
			return true
		}
	}

	return false
}

func (e *Crawler) RequestRobots(ctx *gocrawl.URLContext, robotAgent string) (data []byte, doRequest bool) {
	return nil, false
}
