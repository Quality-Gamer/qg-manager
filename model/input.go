package model

import "github.com/cheekybits/genny/generic"

type Input struct {
	Status   string `json:"status"`
	Body map[string]generic.Type `json:"body"`
	Message string `json:"message"`
}

type InputModel struct {
	Body BodyModel `json:"body"`
}

type BodyModel struct {
	GameModel GameModel `json:"model"`
}

type InputGame struct {
	Body BodyGame `json:"body"`
}

type BodyGame struct {
	ManagerMatch ManagerMatch `json:"match"`
}