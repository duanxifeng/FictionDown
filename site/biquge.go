package site

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/ma6254/FictionDown/store"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// Biquge_1 笔趣阁标准页面
type Biquge_1 struct {
}

func (b *Biquge_1) BookInfo(body io.Reader) (s *store.Store, err error) {
	body = transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
	doc, err := htmlquery.Parse(body)
	if err != nil {
		return
	}

	s = &store.Store{}

	node_title := htmlquery.Find(doc, `//*[@id="info"]/h1`)
	if len(node_title) == 0 {
		err = fmt.Errorf("No matching title")
		return
	}
	s.BookName = htmlquery.InnerText(node_title[0])

	node_desc := htmlquery.Find(doc, `//*[@id="intro"]/p`)
	if len(node_desc) == 0 {
		err = fmt.Errorf("No matching desc")
		return
	}
	s.Description = strings.Replace(
		htmlquery.OutputHTML(node_desc[0], false),
		"<br/>", "\n",
		-1)

	var author = htmlquery.Find(doc, `//*[@id="info"]/p[1]`)
	s.Author = strings.TrimLeft(htmlquery.OutputHTML(author[0], false), "作\u00a0\u00a0\u00a0\u00a0者：")

	node_content := htmlquery.Find(doc, `//*[@id="list"]/dl/dd/a`)
	if len(node_desc) == 0 {
		err = fmt.Errorf("No matching contents")
		return
	}

	var vol = store.Volume{
		Name:     "正文",
		Chapters: make([]store.Chapter, 0),
	}

	for _, v := range node_content[9:] {
		//fmt.Printf("href: %v\n", chapter_u)
		chapterURL, err := url.Parse(htmlquery.SelectAttr(v, "href"))
		if err != nil {
			return nil, err
		}

		vol.Chapters = append(vol.Chapters, store.Chapter{
			Name: strings.TrimSpace(htmlquery.InnerText(v)),
			URL:  chapterURL.String(),
		})
	}
	s.Volumes = append(s.Volumes, vol)

	s.CoverURL = htmlquery.SelectAttr(htmlquery.FindOne(doc, `//*[@id="fmimg"]/img`), "src")

	return
}

func (b *Biquge_1) Chapter(body io.Reader) ([]string, error) {
	body = transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())
	doc, err := htmlquery.Parse(body)
	if err != nil {
		return nil, err
	}

	M := []string{}
	//list
	// nodeContent := htmlquery.Find(doc, `//div[@id="content"]/text()`)
	nodeContent := htmlquery.Find(doc, `//div[@id="content"]/p`)
	if len(nodeContent) == 0 {
		err = fmt.Errorf("No matching content")
		return nil, err
	}
	for _, v := range nodeContent {
		t := htmlquery.InnerText(v)
		t = strings.TrimSpace(t)
		M = append(M, t)
	}

	return M, nil
}
