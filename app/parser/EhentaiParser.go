package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	mo "github.com/enify/multi-downloader/app/model"
	"github.com/enify/multi-downloader/app/request"
	"github.com/enify/multi-downloader/app/util"
	"github.com/PuerkitoBio/goquery"
)

// EhentaiParser support Ehentai Album download
type EhentaiParser struct{ BaseParser }

// GetMeta get meta of this parser
func (parser EhentaiParser) GetMeta() Meta {
	return Meta{
		URLRgx:   `^https://e-hentai.org/g/\w+/\w+/?$`,
		Priority: 9,

		Name:         "Ehentai下载",
		InternalName: "ehentai-parser",
		Version:      "0.1",
		Description:  "可解析下载Ehentai相册",
		Author:       "",
		Link:         "",
	}
}

// Prepare task
func (parser EhentaiParser) Prepare(task *mo.Task, client *request.HTTPClient) (err error) {
	client.SetTimeout(90 * time.Second)

	mainPage, err := client.Req("GET", task.URL, nil, "", nil)
	if err != nil {
		return fmt.Errorf("request task url: E:%w", err)
	}

	if mainPage.StatusCode != http.StatusOK {
		mainPage.Body.Close()
		return fmt.Errorf("request task url: E:status:%s", mainPage.Status)
	}

	mainDom, err := goquery.NewDocumentFromResponse(mainPage)
	if err != nil {
		return fmt.Errorf("gen task dom: E:%w", err)
	}

	task.Title = mainDom.Find("#gn").Text()
	task.FileSize = -1
	task.Path = filepath.Join(task.Path, util.FollowPathRule(task.Title))
	err = os.MkdirAll(task.Path, 0755)
	if err != nil {
		return fmt.Errorf("create task path: path:%s, E:%w", task.Path, err)
	}

	coverURL, _ := mainDom.Find("#gd1 > div").Attr("style")
	exp := regexp.MustCompile(`http.*(.jpg|.jpeg|.png|.bmp|.gif)`)
	coverURL = exp.FindString(coverURL)
	cover, err := client.Fetch("GET", coverURL, nil, "", nil)
	if err != nil {
		return fmt.Errorf("request task cover: E:%w", err)
	}

	task.Preview = filepath.Join(task.Path, fmt.Sprintf("cover%s", filepath.Ext(coverURL)))
	err = ioutil.WriteFile(task.Preview, cover, 0664)
	if err != nil {
		return fmt.Errorf("write task cover file: E:%w", err)
	}

	task.Meta = map[string]string{
		"Title":     task.Title,
		"URL":       task.URL,
		"Category":  mainDom.Find("#gdc > .cs").Text(),
		"Uploader":  mainDom.Find("#gdn > a").Text(),
		"Posted":    mainDom.Find("#gdd tr:nth-child(1) > .gdt2").Text(),
		"Parent":    mainDom.Find("#gdd tr:nth-child(2) > .gdt2").Text(),
		"Visible":   mainDom.Find("#gdd tr:nth-child(3) > .gdt2").Text(),
		"Language":  mainDom.Find("#gdd tr:nth-child(4) > .gdt2").Text(),
		"FileSize":  mainDom.Find("#gdd tr:nth-child(5) > .gdt2").Text(),
		"Length":    mainDom.Find("#gdd tr:nth-child(6) > .gdt2").Text(),
		"Favorited": mainDom.Find("#gdd tr:nth-child(7) > .gdt2").Text(),
		"Rating":    mainDom.Find("#rating_label").Text(),
	}

	type imgPage struct {
		Title string
		URL   string
	}
	var imgPages = []imgPage{}
	var titleExp = regexp.MustCompile(`Page\s\d+:\s`)
	nowDom := mainDom
	for {
		nowDom.Find("#gdt > .gdtm a").Each(func(index int, el *goquery.Selection) {
			title, _ := el.Find("img").Attr("title")
			url, _ := el.Attr("href")
			title = titleExp.ReplaceAllString(title, "")
			imgPages = append(imgPages, imgPage{title, url})
		})
		if nextEl := nowDom.Find(".ptt td:last-child > a"); nextEl.Is("a") {
			nextURL, _ := nextEl.Attr("href")
			nextPage, err := client.Req("GET", nextURL, nil, "", nil)
			if err != nil {
				return fmt.Errorf("request page url: url:%s, E:%w", nextURL, err)
			}
			if nextPage.StatusCode != http.StatusOK {
				nextPage.Body.Close()
				return fmt.Errorf("request page url: url:%s, E:status:%s", nextURL, nextPage.Status)
			}
			nowDom, err = goquery.NewDocumentFromResponse(nextPage)
			if err != nil {
				return fmt.Errorf("gen page dom: url:%s, E:%w", nextURL, err)
			}
		} else {
			break
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(imgPages))
	for _, page := range imgPages {
		go func(p imgPage) {
			defer wg.Done()

			subtask := &mo.SubTask{
				FileName: p.Title,
				Status:   mo.StatusPending,
			}

			subtask.URL, err = func() (string, error) {
				resp, err := client.Req("GET", p.URL, nil, "", nil)
				if err != nil {
					return "", fmt.Errorf("request img page: url:%s, E:%w", p.URL, err)
				}
				if resp.StatusCode != http.StatusOK {
					resp.Body.Close()
					return "", fmt.Errorf("request img page: url:%s, E:status:%s", p.URL, resp.Status)
				}
				dom, err := goquery.NewDocumentFromResponse(resp)
				if err != nil {
					return "", fmt.Errorf("gen img page dom: url:%s, E:%w", p.URL, err)
				}
				src, exists := dom.Find("#img").Attr("src")
				if !exists {
					return "", fmt.Errorf("find img url: from dom of:%s, E:not found", p.URL)
				}
				return src, nil
			}()
			if err != nil {
				subtask.Status = mo.StatusError
				subtask.Err = err
			}

			task.AddSubTask(subtask)
		}(page)
	}
	wg.Wait()
	return
}
