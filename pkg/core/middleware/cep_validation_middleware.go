package middleware

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"regexp"
	"sync"
)

var locy sync.Mutex
var CepValidationMiddlewareInstance CepValidationMiddleware

type CepValidationMiddleware struct{}

func NewCepValidationMiddleware() CepValidationMiddleware {
	return CepValidationMiddleware{}
}

func (c *CepValidationMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cep := chi.URLParam(r, "cep")
		re := regexp.MustCompile(`^\d{5}\d{3}$`)
		valid := re.MatchString(cep)

		if !valid {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
