package main

import (
	"fmt"
	"log"
	"os"

	"github.com/osalomon89/test-crud-api/cmd/api/app"
	"github.com/osalomon89/test-crud-api/internal/platform/environment"
)

func main() {
	env := environment.GetFromString(os.Getenv("GO_ENVIRONMENT"))

	dependencies, err := app.BuildDependencies(env)
	if err != nil {
		log.Fatal("error at dependencies building", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := app.Build(dependencies)
	if err := app.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("error at furyapp startup", err)
	}
}
