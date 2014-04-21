package hut

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
)

type Multiparts struct {
	buffer bytes.Buffer
	writer *multipart.Writer
}

type Part interface {
	AddPart(parts *Multiparts) error
}

type FilePart struct {
	Fieldname string
	Filename  string
	Reader    io.Reader
}

type FormPart struct {
	Fieldname string
	Value     string
}

func AddMultiparts(req *http.Request, parts ...Part) error {
	multiparts := &Multiparts{}

	for _, part := range parts {
		err := part.AddPart(multiparts)
		if err != nil {
			return err
		}
	}

	err := multiparts.writer.Close()
	if err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(&multiparts.buffer)
	req.ContentLength = int64(multiparts.buffer.Len())
	req.Header.Add("Content-Type", multiparts.writer.FormDataContentType())

	return nil
}

func (part FilePart) AddPart(parts *Multiparts) error {
	writer := parts.getMultipartWriter()
	formFile, err := writer.CreateFormFile(part.Fieldname, part.Filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(formFile, part.Reader)
	if err != nil {
		return err
	}

	return nil
}

func (part FormPart) AddPart(parts *Multiparts) error {
	writer := parts.getMultipartWriter()
	formField, err := writer.CreateFormField(part.Fieldname)
	if err != nil {
		return err
	}

	_, err = formField.Write([]byte(part.Value))
	if err != nil {
		return err
	}

	return nil
}

func (parts *Multiparts) getMultipartWriter() *multipart.Writer {
	if parts.writer == nil {
		parts.writer = multipart.NewWriter(&parts.buffer)
	}
	return parts.writer
}

func DebugResponse(resp *http.Response, includeBody bool) string {
	dump, _ := httputil.DumpResponse(resp, includeBody)
	return string(dump)
}

func DebugRequest(req *http.Request, includeBody bool) string {
	dump, _ := httputil.DumpRequest(req, includeBody)
	return string(dump)
}
