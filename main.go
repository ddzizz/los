package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

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

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	for k, b := range cfg.Buckets {
		pattern := path.Join("/s/", k, "/")
		log.Printf("Handle pattern '%s' to %s\n", pattern, b)
		http.Handle("/s/"+k+"/", http.StripPrefix(pattern, http.FileServer(http.Dir(b))))
	}

	log.Print("Listening on :3000...")
	err = http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
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
