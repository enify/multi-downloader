package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/enify/multi-downloader/app/util"

	"github.com/PuerkitoBio/goquery"
	mo "github.com/enify/multi-downloader/app/model"
	"github.com/enify/multi-downloader/app/request"
)

// NhentaiParser support Nhentai Album download
type NhentaiParser struct{ BaseParser }

// GetMeta return meta of this parser
func (NhentaiParser) GetMeta() Meta {
	return Meta{
		URLRgx:   `^https://nhentai.net/g/\w+/?$`,
		Priority: 9,

		Name:         "Nhentai下载",
		InternalName: "nhentai-parser",
		Version:      "0.1",
		Description:  "可解析下载Nhentai图册",
		Author:       "",
		Link:         "",
	}
}

// Prepare task
func (p NhentaiParser) Prepare(task *mo.Task, client *request.HTTPClient) (err error) {
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

	task.Title = mainDom.Find("#info > h1").Text()
	task.FileSize = -1
	task.Path = filepath.Join(task.Path, util.FollowPathRule(task.Title))
	err = os.MkdirAll(task.Path, 0755)
	if err != nil {
		return fmt.Errorf("create task path: path:%s, E:%w", task.Path, err)
	}

	coverURL, _ := mainDom.Find("#cover img").Attr("data-src")
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
		"Title":  task.Title,
		"URL":    task.URL,
		"Length": mainDom.Find("#info > div:nth-of-type(1)").Text(),
		"Posted": mainDom.Find("#info > div:nth-of-type(2) > time").AttrOr("datetime", ""),
	}
	mainDom.Find("#tags > .tag-container").Each(func(idx int, el *goquery.Selection) {
		k := strings.TrimRight(strings.TrimSpace(el.Contents().First().Text()), ":")
		v := ""
		el.Find(".tags > .tag").Each(func(i int, e *goquery.Selection) {
			v += fmt.Sprintf("%s,", e.Contents().First().Text())
		})
		task.Meta[k] = v
	})

	task.SubTasks = []*mo.SubTask{}
	mainDom.Find("#thumbnail-container > .thumb-container > a > img").Each(func(idx int, el *goquery.Selection) {
		subtask := &mo.SubTask{
			Status: mo.StatusPending,
		}

		subtask.URL, err = func() (string, error) {
			s, exists := el.Attr("data-src")
			if !exists {
				return "", fmt.Errorf("find img url: from dom of:%s, E:not found", task.URL)
			}

			u, err := url.Parse(s)
			if err != nil {
				return "", fmt.Errorf("parse img url:%s, E:%w", s, err)
			}

			u.Host = "i.nhentai.net"
			u.Path = strings.Replace(u.Path, "t", "", 1)
			return u.String(), nil
		}()
		subtask.FileName = filepath.Base(subtask.URL)
		if err != nil {
			subtask.Status = mo.StatusError
			subtask.Err = err
		}

		task.AddSubTask(subtask)
	})

	task.ExternalFiles = []string{}
	var infoFile = filepath.Join(task.Path, "info.txt")
	var infoContent = ""
	for k, v := range task.Meta {
		infoContent += fmt.Sprintf("%s: %s\n", k, v)
	}
	err = ioutil.WriteFile(infoFile, []byte(infoContent), 0664)
	if err != nil {
		return fmt.Errorf("write task info file: E:%w", err)
	}

	task.ExternalFiles = append(task.ExternalFiles, infoFile)

	return
}
