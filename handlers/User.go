package handlers

import (
  "fmt"
  "fileshift/models"
)

// CreateRecruit endpoint
func CreateUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

    user := models.User{}
		c.Bind(&user)
		user, err := models.CreateUser(db, user)

    if err != nil {
      return c.JSON(http.StatusInternalServerError, H{
        "error": err.error()
      })
    }

    return c.JSON(http.StatusCreated, H{
      "user": user
    })
	}
}

// GetUser endpoint
func GetUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

    id, _ := strconv.Atoi(c.Param("id"))
		user, err := models.GetUser(db, id)
    if err != nil {
      return c.JSON(http.StatusInternalServerError, H{
        "error": err.error()
      })
    }

		return c.JSON(http.StatusOK, H{
			"User": user
		})
  }
}
