package app

import (
	"log"
	"strings"

	"github.com/revel/revel"
)

type AppConfig struct {
	MainDir  string `json:"main_dir"`
	TrashDir string `json:"trash_dir"`
	Trash    bool   `json"enable_trash"`
	Host     string `json:"host"`
}

func InitAppConfig() {
	d, set := revel.Config.String("app.main_dir")
	if !set {
		log.Fatal("No direcory set")
	}
	h, set := revel.Config.String("app.host")
	if !set {
		log.Fatal("No host set")
	}
	if !strings.HasPrefix(h, "http://") || !strings.HasPrefix(h, "https://") {
		log.Fatal("Invalid host format")
	}
	t, trash := revel.Config.String("app.trash_dir")
	if !set {
		log.Println("WARNING: Trash disabled")
	}
	conf = AppConfig{d, t, trash, h}
}
