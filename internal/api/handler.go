package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/alexadhy/wfreq/internal/model"
	"github.com/alexadhy/wfreq/internal/store"
)

const (
	wlimit = 10
)

// handleUpload handles file upload, on done, it will count the top X most words and outputs it as JSON.
func (a *API) handleUpload(w http.ResponseWriter, r *http.Request) {
	minWordLength := 2

	if mw, ok := r.URL.Query()["min_word"]; ok {
		if m, err := strconv.Atoi(mw[0]); err == nil {
			minWordLength = m
		}
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		a.renderError(w, r, err, http.StatusBadRequest)
		return
	}

	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			a.log(r).Error(err)
		}
		a.s.Clear()
	}(file)

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		a.renderError(w, r, err, http.StatusBadRequest)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "text/plain; charset=utf-8" {
		a.renderError(w, r,
			errors.New(
				"the provided file format is not allowed, Please upload a txt file (plain) with utf-8 charset",
			), http.StatusBadRequest)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		a.renderError(w, r, err, http.StatusBadRequest)
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		a.renderError(w, r, err, http.StatusInternalServerError)
		return
	}

	rid := middleware.GetReqID(r.Context())
	rid = strings.ReplaceAll(rid, "/", "_")
	fileName := fmt.Sprintf("./uploads/%s.txt", rid)
	f, err := os.Create(fileName)
	if err != nil {
		a.renderError(w, r, err, http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(f, file)
	if err != nil {
		a.renderError(w, r, err, http.StatusInternalServerError)
		return
	}

	_ = f.Close() //nolint:errcheck
	f, err = os.Open(fileName)
	if err != nil {
		a.renderError(w, r, err, http.StatusInternalServerError)
		return
	}

	defer func() {
		_ = f.Close()              //nolint:errcheck
		_ = os.RemoveAll(fileName) //nolint:errcheck
	}()

	plist, err := readInput(f, a.s, minWordLength, wlimit)
	if err != nil {
		a.renderError(w, r, err, 0)
		return
	}

	b, err := json.Marshal(plist)
	if err != nil {
		a.renderError(w, r, err, 0)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		a.renderError(w, r, err, 0)
		return
	}
	a.s.Clear()
}

func (a *API) handleWordFrequencies(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			a.log(r).Error(err)
		}
		a.s.Clear()
	}(r.Body)

	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		a.renderError(w, r, errors.New("content-type is not application/json"), http.StatusUnsupportedMediaType)
		return
	}

	var rb model.RequestBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&rb); err != nil {
		a.renderError(w, r, err, http.StatusBadRequest)
		return
	}

	buf := bytes.NewBufferString(rb.Content)
	plist, err := readInput(buf, a.s, rb.MinWordLength, wlimit)
	if err != nil {
		a.renderError(w, r, err, http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(plist)
	if err != nil {
		a.renderError(w, r, err, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out) //nolint:errcheck
}

// readInput function scans for words in the file / input string, filter by the word length, and limit to top-X number
// of the output
func readInput(reader io.Reader, s *store.Store, minWordLength, limit int) (*model.PairList, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		t := scanner.Text()
		if len(t) >= minWordLength {
			v := s.Load(t)
			s.Store(t, v+1)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	plist := s.LoadAll(true, limit)
	return &plist, nil
}
