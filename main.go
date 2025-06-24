package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type QuoteResponse struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

func fetchQuote() (string, error) {
	resp, err := http.Get("https://api.quotable.io/random")
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var qr struct {
			Content string `json:"content"`
			Author  string `json:"author"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&qr); err == nil && qr.Content != "" {
			return fmt.Sprintf("%s — %s", qr.Content, qr.Author), nil
		}
	}

	resp, err = http.Get("https://zenquotes.io/api/random")
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var zq []struct {
			Q string `json:"q"`
			A string `json:"a"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&zq); err == nil && len(zq) > 0 && zq[0].Q != "" {
			return fmt.Sprintf("%s — %s", zq[0].Q, zq[0].A), nil
		}
	}

	resp, err = http.Get("https://quote-garden.herokuapp.com/api/v3/quotes/random")
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var qg struct {
			Data []struct {
				Quote  string `json:"quote"`
				Author string `json:"author"`
			}
		}
		if err := json.NewDecoder(resp.Body).Decode(&qg); err == nil && len(qg.Data) > 0 && qg.Data[0].Quote != "" {
			return fmt.Sprintf("%s — %s", qg.Data[0].Quote, qg.Data[0].Author), nil
		}
	}

	resp, err = http.Get("https://type.fit/api/quotes")
	if err == nil && resp.StatusCode == 200 {
		defer resp.Body.Close()
		var tf []struct {
			Text   string `json:"text"`
			Author string `json:"author"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&tf); err == nil && len(tf) > 0 {
			// Pick a random quote from the slice
			rand.Seed(time.Now().UnixNano())
			q := tf[rand.Intn(len(tf))]
			return fmt.Sprintf("%s — %s", q.Text, q.Author), nil
		}
	}

	return "Could not fetch quote.", fmt.Errorf("")
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		quote, err := fetchQuote()
		if err != nil {
			quote = "Could not fetch quote."
		}
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Random Quote Generator</title></head><body style='font-family:sans-serif;text-align:center;margin-top:10%%;'>`)
		fmt.Fprintf(w, "<h1>Random Quote Generator</h1>")
		fmt.Fprintf(w, "<blockquote style='font-size:1.5em;margin:2em;'>%s</blockquote>", quote)
		fmt.Fprintf(w, "<p>Credit: <a href='https://github.com/1rhino2' target='_blank'>1rhino2</a></p>")
		// Open GitHub in a new tab after 10 seconds (browser JS)
		fmt.Fprintf(w, `<script>setTimeout(function(){ window.open('https://github.com/1rhino2', '_blank'); }, 10000);</script>`)
		fmt.Fprintf(w, "</body></html>")
	} else {
		http.Redirect(w, r, "https://github.com/1rhino2", http.StatusFound)
	}
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Random Quote Generator running on http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}
