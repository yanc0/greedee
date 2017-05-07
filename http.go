package main

import (
	"encoding/json"
	"github.com/yanc0/collectd-http-server/collectd"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.BasicAuth != nil && config.BasicAuth.Active == true {
			user, pass, _ := r.BasicAuth()
			if !check(user, pass) {
				http.Error(w, "401, Unauthorized", 401)
				return
			}
		}
		fn(w, r)
	}
}

func check(user string, pass string) bool {
	for _, pair := range config.BasicAuth.Accounts {
		t := strings.SplitN(pair, ":", 2)
		if t[0] == user && t[1] == pass {
			return true
		}
	}
	return false
}

func handlerMetricPost(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(w, "405 Method Not Allowed - POST Only", http.StatusMethodNotAllowed)
		return
	}

	post, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer req.Body.Close()

	var metrics []collectd.CollectDMetric
	err = json.Unmarshal(post, &metrics)
	if err != nil {
		log.Println("[WARN]", err.Error())
		http.Error(w, "400, Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, p := range pluginList {
		err := p.Send(metrics)
		if err != nil {
			log.Println("[WARN]", err.Error())
		}
	}
	return
}
