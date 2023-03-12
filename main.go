package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port         string `yaml:"port"`
	StoragePaths string `yaml:"storage-paths"`
	Debug        bool   `yaml:"debug"`
}

func main() {
	defaultCfg := Config{
		Port:         "3000",
		StoragePaths: "./storage",
		Debug:        false,
	}
	cfg, err := LoadConfig("cfg.yaml")
	if err != nil {
		fmt.Printf("Load cfg.yaml failed, err=%s\n", err)
		fmt.Println("Use default cfg instead!")
		fmt.Println(defaultCfg)
		cfg = &defaultCfg
	}

	engine := html.New("./views", ".html")

	// Fiber instance
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./static/html/welcome.html")
	})

	// app.Get("/user", func(c *fiber.Ctx) error {
	// 	return c.SendString("Hello")
	// })
	app.Get("/storage/*", func(c *fiber.Ctx) error {
		params, err := url.PathUnescape(c.Params("*"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendFile("./static/html/404.html")
		}
		if cfg.Debug == true {
			return c.SendString(path.Join(cfg.StoragePaths, params))
		}
		return c.SendFile(path.Join(cfg.StoragePaths, params))
	})

	app.Static("/static", "./static")
	app.Use(logger.New(logger.Config{
		Format: "${time} ${status} - ${latency} ${method} ${path}\n",
		// TimeFormat: "02-Jan-2006",
		TimeFormat: time.RFC3339Nano,
		TimeZone:   "Asia/Shanghai",
	}))
	app.Use(favicon.New(favicon.Config{
		File: "./static/icon/favicon.ico",
		URL:  "/favicon.ico",
	}))

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendFile("./static/html/404.html")
	})

	// app.Static("/storage", "./storage/")

	// app.Get("/versions/api", controllers.VersionsGetAll)
	// app.Post("/versions/api", controllers.VersionsInsert)
	// app.Get("/versions/api/:version", controllers.VersionsFind)
	// app.Put("/versions/api/:version", controllers.VersionsUpdate)
	// app.Delete("/versions/api/:version", controllers.VersionsDelete)
	//
	// app.Get("/versions/view", controllers.VersionsViewGetAll)
	// app.Get("/versions/view/edit", controllers.VersionsViewEditFind)
	// app.Post("/versions/view/edit", controllers.VersionsViewEditInsert)
	// app.Get("/versions/view/edit/:version", controllers.VersionsViewEditFind)
	// app.Put("/versions/view/edit/:version", controllers.VersionsViewEditUpdate)
	// app.Get("/versions/:id/page", controllers.VersionsGetAll)

	// Start server
	log.Fatal(app.Listen(":" + cfg.Port))
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
