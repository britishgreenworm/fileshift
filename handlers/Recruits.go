package handlers

import (
	"fmt"
	"net/http"
	"recreview/models"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

//type H map[string]interface{}

// GetRecruits endpoint
func GetUsers(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.JSON(http.StatusOK, models.GetUsers(db, id))

	}
}

// GetRecruit endpoint
func GetUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		// Use our new model to delete a recruit
		Recruit := models.GetRecruit(db, id)
		// Return a JSON response on success
		return c.JSON(http.StatusOK, H{
			"Recruit": Recruit,
		})
	}
}

//UpdateRecruit endPoint
func UpdateUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var Recruit models.Recruit

		c.Bind(&Recruit)

		id, err := models.UpdateRecruit(db, Recruit)
		// Return a JSON response if successful
		if err == nil {
			return c.JSON(http.StatusCreated, H{
				"updated": id,
			})
			// Handle any errors
		} else {
			return err
		}
	}
}

// CreateRecruit endpoint
func CreateUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Instantiate a new task
		var recruit models.Recruit
		// Map imcoming JSON body to the new Recruit
		c.Bind(&recruit)
		// Add a task using our new model
		id, err, recruit := models.CreateRecruit(db, recruit)
		// Return a JSON response if successful
		if err == nil {
			return c.JSON(http.StatusCreated, H{
				"created": id,
				"recruit": recruit,
			})
			// Handle any errors
		} else {
			return err
		}
	}
}

// DeleteRecruit endpoint
func DeleteUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		_, err := models.DeleteRecruit(db, id)
		if err != nil {
			return fmt.Errorf("Delete Recruit: %v", err)
		}

		return c.JSON(http.StatusOK, H{
			"deleted": id,
		})
	}
}
