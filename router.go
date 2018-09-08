package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ResponseStruct struct {
	Status	string	`json:"status"`
	Message	string	`json:"message"`
}

type Response []ResponseStruct

type vhostTemplate struct {
	Domain string
}

var password = ""

func check(e error) error{
	if e != nil {
		return e
	}
	return nil
}

func parse(path string, domain string) {

	t, err := template.ParseFiles(path)
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	// A sample config
	config := map[string]string{
		"domain": domain,
	}

	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
	f.Close()
}


func valdidatePassword(w http.ResponseWriter, r *http.Request) bool{
	if (password == "") {
		return true
	} else if (r.Header.Get("password") != password) {
		result := Response{
			ResponseStruct{
				Status: "error",
				Message: "Invalid password, unauthorized."},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return false
	}
	return true
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/vhost/add", vhostAdd)
	router.HandleFunc("/vhost/delete", vhostDelete)
	log.Fatal(http.ListenAndServe(":7601", router))
}

func vhostAdd(w http.ResponseWriter, r *http.Request) {
	if valdidatePassword(w, r) == false {
		return
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	var domain = r.Form.Get("domain")
	var template = r.Form.Get("template")
	templateContent, err := ioutil.ReadFile("templates/" + template + ".conf")

	f, _ := os.Create("/etc/apache2/sites-available/" + domain + ".conf")
	f.Write(templateContent)
	f.Close()

	parse("/etc/apache2/sites-available/" + domain + ".conf", domain)

	result := Response{
		ResponseStruct{
			Status: "success",
			Message: "Successfully added vhost",},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
func vhostDelete(w http.ResponseWriter, r *http.Request) {
	if valdidatePassword(w, r) == false {
		return
	}

	result := Response{
		ResponseStruct{
			Status: "success",
			Message: "Successfully deleted vhost",},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
