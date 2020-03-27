package apolloDb

import (
	"fmt"
	"io/ioutil"
	"path"
)

type SqlApp struct {
	cs               string
	dbName           string
	scripts          string
	deleteExistingDb bool
	repo             Repository
}

func NewApp(cs string, dbName string, scripts string, provider string, deleteExistingDb bool) (*SqlApp, error) {
	repo, err := NewRepository(provider)
	if err != nil {
		return nil, fmt.Errorf("unable to create repository: %w", err)
	}

	return &SqlApp {
		cs:               cs,
		dbName:           dbName,
		scripts:          scripts,
		repo:             repo,
		deleteExistingDb: deleteExistingDb,
	}, nil
}

func Run(app *SqlApp) error {
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
