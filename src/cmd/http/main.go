package main

import (
	api "github.com/jeielmosi/peq.nu-backend/src/api"
	config "github.com/jeielmosi/peq.nu-backend/src/config"
)

func main() {
	config.Load()
	api.Serve()
}
