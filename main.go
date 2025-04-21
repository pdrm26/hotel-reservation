package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/pdrm26/hotel-reservation/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the API server")
	flag.Parse()

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)

	app.Listen(*listenAddr)
}

func handleRoot(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"message": "All things in my life going very nicly."})
}
