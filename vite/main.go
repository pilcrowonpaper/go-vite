package vite

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//go:embed .html/*
//go:embed .assets/*
var embedded embed.FS
var assetFS, _ = fs.Sub(embedded, ".assets")
var htmlFS, _ = fs.Sub(embedded, ".html")

func GetHTML(name string) ([]byte, error) {
	if serverEnv := os.Getenv("ENV"); serverEnv == "PROD" {
		id := base64.RawURLEncoding.EncodeToString([]byte(name))
		htmlFile, err := htmlFS.Open(id + ".html")
		if err != nil {
			return []byte{}, err
		}
		html, err := io.ReadAll(htmlFile)
		if err != nil {
			return []byte{}, err
		}
		return html, nil
	}
	response, err := http.Get(fmt.Sprintf("http://localhost:%v", vitePort()) + "/" + name)
	if err != nil {
		return []byte{}, err
	}
	if response.StatusCode == 404 {
		return []byte{}, os.ErrNotExist
	}
	node, err := html.Parse(response.Body)
	if err != nil {
		return []byte{}, err
	}
	walk(node, name)
	var transformedHTML bytes.Buffer
	err = html.Render(&transformedHTML, node)
	if err != nil {
		return []byte{}, err
	}
	return transformedHTML.Bytes(), nil
}

func GetStaticAsset(name string) (fs.File, error) {
	asset, err := assetFS.Open(name)
	return asset, err
}

func walk(node *html.Node, htmlFilename string) {
	if node.Type == html.ElementNode && node.Data == "head" {
		for child := node.FirstChild; child != nil; {
			handleHeadChild(child, htmlFilename)
			child = child.NextSibling
		}
		return
	}
	for child := node.FirstChild; child != nil; {
		walk(child, htmlFilename)
		child = child.NextSibling
	}
}

func handleHeadChild(headChild *html.Node, htmlFilename string) {
	if headChild.Type == html.ElementNode && headChild.Data == "script" {
		for i := range headChild.Attr {
			if headChild.Attr[i].Key == "src" {
				src := headChild.Attr[i].Val
				isLocalImport := strings.HasPrefix(src, ".") || strings.HasPrefix(src, "/")
				if isLocalImport {
					importPath := src
					isRelativeImport := strings.HasPrefix(importPath, "../") || strings.HasPrefix(importPath, "./")
					if isRelativeImport {
						importPath = "/" + filepath.Join(filepath.Dir(htmlFilename), importPath)
					}
					headChild.Attr[i].Val = fmt.Sprintf("http://localhost:%v", vitePort()) + importPath
				}
			}
		}
	}
}

func vitePort() int {
	portEnv := os.Getenv("VITE_PORT")
	if portEnv == "" {
		return 5173
	}
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		panic("VITE_PORT must be an int")
	}
	return port
}
