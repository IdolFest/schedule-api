package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type config struct {
	bind           string
	sheetUrl       string
	cacheTimeout   int
	allowedOrigins string
	timezone       string
}

func parseFlags() (config, error) {
	c := config{}
	flag.IntVar(&c.cacheTimeout, "cache-timeout", 300, "The timeout when a new copy of the schedule should be fetched. This applies also when the schedule cannot be fetched.")
	flag.StringVar(&c.bind, "bind", "0.0.0.0:8080", "The host:port to bind to.")
	flag.StringVar(&c.allowedOrigins, "allowed-origins", "", "The set of Origins that should be returned for requests.")
	flag.StringVar(&c.sheetUrl, "sheet-url", "", "The URL of the published Schedule Spreadsheet. Expected response is in CSV format.")
	flag.StringVar(&c.timezone, "timezone", "America/New_York", "The timezome to assume for the spreadsheet")
	flag.Parse()

	if c.sheetUrl == "" {
		return c, fmt.Errorf("you must specify --sheet-url")
	}
	return c, nil
}

func main() {
	c, err := parseFlags()
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}
	http.HandleFunc("/schedule", func(writer http.ResponseWriter, request *http.Request) {
		sched, times, order, err := readEventCache(c.sheetUrl, c.cacheTimeout, c.timezone)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			io.WriteString(writer, "Unable to complete request")
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.Header().Add("Access-Control-Allow-Origin", c.allowedOrigins)

		enc := json.NewEncoder(writer)
		_ = enc.Encode(Response{Rooms: sched, Times: times, RoomOrder: order})
	})
	http.HandleFunc("/schedule/", func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/schedule", http.StatusMovedPermanently)
	})
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	log.Printf("Serving on %s\n", c.bind)
	log.Fatal(http.ListenAndServe(c.bind, nil))
}
