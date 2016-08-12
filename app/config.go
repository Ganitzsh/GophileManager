package app

import (
	"log"

	"github.com/revel/revel"
)

type AppConfig struct {
	MainDir  string `json:"main_dir"`
	TrashDir string `json:"trash_dir"`
	Trash    bool   `json"enable_trash"`
}

func InitAppConfig() {
	d, set := revel.Config.String("app.main_dir")
	if !set {
		log.Fatal("No direcory set")
	}
	t, trash := revel.Config.String("app.trash_dir")
	if !set {
		log.Println("WARNING: Trash disabled")
	}
	conf = AppConfig{d, t, trash}
}
