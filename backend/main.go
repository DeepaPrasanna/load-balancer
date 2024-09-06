// Just a test backend server to see how stuff works. Not using it right now
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", reqhandler)
	fmt.Printf("Backend Server is running at http://localhost:8081")

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func reqhandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Recieved request from %q\n", r.RemoteAddr)
	log.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
	log.Printf("Host: %s\n", r.Host)
	for key, value := range r.Header {
		for _, val := range value {
			log.Println(key, ": ", val)
		}
	}
	fmt.Fprintln(w, "Hello from backend server")

}
