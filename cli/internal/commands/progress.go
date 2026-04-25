package commands

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func printProgress(label string, done, total int64) {
	if total <= 0 {
		fmt.Fprintf(os.Stderr, "\r%s...", label)
		return
	}
	pct := float64(done) / float64(total)
	const width = 20
	filled := int(pct * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	fmt.Fprintf(os.Stderr, "\r%s [%s] %3d%%", label, bar, int(pct*100))
}

func clearProgress() {
	fmt.Fprint(os.Stderr, "\r"+strings.Repeat(" ", 60)+"\r")
}

type progressReader struct {
	r     io.Reader
	total int64
	read  int64
	label string
}

func newProgressReader(r io.Reader, total int64, label string) *progressReader {
	return &progressReader{r: r, total: total, label: label}
}

func (p *progressReader) Read(buf []byte) (int, error) {
	n, err := p.r.Read(buf)
	p.read += int64(n)
	printProgress(p.label, p.read, p.total)
	return n, err
}
