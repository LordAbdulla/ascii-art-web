package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ascii-art-web/internal/ascii"
)

// var to declare the template
var (
	tmpl      *template.Template
	errorTmpl *template.Template
)

// structs for the result
type PageData struct {
	Result string
}

// struct for the error
type ErrorData struct {
	Code    int
	Title   string
	Message string
	Gif     string
	Alt     string
}

func main() {
	var err error
	tmpl, err = template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		log.Fatalf("Error loading index.html: %v", err)
	}

	errorTmpl, err = template.ParseFiles(filepath.Join("templates", "error.html"))
	if err != nil {
		log.Fatalf("Error loading error.html: %v", err)
	}

	// To run the server
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/ascii-art", handleAsciiArt)
	http.Handle("/styles.css", http.FileServer(http.Dir("static")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Println("listening on : http://localhost:8080/ \nTo ShutDown The Server Ctrl+C")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Handles the server errors
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" { // Handles if page not found
		renderError(w, http.StatusNotFound, "Page Not Found : Invalid URL")
		return
	}

	if r.Method != http.MethodGet { // if the index method is not get
		renderError(w, http.StatusMethodNotAllowed, "Method Not Allowed : Only GET is supported")
		return
	}
	// if the template is not found
	if err := tmpl.Execute(w, PageData{}); err != nil {
		log.Println("template error:", err)
		renderError(w, http.StatusInternalServerError, "Internal Server Error : Failed to render page")
		return
	}
}

// Handles doing the AsciiArt
func handleAsciiArt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // if the method is not post
		renderError(w, http.StatusMethodNotAllowed, "Method Not Allowed : Only POST is supported")
		return
	}

	if err := r.ParseForm(); err != nil { // To read the HTML form
		renderError(w, http.StatusBadRequest, "Bad Request : Failed to parse form")
		return
	}
	vals, ok := r.Form["text"]
	if !ok {
		renderError(w, http.StatusBadRequest, "Bad Request: Missing 'text' parameter")
		return

	}
	text := vals[0]                                               // to declare the text in html
	banner := r.FormValue("banner")                                           // to declare the banner in html
	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" { // if the banners given is unavailable
		renderError(w, http.StatusBadRequest, "Bad Request : Invalid banner")
		return

	}
	if banner == "" {
		banner = "standard"
	}

	if strings.TrimSpace(text) == "" {
		renderError(w, http.StatusBadRequest, "Bad Request : Invalid Character ")
		return

	}
	if len(text) >= 1000 {
		renderError(w, http.StatusBadRequest, "Bad Request : length of text is out of range ")
		return

	}

	if err := ascii.ValidatePrintable(text); err != nil { // to validate if the text in printable or no
		renderError(w, http.StatusBadRequest, "Bad Request : "+err.Error())
		return
	}

	path := filepath.Join("banners", banner+".txt")
	if err := ascii.CheckHash(path); err != nil { // checks hash if it is changed or no
		renderError(w, http.StatusInternalServerError, "Internal Server Error : Invalid banner file")
		return
	}

	data, err := os.ReadFile(path)
	if err != nil { // reads the file
		renderError(w, http.StatusInternalServerError, "Internal Server Error : Cannot read banner file")
		return
	}
	var lines []string
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			line := data[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, string(line))
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, string(data[start:]))
	}

	result := ascii.DoAsciiArt(text, lines)                           // renders the result
	if err := tmpl.Execute(w, PageData{Result: result}); err != nil { // Result is saved as .Result
		renderError(w, http.StatusInternalServerError, "Internal Server Error : Failed to render result")
		return
	}
}

// reads and declares the errors in the html
func renderError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(code)
	data := ErrorData{Code: code}
	switch code {
	case http.StatusBadRequest: // 400
		data.Title = "400 — Bad Request"
		data.Gif = "/static/img/jim.webp"
		data.Alt = "Jim Halpert shrugging from The Office"
		if msg == "" {
			msg = "That request did not  look quite right. Try again!"
		}

	case http.StatusNotFound: // 404
		data.Title = "404 — Page Not Found"
		data.Gif = "/static/img/the-office-no.webp"
		data.Alt = "Michael Scott yelling NO from The Office"
		if msg == "" {
			msg = "We could not  find the page you were looking for."
		}

	case http.StatusInternalServerError: // 500
		data.Title = "500 — Internal Server Error"
		data.Gif = "/static/img/walter.webp"
		data.Alt = "Walter White  face from Breaking Bad"
		if msg == "" {
			msg = "Something went wrong on our side. Please try again later."
		}

	default:
		data.Title = fmt.Sprintf("Error %d", code)
		data.Gif = "/static/img/breaking.webp"
		data.Alt = "Walter White  face from Breaking Bad"
		if msg == "" {
			msg = "An unexpected error occurred."
		}
	}
	if msg == "" {
		msg = "We could not find the page you were looking for."
	}
	data.Message = msg

	if err := errorTmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
