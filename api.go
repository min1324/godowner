package main

import (
	"path/filepath"
)

// Clean tmp dir.
func Clean(dir string) {
	CleanDir(dir)
}

// download url's file into path
func FromUrl(url, path string) error {
	name, size, err := checkHead(url)
	if err != nil {
		return err
	}
	path = filepath.Join(path, name)
	c, err := NewCache(url, path, size)
	if err != nil {
		return err
	}
	err = down(c)
	if err != nil {
		return err
	}

	err = c.Merge()
	if err != nil {
		return err
	}
	CleanTmp(path)
	return nil
}

// FromDir download local tmp file.
func FromTmp(path string) error {
	c, err := NewCache("", path, 0)
	if err != nil {
		return err
	}
	_, _, err = checkHead(c.Url())
	if err != nil {
		return err
	}
	err = down(c)
	if err != nil {
		return err
	}
	err = c.Merge()
	if err != nil {
		return err
	}
	CleanTmp(path)
	return nil
}

// FromDir download dir all tmp's.
func FromDir(dir string) error {
	tmp := DirInfo(dir)
	var err error
	for _, p := range tmp {
		err = FromTmp(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func down(c *Cache) error {

	// progress bar
	bar := NewProcessBar(c.info.Size, 0)
	defer bar.Close()
	bar.Add(c.info.cur)

	parts := c.GetUnDonePart()
	for _, p := range parts {
		b, err := downPart(c.Url(), p, bar)
		if err != nil {
			return err
		}
		s := Store{
			ID:   p.Id,
			Data: b,
		}
		err = c.Save(s)
		if err != nil {
			return err
		}
	}
	return nil
}
