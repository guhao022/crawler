package analyze

import (
	"fmt"
	"io"
	"crawler/utils"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/num5/ider"
	"github.com/num5/logger"
	"crawler/analyze/db"
)

type HaHaMx struct {
	doc   *goquery.Document
	name  string
	host  string
	start time.Time
}

func NewHaHaMx() Handler {
	return &HaHaMx{
		name:  "hahamx",
		host:  "www.haha.mx",
		start: time.Now(),
	}
}

func (h *HaHaMx) Prepare(doc *goquery.Document) error {
	h.doc = doc
	return nil
}

func (h *HaHaMx) Process() (bool, error) {
	defer func() {

		if err := recover(); err != nil {
			logger.Fatalf("panic 错误：%s ", err)
		}
	}()
	
	doc := h.doc

	host := doc.Url.Host

	//var reg = regexp.MustCompile(`http:\/\/www.gbfzh.com\/[\w]+\/[\d]+\.html`)

	if host == h.host {

		source_url := doc.Url.String()

		var cate = "笑话"

		main := doc.Find(".joke-main")

		topic := main.Find(".joke-main-topic").Text()
		if topic == "小编说" || topic == "蘑菇广播" {
			return false, fmt.Errorf("不是笑话，跳过...")
		}

		content := main.Find(".joke-main-content")
		text := content.Find(".joke-main-content-text").Text()
		like := main.Find(".joke-main-footer").Find(".btn-icon-good").Text()
		comment_text := main.Find(".joke-main-footer").Find(".btn-icon-comment").Text()
		comment := strings.Trim(comment_text, "评论 (|)")

		_, fond := content.Find(".joke-main-content-img").Attr("src")
		if fond {
			// 保存图片
			/*content.Find(".joke-main-content-img").Each(func(i int, sel *goquery.Selection) {
				img_src, fond := sel.Attr("src")
				if fond {
					if strings.HasSuffix(img_src, ".gif") {
						imgUrl, err := h.SaveImage("https:"+img_src)
						if err != nil {
							panic("保存图片失败，错误原因：" + err.Error())
						}

						img := new(db.Image)
						img.SourceUrl = "https:"+img_src
						img.Path = imgUrl
						img.ReadNum = 0
						img.LikeNum, _ = strconv.Atoi(like)
						img.CommentNum, _ = strconv.Atoi(comment)
						img.CreatedAt = time.Now()

						img.Store()
					}
				}
			})*/

			return true, fmt.Errorf("图片保存成功...")
		}

		if len(text) <= 0 {
			return false, nil
		}

		joke := new(db.Joker)
		joke.SourceUrl = source_url
		joke.Title = "哈哈mx"
		joke.Category = cate
		joke.Topic = topic
		joke.Content = text
		joke.ReadNum = 0
		joke.LikeNum, _ = strconv.Atoi(like)
		joke.CommentNum, _ = strconv.Atoi(comment)
		joke.CreatedAt = time.Now()

		joke.Store()

		return true, nil

	}

	return false, fmt.Errorf("获取内容失败")
}

func (h *HaHaMx) SaveImage(url string) (string, error) {
	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		return "", err
	}

	suffix := path.Ext(url)

	fpath := path.Join(h.name, strconv.Itoa(time.Now().Year()), strconv.Itoa(utils.Month(time.Now().Month().String())), strconv.Itoa(time.Now().Day()))

	filepath := path.Join(os.Getenv("IMG_PATH"), fpath)

	if _, err := os.Stat(filepath); err != nil {

		if os.IsNotExist(err) {

			err = os.MkdirAll(filepath, os.ModePerm)

			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	ide := ider.NewID(1)

	filename := strconv.FormatInt(ide.Next(), 10) + "." + suffix

	dst, err := os.Create(path.Join(filepath, filename))
	if err != nil {
		logger.Fatal(err.Error())
		return "", err
	}

	_, err = io.Copy(dst, res.Body)

	if err != nil {
		return "", err
	}

	var img_url = path.Join(fpath, filename)

	return img_url, nil
}

func (h *HaHaMx) SourceUrl() string {
	return h.doc.Url.String()
}

func (h *HaHaMx) Host() string {
	return h.host
}

func (h *HaHaMx) Name() string {
	return h.name
}

func (h *HaHaMx) Close() float64 {
	now := time.Now()
	start := h.start

	d := now.Sub(start).Nanoseconds()

	return float64(d) / 1e9
}

func init() {
	Register("hahamx", NewHaHaMx)
}
