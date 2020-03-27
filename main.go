package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
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

type SqlApp struct {
	cs               string
	dbName           string
	scripts          string
	deleteExistingDb bool
	repo             Repository
killer 				AppKiller
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

	repo, err := NewRepository("postgres")
	if err != nil {
		killer.Kill("Invalid provider: 'postgres'")
	}

	app := &SqlApp{
		cs:               *cs,
		dbName:           strings.ToLower(*dbName),
		scripts:          *scripts,
		repo:             repo,
		deleteExistingDb: true,
		killer:killer,
	}

	err = run(app)
	if err != nil {
		killer.Kill(fmt.Sprintf("Error running app: %s", err))
	}
}

func run(app *SqlApp) error {
	//cwd, _ := os.Getwd()
	//root := path.Join(cwd, app.scripts)
	root := app.scripts
	directories, err := ioutil.ReadDir(root)
	if err != nil {
		return fmt.Errorf("unable to process scripts folder: %w", err)
	}

	if app.deleteExistingDb {
		fmt.Println(fmt.Sprintf("dropping database %s", app.dbName))
		err = app.repo.DropDatabase(app.cs, app.dbName)
		if err != nil {
			return fmt.Errorf("unable to drop database: %w", err)
		}
		fmt.Println(fmt.Sprintf("dropped successfully"))
		fmt.Println()

		fmt.Println(fmt.Sprintf("creating database %s", app.dbName))
		err = app.repo.CreateDatabase(app.cs, app.dbName)
		if err != nil {
			return fmt.Errorf("unable to create database: %w", err)
		}
		fmt.Println("created successfully")
		fmt.Println()
	}

	for _, dir := range directories {
		if !dir.IsDir() {
			continue
		}

		dirPath := path.Join(root, dir.Name())
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			return fmt.Errorf("unable to get directory contents: %w", err)
		}

		for _, file := range files {
			filePath := path.Join(dirPath, file.Name())
			contents, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("unable to read file contents '%s': %w", filePath, err)
			}

			sql := string(contents)
			fmt.Println(fmt.Sprintf("executing %s/%s", dir.Name(), file.Name()))

			err = app.repo.Execute(app.cs, app.dbName, sql)
			if err != nil {
				return fmt.Errorf("unable to execute: %w", err)
			}
			fmt.Println(fmt.Sprintf("finshed executing %s/%s", dir.Name(), file.Name()))
		}
	}
	fmt.Println()

	return nil
}
