package main

import (
	"net/http"
	"fmt"
	"log"
	"context"
	"io"
	"time"
)

type flushWriter struct {
	f http.Flusher
	w io.Writer
}

func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if fw.f != nil {
		fw.f.Flush()
	}
	return
}

func main() {
	srv := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}

	closeChan := make(chan struct{})
	go func(closeChan chan struct{}) {
		fmt.Println("Press enter to shutdown server")
		fmt.Scanln()
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("could not shutdown: %v", err)
		}
		close(closeChan)
	}(closeChan)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fw := flushWriter{w: w}
		if f, ok := w.(http.Flusher); ok {
			fw.f = f
		}

		for i := 0; i < 100; i++ {
			fmt.Fprintln(&fw, "Happy Go 1.8'th")
			time.Sleep(time.Millisecond * 200)
		}
	})
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
	<-closeChan
}
