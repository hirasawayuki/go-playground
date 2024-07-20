package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

const (
	baseURL     = "https://www.invoice-kohyo.nta.go.jp"
	targetURL   = baseURL + "/download/zenken#csv-filtyp"
	downloadURL = baseURL + "/download/zenken/dlfile"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
		return
	}
	os.Exit(0)
}

func run() error {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	listCSVInfo, err := scrapeCSVInfo(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("unepxected error: %v", err))
	}

	for _, info := range listCSVInfo {
		err := func() error {
			u, err := url.Parse(downloadURL)
			if err != nil {
				return err
			}
			q := u.Query()
			q.Set("dlFilKanriNo", info.FileNum)
			q.Set("jinkakukbn", info.Section)
			q.Set("type", "csv")
			u.RawQuery = q.Encode()

			req, err := http.NewRequest(http.MethodGet, u.String(), nil)
			if err != nil {
				return err
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("bad status: %s", resp.Status)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if err := uploadToS3(body); err != nil {
				return err
			}
			return nil
		}()

		if err != nil {
			return err
		}
	}
	return nil
}

type CSVInfo struct {
	FileNum string
	Section string
}

func (info CSVInfo) DownloadURL() (string, error) {
	u, err := url.Parse(downloadURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Set("dlFilKanriNo", info.FileNum)
	q.Set("jinkakukbn", info.Section)
	q.Set("type", "csv")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func scrapeCSVInfo(r io.Reader) ([]CSVInfo, error) {
	htmlNode, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var tableNode *html.Node
	var findTable func(*html.Node)
	findTable = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "DLtableCsv" {
					tableNode = n
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTable(c)
		}
	}
	findTable(htmlNode)
	if tableNode == nil {
		return nil, err
	}

	var listCSVInfo []CSVInfo
	var extractCSVInfo func(*html.Node)
	extractCSVInfo = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "onclick" && strings.HasPrefix(a.Val, "return doDownload(") {
					start := strings.Index(a.Val, "(")
					end := strings.Index(a.Val, ")")
					if start != -1 && end != -1 && start < end {
						attrVal := a.Val[start+1 : end]
						params := strings.Split(strings.Trim(attrVal, "'"), "','")
						if len(params) >= 2 {
							listCSVInfo = append(listCSVInfo, CSVInfo{FileNum: params[0], Section: params[1]})
						}
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractCSVInfo(c)
		}
	}

	extractCSVInfo(tableNode)
	return listCSVInfo, nil
}

func uploadToS3(data []byte) error {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		if filepath.Ext(file.Name) != ".csv" {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join("extracted", file.Name)
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
