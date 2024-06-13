package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := "7540"
	webDir := "./web"
	fmt.Printf("Starting server at port %s", port)
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finished server at port " + port) // наверное не нужно, так как не выполнится из-за ListenAndServe
	/**if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}**/
}
