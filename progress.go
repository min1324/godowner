package main

import (
	"fmt"
	"strings"
)

type ProcessBar struct {
	Total int
	Cur   int
}

func NewProcessBar(max, cur int) *ProcessBar {
	return &ProcessBar{Total: max, Cur: cur}
}

func (wc *ProcessBar) Write(p []byte) (int, error) {
	n := len(p)
	wc.Cur += n
	wc.render()
	return n, nil
}

func (wc *ProcessBar) Close() {
	fmt.Println("")
}

func (wc *ProcessBar) Add(n int) {
	wc.Cur += n
}

func (wc *ProcessBar) render() {
	cur := wc.Cur
	max := wc.Total

	present := float64(cur) / float64(max) * 100
	i := int(present / 4)
	if i > 25 {
		i = 25
	}
	h := strings.Repeat("▅", i) + strings.Repeat(" ", 25-i)

	fmt.Printf("\r%s", strings.Repeat(" ", 78))
	fmt.Printf("\r%-8s[%s]%3.0f%%",
		format(max), h, present)
}

func format(s int) string {
	var tmp = []string{"B", "KB", "MB", "GB", "TB"}
	i, p, q := 0, 0.0, float64(s)
	for ; i < len(tmp); i++ {
		p = q / 1024
		if p < 1 {
			break
		}
		q = p
	}
	return fmt.Sprintf("%.2f%s", q, tmp[i])
}
