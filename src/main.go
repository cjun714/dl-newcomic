package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"./log"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var siteURL = "https://newcomic.info/"
var baseURL = "https://newcomic.info/page/"
var baseDir = "z:/"

func main() {
	start, e := strconv.Atoi(os.Args[1])
	if e != nil {
		panic(e)
	}
	end, e := strconv.Atoi(os.Args[2])
	if e != nil {
		panic(e)
	}

	e = downloadAll(baseURL, start, end)
	if e != nil {
		log.E(e)
	}
}

func downloadAll(baseURL string, startNum, endNum int) error {
	start := strconv.Itoa(startNum)
	end := strconv.Itoa(endNum)
	baseDir = baseDir + "/" + start + "-" + end
	e := os.Mkdir(baseDir, os.ModeDir) // create dir like z:/1-100
	if e != nil {
		return e
	}
	pageDir := baseDir + "/" + "pages"
	e = os.Mkdir(pageDir, os.ModeDir) // create dir like z:/1-100/pages
	if e != nil {
		return e
	}
	imageDir := baseDir + "/" + "images"
	e = os.Mkdir(imageDir, os.ModeDir) // create dir like z:/1-100/images
	if e != nil {
		return e
	}

	for i := startNum; i <= endNum; i++ {
		n := strconv.Itoa(i)
		url := baseURL + n

		// download index page
		indexPagePath := baseDir + "/" + strconv.Itoa(i) + ".html"
		log.I("download:", url)
		e := downloadHTML(url, indexPagePath)
		if e != nil {
			log.E("download index page failed:", url, ",", e)
			continue
		}

		// parse index page
		infoList, e := parseIndexPage(indexPagePath)
		if e != nil {
			log.E("parse index page failed:", indexPagePath, ",", e)
			continue
		}

		var wg sync.WaitGroup
		for _, mi := range infoList {
			// donwload cover image into /images
			wg.Add(1)
			coverPath := imageDir + "/" + getNameFromURL(mi.ImageURL)
			go func(url, path string) {
				defer wg.Done()
				e := downloadImage(url, path)
				if e != nil {
					log.E("download image failed:", url, ",", e)
				}
			}(mi.ImageURL, coverPath)

			// download detail page into /pages
			wg.Add(1)
			detailPagePath := pageDir + "/" + getNameFromURL(mi.DetailPageURL)
			go func(u, p string) {
				defer wg.Done()
				e = downloadHTML(u, p)
				if e != nil {
					log.E("download detail page failed:", u, ",", e)
				}
			}(mi.DetailPageURL, detailPagePath)

			time.Sleep(300 * time.Millisecond)
		}
		wg.Wait()

	}

	log.I("done")

	return nil
}

func downloadHTML(url string, targetPath string) error {
	// get html into bytes[]
	res, e := http.Get(url)
	if e != nil {
		return e
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("get html faild,resp code:" + strconv.Itoa(res.StatusCode))
	}
	bs, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return e
	}

	// save html content into file
	f, e := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0660)
	if e != nil {
		return e
	}
	defer f.Close()

	wr := bufio.NewWriter(f)
	_, e = wr.Write(bs)
	if e != nil {
		return e
	}

	return nil
}

// MangaInfo manga information
type MangaInfo struct {
	Title         string
	ImageURL      string
	DetailPageURL string
	DownloadURL   string
	Tags          string
}

func parseIndexPage(path string) ([]MangaInfo, error) {
	bs, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}

	r := bytes.NewReader(bs)
	doc, e := goquery.NewDocumentFromReader(r)
	if e != nil {
		return nil, e
	}

	infoList := make([]MangaInfo, 0, 1)

	doc.Find(".newcomic-short").Each(func(i int, s *goquery.Selection) {
		var mi MangaInfo

		// read pages + size
		fileInfo := s.Find(".newcomic-mask-top").Clone().Children().Remove().End().Text()
		fileInfo = strings.ReplaceAll(fileInfo, "\n", "")
		fileInfo = strings.Trim(fileInfo, "	")

		// read tags
		s.Find(".newcomic-mask-top a").Each(func(i int, s *goquery.Selection) {
			tag := s.Text()
			mi.Tags = mi.Tags + "|" + tag
		})

		s.Find(".newcomic-mask-bottom a").Each(func(i int, s *goquery.Selection) {
			title, _ := s.Attr("title")
			url, _ := s.Attr("href")
			mi.Title = title
			mi.DetailPageURL = url
		})

		// img
		imgURL, _ := s.Find("img").Attr("src")
		if strings.Index(imgURL, "/") == 0 {
			imgURL = siteURL + imgURL
		}
		mi.ImageURL = imgURL

		infoList = append(infoList, mi)
	})

	return infoList, e
}

func getNameFromURL(url string) string {
	idx := strings.LastIndex(url, "/")
	if idx == -1 {
		return "errorname"
	}
	return url[idx:]
}

func downloadImage(url string, targetPath string) error {
	resp, e := http.Get(url)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	f, e := os.Create(targetPath)
	if e != nil {
		return e
	}
	defer f.Close()

	_, e = io.Copy(f, resp.Body)
	if e != nil {
		return e
	}

	return nil
}
