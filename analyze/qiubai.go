package analyze

import (
	"time"
	"fmt"
	"crawler/analyze/db"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"github.com/num5/logger"
)

type QiuBai struct {
	doc   *goquery.Document
	name  string
	host  string
	start time.Time
}

func NewQiuBai() Handler {
	return &HaHaMx{
		name:  "qiubai",
		host:  "www.qiushibaike.com",
		start: time.Now(),
	}
}

func (q *QiuBai) Prepare(doc *goquery.Document) error {
	q.doc = doc
	return nil
}

func (q *QiuBai) Process() (bool, error) {
	defer func() {

		if err := recover(); err != nil {
			logger.Fatalf("panic 错误：%s ", err)
		}
	}()

	doc := q.doc

	host := doc.Url.Host

	if host == q.host {

		source_url := doc.Url.String()

		main := doc.Find("#content")

		cate := main.Find(".source-column").Text()

		block := main.Find(".content-block")
		content := block.Find(".content").Text()
		stats := block.Find(".stats")
		like := stats.Find(".stats-vote").Find(".number").Text()
		comment := stats.Find(".stats-comments").Find(".number").Text()

		_, fond := block.Find(".thumb").Find("img").Attr("src")
		if fond {
			return false, fmt.Errorf("图片暂不保存...")
		}

		if len(content) <= 0 {
			return false, nil
		}

		joke := new(db.Joker)
		joke.SourceUrl = source_url
		joke.Title = "糗事百科"
		joke.Category = cate
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

func (q *QiuBai) SourceUrl() string {
	return q.doc.Url.String()
}

func (q *QiuBai) Host() string {
	return q.host
}

func (q *QiuBai) Name() string {
	return q.name
}

func (q *QiuBai) Close() float64 {
	now := time.Now()
	start := q.start

	d := now.Sub(start).Nanoseconds()

	return float64(d) / 1e9
}

func init() {
	Register("qiubai", NewQiuBai)
}
