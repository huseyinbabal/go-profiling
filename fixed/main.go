package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"regexp"
	"strings"
	"sync"
)

var logs = makeLogs(20000)

// recent keeps every search result for the "recent searches" page.
var (
	mu     sync.Mutex
	recent [][]byte
)

func search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	re, err := regexp.Compile(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var out strings.Builder
	for _, line := range logs {
		if re.MatchString(line) {
			out.WriteString(line)
			out.WriteByte('\n')
		}
	}
	body := []byte(out.String())

	mu.Lock()
	recent = append(recent, body)
	if len(recent) > 50 {
		recent = recent[1:]
	}
	mu.Unlock()

	w.Write(body)
}

func main() {
	http.HandleFunc("/search", search)
	log.Println("listening on :9092")
	log.Fatal(http.ListenAndServe(":9092", nil))
}

func makeLogs(n int) []string {
	levels := []string{"INFO", "WARN", "ERROR"}
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("2026-10-25T10:%02d:%02d %s order=%d user=%d took=%dms",
			i/60%60, i%60, levels[i%7%3], 100000+i, i%977, i%450)
	}
	return lines
}
