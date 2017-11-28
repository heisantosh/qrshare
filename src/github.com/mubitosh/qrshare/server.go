package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"

	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type fileInfo struct {
	Name       string
	Icon  string
	ChildDirs  []fileInfo
	ChildFiles []fileInfo
}

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

// fileServer serves a file on a random port number. It shuts down if there
// is no download from the server within a period of time.
type fileServer struct {
	http.Server
	port     int
	listener net.Listener
}

var rootSelectedFiles map[string]bool

func fileServerNew() (*fileServer, error) {
	fs := &fileServer{}
	fs.Server.Addr = ":"
	listener, err := net.Listen("tcp", fs.Server.Addr)
	fs.listener = listener
	if err != nil {
		return nil, err
	}
	fs.port = fs.listener.Addr().(*net.TCPAddr).Port
	return fs, nil
}

func (fs *fileServer) start(app *QrShare, qrWindow *gtk.ApplicationWindow) error {
	rootSelectedFiles := make(map[string]bool)
	for _, s := range app.files {
		rootSelectedFiles[s] = true
	}

	serving, justServed := new(srvFlag), new(srvFlag)
	serving.set(false)
	justServed.set(false)

	fi, err := os.Stat(app.files[0])
	if err != nil {
		log.Println("Error starting server:", err)
		return err
	}

	baseDir := app.files[0]
	if !fi.IsDir() {
		baseDir = path.Dir(app.files[0])
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serving.set(true)

		name := path.Join(baseDir, path.Clean(r.URL.Path))

		// TODO: If URL path contains a filename with parent directory as baseDir
		// and the filename is not in the files list, return http.Error

		f, err := os.Open(name)
		if err != nil {
			log.Println("Error opening requested file:", err)
			http.Error(w, "StatusNotFound", 404)
			return
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			log.Println(err)
			http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
			return
		}

		// A directory
		if stat.IsDir() {
			serveDir(w, r, f)
		}

		// A socket
		if (stat.Mode() &^ 07777) == os.ModeSocket {
			log.Println("Not allowed to serve socket")
			http.Error(w, "StatusForbidden", http.StatusForbidden)
			return
		}

		// A file
		serveFile(w, r, f.Name())

		serving.set(false)
		justServed.set(true)
	})

	fs.Server.Handler = mux

	// Stop sharing when no activity is there.
	go func() {
		for {
			justServed.set(false)
			time.Sleep(time.Duration(*app.inActive) * time.Second)
			if !serving.get() && !justServed.get() {
				log.Println("Exceeded inactive time of", *app.inActive, "seconds")
				log.Println("Stopping file sharing")
				if app.isContractor {
					log.Println("App was started to display QR window only, exiting app")
					os.Exit(0)
				}
				log.Println("App was started with main window, back to main window")
				glib.IdleAdd(qrWindow.Destroy)
				return
			}
		}
	}()

	log.Println("Starting file sharing")
	return fs.Serve(tcpKeepAliveListener{fs.listener.(*net.TCPListener)})
}

// serveDir serves a directory content.
func serveDir(w http.ResponseWriter, r *http.Request, f *os.File) {
	fStats, err := f.Readdir(-1)
	if err != nil {
		log.Println("Error reading directory:", err)
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}

	fis := fileInfo{
		ChildFiles: []fileInfo{},
		ChildDirs: []fileInfo{},
	}

	// TODO: Sort the filenames
	
	for _, fStat := range fStats {
		if fStat.Name()[0] == '.' {
			continue
		}

		// If path is root, filter files not in app.files.
		if _, ok := rootSelectedFiles[f.Name()]; !ok && r.URL.Path == "/" {
			continue
		} else {
			if fStat.IsDir() {
				child := fileInfo {
					Name: f.Name(),
					Icon: "",
				}
				fis.ChildDirs = append(fis.ChildDirs, child)
			} else {
				child := fileInfo {
					Name: f.Name(),
					Icon: "",
				}
				fis.ChildFiles = append(fis.ChildFiles, child)
			}
		}
	}

	tpl, err := template.New("t").Parse(listingHTML)
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, fis)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}
}

// serveFile serves a file using standard library http.ServeFile.
func serveFile(w http.ResponseWriter, r *http.Request, name string) {
	http.ServeFile(w, r, name)
}

var listingHTML = `<html>
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style type="text/css">
      a:hover, a:visited, a:link, a:active {
    	text-decoration: none!important;
    	-webkit-box-shadow: none!important;
    	box-shadow: none!important;
	  }
      .file {
        max-width: 90px;
        word-wrap: break-word;
        display: inline-block;
        margin: 10px;
        vertical-align: top;
      }
      .icon {
        display: flex;
        justify-content: center;
        margin-bottom: 5px;
      }
      .icon-image {
        max-width: 100%;
      }
      .file-name {
        text-align: center;
        color: black;
      }
    </style>
  </head>
  <body>
    <div>

      {{range .ChildDirs}}
      	<div class="file">
        <a class="file-url" href="{{.Name}}">
         	<div class="icon">
            	<img class="icon-image" src="data:image/svg+xml;base64,{{.Icon}}">
          	</div>
          	<div class="file-name">{{.Name}}</div>
        </a>
      	</div>
      {{end}}

      {{range .ChildFiles}}
      	<div class="file">
        <a class="file-url" href="{{.Name}}">
         	<div class="icon">
            	<img class="icon-image" src="data:image/svg+xml;base64,{{.Icon}}">
          	</div>
          	<div class="file-name">{{.Name}}</div>
        </a>
      	</div>
      {{end}}

    </div>
  </body>
</html>`
