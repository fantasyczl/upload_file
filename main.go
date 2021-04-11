package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const DataRoot = "./data"

func main() {
	port := ":10001"
	fmt.Printf("addr: http://localhost%s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := "./resource/view/index.html"

		bs, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatalf("path:%s read failed, %s", path, err)
		}

		w.Write(bs)
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		f, h, err := r.FormFile("file")
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		defer f.Close()

		if h.Filename == "" {
			w.Write([]byte("File name is empty"))
			return
		}

		targetPath := filepath.Join(DataRoot, h.Filename)

		//log.Printf("f:%+v\n", f)
		log.Printf("name:%+v\n", h.Filename)
		wf, err := os.Create(targetPath)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("create file failed. path:%s, err:%s", targetPath, err)))
			return
		}
		defer wf.Close()

		var buf = make([]byte, 1024)
		for {
			n, err := f.Read(buf)
			if err != nil {
				if err != io.EOF {
					w.Write([]byte(fmt.Sprintf("read file failed. err:%s", err)))
				}

				break
			}

			log.Printf("read num:%d\n", n)
			wf.Write(buf[:n])
		}
	})

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
