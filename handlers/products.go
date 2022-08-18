package handlers

import (
	"api/data"
	"log"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.l.Println("receive request for products")

	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}

	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Products) getProducts(w http.ResponseWriter, r *http.Request) {
	h.l.Println("receive request for products")

	lp := data.GetProducts()
	err := lp.ToJSON(w)

	if err != nil {
		http.Error(w, "Unable to marshal json ", http.StatusInternalServerError)
		return
	}
}
