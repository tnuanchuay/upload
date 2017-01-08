package main

import (
	"flag"
	"net/http"
	"fmt"
	"log"
	"io"
	"os"
	"io/ioutil"
	"github.com/parnurzeal/gorequest"
	"path/filepath"
)

var(
	path		string
)

func main(){
	serveMode := flag.Bool("serve", false, "serve mode")
	pathPtr := flag.String("path", "/var/www/html", "specific path to save file")
	server := flag.String("host", "", "a server to upload")
	flag.Parse()

	if *serveMode {
		f, err := os.OpenFile("/var/log/upload.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		log.SetOutput(f)

		path = *pathPtr
		if path[len(path) - 1] != '/'{
			path += "/"
		}

		http.HandleFunc("/", handler)
		http.ListenAndServe(":54321", nil)
	}else{
		args := os.Args
		if len(args) == 3{
			filename := args[2]
			host := fmt.Sprintf("http://%s:54321/", *server)

			f, _ := filepath.Abs(filename)
			byteOfFile, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatal("file err : ", err)
			}

			resp, body, errs := gorequest.New().
				Post(host).
				Type("multipart").
				SendFile(byteOfFile, filename, "upload").
				End()

			if errs != nil{
				for _, err := range errs{
					log.Println("gorequest", err)
				}
				os.Exit(1)
			}

			if resp.StatusCode == 200{
				fmt.Println(body)
			}else{
				fmt.Println("err", resp.StatusCode)
			}
		}else{
			log.Fatal("usage : upload --host=127.0.0.1 file.jpg")
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET"{
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("upload")
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()
		newpath := fmt.Sprintf("%s%s", path, handler.Filename)
		f, err := os.OpenFile(newpath, os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		w.Write([]byte(fmt.Sprintf("ok /%s", handler.Filename)))
		log.Println(fmt.Sprintf("/%s", handler.Filename))
	}else{
		w.Write([]byte("not use in method get"))
	}
}
