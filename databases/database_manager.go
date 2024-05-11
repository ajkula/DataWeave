package databases

import (
	mysqlConnector "db_meta/databases/mysql"
	postgresConnector "db_meta/databases/postgres"
	sqliteConnector "db_meta/databases/sqlite"
	sqlserverConnector "db_meta/databases/sqlserver"
	"db_meta/dbstructs"
	"errors"
	"log"
	"strconv"

	"gorm.io/gorm"
)

var (
	instance *DatabaseManager
	// once     sync.Once
)

type DatabaseManager struct {
	connector DatabaseConnector
	DB        *gorm.DB
	Tables    []*dbstructs.TableMetadata
	Nodes     []*dbstructs.NodeElement
	Edges     []*dbstructs.RelationshipEdge
}

// func GetDatabaseManagerInstance() *DatabaseManager {
// 	once.Do(func() {
// 		instance = &DatabaseManager{}
// 	})
// 	return instance
// }

func GetDatabaseManagerInstance() *DatabaseManager {
	if instance == nil {
		instance = &DatabaseManager{}
	}
	return instance
}

func SetInstance(mockInstance *DatabaseManager) {
	instance = mockInstance
}

func (dbm *DatabaseManager) SetConnector(con DatabaseConnector) {
	dbm.connector = con
}

func (dbm *DatabaseManager) GetTablesList() []*dbstructs.TableMetadata {
	return dbm.Tables
}

// GetTablesListFunc is a function variable to get table metadata.
var GetTablesListFunc = DefaultGetTablesList

func DefaultGetTablesList() []*dbstructs.TableMetadata {
	return GetDatabaseManagerInstance().GetTablesList()
}

func (dbm *DatabaseManager) Connect(dbType, host, port, database, user, password string) (*gorm.DB, error) {
	log.Println("Trying to connect to DB...")
	switch dbType {
	case "postgres":
		dbm.SetConnector(&postgresConnector.PostgresConnector{})
	case "mysql":
		dbm.SetConnector(&mysqlConnector.MySQLConnector{})
	case "sqlite":
		dbm.SetConnector(&sqliteConnector.SQLiteConnector{})
	case "sqlserver":
		dbm.SetConnector(&sqlserverConnector.SQLServerConnector{})
	default:
		return nil, errors.New("unsopported database")
	}

	if dbm.connector == nil {
		return nil, errors.New("DB connector not initialized")
	}

	db, err := dbm.connector.Connect(host, port, database, user, password)
	if err != nil {
		log.Println("database_manager.go:[1]", err)
		return nil, err
	}
	dbm.DB = db
	return dbm.DB, nil
}

func (dbm *DatabaseManager) GetTableMetadata() ([]*dbstructs.TableMetadata, error) {
	if dbm.DB == nil {
		return nil, errors.New("DB not connected")
	}

	if dbm.connector == nil {
		return nil, errors.New("DB connector not initialized")
	}

	tables, err := dbm.connector.GetTableMetadata(dbm.DB)
	if err != nil {
		log.Println("database_manager.go:[2]", err)
		return nil, err
	}
	dbm.Tables = tables
	log.Printf("%#v\n", dbm.Tables)
	dbm.TransformToGraph()
	return dbm.Tables, nil
}

func (dbm *DatabaseManager) TransformToGraph() {
	dbm.Nodes = []*dbstructs.NodeElement{}
	dbm.Edges = []*dbstructs.RelationshipEdge{}

	for index, table := range dbm.Tables {
		// Add Nodes
		dbm.Nodes = append(dbm.Nodes, &dbstructs.NodeElement{
			Data: &dbstructs.NodeData{
				ID:         strconv.Itoa(index),
				Name:       table.TableName,
				Columns:    table.Columns,
				PrimaryKey: table.PrimaryKey,
				Indexes:    table.Indexes,
			},
		})

		// Add relations
		for _, rel := range table.Relationships {
			dbm.Edges = append(dbm.Edges, &dbstructs.RelationshipEdge{
				Data: &dbstructs.EdgeData{
					ID:     rel.Conname,
					Source: rel.SourceTableName,
					Target: rel.RelatedTableName,
				},
			})
		}
	}
}
