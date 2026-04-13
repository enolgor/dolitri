package utils

import (
	"archive/zip"
	"compress/flate"
	"fmt"
	"io"
	"os"
)

type Zip struct {
	entries []zipEntry
}

type zipEntry struct {
	path   string
	reader io.Reader
}

func NewZip() *Zip {
	return &Zip{
		entries: []zipEntry{},
	}
}

func (z *Zip) AddEntry(path string, reader io.Reader) {
	fmt.Println("adding entry: " + path)
	z.entries = append(z.entries, zipEntry{path: path, reader: reader})
}

func (z *Zip) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	for _, entry := range z.entries {
		fw, err := w.CreateHeader(&zip.FileHeader{
			Name:   entry.path,
			Method: zip.Deflate,
		})
		if err != nil {
			return err
		}
		if _, err := io.Copy(fw, entry.reader); err != nil {
			return err
		}
	}
	return nil
}
