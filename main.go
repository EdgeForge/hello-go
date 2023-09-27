package main

import (
	"embed"
	"flag"
	"fmt"
	"html"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	//go:embed static/*
	content embed.FS
)

func main() {
	var maintenance, stdout bool
	port := "8080"

	flag.BoolVar(&maintenance, "maintenance", false, "show maintenance page")
	flag.BoolVar(&stdout, "stdout", false, "show page its on stdout")
	flag.StringVar(&port, "port", port, "port to listen on")

	flag.Parse()

	if ok, _ := strconv.ParseBool(os.Getenv("MAINTENANCE")); ok {
		maintenance = true
	}

	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "0.0.0.0"
	}

	if isout := os.Getenv("STDOUT"); isout != "" {
		if ok, err := strconv.ParseBool(isout); err == nil {
			stdout = ok
		} else {
			log.Printf("STDOUT env was set to %q but not parsable as bool: %v", isout, err)
		}
	}
	if maintenance {
		log.Println("MAINTAIN!")
		dir, err := fs.Sub(content, "static")
		if err != nil {
			log.Fatalf("failed to get static dir: %v", err)
		}
		http.Handle("/", http.FileServer(http.FS(dir)))
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "[%-6s] %q\n", r.Method, html.EscapeString(r.URL.Path))
			if stdout {
				fmt.Fprintf(os.Stdout, "[%-6s] %q\n", r.Method, html.EscapeString(r.URL.Path))
			}
		})
	}

	listen := addr + ":" + port
	log.Println("listening on: ", listen)

	log.Fatal(http.ListenAndServe(listen, nil))
}
