//extractor is a tool used to extarcts words and its frequency from a given url..
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"golang.org/x/net/html"
)

var workers = 1

func main() {
	fmt.Println("Starting the server on port : 8087")

	//handle to serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	//handle for word extractor
	http.HandleFunc("/", getLandingPage)

	//handle for excel column finder
	http.HandleFunc("/finder", finder)

	err := http.ListenAndServe(":8087", nil)
	if err != nil {
		log.Fatal(err)
	}

}

//handler
func getLandingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("static/tpl/layout.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		start := time.Now()
		if r.FormValue("url") == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s", "Invalid url!")
			return
		}
		resp, err := http.Get(r.FormValue("url"))
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprintf(w, "%s", "Invalid url")
			return
		}
		defer resp.Body.Close()
		data := make(chan string, 100)
		words := make(chan string)
		minionResults := make([]<-chan bool, workers)
		for i := 0; i < workers; i++ {
			minionResults[i] = minion(data, words)

		}

		var waitr sync.WaitGroup
		waitr.Add(2)
		m := make(map[string]int)
		go func(m map[string]int, words chan string) {
			for k := range words {
				if k != "" {
					if _, ok := m[k]; !ok {
						m[k] = 1
					} else {
						m[k]++
					}
				}
			}
			waitr.Done()
		}(m, words)
		go func(words chan string) {
			for _, v := range minionResults {
				for range v {

				}
			}
			close(words)
			waitr.Done()
		}(words)

		extrator(resp.Body, data)
		close(data)

		waitr.Wait()
		fmt.Println("The time elapsed :", time.Now().Sub(start))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", "Invalid request!")
	}

}

//minion cosume texts from the jobs channel , after processing sends the words to the out channel and return a done channel
func minion(jobs <-chan string, out chan string) <-chan bool {
	m := make(chan bool)

	go func() {
		for v := range jobs {
			text := mapper(v)
			words := strings.Split(text, " ")
			for _, k := range words {
				out <- k
			}
		}
		close(m)
	}()
	return m
}

//extrator tokenise the html and puts the content in to the out channel
func extrator(body io.Reader, out chan string) {
	z := html.NewTokenizer(body)
	for {
		n, _ := z.TagName()
		if string(n) == "script" || string(n) == "style" {
			z.Next()
		} else {
			tt := z.Next()
			switch {
			case tt == html.ErrorToken:
				// End of the document, we're done
				return
			case tt == html.TextToken:
				s := z.Text()
				out <- string(s)

			}
		}

	}
}

func mapper(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSymbol(r) {
			return -1
		}
		if unicode.IsSpace(r) {
			return ' '
		}

		if unicode.IsPunct(r) {
			return ' '
		}
		if unicode.IsNumber(r) {
			return -1
		}
		return unicode.ToLower(r)
	}, str)
}

//handler excel column finder
func finder(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl := template.Must(template.ParseFiles("static/tpl/finder.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == "POST" {
		m := make(map[string]interface{})
		status := ""
		if !validateColumn(r.FormValue("start")) {
			status = status + " Invalid start column !\n"
		}
		if !validateInt(r.FormValue("row")) {
			status = status + " Invalid row count !\n"
		}
		if !validateInt(r.FormValue("column")) {
			status = status + " Invalid column count !\n"
		}

		if status != "" {
			m["error"] = status
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s", status)
			return
		}
		row, _ := strconv.Atoi(r.FormValue("row"))
		column, _ := strconv.Atoi(r.FormValue("column"))
		data := make([]string, 0, row*column+1)
		data = append(data, r.FormValue("start"))
		initial := r.FormValue("start")
		for i := 0; i < row*column; i++ {
			res := increment26(initial)
			data = append(data, res)
			initial = res

		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		m["data"] = data
		json.NewEncoder(w).Encode(m)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", "Invalid request!")
	}

	return
}

func validateInt(id string) bool {
	if id == "" {
		return false
	}
	valid, err := strconv.Atoi(id)
	if err != nil {
		return false
	}
	if valid <= 0 {
		return false
	}
	return true
}
func validateColumn(val string) bool {
	if val == "" {
		return false
	}
	d := []rune(val)
	for _, k := range d {
		if k < 65 || k > 91 {
			return false
		}
	}
	return true
}

//increment26 increment the give excel column
func increment26(s string) string {
	d := []rune(s)
	n := len(d)
	flag := false
	for n > 0 {
		d[n-1] = d[n-1] + 1
		if int(d[n-1]-'A') >= 26 {
			d[n-1] = 'A'
		} else {
			flag = true
			break
		}
		n--
	}
	//check to handle ZZ case
	if !flag {
		t := []rune("A")
		d = append(t, d...)
	}
	return string(d)
}
