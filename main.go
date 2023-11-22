package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Port    string            `yaml:"port"`
	Buckets map[string]string `yaml:"buckets"`
	Debug   bool              `yaml:"debug"`
}

func main() {
	defaultCfg := Config{
		Port:  "3000",
		Debug: false,
	}
	cfg, err := LoadConfig("cfg.yaml")
	if err != nil {
		fmt.Printf("Load cfg.yaml failed, err=%s\n", err)
		fmt.Println("Use default cfg instead!")
		fmt.Println(defaultCfg)
		cfg = &defaultCfg
	}

	e := echo.New()

	// Routes
	e.GET("/", func(c echo.Context) error {
		bs, err := os.ReadFile("./static/html/welcome.html")
		if err != nil {
			return err
		}

		return c.HTML(http.StatusOK, string(bs))
	})

	for k, b := range cfg.Buckets {
		e.Static("/s/"+k, b)
	}

	// app.Static("/assets", "./assets")
	e.Static("/static", "./static")
	// app.Use(logger.New(logger.Config{
	// 	Format: "${time} ${status} - ${latency} ${method} ${path}\n",
	// 	// TimeFormat: "02-Jan-2006",
	// 	TimeFormat: time.RFC3339Nano,
	// 	TimeZone:   "Asia/Shanghai",
	// }))
	// app.Use(favicon.New(favicon.Config{
	// 	File: "./static/icon/favicon.ico",
	// 	URL:  "/favicon.ico",
	// }))
	//
	// app.Use(func(c *fiber.Ctx) error {
	// 	return c.Status(fiber.StatusNotFound).SendFile("./static/html/404.html")
	// })
	//
	// Start server
	log.Fatal(e.Start(":" + cfg.Port))
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
