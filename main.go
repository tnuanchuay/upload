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
	"strings"
	"runtime"
)

const(
	UPLOAD_HOST_ENV		=	"UPLOAD_HOST"
)

var(
	path		string
)

func main(){
	serveMode := flag.Bool("serve", false, "serve mode")
	pathPtr := flag.String("path", "/var/www/html", "specific path to save file")
	server := flag.String("host", "", "a server to upload")
  port := flag.String("port", "54321", "a TCP port for listening or establishing")

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

    listeningAddress := fmt.Sprintf(":%s", *port)

		http.HandleFunc("/", handler)
		http.ListenAndServe(listeningAddress, nil)
	}else{
		args := os.Args
<<<<<<< HEAD
=======
		if len(args) == 3{
			filename := args[2]
			host := fmt.Sprintf("http://%s:%s/", *server, *port)

			f, _ := filepath.Abs(filename)
			byteOfFile, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatal("file err : ", err)
			}
>>>>>>> 21d3a7f2683647a2e3dea0398a69b823543a516d

		filename := args[len(args) -1]
		if *server == "" {
			for _, e := range os.Environ(){
				pair := strings.Split(e, "=")
				if pair[0] == UPLOAD_HOST_ENV{
					server = &pair[1]
					break
				}

			}
		}

		if *server == ""{
			usageString := ""
			if runtime.GOOS == "windows"{
				usageString = fmt.Sprintf("SET %s=example.com", UPLOAD_HOST_ENV)
			}else{
				usageString = fmt.Sprintf("export %s=example.com", UPLOAD_HOST_ENV)
			}

			log.Fatal("need to define host to upload ", usageString, " or use upload --host=example.com filename")

		}

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
