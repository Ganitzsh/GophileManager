package controllers

import (
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ganitzsh/WebManager/app"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Check() revel.Result {
	if !app.Context.LoggedIn {
		log.Println("Not logged in")
		return c.Redirect("/auth")
	}
	return c.Redirect("/app")
}

func (c App) Serve(prefix, filepath string) revel.Result {
	file := c.Params.Get("target")
	fPath := app.Context.Config.MainDir + "/" + file
	f, err := os.Open(fPath)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}
	return c.RenderFile(f, revel.Attachment)
}

func (c App) Index() revel.Result {
	tmp := c.Params.Get("folder")
	log.Println("Folder:", tmp)
	c.RenderArgs["test"] = app.Context.GetMainCategories()
	c.RenderArgs["content"] = app.Context.Categories
	return c.Render()
}

func (c App) Compress() revel.Result {
	name := c.Params.Get("name")
	file := c.Params.Get("target")
	target := app.Context.Config.MainDir
	log.Println(target + "/" + file)
	fName := name + ".tar"
	log.Println(fName)
	_, err := app.CreateArchive(target+"/"+file, target, name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(target + "/" + fName)
	f, err := os.Open(target + "/" + fName)
	if err != nil {
		log.Fatal(err)
	}
	return c.RenderFile(f, revel.Attachment)
}

func (c App) Delete() revel.Result {
	file := c.Params.Get("target")
	fPath := app.Context.Config.MainDir + "/" + file
	if err := os.Remove(fPath); err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}
	time.Sleep(1 * time.Second)
	return c.RenderJson(map[string]interface{}{
		"message": "Deleted successfuly",
		"status":  http.StatusOK,
	})
}

func (c App) GetFiles() revel.Result {
	return c.RenderJson(app.Context.Categories)
}

func (c App) Download() revel.Result {
	file := c.Params.Get("target")
	fPath := app.Context.Config.MainDir + "/" + file
	f, err := os.Open(fPath)
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}
	return c.RenderFile(f, revel.Inline)
}

func toMP4(target string) error {
	ext := filepath.Ext(target)
	base := app.Context.Config.MainDir + "/" + strings.TrimSuffix(target, ext) + ".mp4"
	fPath := app.Context.Config.MainDir + "/" + target
	cmd := exec.Command("ffmpeg", "-i", fPath, "-vcodec", "copy", "-acodec", "copy", base)
	err := cmd.Start()
	if err != nil {
		return err
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (c App) Convert() revel.Result {
	file := c.Params.Get("target")
	if err := toMP4(file); err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}
	return c.RenderJson(map[string]interface{}{
		"message": "Successfully converted",
	})
}

func (c App) Video() revel.Result {
	file := c.Params.Get("target")
	ext := filepath.Ext(file)
	mime := mime.TypeByExtension(ext)
	c.RenderArgs["video"] = file
	c.RenderArgs["mime"] = mime
	return c.Render()
}
