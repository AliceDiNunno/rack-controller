package rest

import e "github.com/AliceDiNunno/go-nested-traced-error"

type Status struct {
	Success bool
	Message string   `json:",omitempty"`
	Error   *e.Error `json:",omitempty"`
	Data    interface{}
	Host    string
}
