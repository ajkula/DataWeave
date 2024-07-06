package main

import (
	"context"
	"db_meta/api"
	"db_meta/apigen"
	"db_meta/databases"
	"db_meta/dbstructs"
	"encoding/json"
	"log"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ConfigureGorm(dbType, host, port, database, user, password string) (string, error) {
	var tableMetadata []*dbstructs.TableMetadata
	connector := databases.GetDatabaseManagerInstance()

	_, err := connector.Connect(dbType, host, port, database, user, password)
	if err != nil {
		log.Println("app.go:38", err)

		return "", err
	}

	tableMetadata, err = connector.GetTableMetadata()
	if err != nil {
		log.Println("app.go:45 - Erreur lors de la récupération des métadonnées des tables :", err)
		return "", err
	}

	jsonData, err := json.Marshal(tableMetadata)
	if err != nil {
		log.Println("app.go:51", err)
		return "", err
	}

	return string(jsonData), nil
}

func (a *App) GetTablesList() []*dbstructs.TableMetadata {
	tables := databases.GetDatabaseManagerInstance().GetTablesList()
	return tables
}

func (a *App) GraphTransform() (string, error) {
	connector := databases.GetDatabaseManagerInstance()
	response := &dbstructs.GraphResponse{
		Edges: connector.Edges,
		Nodes: connector.Nodes,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	return string(jsonResponse), nil
}

func (a *App) PerformAllVerifications() (string, error) {
	connector := databases.GetDatabaseManagerInstance()
	var verifications *dbstructs.SchemaVerificationResults
	var jsonResponse []byte
	var err error
	if verifications, err = connector.PerformAllVerifications(); err != nil {
		return "", err
	}
	if jsonResponse, err = json.Marshal(verifications); err != nil {
		return "", err
	}
	return string(jsonResponse), err
}

func (a *App) GenerateOpenApi(config *api.APIConfig) (string, error) {
	tables := databases.GetDatabaseManagerInstance().GetTablesList()
	var bytesArray []byte
	var err error
	if bytesArray, err = apigen.GenerateOpenAPI(tables, config); err != nil {
		return "", err
	}
	return string(bytesArray), err
}
