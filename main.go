package main

import (
	"flag"
	"net/http"
	"fmt"
	"log"
	"io"
	"os"
)

var(
	pathPtr		*string
	path		string

)

func main(){
	serveMode := flag.Bool("serve", false, "a bool")
	pathPtr := flag.String("path", "/var/www/html", "a string")

	flag.Parse()
	f, err := os.OpenFile("log.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)

	if *serveMode {
		path = *pathPtr
		if path[len(path) - 1] != '/'{
			path += "/"
		}

		http.HandleFunc("/", handler)
		http.ListenAndServe(":54321", nil)
	}
}

func handler(w http.ResponseWriter, r *http.Request){
	if r.Method != "GET"{
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		f, err := os.OpenFile(fmt.Sprintf("%s%s", path, handler.Filename), os.O_WRONLY | os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		w.Write([]byte(fmt.Sprintf("ok /%s", handler.Filename)))
		log.Println("/", handler.Filename)
	}

	w.Write([]byte("not use in method get"))
}