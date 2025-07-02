package main

import (
	"crypto/md5"
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
	// Try APIs in order, fallback to local quotes if all fail
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

	// Local fallback quotes
	localQuotes := []string{
		"The best way to get started is to quit talking and begin doing. — Walt Disney",
		"Don't let yesterday take up too much of today. — Will Rogers",
		"It's not whether you get knocked down, it's whether you get up. — Vince Lombardi",
		"If you are working on something exciting, it will keep you motivated. — Steve Jobs",
		"Success is not in what you have, but who you are. — Bo Bennett",
		"The harder you work for something, the greater you'll feel when you achieve it. — Unknown",
		"Dream bigger. Do bigger. — Unknown",
		"Don't watch the clock; do what it does. Keep going. — Sam Levenson",
		"Great things never come from comfort zones. — Unknown",
		"Push yourself, because no one else is going to do it for you. — Unknown",
	}
	rand.Seed(time.Now().UnixNano())
	return localQuotes[rand.Intn(len(localQuotes))], nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		quote, err := fetchQuote()
		if err != nil {
			quote = "Could not fetch quote."
		}
		fmt.Fprintf(w, `<!DOCTYPE html><html><head><title>Random Quote Generator</title>
		<meta name='viewport' content='width=device-width,initial-scale=1'>
		<link rel='icon' href='/favicon.ico'>
		<style>
		:root { --bg: #fff; --fg: #222; --btn: #eee; --btnfg: #222; }
		body.dark { --bg: #181818; --fg: #eee; --btn: #333; --btnfg: #eee; }
		body { background: var(--bg); color: var(--fg); transition: background .3s, color .3s; }
		button { background: var(--btn); color: var(--btnfg); border: none; padding: 0.5em 2em; font-size: 1em; border-radius: 5px; cursor: pointer; margin: 1em; }
		.toggle { float:right; margin:1em; }
		.avatar { width:48px; height:48px; border-radius:50%%; display:inline-block; vertical-align:middle; margin-right:0.5em; }
		</style>
		</head><body style='font-family:sans-serif;text-align:center;margin-top:10%%;'>`)
		fmt.Fprintf(w, `<button class='toggle' onclick='toggleDark()' style='background:#222;color:#fff;font-size:1.5em;width:2.5em;height:2.5em;border-radius:50%%;border:2px solid #fff;box-shadow:0 2px 8px #0002;'>🌙</button>`) // more visible dark mode toggle
		fmt.Fprintf(w, "<h1>Random Quote Generator</h1>")
		fmt.Fprintf(w, "<blockquote id='quote' style='font-size:1.5em;margin:2em;'>%s</blockquote>", quote)
		fmt.Fprintf(w, `<div id='author'>%s</div>`, authorHTML(quote))
		fmt.Fprintf(w, `<button onclick='refreshQuote()' id='refreshBtn'>New Quote</button>`)
		fmt.Fprintf(w, `<script src="https://cdnjs.cloudflare.com/ajax/libs/blueimp-md5/2.19.0/js/md5.min.js"></script>`)
		fmt.Fprintf(w, `<script>
		let loading = false;
		function refreshQuote() {
			if(loading) return;
			loading = true;
			document.getElementById('refreshBtn').disabled = true;
			fetch('/api/quote').then(r=>r.json()).then(d=>{
				document.getElementById('quote').textContent = d.quote;
				if(d.author) document.getElementById('author').innerHTML = authorHTML(d.author);
			}).catch(()=>{
				document.getElementById('quote').textContent = 'Could not fetch quote.';
				document.getElementById('author').innerHTML = '';
			}).finally(()=>{
				loading = false;
				document.getElementById('refreshBtn').disabled = false;
			});
		}
		function authorHTML(author) {
			if(!author) return '';
			return '<img class="avatar" src="https://www.gravatar.com/avatar/' + md5(author.trim().toLowerCase()) + '?d=identicon"/>' + '— ' + author;
		}
		function toggleDark() {
			document.body.classList.toggle('dark');
			localStorage.setItem('dark', document.body.classList.contains('dark'));
		}
		if(localStorage.getItem('dark')==='true') document.body.classList.add('dark');
		</script>`)
		fmt.Fprintf(w, "</body></html>")
	case "/favicon.ico":
		http.ServeFile(w, r, "favicon.ico")
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 not found")
	}
}

// Helper for author avatar HTML (server fallback)
func authorHTML(quote string) string {
	// Try to extract author from quote string
	sep := " — "
	if i := len(quote); i > 0 {
		parts := []rune(quote)
		for j := 0; j < len(parts)-2; j++ {
			if string(parts[j:j+3]) == sep {
				author := string(parts[j+3:])
				if author != "" {
					hash := fmt.Sprintf("%x", md5sum([]byte(author)))
					return fmt.Sprintf("<img class='avatar' src='https://www.gravatar.com/avatar/%s?d=identicon'/> — %s", hash, author)
				}
			}
		}
	}
	return ""
}

// Minimal md5sum for gravatar
func md5sum(data []byte) [16]byte {
	return md5.Sum(data)
}

func apiQuoteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	quote, err := fetchQuote()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not fetch quote."})
		return
	}
	// Try to split quote and author for API consumers
	var text, author string
	sep := " — "
	if i := len(quote); i > 0 {
		parts := []rune(quote)
		for j := 0; j < len(parts)-2; j++ {
			if string(parts[j:j+3]) == sep {
				text = string(parts[:j])
				author = string(parts[j+3:])
				break
			}
		}
	}
	if text == "" {
		text = quote
	}
	json.NewEncoder(w).Encode(map[string]string{"quote": text, "author": author})
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/quote", apiQuoteHandler)
	fmt.Println("Random Quote Generator running on http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}

// To run the server, use the command: go run main.go
// Access the web interface at http://localhost:8080
// Fuck this shit project....
