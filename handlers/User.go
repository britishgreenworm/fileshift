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
