package handlers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"recreview/models"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

//type H map[string]interface{}

//GetAppDocs endpoint
func GetAppDocs(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		AppDocs := models.GetAppDocs(db, id)

		return c.JSON(http.StatusOK, AppDocs)
	}
}

func DeleteAppDoc(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return fmt.Errorf("ID Atoi: %v: ", err)
		}

		// Use our new model to delete a task
		_, err = models.DeleteAppDoc(db, id)
		if err != nil {
			return fmt.Errorf("Delete AppDoc: %v", err)
		}

		return c.JSON(http.StatusOK, H{
			"deleted": id,
		})
	}
}

// CreateAppDoc endpoint
func CreateAppDoc(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Instantiate a new task
		var AppDoc models.AppDoc
		// Map imcoming JSON body to the new Recruit

		file, err := c.FormFile("FileUpload")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		//prep for db
		AppDoc.Name = file.Filename
		AppDoc.RecruitID, _ = strconv.ParseInt(c.FormValue("RecruitID"), 10, 64)

		//create db row
		id, err := models.CreateAppDoc(db, AppDoc)
		if err != nil {
			return err
		}

		//thumbnail of pdf file
		b64data := strings.Split(c.FormValue("preview"), "base64,")[1]
		decodedString, err := base64.StdEncoding.DecodeString(b64data)
		if err != nil {
			panic(err)
		}

		f, err := os.Create("./public/media/t" + fmt.Sprint(id))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(decodedString); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}

		//pdf file
		dst, err := os.Create("./public/media/" + fmt.Sprint(id))
		if err != nil {
			panic(err)
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		// Return a JSON response if successful
		if err == nil {
			return c.JSON(http.StatusCreated, H{
				"created": id,
			})
			// Handle any errors
		} else {
			return err
		}
	}
}

//GetPdf endpoint
func GetPdf(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		id, _ := strconv.Atoi(c.Param("id"))
		pdf := models.GetPdf(db, id)

		return c.JSON(http.StatusOK, pdf)
	}
}
