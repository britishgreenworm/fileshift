package handlers

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

//do not make this key public
var CookieStore = GetCookieSigningKey()

func GetCookieSigningKey() *sessions.CookieStore {
	SigningKey, err := ReadConfig("config.ini", "signingkey")
	if err != nil {
		log.Println("Read Cookie Signing Key: " + err.Error())
	}
	CookieSigningKey := sessions.NewCookieStore([]byte(SigningKey))
	return CookieSigningKey
}

func ReadConfig(filePath string, arg string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Opening Config File: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strAry := strings.Split(scanner.Text(), "=")

		if strAry[0] == arg {
			return strAry[1], nil
		}

	}
	return "", fmt.Errorf("Reading Config: %v", err)
}

func CheckCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		sess, err := CookieStore.Get(c.Request(), "fileshift")
		if err != nil {
			log.Println(err.Error())
			return c.Redirect(http.StatusFound, "/login?loginError=cookie")
		}

		if sess.IsNew {
			log.Println("Cookie not found...")
			//return c.String(http.StatusUnauthorized, "You do not have permission to view this page")
			return c.Redirect(302, "/login")
		}

		return next(c)
	}
}

func CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := CookieStore.Get(c.Request(), "fileshift")
		if err != nil {
			log.Println(err.Error())
			return c.Redirect(http.StatusFound, "/login?loginError=admin")
		}

		if sess.IsNew {
			return c.String(http.StatusUnauthorized, "You do not have permission to view this content")
		}

		if sess.Values["role"] != "admin" {
			return c.Redirect(http.StatusFound, "/login?loginError=admin")
		}

		return next(c)
	}
}

func Login(db *gorm.DB, cookieStore *sessions.CookieStore) echo.HandlerFunc {
	return func(c echo.Context) error {

		password := c.FormValue("password")

		ConfigAdminPassword, err := ReadConfig("config.ini", "adminpassword")
		if err != nil {
			log.Printf("Login: Read Admin Config: %v", err)
			return c.JSON(http.StatusInternalServerError, "Internal Error Reading Config File, Contact server administrator.")
		}

		if password == ConfigAdminPassword {
			log.Println("admin logged in...")

			// Create session
			sess, err := cookieStore.Get(c.Request(), "fileshift")
			if err != nil {
				log.Println(err.Error())
				return c.Redirect(http.StatusFound, "/login?loginError=admin")
			}

			log.Println("Got session")

			sess.Values["role"] = "admin"
			sess.Options.MaxAge = 0
			err = sess.Save(c.Request(), c.Response())

			//sess.Store()
			if err != nil {
				log.Println(err.Error())
			}

			log.Println("Added role and saved session`")

			return c.Redirect(http.StatusFound, "/")

		}

		ConfigGeneralPassword, err := ReadConfig("config.ini", "generalpassword")
		if err != nil {

			log.Printf("Login: Read General Config: %v", err)
			return c.JSON(http.StatusInternalServerError, "Internal Error Reading Config File, Contact server administrator.")
		}

		if password == ConfigGeneralPassword {
			log.Println("general user logged in...")

			// Create session
			sess, err := cookieStore.Get(c.Request(), "fileshift")
			if err != nil {
				log.Println(err.Error())
				return c.Redirect(http.StatusFound, "/login?loginError=admin")
			}

			log.Println("Got session")

			sess.Values["role"] = "general"
			sess.Options.MaxAge = 0
			err = sess.Save(c.Request(), c.Response())

			//sess.Store()
			if err != nil {
				log.Println(err.Error())
			}

			log.Println("Added role and saved session`")

			return c.Redirect(http.StatusFound, "/")
		}

		return c.Redirect(http.StatusFound, "/login?loginError=password")

	}
}

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}
