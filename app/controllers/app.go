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
	fPath := c.Session["pwd"] + "/" + file
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
	return c.Render()
}

func (c App) Compress() revel.Result {
	reload := false
	oldPwd := c.Session["pwd"]
	name := c.Params.Get("name")
	file := c.Params.Get("target")
	target := c.Session["pwd"]
	_, err := app.CreateArchive(target+"/"+file, target, name)
	if err != nil {
		app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
			"message": "<strong>" + file + "</strong> commpression error:<br/>" + err.Error(),
			"alert":   c.Params.Get("alert_id"),
			"reload":  reload,
		})
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err.Error(),
		})
	}
	log.Println("After compression:", c.Session["pwd"])
	reload = c.Session["pwd"] == oldPwd
	app.Context.SocketIO.BroadcastTo("notif", "notif action done", map[string]interface{}{
		"message": "<strong>" + file + "</strong> commpressed successfully!",
		"alert":   c.Params.Get("alert_id"),
		"reload":  reload,
	})
	return c.RenderJson(map[string]interface{}{
		"message": "Success!",
	})
	// fName := name + ".tar"
	// _, err = os.Open(target + "/" + fName)
	// if err != nil {
	// 	c.RenderJson(map[string]interface{}{
	// 		"message": err,
	// 	})
	// }
	// return c.RenderFile(f, revel.Attachment)
}

func (c App) Delete() revel.Result {
	reload := false
	oldPwd := c.Session["pwd"]
	file := c.Params.Get("target")
	fPath := c.Session["pwd"] + "/" + file
	if err := os.Remove(fPath); err != nil {
		app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
			"message": "<strong>" + file + "</strong> could not be deleted:<br?/>" + err.Error(),
			"alert":   c.Params.Get("alert_id"),
			"reload":  reload,
		})
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}
	reload = (c.Session["pwd"] == oldPwd)
	time.Sleep(1 * time.Second)
	app.Context.SocketIO.BroadcastTo("notif", "notif action done", map[string]interface{}{
		"message": "<strong>" + file + "</strong> deleted successfully!",
		"alert":   c.Params.Get("alert_id"),
		"reload":  reload,
	})
	return c.RenderJson(map[string]interface{}{
		"message": "Deleted successfuly",
		"status":  http.StatusOK,
	})
}

func (c App) GetFiles() revel.Result {
	tmp := c.Params.Get("dir")
	path := c.Session["pwd"]
	if path == "" {
		path = app.Context.Config.MainDir
	}
	if tmp != "" {
		if tmp == "up" && c.Session["pwd"] != app.Context.Config.MainDir {
			path = filepath.Dir(c.Session["pwd"])
		} else if tmp != "up" && tmp != "current" {
			path += "/" + tmp
		}
	}
	content, err := app.ProcessDir(path)
	if err != nil {
		return c.RenderJson(map[string]interface{}{
			"message": "Directory does not exist: " + strings.TrimPrefix(path, app.Context.Config.MainDir),
			"status":  http.StatusBadRequest,
		})
	}
	c.Session["pwd"] = path
	c.RenderArgs["isRoot"] = (path == app.Context.Config.MainDir)
	c.RenderArgs["content"] = content
	log.Println("New PWD:", c.Session["pwd"])
	return c.RenderTemplate("App/files.html")
}

func (c App) Download() revel.Result {
	file := c.Params.Get("target")
	fPath := c.Session["pwd"] + "/" + file
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

func toMP4(pwd, target string) error {
	ext := filepath.Ext(target)
	base := pwd + "/" + strings.TrimSuffix(target, ext) + ".mp4"
	fPath := pwd + "/" + target
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
	if err := toMP4(c.Session["pwd"], file); err != nil {
		log.Println(err)
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
