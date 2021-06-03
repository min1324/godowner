package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	partSize = 1 << 20

	cacheTmpSuffix  = ".tmp"
	cacheInfoSuffix = ".info"
)

// CleanSuffix clean path suffix
func CleanSuffix(path string) string {
	if strings.HasSuffix(path, cacheInfoSuffix) {
		path = strings.TrimSuffix(path, cacheInfoSuffix)
	} else if strings.HasSuffix(path, cacheTmpSuffix) {
		path = strings.TrimSuffix(path, cacheTmpSuffix)
	}
	path = filepath.Clean(path)
	return path
}

// CleanTmp
func CleanTmp(path string) {
	path = CleanSuffix(path)
	tmp := path + cacheTmpSuffix
	info := path + cacheInfoSuffix
	os.Remove(tmp)
	os.Remove(info)
}

// CleanTmp
func CleanDir(dir string) {
	filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(d.Name(), cacheInfoSuffix) ||
			strings.HasSuffix(d.Name(), cacheTmpSuffix) {
			os.Remove(path)
		}
		return nil
	})
}

// DirInfo return dir's .info file's path
func DirInfo(dir string) []string {
	var tmp []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		dir = filepath.Dir(path)
		if strings.HasSuffix(info.Name(), cacheInfoSuffix) {
			name := strings.Trim(info.Name(), cacheInfoSuffix)
			name = filepath.Join(dir, name)
			tmp = append(tmp, name)
		}
		return nil
	})
	return tmp
}

// convert file size int to string
func sizeConvert(s int64) string {
	t := []string{"B", "KB", "MB", "GB"}
	i := 0
	p := float64(s)
	for ; i < 4; i++ {
		q := p * 1.0 / 1024
		if q > 1 {
			p = q
			continue
		}
		break
	}
	return strconv.FormatFloat(p, 'f', 1, 64) + t[i]
}

// Store cache file part store in .tmp
type Store struct {
	ID   int    // 文件分片的序号
	Data []byte // http下载得到的文件内容
}

// FilePart seperate file into Total path
// each part msg
type FilePart struct {
	Id    int
	Start int
	End   int
	done  bool
}

// cache infomation store in .info
type FileInfo struct {
	Url   string     // file name
	Path  string     // cache path
	Size  int        // file size
	Total int        // total FilePart num
	Parts []FilePart // len=Total
	cur   int        // cur is done part size.
}

// Init initialize fileinfo
func (f *FileInfo) Init(url, path string, size int) {
	f.Url = url
	f.Path = path
	f.Size = size
	f.makePart()
}

// URL
func (f *FileInfo) URL() string {
	return f.Url
}

// sep file part
func (f *FileInfo) makePart() {
	f.Total = f.Size / partSize
	if f.Total < 1 {
		f.Total = 1
	}
	f.Parts = make([]FilePart, f.Total)
	for i := 0; i < f.Total; i++ {
		f.Parts[i].Id = i
		f.Parts[i].Start = i * partSize
		f.Parts[i].End = (i+1)*partSize - 1
		f.Parts[i].done = false
	}
	f.Parts[f.Total-1].End = f.Size - 1
}
