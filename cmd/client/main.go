package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/alexadhy/wfreq/internal/logging"
	"github.com/alexadhy/wfreq/internal/model"
)

var (
	inputString   string
	minWordLength int
	requestURL    string
)

const (
	requestTimeout    = 5 * time.Second
	defaultRequestURL = "http://localhost:3334/"
)

type requestForm struct {
	body        *bytes.Buffer
	contentType string
	contentLen  int
}

func main() {
	l := logging.New()

	flag.StringVar(&inputString, "i", "some input string", "input string or input file in plain text utf-8 format")
	flag.IntVar(&minWordLength, "wl", 2, "minimum word length to be counted")
	flag.StringVar(&requestURL, "u", defaultRequestURL, "wfreq-svc url (default to http://localhost:3334)")
	flag.Parse()

	var request http.Request

	// check if input string is a file
	stat, err := os.Stat(inputString)
	if err == nil {
		u, err := url.Parse(requestURL + "upload")
		if err != nil {
			l.Fatal(err)
		}
		query := u.Query()
		query.Set("min_word", strconv.Itoa(minWordLength))
		u.RawQuery = query.Encode()

		rf, err := newRequestForm(inputString, stat)
		if err != nil {
			l.Fatal(err)
		}
		header := make(http.Header)
		header.Set("Content-Type", rf.contentType)
		request = http.Request{
			Method:        http.MethodPost,
			URL:           u,
			Header:        header,
			Body:          ioutil.NopCloser(rf.body),
			ContentLength: int64(rf.contentLen),
		}
	} else {
		u, err := url.Parse(requestURL)
		if err != nil {
			l.Fatal(err)
		}
		rb := model.RequestBody{
			MinWordLength: minWordLength,
			Content:       inputString,
		}
		b, err := json.Marshal(&rb)
		if err != nil {
			l.Fatal(err)
		}
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		request = http.Request{
			Method: http.MethodPost,
			URL:    u,
			Body:   ioutil.NopCloser(bytes.NewBuffer(b)),
		}
	}

	hc := http.Client{
		Timeout: requestTimeout,
	}
	res, err := hc.Do(&request)
	if err != nil {
		l.Fatal(err)
	}

	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			l.Error(err)
		}
	}(res.Body)

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		l.Fatal(err)
	}

	var plist []model.Pair
	if err := json.Unmarshal(b, &plist); err != nil {
		l.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err = enc.Encode(&plist); err != nil {
		l.Fatal(err)
	}
}

func newRequestForm(inputString string, stat os.FileInfo) (*requestForm, error) {
	hdr := make(textproto.MIMEHeader)
	cd := mime.FormatMediaType("form-data", map[string]string{
		"name":     "file",
		"filename": inputString,
	})
	hdr.Set("Content-Disposition", cd)
	hdr.Set("Content-Type", "text/plain; charset=utf-8")
	hdr.Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)

	part, err := mw.CreatePart(hdr)
	if err != nil {
		return nil, fmt.Errorf("failed to create new form part: %v", err)
	}

	fd, err := os.Open(inputString)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}

	n, err := io.Copy(part, fd)
	if err != nil {
		return nil, fmt.Errorf("failed to write form part: %w", err)
	}

	if n != stat.Size() {
		return nil, fmt.Errorf("file size changed while writing: %s", fd.Name())
	}

	err = mw.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare form: %w", err)
	}

	return &requestForm{
		body:        &buf,
		contentType: mw.FormDataContentType(),
		contentLen:  buf.Len(),
	}, nil
}
