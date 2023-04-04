/*
Copyright 2023 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// This is a very simple server that stores and retrieves robots to a map that
// will be lost when the server is restarted.

type Robot struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

var robots = make(map[string]*Robot)

func main() { //nolint:gocyclo
	http.HandleFunc("/robots", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			name := r.URL.Query().Get("name")
			robot, ok := robots[name]
			if !ok {
				http.Error(w, "Robot not found", http.StatusNotFound)
				return
			}
			if err := json.NewEncoder(w).Encode(robot); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case "POST":
			robot := new(Robot)
			if err := json.NewDecoder(r.Body).Decode(robot); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if _, ok := robots[robot.Name]; ok {
				http.Error(w, "Robot already exists", http.StatusBadRequest)
				return
			}
			robots[robot.Name] = robot
			w.WriteHeader(http.StatusCreated)
		case "PUT":
			name := r.URL.Query().Get("name")
			if name == "" {
				http.Error(w, "Missing name parameter", http.StatusBadRequest)
				return
			}
			robot := new(Robot)
			err := json.NewDecoder(r.Body).Decode(robot)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if _, ok := robots[name]; !ok {
				http.Error(w, "Robot not found", http.StatusNotFound)
				return
			}
			delete(robots, name)
			robots[robot.Name] = robot
			w.WriteHeader(http.StatusOK)
		case "DELETE":
			name := r.URL.Query().Get("name")
			if name == "" {
				http.Error(w, "Missing name parameter", http.StatusBadRequest)
				return
			}
			if _, ok := robots[name]; !ok {
				http.Error(w, "Robot not found", http.StatusNotFound)
				return
			}
			delete(robots, name)
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	log.Println("Listening on :8080")
	log.Printf("Example commands:\n")
	log.Printf("Create:\n\tcurl -XPOST -H \"Content-type: application/json\" -d '{\"name\": \"myrobot\", \"color\": \"green\"}' 'http://127.0.0.1:8080/robots'\n")
	log.Printf("Get:\n\tcurl 'http://127.0.0.1:8080/robots?&name=myrobot'\n")
	log.Printf("Update:\n\tcurl -XPUT -H \"Content-type: application/json\" -d '{\"name\": \"myrobot\", \"color\": \"yellow\"}' 'http://127.0.0.1:8080/robots?&name=myrobot'\n")
	log.Printf("Delete:\n\tcurl -XDELETE 'http://127.0.0.1:8080/robots?&name=myrobot'\n")
	s := &http.Server{
		Addr:              ":8080",
		Handler:           http.TimeoutHandler(http.DefaultServeMux, 10*time.Second, "timeout"),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
