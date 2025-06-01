package application

import (
	"sync"

	"github.com/ahmadabdelrazik/jasad/pkg/config"
)

type Application struct {
	cfg config.Config

	wg sync.WaitGroup
}

func New(cfg config.Config) Application {
	return Application{
		cfg: cfg,
	}
}
