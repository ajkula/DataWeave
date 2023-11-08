package databases

import (
	mysqlConnector "db_meta/databases/mysql"
	postgresConnector "db_meta/databases/postgres"
	sqlserverConnector "db_meta/databases/sqlserver"
	"db_meta/dbstructs"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"gorm.io/gorm"
)

var (
	instance *DatabaseManager
	once     sync.Once
)

type DatabaseManager struct {
	connector DatabaseConnector
	DB        *gorm.DB
	Tables    []*dbstructs.TableMetadata
	Nodes     []*dbstructs.NodeElement
	Edges     []*dbstructs.RelationshipEdge
}

func GetDatabaseManagerInstance() *DatabaseManager {
	once.Do(func() {
		instance = &DatabaseManager{}
	})
	return instance
}

func (dbm *DatabaseManager) Connect(dbType, host, port, database, user, password string) (*gorm.DB, error) {
	log.Println("Trying to connect to DB...")
	switch dbType {
	case "postgres":
		dbm.connector = &postgresConnector.PostgresConnector{}
	case "mysql":
		dbm.connector = &mysqlConnector.MySQLConnector{}
	case "sqlserver":
		dbm.connector = &sqlserverConnector.SQLServerConnector{}
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

		dbm.Edges = []*dbstructs.RelationshipEdge{}
		// Add relations
		for _, rel := range table.Relationships {
			dbm.Edges = append(dbm.Edges, &dbstructs.RelationshipEdge{
				Data: &dbstructs.EdgeData{
					ID:     rel.Conname,
					Source: table.TableName,
					Target: rel.RelatedTableName,
				},
			})
		}
	}
}

// PerformAllVerifications scans the schema for potential issues.
func (dbm *DatabaseManager) PerformAllVerifications() (*dbstructs.SchemaVerificationResults, error) {
	results := &dbstructs.SchemaVerificationResults{}
	var mu sync.Mutex // Mutex to protect append operations

	done := make(chan bool, 5) // Channel to signal when a goroutine is done

	// Check for missing primary keys
	go func() {
		for _, table := range dbm.Tables {
			if len(table.PrimaryKey) == 0 {
				issue := &dbstructs.PrimaryKeyIssue{
					TableName:        table.TableName,
					IssueDescription: "Missing primary key",
				}
				mu.Lock()
				results.MissingPrimaryKeys = append(results.MissingPrimaryKeys, issue)
				mu.Unlock()
			}
		}
		done <- true
	}()

	// Check for columns that should be NOT NULL
	go func() {
		for _, table := range dbm.Tables {
			for _, column := range table.Columns {
				if !column.NotNull && column.DataType != "serial" {
					issue := &dbstructs.NullableColumnIssue{
						TableName:        table.TableName,
						ColumnName:       column.ColumnName,
						IssueDescription: "Column should be NOT NULL",
					}
					mu.Lock()
					results.NullableColumns = append(results.NullableColumns, issue)
					mu.Unlock()
				}
			}
		}
		done <- true
	}()

	// Check for uniqueness constraints without indexes
	go func() {
		for _, table := range dbm.Tables {
			for _, column := range table.Columns {
				if column.Unique && !dbm.columnHasIndex(table, column.ColumnName) {
					issue := &dbstructs.UniqueIndexIssue{
						TableName:        table.TableName,
						ColumnName:       column.ColumnName,
						IssueDescription: "Missing unique index",
					}
					mu.Lock()
					results.MissingUniqueIndexes = append(results.MissingUniqueIndexes, issue)
					mu.Unlock()
				}
			}
		}
		done <- true
	}()

	// Check for foreign key issues
	go func() {
		for _, table := range dbm.Tables {
			for _, relationship := range table.Relationships {
				relatedTable := dbm.findTableByName(relationship.RelatedTableName)
				if relatedTable == nil {
					issue := &dbstructs.ForeignKeyIssue{
						TableName:        table.TableName,
						ColumnName:       relationship.Conname,
						IssueDescription: fmt.Sprintf("Linked table not found: %s", relationship.RelatedTableName),
					}
					mu.Lock()
					results.ForeignKeyIssues = append(results.ForeignKeyIssues, issue)
					mu.Unlock()
					continue
				}

				hasIndex := dbm.columnHasIndex(relatedTable, relationship.Conname)
				if !hasIndex {
					issue := &dbstructs.ForeignKeyIssue{
						TableName:        relatedTable.TableName,
						ColumnName:       relationship.Conname,
						IssueDescription: "Missing index for foreign key",
					}
					mu.Lock()
					results.ForeignKeyIssues = append(results.ForeignKeyIssues, issue)
					mu.Unlock()
				}
			}
		}
		done <- true
	}()

	// Check for redundant indexes
	go func() {
		for _, table := range dbm.Tables {
			indexMap := make(map[string][]string)
			for _, index := range table.Indexes {
				indexMap[index.Name] = index.Columns
			}

			for indexName, columns := range indexMap {
				for otherIndexName, otherColumns := range indexMap {
					if indexName != otherIndexName && dbm.isSubset(columns, otherColumns) {
						issue := &dbstructs.RedundantIndexIssue{
							TableName:        table.TableName,
							IndexName:        indexName,
							RedundantWith:    otherIndexName,
							IssueDescription: "Redundant index",
						}
						mu.Lock()
						results.RedundantIndexes = append(results.RedundantIndexes, issue)
						mu.Unlock()
					}
				}
			}
		}
		done <- true
	}()

	for i := 0; i < 5; i++ {
		<-done
	}

	return results, nil
}

// Toolbox methods
func (dbm *DatabaseManager) findTableByName(tableName string) *dbstructs.TableMetadata {
	for _, table := range dbm.Tables {
		if table.TableName == tableName {
			return table
		}
	}
	return nil
}

func (dbm *DatabaseManager) columnHasIndex(table *dbstructs.TableMetadata, columnName string) bool {
	for _, index := range table.Indexes {
		for _, idxColumn := range index.Columns {
			if idxColumn == columnName {
				return true
			}
		}
	}
	return false
}

func (dbm *DatabaseManager) isSubset(subset, set []string) bool {
	setMap := make(map[string]struct{})
	for _, item := range set {
		setMap[item] = struct{}{}
	}
	for _, item := range subset {
		if _, found := setMap[item]; !found {
			return false
		}
	}
	return true
}
