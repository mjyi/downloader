package ooxx

import (
	"github.com/mjyi/downloader"
	"net/http"
	"encoding/json"
	"log"
	"io/ioutil"
	"os"
	"io"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

const OOXXURL = "http://i.jandan.net/?oxwlxojflwblxbsapi=jandan.get_ooxx_comments&page=%d"
const dir = "ooxx_images"
// init 初始化操作
func init() {
	InitDB()
	initDL()
	os.Mkdir(dir, 0777)
}

type OOXXDownLoader struct {
	PageDownloader *downloader.Downloader
	Max uint32
}

var OXDOwnloader *OOXXDownLoader


func initDL() {
	pageDL := downloader.NewDownloader(
		downloader.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"),
		downloader.Async(true),
	)

	pageDL.OnResponse(func(response *http.Response) {
		url := response.Request.URL.String()
		if strings.HasPrefix(url, "http://i.jandan.net") {
			parsingResponse(response)
		} else {
			saveImageData(response)
		}
	})

	pageDL.OnError(func(response *http.Response, e error) {
		fmt.Println(response.Request.URL, e)
	})

	OXDOwnloader = &OOXXDownLoader{
		PageDownloader:pageDL,
	}
}

func parsingResponse(resp *http.Response) {

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	var ret OOXXResult
	if err := json.Unmarshal(data, &ret); err != nil {
		log.Println(err)
		return
	}

	if ret.CurrentPage <= ret.PageCount {
		go StartPage(ret.CurrentPage + 1)
	}

	oxs := ret.Comments
	InsertModels(oxs)
	for _, ox := range oxs {
		for _, pic := range ox.Pics {
			go downloadFile(pic)
		}
	}
}


func Start() {
	StartPage(1)
	OXDOwnloader.PageDownloader.Wait()
	CloseDB()
}

func StartPage(page int32)  {
	url := fmt.Sprintf(OOXXURL, page)
	OXDOwnloader.PageDownloader.Get(url)
}

func downloadFile(filePath string) {
	url := strings.Replace(filePath, "mw600", "large", 1)
	time.Sleep(50 * time.Millisecond)
	OXDOwnloader.PageDownloader.Get(url)
}

func saveImageData(resp *http.Response)  {
	defer resp.Body.Close()
	fpath := resp.Request.URL.String()
	name := filepath.Base(fpath)
	name = filepath.Join(dir, name)
	if PathExist(name) {
		return
	}
	
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	size, err := io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(name, size, resp.Status)
}

func PathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}