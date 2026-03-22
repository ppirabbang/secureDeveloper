package main

import (
	"gosecureskeleton/pkg/ext/db/sqlite"
	"gosecureskeleton/pkg/handler"
	"gosecureskeleton/pkg/session"
	"gosecureskeleton/pkg/util"
)

const (
	databasePath = "./app.db"
	schemaFile   = "./schema.sql"
	seedFile     = "./seed.sql"

	defaultServerPort = ":8080"
)

func main() {
	util.SetDefaultLogger()

	store, err := sqlite.New(databasePath, schemaFile, seedFile)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	sessions := session.NewStore()
	router := handler.SetupRouter(store, sessions)

	if err = router.Run(defaultServerPort); err != nil {
		panic(err)
	}
}
