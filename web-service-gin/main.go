package main

import (
	"example/web-service-gin/routers"
	"os"
)

func main() {
	address := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	routers.Init(address, port)

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	log.Fatal("PORT env must be set")
	// }

	// instanceId := os.Getenv("INSTANCE_ID")
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	// 	if r.Method != http.MethodGet {
	// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	text := "Hello, World!"
	// 	if instanceId != "" {
	// 		text = text + " from instance " + instanceId
	// 	}

	// 	w.Write([]byte(text))
	// })

	// server := new(http.Server)
	// server.Handler = mux
	// server.Addr = "0.0.0.0:" + port

	// log.Println("Listening on", server.Addr)
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
