package app

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/googollee/go-socket.io"
	"github.com/revel/revel"

	"gopkg.in/fsnotify.v1"
)

type Category string

type WebManager struct {
	Config      *AppConfig
	LoggedIn    bool
	PWD         string
	Categories  map[Category]interface{}
	Watcher     *fsnotify.Watcher
	SocketIO    *socketio.Server
	RevelHandle http.Handler
	Trash       bool
	CanConvert  bool
}

func handle(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/socket.io/") {
		Context.SocketIO.ServeHTTP(w, r)
	} else {
		Context.RevelHandle.ServeHTTP(w, r)
	}
}

func (wm WebManager) GetMainCategories() []string {
	var c []string

	for cat, _ := range wm.Categories {
		c = append(c, strings.Title(string(cat)))
	}
	sort.Strings(c)
	return c
}

func (wm WebManager) GetSubCategories(main string) []string {
	var c []string

	ptr := wm.Categories[Category(main)]
	if ptr == nil {
		return nil
	}
	t := reflect.TypeOf(ptr)
	if t.Kind() != reflect.Map {
		return nil
	}
	tmp := ptr.(map[Category]interface{})
	for cat, _ := range tmp {
		c = append(c, strings.Title(string(cat)))
	}
	return c
}

func watchDirActivity(watcher *fsnotify.Watcher) {
	for {
		select {
		case ev := <-watcher.Events:
			if ev.Op != fsnotify.Write {
				log.Println("event:", ev)
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func binExist(bin string) bool {
	sysPath := os.Getenv("PATH")
	log.Println(sysPath)
	paths := strings.Split(sysPath, ":")
	for _, path := range paths {
		f, err := os.Open(path + "/" + bin)
		if err == nil {
			return true
		}
		f.Close()
	}
	return false
}

func NewWebManager() *WebManager {
	ffmpegInstalled := binExist("ffmpeg")
	if ffmpegInstalled {
		log.Println("ffmpeg is installed")
	} else {
		log.Println("ffmpeg is not installed")
	}
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("Client connected")
		so.Join("notif")
		so.On("disconnection", func() {
			log.Println("Client disconnected")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	return &WebManager{
		Config:      &conf,
		LoggedIn:    true,
		Categories:  make(map[Category]interface{}),
		PWD:         conf.MainDir,
		SocketIO:    server,
		RevelHandle: revel.Server.Handler,
		Trash:       conf.Trash,
		CanConvert:  ffmpegInstalled,
	}
}

func NewWatcher(dir string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go watchDirActivity(watcher)
	err = Context.Watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func InitApp() {
	Context = NewWebManager()
	revel.Server.Handler = http.HandlerFunc(handle)
}
