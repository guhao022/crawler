package analyze

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
)

type Handler interface {
	Prepare(doc *goquery.Document) error
	Process() (bool, error)
	SourceUrl() string
	Host() string
	Name() string
	Close() float64
}

type Analyze struct {
	doc     *goquery.Document
	handler map[string]Handler
	lock    sync.Mutex
}

func NewAnalyze(doc *goquery.Document) *Analyze {
	a := &Analyze{
		doc:     doc,
		handler: make(map[string]Handler),
	}

	return a
}

// 定义处理引擎字典
type engineType func() Handler

var engines = make(map[string]engineType)

// 注册引擎
func Register(name string, a engineType) {
	if a == nil {
		panic("analyze engine: Register provide is nil")
	}
	if _, dup := engines[name]; dup {
		panic("analyze engine: Register called twice for provider " + name)
	}

	engines[name] = a
}

func (p *Analyze) SetHandler(name, conf string) *Analyze {
	p.lock.Lock()
	defer p.lock.Unlock()

	//获取引擎
	if handle, ok := engines[name]; ok {
		h := handle()
		err := h.Prepare(p.doc)
		if err != nil {
			errmsg := fmt.Errorf("SetEngine error: %s", err)
			fmt.Println(errmsg.Error())
			return nil
		}

		p.handler[name] = h
	} else {
		fmt.Printf("unknown Engine % ", name)
		return nil
	}

	return p
}

func (a *Analyze) Hosts() []string {
	var hosts []string

	for _, e := range a.handler {
		hosts = append(hosts, e.Host())
	}

	return hosts
}

func (a *Analyze) Process() (bool, float64, error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	var iscr bool
	var used float64
	var err error

	for _, e := range a.handler {

		if a.doc.Url.Host == e.Host() {
			iscr, err = e.Process()
			used = e.Close()
			break
		}
	}

	return iscr, used, err
}
