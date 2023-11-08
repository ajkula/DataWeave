package sqliteConnector

import (
	"db_meta/dbstructs"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteConnector struct{}

func (conn SQLiteConnector) Connect(host, port, database, user, password string) (*gorm.DB, error) {
	// only the 'database' parameter is used...
	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{})
	if err != nil {
		log.Println("sqlite.go: Connect error", err)
		return nil, err
	}
	return db, nil
}

func (conn SQLiteConnector) GetTableMetadata(db *gorm.DB) ([]*dbstructs.TableMetadata, error) {
	tableNames, err := conn.GetTableNames(db)
	if err != nil {
		log.Println("Error fetching table names:", err)
		return nil, err
	}

	var tables []*dbstructs.TableMetadata
	for _, tableName := range tableNames {
		table := &dbstructs.TableMetadata{TableName: tableName}

		// Get columns
		columns, err := conn.GetColumns(db, tableName)
		if err != nil {
			log.Println("Error fetching columns:", err)
			return nil, err
		}
		table.Columns = columns

		// Get primary keys
		primaryKeys, err := conn.GetPrimaryKeys(db, tableName)
		if err != nil {
			log.Println("Error fetching primary keys:", err)
			return nil, err
		}
		table.PrimaryKey = primaryKeys

		// Get relationships
		relationships, err := conn.GetRelationships(db, tableName)
		if err != nil {
			log.Println("Error fetching relationships:", err)
			return nil, err
		}
		table.Relationships = relationships

		// Get indexes
		indexes, err := conn.GetIndexes(db, tableName)
		if err != nil {
			log.Println("Error fetching indexes:", err)
			return nil, err
		}
		table.Indexes = indexes

		tables = append(tables, table)
	}

	return tables, nil
}

func (conn SQLiteConnector) GetColumns(db *gorm.DB, tableName string) ([]*dbstructs.Column, error) {
	var columns []*dbstructs.Column
	rows, err := db.Raw(fmt.Sprintf("PRAGMA table_info('%s');", tableName)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			dataType   string
			notNullInt int
			dfltValue  interface{}
			pkInt      int
		)
		if err := rows.Scan(&cid, &name, &dataType, &notNullInt, &dfltValue, &pkInt); err != nil {
			return nil, err
		}
		column := &dbstructs.Column{
			ColumnName: name,
			DataType:   dataType,
			NotNull:    notNullInt != 0,
			Unique:     false, // Will be updated after fetching unique indexes
		}
		columns = append(columns, column)
	}

	// Check for unique columns
	uniqueIndexes, err := conn.GetUniqueIndexes(db, tableName)
	if err != nil {
		log.Println("Error fetching unique indexes:", err)
		return nil, err
	}
	for _, column := range columns {
		if _, exists := uniqueIndexes[column.ColumnName]; exists {
			column.Unique = true
		}
	}

	return columns, nil
}

func (conn SQLiteConnector) GetPrimaryKeys(db *gorm.DB, tableName string) ([]string, error) {
	var primaryKeys []string
	rows, err := db.Raw(fmt.Sprintf("PRAGMA table_info('%s');", tableName)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			cid        int
			name       string
			dataType   string
			notNullInt int
			dfltValue  interface{}
			pkInt      int
		)
		if err := rows.Scan(&cid, &name, &dataType, &notNullInt, &dfltValue, &pkInt); err != nil {
			return nil, err
		}
		if pkInt == 1 {
			primaryKeys = append(primaryKeys, name)
		}
	}

	return primaryKeys, nil
}

func (conn SQLiteConnector) GetRelationships(db *gorm.DB, tableName string) ([]*dbstructs.RelationshipMetadata, error) {
	var relationships []*dbstructs.RelationshipMetadata
	rows, err := db.Raw(fmt.Sprintf("PRAGMA foreign_key_list('%s');", tableName)).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id       int
			seq      int
			table    string
			from     string
			to       string
			onUpdate string
			onDelete string
			match    string
		)
		if err := rows.Scan(&id, &seq, &table, &from, &to, &onUpdate, &onDelete, &match); err != nil {
			return nil, err
		}
		relationship := &dbstructs.RelationshipMetadata{
			Conname:          from,
			RelatedTableName: table,
		}
		relationships = append(relationships, relationship)
	}

	return relationships, nil
}

func (conn SQLiteConnector) GetIndexes(db *gorm.DB, tableName string) ([]*dbstructs.Index, error) {
	var indexes []*dbstructs.Index
	indexRows, err := db.Raw(fmt.Sprintf("PRAGMA index_list('%s');", tableName)).Rows()
	if err != nil {
		return nil, err
	}
	defer indexRows.Close()

	for indexRows.Next() {
		var index dbstructs.Index
		if err := indexRows.Scan(&index.Name, nil, nil); err != nil {
			return nil, err
		}

		colRows, err := db.Raw(fmt.Sprintf("PRAGMA index_info('%s');", index.Name)).Rows()
		if err != nil {
			return nil, err
		}
		defer colRows.Close()

		for colRows.Next() {
			var colName string
			if err := colRows.Scan(nil, &colName); err != nil {
				return nil, err
			}
			index.Columns = append(index.Columns, colName)
		}

		indexes = append(indexes, &index)
	}

	return indexes, nil
}

func (conn SQLiteConnector) GetUniqueIndexes(db *gorm.DB, tableName string) (map[string]bool, error) {
	uniqueIndexes := make(map[string]bool)
	indexRows, err := db.Raw(fmt.Sprintf("PRAGMA index_list('%s');", tableName)).Rows()
	if err != nil {
		return nil, err
	}
	defer indexRows.Close()

	for indexRows.Next() {
		var indexName string
		var unique bool
		if err := indexRows.Scan(&indexName, &unique, nil); err != nil {
			return nil, err
		}
		if unique {
			colRows, err := db.Raw(fmt.Sprintf("PRAGMA index_info('%s');", indexName)).Rows()
			if err != nil {
				return nil, err
			}
			defer colRows.Close()

			for colRows.Next() {
				var colName string
				if err := colRows.Scan(nil, &colName); err != nil {
					return nil, err
				}
				uniqueIndexes[colName] = true
			}
		}
	}

	return uniqueIndexes, nil
}

func (conn SQLiteConnector) GetTableNames(db *gorm.DB) ([]string, error) {
	var tableNames []string
	rows, err := db.Raw("SELECT name FROM sqlite_master WHERE type='table';").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}

	return tableNames, nil
}
