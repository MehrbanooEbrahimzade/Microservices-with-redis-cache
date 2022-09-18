package main

import (
	"encoding/json"
	"fmt"
	"github.com/MehrbanooEbrahimzade/golang-redis-example/cache"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/products", httpHandler)
	http.ListenAndServe(":8080", nil)
}

func httpHandler(w http.ResponseWriter, req *http.Request) {
	for i := 0; i < 10; i++ {

		t := time.Now()
		response, err := cache.GetProducts()

		if err != nil {

			fmt.Fprintf(w, err.Error()+"\r\n")

		} else {

			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")

			if err := enc.Encode(response); err != nil {
				fmt.Println(err.Error())
			}

		}
		t2 := time.Since(t)
		fmt.Println(t2.String())
		fmt.Println("-------------------------------------------------")
		time.Sleep(3 * time.Second)
	}

}
