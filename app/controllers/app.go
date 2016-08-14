package controllers

import (
	"log"
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

func (c App) Serve() revel.Result {
	file := app.Context.Config.MainDir + "/" + c.Params.Get("filepath")
	log.Println("Serving:", file)
	f, err := os.Open(file)
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
	pchan := make(chan uint64) // Follow progress
	oldPwd := c.Session["pwd"]
	name := c.Params.Get("name")
	file := c.Params.Get("target")
	target := c.Session["pwd"]

	go func() {
		for {
			select {
			case progress := <-pchan:
				if progress == 100 {
					return
				}
				app.Context.SocketIO.BroadcastTo("notif", "notif action progress", map[string]interface{}{
					"message":  "Compressing <strong>" + file + "</strong>...",
					"progress": progress,
					"alert":    c.Params.Get("alert_id"),
					"reload":   reload,
				})
			}
		}
	}()
	_, err := app.CreateArchive(target+"/"+file, target, name, pchan)
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
	close(pchan)
	return c.RenderJson(map[string]interface{}{
		"message": "Success!",
	})
}

func (c App) EmptyTrash() revel.Result {
	if !app.Context.Trash {
		app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
			"message": "There is no trash set on this server",
			"alert":   c.Params.Get("alert_id"),
			"reload":  false,
		})
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": "No trash set",
			"status":  http.StatusBadRequest,
		})
	}
	if err := os.RemoveAll(app.Context.Config.TrashDir); err != nil {
		app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
			"message": "Could not empty the trash: " + err.Error(),
			"alert":   c.Params.Get("alert_id"),
			"reload":  false,
		})
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err,
			"status":  http.StatusBadRequest,
		})
	}
	if err := os.MkdirAll(app.Context.Config.TrashDir, os.ModePerm); err != nil {
		app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
			"message": "Could not create the trash: " + err.Error(),
			"alert":   c.Params.Get("alert_id"),
			"reload":  false,
		})
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err,
			"status":  http.StatusBadRequest,
		})
	}
	app.Context.SocketIO.BroadcastTo("notif", "notif action done", map[string]interface{}{
		"message": "Trash emptied successfully!",
		"alert":   c.Params.Get("alert_id"),
		"reload":  true,
	})
	return c.RenderJson(map[string]interface{}{
		"message": "Trash emptied",
		"status":  http.StatusOK,
	})
}

func (c App) MoveToTrash() revel.Result {
	file := c.Params.Get("target")
	fPath := c.Session["pwd"] + "/" + file
	newPath := app.Context.Config.TrashDir + "/" + file
	if err := os.Rename(fPath, newPath); err != nil {
		app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
			"message": "<strong>" + file + "</strong> could not be moved to trash:<br?/>" + err.Error(),
			"alert":   c.Params.Get("alert_id"),
			"reload":  false,
		})
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}
	time.Sleep(1 * time.Second)
	app.Context.SocketIO.BroadcastTo("notif", "notif action done", map[string]interface{}{
		"message": "<strong>" + file + "</strong> moved to trash!",
		"alert":   c.Params.Get("alert_id"),
		"reload":  true,
	})
	return c.RenderJson(map[string]interface{}{
		"message": "Deleted successfuly",
		"status":  http.StatusOK,
	})
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
	var trashCount uint64
	tmp := c.Params.Get("dir")
	path := c.Session["pwd"]

	c.RenderArgs["noTrash"] = !app.Context.Trash
	c.RenderArgs["isTrash"] = false
	if path == "" {
		path = app.Context.Config.MainDir
	}
	switch tmp {
	case "trash":
		if !app.Context.Trash {
			app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
				"message": "Trash is disabled",
				"alert":   c.Params.Get("alert_id"),
				"reload":  false,
			})
			path = app.Context.Config.MainDir
		} else {
			path = app.Context.Config.TrashDir
		}
	case "home":
		path = app.Context.Config.MainDir
	case "current":
		path = path
	case "up":
		isAllowed := strings.HasPrefix(path, app.Context.Config.MainDir)
		if path != app.Context.Config.MainDir && isAllowed {
			path = filepath.Dir(path)
		} else {
			path = app.Context.Config.MainDir
		}
	default:
		path += "/" + tmp
	}
	c.Session["pwd"] = path
	content, err := app.ProcessDir(path)
	if err != nil {
		app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
			"message": "Could not get files: " + err.Error(),
			"alert":   c.Params.Get("alert_id"),
			"reload":  false,
		})
		return c.RenderJson(map[string]interface{}{
			"message": "Directory does not exist: " + strings.TrimPrefix(path, app.Context.Config.MainDir),
			"status":  http.StatusBadRequest,
		})
	}
	if app.Context.Trash {
		trashCount, err = app.CountFilesInDir(app.Context.Config.TrashDir)
		if err != nil {
			app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
				"message": "Could not access trash: " + err.Error(),
				"alert":   c.Params.Get("alert_id"),
				"reload":  false,
			})
		}
	}
	c.RenderArgs["isTrash"] = (path == app.Context.Config.TrashDir)
	c.RenderArgs["trashCount"] = trashCount
	c.RenderArgs["empty"] = (len(content) == 0)
	c.RenderArgs["isRoot"] = (path == app.Context.Config.MainDir)
	c.RenderArgs["content"] = content
	c.RenderArgs["relPwd"] = strings.TrimLeft(path, app.Context.Config.MainDir+"/")
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
	if !app.Context.CanConvert {
		app.Context.SocketIO.BroadcastTo("notif", "notif action error", map[string]interface{}{
			"message": "<strong>Cannot convert this file</strong>ffmpeg not installed on the server",
			"alert":   c.Params.Get("alert_id"),
			"reload":  false,
		})
		return c.RenderJson(map[string]interface{}{
			"message": "Cannot convert this file",
			"status":  http.StatusBadRequest,
		})
	}
	file := c.Params.Get("target")
	if err := toMP4(c.Session["pwd"], file); err != nil {
		log.Println(err)
		c.Response.Status = http.StatusBadRequest
		return c.RenderJson(map[string]interface{}{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}
	app.Context.SocketIO.BroadcastTo("notif", "notif action done", map[string]interface{}{
		"message": "<strong>" + file + "</strong> Converted successfully!",
		"alert":   c.Params.Get("alert_id"),
		"reload":  false,
	})
	return c.RenderJson(map[string]interface{}{
		"message": "Successfully converted",
	})
}

func (c App) Video() revel.Result {
	file := c.Params.Get("target")
	relPwd := strings.TrimPrefix(c.Session["pwd"], app.Context.Config.MainDir+"/")
	path := relPwd
	if relPwd != "" {
		path += "/"
	}
	log.Println("Watching:", file)
	c.RenderArgs["videoPath"] = path + file
	c.RenderArgs["videoName"] = file
	return c.Render()
}
