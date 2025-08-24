package encoding

import (
	"compress/gzip"
	"io"
)

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewGzipReader(r io.ReadCloser) (io.ReadCloser, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		zr: zr,
	}, nil
}

func (g *gzipReader) Read(p []byte) (n int, err error) {
	return g.zr.Read(p)
}

func (g *gzipReader) Close() error {
	if err := g.zr.Close(); err != nil {
		return err
	}
	return g.r.Close()
}
