package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type srvFlag struct {
	value bool
	mutex sync.Mutex
}

func (f *srvFlag) set(v bool) {
	f.mutex.Lock()
	f.value = v
	f.mutex.Unlock()
}

func (f *srvFlag) get() bool {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.value
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

type FileServer struct {
	http.Server
	port     int
	listener net.Listener
}

func FileServerNew() (*FileServer, error) {
	fs := &FileServer{}
	fs.Server.Addr = ":"
	listener, err := net.Listen("tcp", fs.Server.Addr)
	fs.listener = listener
	if err != nil {
		return nil, err
	}
	fs.port = fs.listener.Addr().(*net.TCPAddr).Port
	return fs, nil
}

func (fs *FileServer) start(app *App, window *gtk.ApplicationWindow) error {
	serving, justServed := new(srvFlag), new(srvFlag)
	serving.set(false)
	justServed.set(false)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serving.set(true)
		log.Println("Serving file:", *app.file)
		http.ServeFile(w, r, *app.file)
		log.Println("File served")
		serving.set(false)
		justServed.set(true)
	})

	fs.Server.Handler = mux

	go func() {
		for {
			justServed.set(false)
			time.Sleep(time.Duration(*app.inactive) * time.Second)
			if !serving.get() && !justServed.get() {
				log.Println("Exceeded inactive time of", *app.inactive, "seconds")
				log.Println("Stopping file sharing")
				if app.isCmdLine {
					log.Println("App was started to display QR window only, exiting app")
					os.Exit(0)
				}
				log.Println("App was started with main window, back to main window")
				glib.IdleAdd(window.Destroy)
				return
			}
		}
	}()

	log.Println("Starting file sharing")
	return fs.Serve(tcpKeepAliveListener{fs.listener.(*net.TCPListener)})
}
