package scm

// This file was automatically generated by go-binary (github.com/jimmysawczuk/go-binary)
// at 2016-01-07T13:56:38-0500; compression = 9, -Inf% saved

import (
	"bytes"
	"compress/gzip"
	"io"
)

func getVersion() ([]byte, error) {

	in := bytes.NewBuffer([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x01, 0x00,
		0x00, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	out := bytes.NewBuffer([]byte{})

	gz, err := gzip.NewReader(in)
	if err != nil {
		return []byte{}, err
	}
	io.Copy(out, gz)

	return out.Bytes(), nil
}
