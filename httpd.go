package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

// Main webserver structure
type Server struct {
	Host string
	Port string
	Dir  string
}

// Initializes the server
func (s *Server) Init(Host, Port, Dir string) error {

	// Save config
	s.Host = Host
	s.Port = Port
	s.Dir = Dir

	// This is just to notify the user via log message later
	host := fmt.Sprintf("%s:%s", s.Host, s.Port)

	// Configure handlers
	http.HandleFunc("/", s.HandleHTMLRequest)
	http.HandleFunc("/json/", s.HandleJSONRequest)

	// And eventually start the webserver
	INFO.Println("Starting webserver:", host)
	err := http.ListenAndServe(host, nil)
	return err
}

func (s *Server) HandleHTMLRequest(w http.ResponseWriter, r *http.Request) {

	funcMap := template.FuncMap{
		"formatDate": func(date time.Time) (res string) {
			res = date.Format("01 Jan")
			return res
		},
	}

	url_path := r.URL.Path
	if url_path == "/" {
		fp := path.Join(s.Dir, "index.tpl")
		tmpl := template.Must(template.New("index.tpl").Funcs(funcMap).ParseFiles(fp))
		if err := tmpl.Execute(w, *cache); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			ERROR.Println("HTTP Error", err.Error())
		}
	} else {
		url_path = path.Join(s.Dir, url_path)
		http.ServeFile(w, r, url_path)
	}
}

// HandleApiRequest handles REST API cals from the outer space
func (s *Server) HandleJSONRequest(w http.ResponseWriter, r *http.Request) {

	var rMax = returnMax
	var js []byte

	if r.Method == "GET" {
		r.ParseForm()
		var vals *url.Values = &r.Form
		if vals.Get("max") != "" {
			max, err := strconv.Atoi(vals.Get("max"))
			if err != nil {
				// handle error
				ERROR.Println(err)
			}
			rMax = max
		}
	}

	for index, tweet := range *cache {

		if rMax < index+1 {
			break
		}

		js_part, err := json.Marshal(tweet)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			ERROR.Println("HTTP Error", err.Error())
			return
		}
		js = append(js, js_part...)

	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
