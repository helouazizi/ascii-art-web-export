package server

import (
	"fmt"
	"html/template"
	"net/http"

	"ascii-art-web-export/functions"
)

type PageData struct {
	Message string
}

// store this assci here for thee export needed
var asciiArt string

const maxInputTextLength = 500

var temple01 *template.Template

func Home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "Error 404: NOT FOUND", http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("template/index.html")
	temple01 = tmpl
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ParseForm(r *http.Request) (string, string, error) {
	inputText := r.FormValue("inputText")
	if len(inputText) > maxInputTextLength {
		return "", "", fmt.Errorf("input text exceeds %d characters", maxInputTextLength)
	}
	banner := r.FormValue("choice")
	return inputText, banner, nil
}

func ReadBannerTemplate(banner string) ([]string, error, bool) {
	switch banner {
	case "standard", "shadow", "thinkertoy":
		return functions.ReadFile("banners/" + banner + ".txt")
	default:
		return nil, fmt.Errorf("error: 300 invalid banner choice: %s", banner), false
	}
}

func TreatData(templ []string, inputText string) string {
	asciiArt = functions.TraitmentData(templ, inputText)
	return functions.TraitmentData(templ, inputText)
}

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error 405: Method not allowed", http.StatusBadRequest)
		return
	}
	inputText, banner, err := ParseForm(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	templ, err, bol := ReadBannerTemplate(banner)
	if err != nil {
		if bol {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	treatedText := TreatData(templ, inputText)

	test := temple01
	err = test.Execute(w, PageData{Message: treatedText})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* add a function hundler to hundle the export  of the data */
func ExportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error 405: Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Create a file with the ASCII art text
	fileContent := []byte(asciiArt)

	// Set the Content-Disposition header to force the browser to download the file
	w.Header().Set("Content-Disposition", "attachment; filename=ascii-art.txt")

	// Set the Content-Type header to text/plain
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Write the file content to the response writer
	w.Write(fileContent)

	// or
	// Serve the file using http.ServeContent
	// http.ServeContent(w, r, fileName, time.Now(), bytes.NewReader(fileContent))
}
