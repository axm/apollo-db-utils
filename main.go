package main

import (
	"db-setup/apolloDb"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type AppKiller interface {
	Kill(message string)
}

type killer struct{}

func (killer *killer) Kill(message string) {
	log.Fatal(message)

	os.Exit(1)
}

func validateCli(killer AppKiller, cs string, dbName string, scripts string) {
	if cs == "" {
		flag.Usage()
		killer.Kill("Invalid arguments: missing cs.")
	}

	if dbName == "" {
		flag.Usage()
		killer.Kill("Invalid arguments: missing dbName.")
	}

	if scripts == "" {
		flag.Usage()
		killer.Kill("Invalid arguments: missing scripts.")
	}
}

var (
	cs      = flag.String("cs", "", "Postgresql connection string")
	dbName  = flag.String("dbName", "", "Postgresql database name")
	scripts = flag.String("scripts", "", "Scripts root folder")
)

func main() {
	flag.Parse()

	killer := &killer{}
	validateCli(killer, *cs, *dbName, *scripts)

	app, err := apolloDb.NewApp(*cs, strings.ToLower(*dbName), *scripts, "postgres", true)
	if err != nil {
		killer.Kill(fmt.Sprintf("error creating db setup app: %v", err))
	}

	err = apolloDb.Run(app)
	if err != nil {
		killer.Kill(fmt.Sprintf("Error running app: %s", err))
	}
}
