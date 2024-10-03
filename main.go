package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func CheckPrintable(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

func CheckNewLine(s []string) bool {
	for _, r := range s {
		if len(r) != 0 {
			return false
		}
	}
	return true
}

func PrintASCIIArt(text string, banner string) (string, error) {
	spaceASCII := int(' ')
	output := ""

	file, err := os.Open(banner + ".txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	linesOfText := strings.Split(text, "\n")

	for i, word := range linesOfText {
		if CheckNewLine(linesOfText) && i != len(linesOfText)-1 {
			output += "\n"
		}

		if len(word) == 0 && !CheckNewLine(linesOfText) {
			output += "\n"
		}

		if len(word) != 0 && !CheckNewLine(linesOfText) {
			for i := 0; i < 8; i++ {
				var lineOutput []string
				for _, letter := range word {
					letterASCII := int(letter)
					start := (letterASCII-spaceASCII)*9 + 1
					if start+i < len(lines) {
						lineOutput = append(lineOutput, lines[start+i])
					}
				}

				output += strings.Join(lineOutput, "") + "\n"
			}
		}
	}

	return output, nil
}

// not found - 404
// internal server error - 500
// bad request - 400
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 - not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "405 - method not allowed", http.StatusMethodNotAllowed)
		return
	}
	tmpl, er := template.ParseFiles("index.html")
	if er != nil {
		http.Error(w, "505 -- internal server error ", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func AssciiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "405 - method not allowed", http.StatusMethodNotAllowed)
		return
	}
	text := r.FormValue("userText")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	banner := r.FormValue("bannerSelect")
	if (banner != "shadow") && (banner != "thinkertoy") && (banner != "standard") {
		http.Error(w, "400 -bad request", http.StatusBadRequest)
		return
	}
	if len(text) > 400 {
		http.Error(w, "400 - mrg", http.StatusBadRequest)

		return
	}

	result, err := PrintASCIIArt(text, banner)
	if err != nil {
		http.Error(w, "500 - internal server error", http.StatusInternalServerError)
		return
	}
	// fmt.Println(result)
	tmpl, _ := template.ParseFiles("index.html")
	// hh := map[string]string{
	// 	"Result": result,
	// }
	tmpl.Execute(w, result)
}

func main() {
	if len(os.Args) != 1 {
		fmt.Println("check args")
		return
	}
	fmt.Println("this is your port : http://localhost:8080/ ")
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/Asscii", AssciiHandler)
	http.ListenAndServe(":8080", nil)
}
