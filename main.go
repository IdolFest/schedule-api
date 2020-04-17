package main

// import (
//     "io"
//     "log"
//     "net/http"
//     "encoding/json"
// )

func main() {

	readEventCache()

	// http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
	//     var sched = readEventCache()

	//     w.Header().Add("Content-Type", "application/json")

	//     enc := json.NewEncoder(w)

	//     io.WriteString(w, "Hello from a HandleFunc!\n")
	// })
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
