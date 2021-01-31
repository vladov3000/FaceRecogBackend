package main

import (
	"fmt"
	"infer"
)

// const PORT = 8000

// func barHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello world, %q", html.EscapeString(r.URL.Path))
// }

// func main() {
// 	http.HandleFunc("/bar", barHandler)

// 	log.Printf("Listening on port %d", PORT)
// 	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil))
// }

func main() {
	inferer := infer.NewInferer("build/model/model.dat", "build/model/shape_predictor.dat")
	res := inferer.GetResults("test-images/obama.jpg")
	fmt.Printf("%+v\n", res)
	inferer.Free()
}
