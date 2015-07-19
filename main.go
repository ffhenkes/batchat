package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		<html>
			<head>
				<title>The Bat Chat</title>
			</head>
			<body>
				Bat Chat!!
			</body>
		</html>
		`))
	})

	// serve
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
