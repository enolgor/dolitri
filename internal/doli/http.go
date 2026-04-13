package doli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type DownloadResponse struct {
	Filename    string    `json:"filename"`
	ContentType string    `json:"content-type"`
	Filesize    int64     `json:"filesize"`
	Content     io.Reader `json:"-"`
	Encoding    string    `json:"encoding"`
}

func doGet(path string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", baseURL, path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("DOLAPIKEY", doliKey)
	return http.DefaultClient.Do(req)
}

func get(path string) (json.RawMessage, error) {
	resp, err := doGet(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func download(path string) (*DownloadResponse, error) {
	resp, err := doGet(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		Filename    string `json:"filename"`
		ContentType string `json:"content-type"`
		Filesize    int64  `json:"filesize"`
		Content     string `json:"content"`
		Encoding    string `json:"encoding"`
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var content io.Reader
	if raw.Encoding == "base64" {
		content = base64.NewDecoder(base64.StdEncoding, strings.NewReader(raw.Content))
	} else {
		content = strings.NewReader(raw.Content)
	}

	return &DownloadResponse{
		Filename:    raw.Filename,
		ContentType: raw.ContentType,
		Filesize:    raw.Filesize,
		Content:     content,
		Encoding:    raw.Encoding,
	}, nil
}

func DownloadInvoice(invoice Invoice) (*DownloadResponse, error) {
	return download("/documents/download?modulepart=facture&original_file=" + url.QueryEscape(invoice.FilePath()))
}

func DownloadExpense(name string) (*DownloadResponse, error) {
	return download("/documents/download?modulepart=expensereport&original_file=" + url.QueryEscape(name))
}
