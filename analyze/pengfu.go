package analyze

import (
	"time"
	"fmt"
	"crawler/analyze/db"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"github.com/num5/logger"
)

type PengFu struct {
	doc   *goquery.Document
	name  string
	host  string
	start time.Time
}

func NewPengFu() Handler {
	return &PengFu{
		name:  "pengfu",
		host:  "www.pengfu.com",
		start: time.Now(),
	}
}

func (p *PengFu) Prepare(doc *goquery.Document) error {
	p.doc = doc
	return nil
}

func (p *PengFu) Process() (bool, error) {

	defer func() {

		if err := recover(); err != nil {
			logger.Fatalf("panic 错误：%s ", err)
		}
	}()

	doc := p.doc

	host := doc.Url.Host

	if host == p.host {

		source_url := doc.Url.String()

		main := doc.Find(".w960").Find(".w645").Find(".list-item")

		content := main.Find(".content-txt").Text()

		action := main.Find(".action")
		like := action.Find(".ding").Find("em").Text()
		comment := action.Find(".det-commentClick").Find("em").Text()

		_, fond := main.Find(".content-txt").Find("img").Attr("src")
		if fond {
			return false, fmt.Errorf("图片暂不保存...")
		}

		if len(content) <= 0 {
			return false, nil
		}

		joke := new(db.Joker)
		joke.SourceUrl = source_url
		joke.Title = "捧腹"
		joke.Category = " "
		joke.Content = content
		joke.ReadNum = 0
		joke.LikeNum, _ = strconv.Atoi(like)
		joke.CommentNum, _ = strconv.Atoi(comment)
		joke.CreatedAt = time.Now()

		joke.Store()

		return true, nil

	}

	return false, fmt.Errorf("获取内容失败")
}

func (p *PengFu) SourceUrl() string {
	return p.doc.Url.String()
}

func (p *PengFu) Host() string {
	return p.host
}

func (p *PengFu) Name() string {
	return p.name
}

func (p *PengFu) Close() float64 {
	now := time.Now()
	start := p.start

	d := now.Sub(start).Nanoseconds()

	return float64(d) / 1e9
}

func init() {
	Register("pengfu", NewPengFu)
}
