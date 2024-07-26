package middleware

import "github.com/rs/zerolog"

type Middleware struct {
	log *zerolog.Logger
}


func New(log *zerolog.Logger) *Middleware {
	return &Middleware{
		log: log,
	}
}
