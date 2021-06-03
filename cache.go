package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

// Cache store file in .tmp and .info
type Cache struct {
	info FileInfo
	fd   *os.File
}

// NewCache initialize cache
func NewCache(url, path string, size int) (p *Cache, err error) {
	var c Cache
	err = c.Init(url, path, size)
	return &c, err
}

// Init initialize cache
func (c *Cache) Init(url, path string, size int) error {
	path = CleanSuffix(path)

	c.info.Init(url, path, size)
	err := c.checkInfo()
	if err != nil {
		c.info.Init(url, path, size)
	}
	err = c.checkPart()
	if err != nil {
		return err
	}
	c.fd, err = os.OpenFile(path+cacheTmpSuffix, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	return nil
}

// checkInfo read *.info, and update c.info.
// if *.info not exit,write c.info into it.
func (c *Cache) checkInfo() error {
	fd, err := os.OpenFile(c.info.Path+cacheInfoSuffix, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	err = readInfo(fd, &c.info)
	if err != nil {
		err = write(fd, c.info)
		return err
	}
	return nil
}

// checkPark read *.tmp and update c.info.part
func (c *Cache) checkPart() error {
	fd, err := os.OpenFile(c.info.Path+cacheTmpSuffix, os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	var s Store
	read(fd, '\n', func(b []byte) error {
		err := json.Unmarshal(b, &s)
		if err != nil {
			return nil
		}
		c.info.Parts[s.ID].done = true
		c.info.cur += len(s.Data)
		return nil
	})
	return nil
}

// Url return cache file download url.
func (c *Cache) Url() string {
	return c.info.URL()
}

// Save save s into *.tmp, update part.
func (c *Cache) Save(s Store) error {
	err := write(c.fd, s)
	if err != nil {
		return err
	}
	c.info.Parts[s.ID].done = true
	return nil
}

// GetUnDonePart return part not download.
func (c *Cache) GetUnDonePart() []FilePart {
	var tmp = make([]FilePart, c.info.Total)
	j := 0
	for _, p := range c.info.Parts {
		if !p.done {
			tmp[j] = p
			j++
		}
	}
	return tmp[:j]
}

// Clean clean file tmp.
func (c *Cache) Clean() {
	c.fd.Close()
	path := CleanSuffix(c.info.Path)
	CleanTmp(path)
}

// Merge merge .tmp store to destination path.
func (c *Cache) Merge() error {
	fc, err := os.Open(c.info.Path + cacheTmpSuffix)
	if err != nil {
		return err
	}
	defer fc.Close()

	fd, err := os.Create(c.info.Path)
	if err != nil {
		return err
	}
	defer fd.Close()

	writer := bufio.NewWriter(fd)
	fs, err := readStore(fc)
	if err != nil {
		return err
	}
	for _, s := range fs {
		_, err = writer.Write(s.Data)
		if err != nil {
			os.Remove(c.info.Path)
			return err
		}
	}
	c.Clean()
	return nil
}

// Write use json marshal i and write into f.
func write(f *os.File, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	f.WriteString("\n")
	return err
}

// read .info fileinfo.
func readInfo(f *os.File, i *FileInfo) error {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return errors.New("empty file:" + f.Name())
	}
	err = json.Unmarshal(b, i)
	if err != nil {
		return err
	}
	return nil
}

// read .tmp store.and sort by id.
// it will use last id replace preview id.
func readStore(f *os.File) (fs []Store, e error) {
	var tmp []Store
	max := 0
	err := read(f, '\n', func(b []byte) error {
		var s Store
		err := json.Unmarshal(b, &s)
		if err != nil {
			return err
		}
		tmp = append(tmp, s)
		if s.ID > max {
			max = s.ID
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	fs = make([]Store, max+1)
	for _, v := range tmp {
		fs[v.ID] = v
	}
	return
}

// read until delim using fn handle each data.
func read(f *os.File, delim byte, fn func([]byte) error) error {
	reader := bufio.NewReader(f)
	for {
		b, err := reader.ReadBytes(delim)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = fn(b)
		if err != nil {
			return err
		}
	}
	return nil
}
