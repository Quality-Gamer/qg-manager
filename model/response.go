package model

import "github.com/cheekybits/genny/generic"

type Response struct {
	Status   string `json:"status"`
	Response []generic.Type `json:"response"`
	Message string `json:"message"`
}
