package app

import (
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/googollee/go-socket.io"

	"gopkg.in/fsnotify.v1"
)

type Category string

type WebManager struct {
	Config     *AppConfig
	LoggedIn   bool
	PWD        string
	Categories map[Category]interface{}
	Watcher    *fsnotify.Watcher
	Server     *socketio.Server
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

func NewWebManager() *WebManager {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("notif")
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})
	return &WebManager{
		Config:     &conf,
		LoggedIn:   true,
		Categories: make(map[Category]interface{}),
		PWD:        conf.MainDir,
		Server:     server,
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
}
