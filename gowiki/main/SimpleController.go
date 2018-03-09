package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"html/template"
)

type Page struct {
	Title string
	Body  [] byte
}

var viewPage = "show.html"
var editPage = "edit.html"
var templates = template.Must(template.ParseFiles(viewPage, editPage))

func showHandler(writer http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/show/"):]

	page, err := loadContent(fileName)

	if err != nil {
		http.Redirect(writer, r, "/edit/"+fileName, http.StatusFound)
		return
	}

	renderTemplate(writer, viewPage, page)
}

func editHandler(writer http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[len("/edit/"):]

	page, _ := loadContent(fileName)

	renderTemplate(writer, editPage, page)
}

func saveHandler(writer http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	content := r.PostFormValue("content")

	page := &Page{title, []byte(content)}

	page.save()

	http.Redirect(writer, r, "/show/"+title, http.StatusFound)
}

func renderTemplate(writer http.ResponseWriter, pageName string, page *Page) {
	err := templates.ExecuteTemplate(writer, pageName, page)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func loadContent(fileName string) (*Page, error) {
	bytes, err := readFile(fileName + ".txt")

	if err != nil {
		return &Page{Title: fileName}, err
	}

	return &Page{fileName, bytes}, err
}

func readFile(fileName string) ([]byte, error) {
	body, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (page Page) save() error {
	fileName := page.Title + ".txt"

	return ioutil.WriteFile(fileName, page.Body, 0600)
}

func main() {
	http.HandleFunc("/show/", showHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
