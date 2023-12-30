package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go-vite/vite"
	"html/template"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
)

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.HandlerFunc(handleStaticAssets)))
	http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/" {
			w.WriteHeader(200)
			templateSource, err := vite.GetHTML("pages/index.html")
			if err != nil {
				panic(err)
			}
			var pageHTML bytes.Buffer
			htmlTemplate, _ := template.New(request.URL.Path).Parse(string(templateSource))
			htmlTemplate.Execute(&pageHTML, Data{
				Message: "Hello from Go!",
			})
			w.Write(pageHTML.Bytes())
			return
		}
		w.WriteHeader(404)
	})
	fmt.Println("Starting server on port 3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}

func handleStaticAssets(w http.ResponseWriter, request *http.Request) {
	asset, err := vite.GetStaticAsset(request.URL.Path)
	if errors.Is(err, fs.ErrNotExist) {
		w.WriteHeader(404)
		return
	}
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer asset.Close()
	contentType := mime.TypeByExtension(filepath.Ext(request.URL.Path))
	if contentType == "" {
		panic(fmt.Sprintf("Unknown file extension: %v", contentType))
	}
	w.Header().Set("Cache-Control", "max-age=31536000,immutable")
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(200)
	scanner := bufio.NewScanner(asset)
	for scanner.Scan() {
		w.Write(scanner.Bytes())
	}
}

type Data struct {
	Message string
}
