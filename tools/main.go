package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
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

	infos, err := scrapeCSVInfos(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("unepxected error: %v", err))
	}

	for _, info := range infos {
		err := func() error {
			req, err := http.NewRequest(http.MethodGet, info.DownloadURL(), nil)
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

func (info CSVInfo) DownloadURL() string {
	return fmt.Sprintf("%s?dlFilKanriNo=%s&jinkakukbn=%s&type=csv", downloadURL, info.FileNum, info.Section)
}

func scrapeCSVInfos(r io.Reader) ([]CSVInfo, error) {
	htmlNode, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parsing HTML: %w", err)
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
		return nil, errors.New("failed to find table#DLtableCsv")
	}

	var csvInfos []CSVInfo
	var findCSVInfos func(*html.Node)
	findCSVInfos = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "onclick" && strings.HasPrefix(a.Val, "return doDownload(") {
					start := strings.Index(a.Val, "(")
					end := strings.Index(a.Val, ")")
					if start != -1 && end != -1 && start < end {
						attrVal := a.Val[start+1 : end]
						params := strings.Split(strings.Trim(attrVal, "'"), "','")
						if len(params) >= 2 {
							csvInfos = append(csvInfos, CSVInfo{FileNum: params[0], Section: params[1]})
						}
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findCSVInfos(c)
		}
	}
	findCSVInfos(tableNode)

	return csvInfos, nil
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
