package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type LogLevel int

const (
	LevelError LogLevel = iota
	LevelInfo
	LevelVerbose
)

type vroomHandler struct {
	address      string
	port         int
	pathToVRoom  string
	vroomOptions []string
	logLevel     LogLevel
}

func main() {
	addr := flag.String("address", "localhost", "address to listen on")
	port := flag.Int("port", 8088, "port to listen on")
	executable := flag.String("vroom", "./vroom", "path to vroom executable")
	verbose := flag.Bool("verbose", false, "log some messages")
	verbosest := flag.Bool("verbosest", false, "log impossibly much messages")
	flag.Usage = func() {
		fmt.Printf("usage: %s [flags] [ -- options-to-vroom]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	h := vroomHandler{
		address:      *addr,
		port:         *port,
		pathToVRoom:  *executable,
		vroomOptions: flag.Args(),
	}
	if *verbose {
		h.logLevel = LevelInfo
	}
	if *verbosest {
		h.logLevel = LevelVerbose
	}

	h.start()
}

func (h vroomHandler) start() {
	//http.Handle("/", h)

	log.Fatal(http.ListenAndServe(h.address+":"+strconv.Itoa(h.port), h))
}

func (h vroomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if h.logLevel == LevelVerbose {
		log.Printf("new request: %s/%s %d", r.RemoteAddr, r.Method, r.ContentLength)
	}
	if r.Method != http.MethodPost {
		log.Print("error: POST method is awaiting")
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}
	command := exec.Command(h.pathToVRoom, h.vroomOptions...)
	stdIn, err := command.StdinPipe()
	if err != nil {
		log.Printf("error: can not get stdin: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer stdIn.Close()
	stdOut, err := command.StdoutPipe()
	if err != nil {
		log.Printf("error: can not get stdout: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer stdOut.Close()
	stdErr, err := command.StderrPipe()
	if err != nil {
		log.Printf("error: can not get stdErr: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer stdErr.Close()
	err = command.Start()
	if err != nil {
		log.Printf("error: when starting: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	count, err := io.Copy(stdIn, r.Body)
	if err != nil {
		log.Printf("error: when copying to stdin: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if h.logLevel == LevelVerbose {
		log.Printf("%d bytes request was passed to vroom", r.ContentLength)
	}
	err = command.Wait()
	if err != nil {
		log.Printf("error: when waiting: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	errors, err := io.ReadAll(stdErr)
	if err != nil {
		log.Printf("error: when getting error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if len(errors) > 0 {
		log.Printf("error: errors got from vroom executable: %s", string(errors))
		http.Error(w, string(errors), http.StatusInternalServerError)
		return
	}
	count, err = io.Copy(w, stdOut)
	if err != nil {
		log.Printf("error: when getting result from vrrom binary: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if h.logLevel >= LevelInfo {
		log.Printf("request was processed: %s %d/%d (%d ms)", r.RemoteAddr, r.ContentLength, count, time.Since(start)/time.Millisecond)
	}
	w.Header().Set("Content-Type", "application/json")
}
