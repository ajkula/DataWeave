package databases

import (
	"db_meta/dbstructs"
	"fmt"
	"sync"
)

// PerformAllVerifications scans the schema for potential issues.
func (dbm *DatabaseManager) PerformAllVerifications() (*dbstructs.SchemaVerificationResults, error) {
	results := &dbstructs.SchemaVerificationResults{}
	var mu sync.Mutex // Mutex to protect append operations

	done := make(chan bool, 6) // Channel to signal when a goroutine is done

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

	// check for SCCs in schema relations
	go func() {
		graph := &dbstructs.GraphResponse{
			Edges: dbm.Edges,
			Nodes: dbm.Nodes,
		}
		results.SCCs = FindSCCs(graph)
		done <- true
	}()

	for i := 0; i < 6; i++ {
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
