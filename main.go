package main

import (
	"fileshift/models"
	"fmt"
	"log"
	"recreview/handlers"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

func main() {

	db := initDB("storage.db")
	migrate(db)

	e := echo.New()

	auth := e.Group("")
	admin := e.Group("/admin")

	auth.Use(session.Middleware(handlers.CookieStore))
	auth.Use(handlers.CheckCookie)
	admin.Use(handlers.CheckCookie)
	admin.Use(handlers.CheckAdmin)

	e.Static("/plugins", "public/plugins")

	//static html pages
	auth.Static("/", "public/index.html")
	e.File("/login", "public/login.html")

	auth.GET("/getSettings/:id", handlers.GetSettings(db))
	admin.PUT("/updateSettings", handlers.UpdateSettings(db))
	e.POST("/login", handlers.Login(db, handlers.CookieStore))

	var port string
	var err error

	port, err = handlers.ReadConfig("config.ini", "port")
	if err != nil {
		log.Printf("Main: Read Port Config: %v", err)
	}

	log.Println("Starting server...")
	log.Println(" *** Use Port: " + port + " ***")
	e.HideBanner = true
	e.HidePort = true

	https, err := handlers.ReadConfig("config.ini", "https")
	if err != nil {
		log.Printf("Main: Read https Config: %v", err)
	}

	if https == "true" {
		err = e.StartTLS(":"+port, "./certs/cert.pem", "./certs/key.pem")
		if err != nil {
			panic(fmt.Errorf("Start Https Server: %v", err))
		}

	} else {
		err = e.Start(":" + port)
		if err != nil {
			panic(fmt.Errorf("Start Http Server: %v", err))
		}
	}

}

func initDB(filepath string) *gorm.DB {
	db, err := gorm.Open("sqlite3", filepath)

	// Here we check for any db errors then exit
	if err != nil {
		log.Println(err.Error())
		panic(err)
	}

	// If we don't get any errors but somehow still don't get a db connection
	// we exit as well
	if db == nil {
		panic("db does not exist.")
	}

	db.Exec("PRAGMA foreign_keys = ON;")

	return db
}

func migrate(db *gorm.DB) {
	//use for db debugging
	//db.LogMode(true)

	db.AutoMigrate(&models.File{})
	db.AutoMigrate(&models.User{})
}
