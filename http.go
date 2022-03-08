package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func downAll(url, dir string, bar *ProcessBar) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	name := filepath.Base(url)
	name = filepath.Join(dir, name)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, io.TeeReader(resp.Body, bar))
	return err
}

func downPart(url string, fp FilePart, bar *ProcessBar) ([]byte, error) {
	r, err := newRequest("GET", url)
	if err != nil {
		return nil, err
	}
	//log.Printf("开始[%d]下载from:%d to:%d\n", c.Id, c.Start, c.End)
	r.Header.Set("Range", fmt.Sprintf("bytes=%v-%v", fp.Start, fp.End))
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		return nil, fmt.Errorf(fmt.Sprintf("服务器错误状态码: %v", resp.StatusCode))
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(io.TeeReader(resp.Body, bar))
	return b, err
}

// check url if support.
func checkHead(url string) (name string, size int, e error) {
	r, err := newRequest("HEAD", url)
	if err != nil {
		return "", 0, err
	}
	rsp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", 0, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode > 299 {
		return "", 0, fmt.Errorf(fmt.Sprintf("Can't process, response is %v", rsp.StatusCode))
	}
	//检查是否支持 断点续传
	if rsp.Header.Get("Accept-Ranges") != "bytes" {
		return "", 0, errors.New("服务器不支持文件断点续传")
	}
	name = parseForm(rsp)
	if name == "" {
		name = filepath.Base(url)
	}
	l, err := strconv.Atoi(rsp.Header.Get("Content-Length"))
	return name, l, err
}

func newRequest(method string, url string) (*http.Request, error) {
	r, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agen", "downloadfile")
	return r, nil
}

func parseForm(resp *http.Response) string {
	cd := resp.Header.Get("Content-Disposition")
	if cd != "" {
		_, params, err := mime.ParseMediaType(cd)
		if err == nil {
			return params["filename"]
		}
	}
	filename := filepath.Base(resp.Request.URL.Path)
	return filename
}
