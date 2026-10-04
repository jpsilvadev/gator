package main

import (
	"github.com/jpsilvadev/gator/internal/config"
	"github.com/jpsilvadev/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
