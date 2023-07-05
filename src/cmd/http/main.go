package main

import (
	api "github.com/jei-el/vuo.be-backend/src/api"
	config "github.com/jei-el/vuo.be-backend/src/config"
)

func main() {
	config.Load()
	api.Serve()
}
