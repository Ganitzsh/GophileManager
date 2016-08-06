package app

import (
	"log"

	"github.com/revel/revel"
)

type AppConfig struct {
	MainDir string `json:"directory"`
}

func InitAppConfig() {
	d, set := revel.Config.String("app.main_dir")
	if !set {
		log.Fatal("No direcory set")
	}
	conf = AppConfig{d}
}
