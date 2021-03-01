package routes

import (
	"fmt"
	"hd-virtual-plus-plus/filefinder"
	"html/template"
	"log"
	"os"
	fp "path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

//Index Load login page, with or without warning of wrong user or password
func Index(c *fiber.Ctx) error {
	if c.FormValue("auth") == "false" {
		log.Output(1, "USER OR PASSWORD WRONG")
	}
	return c.SendFile("./frontend/html/login.html", false)
}

//Files serve files to user based on url path
func Files(c *fiber.Ctx) error {
	filePath := c.Params("*")
	fileNames, err := filefinder.FindFiles("uploads/" + filePath)
	if err != nil {
		log.Fatalf("ERROR: FILE FINDER: %v\n", err)
		return c.SendFile("frontend/html/fileNotFound.html")
	}
	htmlStr := ""
	for _, fileName := range fileNames {

		//Default file is a folder
		fileLink := fp.Join("arquivos", filePath, fileName)
		download := ""
		fileType := "folder"

		//If it has an extension it is a file
		if strings.Contains(fileName, ".") {
			fileLink = fp.Join("download", filePath, fileName)
			download = "download='" + fileName + "'"
			fileType = "description"
		}

		//Transform files into html
		htmlStr = htmlStr + "<a href='/" + fileLink + "' " + download + " class='item'>" +
			"<span class='material-icons'>" +
			fileType +
			"</span>" +
			"<div class='name'>" +
			fileName +
			"</div>" +
			"</a>"

	}
	html := template.HTML(htmlStr)
	return c.Render("files", fiber.Map{
		"Files": html,
		"Path":  filePath,
	})
}

//AddFilesForm Add files or folders to a form and send
func AddFilesForm(c *fiber.Ctx) error {
	pathDir := c.Params("*")
	pathName := pathDir
	if len(pathDir) == 0 {
		pathDir = ""
		pathName = "pasta inicial"
	}
	return c.Render("addFileForm", fiber.Map{
		"PathName": pathName,
		"PathDir":  pathDir,
	})
}

//SaveFiles Save files from form and tag them with their id
func SaveFiles(c *fiber.Ctx) error {
	addpath := c.FormValue("addpath")
	addtype := c.FormValue("addtype")

	if addtype == "dir" {
		dirname := c.FormValue("dirname")

		if dirname == "" {
			c.Request().Header.Add("error", "missing-dirname")
			return c.Redirect("/add/" + addpath)
		}

		savepath := fp.Join("uploads", addpath, strings.ReplaceAll(dirname, " ", "_"))
		if err := os.Mkdir(savepath, 0755); err != nil {
			log.Fatalf("ERROR: SAVE DIR: %v\n", err)
			return c.Redirect("/add/" + addpath)
		}
		log.Output(1, fmt.Sprintf("%v %v %v", dirname, addtype, savepath))
	} else {
		filedata, err := c.FormFile("filedata")
		if err != nil {
			log.Fatalf("ERROR: FILE UPLOAD: %v\n", err)
			c.Request().Header.Add("error", "missing-file")
			return c.Redirect("/add/" + addpath)
		}

		savepath := fp.Join("uploads", addpath, strings.ReplaceAll(filedata.Filename, " ", "_"))
		if err = c.SaveFile(filedata, savepath); err != nil {
			log.Fatalf("ERROR: SAVE UPLOAD: %v\n", err)
			return c.Redirect("/add/" + addpath)
		}
		log.Output(1, fmt.Sprintf("%v %v %v", filedata.Filename, addtype, savepath))
	}

	return c.Redirect("/arquivos/" + addpath)
}

//Login login the user and give access to uploaded files
func Login(c *fiber.Ctx) error {
	username := c.FormValue("username")
	passwd := c.FormValue("password")

	log.Output(0, username)
	log.Output(0, passwd)
	// if username != "john" || passwd != "doe" {
	// 	return c.Redirect("../", fiber.StatusUnauthorized)
	// }

	return c.Redirect("/arquivos")
}
