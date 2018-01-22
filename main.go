package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Rule struct {
	Pattern string `json:"pattern"`
	Target  string `json:"target"`
	Status  int    `json:"status"`
}

type Config struct {
	Port      int
	TargetURL string `json:"target_url"`
	Rules     []Rule `json:"rules"`
}

const (
	configFile = "./config.json"
)

var (
	config Config
)

func redirectFunc(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	for _, rule := range config.Rules {
		match, err := regexp.MatchString(rule.Pattern, r.URL.Path)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(rule.Pattern, match)

		if match {
			fullURL := config.TargetURL + r.URL.Path
			log.Println(fullURL)

			http.Redirect(w, r, fullURL, rule.Status)
			return
		}
	}
}

func init() {
	fileContent, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		panic(err)
	}
}

func main() {
	for _, rule := range config.Rules {
		http.HandleFunc(rule.Pattern, redirectFunc)
	}

	err := http.ListenAndServe(fmt.Sprintf(":%v", config.Port), nil)
	if err != nil {
		panic(err)
	}
}
