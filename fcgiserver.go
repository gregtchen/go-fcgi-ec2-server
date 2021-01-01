//source: https://www.dav-muz.net/blog/2013/09/how-to-use-go-and-fastcgi/

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"runtime"
)

var (
	local = flag.String("local", "", "serve as webserver, example: 0.0.0.0:8000")
	tcp   = flag.String("tcp", "", "serve as FCGI via TCP, example: 0.0.0.0:8000")
	unix  = flag.String("unix", "", "serve as FCGI via UNIX socket, example: /tmp/myprogram.sock")
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func homeView(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	headers.Add("Content-Type", "text/html")
	io.WriteString(w, "<html><head></head><body><p>It works!</p></body></html>")
}

func app2(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "App 2 served")
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("cwd = %s\n", cwd)

	staticHandler := http.NewServeMux()
	staticHandler.Handle("/", http.FileServer(http.Dir(cwd)))
	staticHandler.HandleFunc("/app2", app2)

	flag.Parse()

	fmt.Printf("you have made it this far\n")
	if *local != "" { // Run as a local web server
		err = http.ListenAndServe(*local, staticHandler)
	} else if *tcp != "" { // Run as FCGI via TCP
		listener, err := net.Listen("tcp", *tcp)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		err = fcgi.Serve(listener, staticHandler)
	} else if *unix != "" { // Run as FCGI via UNIX socket
		listener, err := net.Listen("unix", *unix)
		if err != nil {
			log.Fatal(err)
		}
		defer listener.Close()

		err = fcgi.Serve(listener, staticHandler)
	} else { // Run as FCGI via standard I/O
		err = fcgi.Serve(nil, staticHandler)
	}
	if err != nil {
		log.Fatal(err)
	}
}
