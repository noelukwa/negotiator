package negotiator

import "net/http"

type Negotiator struct {
	req *http.Request
}

func New(req *http.Request) *Negotiator {
	return &Negotiator{req: req}
}
