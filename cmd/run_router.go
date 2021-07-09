package main

import "github.com/AndreyBMWX6/raft-db/internal/router"

func main() {
	r := router.NewRouter()
	r.RunRouter()
}
