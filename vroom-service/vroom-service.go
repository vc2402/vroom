package main

import (
	"bytes"
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
	traceIncomig bool
	traceOutgoig bool
}

func main() {
	addr := flag.String("address", "localhost", "address to listen on")
	port := flag.Int("port", 8088, "port to listen on")
	executable := flag.String("vroom", "./vroom", "path to vroom executable")
	verbose := flag.Bool("verbose", false, "log some messages")
	verbosest := flag.Bool("verbosest", false, "log impossibly much messages")
	traceIncomig := flag.Bool("traceIncoming", false, "trace incoming messages")
	traceOutgoig := flag.Bool("traceOutgoig", false, "trace outgoing messages")
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
		traceIncomig: *traceIncomig,
		traceOutgoig: *traceOutgoig,
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
	if h.traceIncomig {
		if h.logLevel == LevelVerbose {
			log.Println("going to read from stdout")
		}
		var buf bytes.Buffer
		tee := io.TeeReader(r.Body, &buf)
		command.Stdin = &buf
		incoming, _ := io.ReadAll(tee)
		log.Print("incoming message: \n\t", string(incoming))
	} else {
		command.Stdin = r.Body
	}
	var stdOut, stdErr bytes.Buffer
	command.Stdout = &stdOut
	command.Stderr = &stdErr
	err := command.Run()
	errors := stdErr.Bytes()
	var count int64
	if err != nil || len(errors) > 0 {
		if err != nil {
			log.Printf("error: when getting error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		} else {
			log.Printf("error: errors got from vroom executable: %s", string(errors))
			http.Error(w, string(errors), http.StatusInternalServerError)
		}
	} else {
		if h.traceOutgoig {
			if h.logLevel == LevelVerbose {
				log.Println("going to read from stdout")
			}
			var buf bytes.Buffer
			count, err = io.Copy(w, io.TeeReader(&stdOut, &buf))
			if err != nil {
				log.Printf("error: when copying output: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			outgoing, _ := io.ReadAll(&buf)
			log.Print("outgoing message: \n\t", string(outgoing))
		} else {
			if h.logLevel == LevelVerbose {
				log.Println("going to read from stdout")
			}
			count, err = io.Copy(w, &stdOut)
			if h.logLevel == LevelVerbose {
				log.Printf("%d bytes were read from stdout", count)
			}
		}

		if err != nil {
			log.Printf("error: when getting result from vrrom binary: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
	//_ = command.Wait()
	if err != nil || len(errors) > 0 {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if h.logLevel >= LevelInfo {
		log.Printf(
			"request was processed: %s %d/%d (%d ms)",
			r.RemoteAddr,
			r.ContentLength,
			count,
			time.Since(start)/time.Millisecond,
		)
	}
}
