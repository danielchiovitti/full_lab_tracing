package full_tracing

import (
	"fmt"
	"net/http"
)

func main() {
	h := full_lab_cep.InitializeHandlers()
	r := h.GetRoutes()
	fmt.Println("Starting")
	http.ListenAndServe(":3500", r)
}
