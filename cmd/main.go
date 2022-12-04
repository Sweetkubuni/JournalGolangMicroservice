package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Sweetkubuni/journal/api"
	"github.com/Sweetkubuni/journal/api/controller"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	echo "github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type Pair struct {
	label         string
	default_value string
}

func ConfigEnv(args []Pair) map[string]string {
	var env map[string]string
	env = make(map[string]string)
	for _, pairVal := range args {
		if val, ok := os.LookupEnv(pairVal.label); ok {
			env[pairVal.label] = val
		} else {
			env[pairVal.label] = pairVal.default_value
		}
	}
	return env
}

func main() {

	env := ConfigEnv([]Pair{
		{"SERVER_PORT", "8080"},
		{"DB_HOST", "127.0.0.1"},
		{"DB_USER", "postgres"},
		{"DB_PASSWORD", "mysecretpassword"},
		{"DB_NAME", "app"},
		{"DB_PORT", "5432"},
	})

	var port = flag.Int("port", 8080, "Port for test HTTP server")
	flag.Parse()

	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create an instance of our handler which satisfies the generated interface
	journalHandler, err := controller.NewJournalHandlers(env["DB_HOST"], env["DB_PORT"], env["DB_USER"], env["DB_NAME"], env["DB_PASSWORD"])
	if err != nil {
		fmt.Printf("something went wrong with NewJournalHandlers")
		log.Fatal("NewJournalHandlers:", err)
		os.Exit(1)
	}

	// This is how you set up a basic Echo router
	e := echo.New()

	//set static file hosting
	e.Static("/static", "audio")

	// Log all requests
	e.Use(echomiddleware.Logger())
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	e.Use(middleware.OapiRequestValidator(swagger))

	// We now register our journalHandler above as the handler for the interface
	api.RegisterHandlers(e, journalHandler)

	// And we serve HTTP until the world ends.
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *port)))
}
