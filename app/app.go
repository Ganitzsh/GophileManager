package app

import (
	"log"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/fsnotify.v1"
)

type Category string

type WebManager struct {
	Config     *AppConfig
	LoggedIn   bool
	Categories map[Category]interface{}
	Watcher    *fsnotify.Watcher
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

func watchDirActivity() {
	for {
		select {
		case ev := <-Context.Watcher.Events:
			if ev.Op != fsnotify.Write {
				log.Println("event:", ev)
				Context.Categories = make(map[Category]interface{})
				ProcessDir(Context.Config.MainDir)
			}
		case err := <-Context.Watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func NewWebManager() *WebManager {
	return &WebManager{
		Config:     &conf,
		LoggedIn:   true,
		Categories: make(map[Category]interface{}),
	}
}

func InitApp() {
	Context = NewWebManager()
	ProcessDir(Context.Config.MainDir)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	Context.Watcher = watcher
	go watchDirActivity()
	err = Context.Watcher.Add(Context.Config.MainDir)
	if err != nil {
		log.Fatal(err)
	}
}
