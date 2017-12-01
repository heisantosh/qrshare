package main

import (
	"github.com/gotk3/gotk3/gtk"

	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"time"
)

type fileInfo struct {
	Name       string
	Icon       string
	ChildDirs  []fileInfo
	ChildFiles []fileInfo
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

// Map of selected files provided as command line args to the application.
var rootSelectedFiles map[string]bool

var sharedPath = "/"

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
	rootSelectedFiles = make(map[string]bool)
	for _, s := range app.files {
		rootSelectedFiles[path.Base(s)] = true
	}

	_, err := os.Stat(app.files[0])
	if err != nil {
		log.Println("Error starting server:", err)
		return err
	}

	ap := getAbsPath(app.files)

	mux := http.NewServeMux()

	// Serve shared files under path sharedPath
	mux.Handle(sharedPath, http.StripPrefix(sharedPath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if path.Clean(r.URL.Path) == "/" {
				// Because there might be only one file to be served at the root.
				// Joining a URL path / at the end of a filename will make the
				// filepath a directory name.
				serve(w, r, ap)
			} else {
				serve(w, r, path.Join(ap, path.Clean(r.URL.Path)))
			}
		})))

	fs.Server.Handler = mux

	log.Println("Starting file sharing")
	return fs.Serve(tcpKeepAliveListener{fs.listener.(*net.TCPListener)})
}

// getParentDir returns the parent directory of the given file.
func getAbsPath(names []string) string {
	p := names[0]

	// Get the absolute path name for relative filename.
	if !path.IsAbs(p) {
		c, _ := os.Getwd()
		p = path.Join(c, p)
	}

	if len(names) == 1 {
		return p
	}

	return path.Dir(p)
}

// serve serves a file or directory request with the given file path.
func serve(w http.ResponseWriter, r *http.Request, filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Println("Error opening file:", err)
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}

	defer f.Close()

	fStat, err := f.Stat()
	if err != nil {
		log.Println("Error getting file info:", err)
		http.Error(w, "StatusInternalServerError", http.StatusInternalServerError)
		return
	}

	if fStat.IsDir() {
		serveDir(w, r, f)
		return
	}

	if (fStat.Mode() &^ 07777) == os.ModeSocket {
		log.Println("file is a socket: not serving it")
		http.Error(w, "StatusForbidden", http.StatusForbidden)
		return
	}

	serveFile(w, r, filePath)
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
		Name:       "",
		ChildFiles: []fileInfo{},
		ChildDirs:  []fileInfo{},
	}

	// If not at root of sharing, url path should equal to relative file path.
	if r.URL.Path != "/" {
		fis.Name = r.URL.Path
	}

	// TODO: Sort the filenames
	// I think the slice return from os.Readdir is sorted already.

	for _, fStat := range fStats {
		if fStat.Name()[0] == '.' {
			continue
		}

		// If path is root, filter files not in app.files.
		if _, ok := rootSelectedFiles[fStat.Name()]; !ok &&
			// http.StripPrefix removes / also. Need to check for that.
			(r.URL.Path == "/" || r.URL.Path == "") &&
			// It's okay if there is only one directory is to be served.
			len(rootSelectedFiles) != 1 {
			continue
		}

		if fStat.IsDir() {
			fis.ChildDirs = append(fis.ChildDirs, fileInfo{Name: fStat.Name(), Icon: iconFolder})
		} else {
			fis.ChildFiles = append(fis.ChildFiles, fileInfo{Name: fStat.Name(), Icon: iconText})
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
func serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	http.ServeFile(w, r, filePath)
}

var listingHTML = `<html>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style type="text/css">
        a:hover,
        a:visited,
        a:link,
        a:active {
            text-decoration: none!important;
            -webkit-box-shadow: none!important;
            box-shadow: none!important;
        }
        
        .file {
            width: 80px;
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
        	font-family: "Arial";
        	font-size: small;
        	color: rgb(80, 80, 80)
        }
    </style>
</head>

<body>
    <div>

        {{$name := .Name}} 

        {{range .ChildDirs}}
        <div class="file">
            <a class="file-url" href="{{$name}}/{{.Name}}">
                <div class="icon">
                    <img class="icon-image" src="data:image/svg+xml;base64,{{.Icon}}">
                </div>
                <div class="file-name">{{.Name}}</div>
            </a>
        </div>
        {{end}} 

        {{range .ChildFiles}}
        <div class="file">
            <a class="file-url" href="{{$name}}/{{.Name}}">
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
