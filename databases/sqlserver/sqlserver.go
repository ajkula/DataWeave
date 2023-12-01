package sqlServerConnector

import (
	"db_meta/dbstructs"
	"fmt"
	"log"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type SQLServerConnector struct{}

func (conn SQLServerConnector) Connect(host, port, database, user, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", user, password, host, port, database)

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("sqlserver.go:[1]", err)
		return nil, err
	}
	return db, nil
}

func (conn SQLServerConnector) GetTableMetadata(db *gorm.DB) ([]*dbstructs.TableMetadata, error) {
	tableNames, err := conn.GetTableNames(db)
	if err != nil {
		log.Println("Error fetching table names:", err)
		return nil, err
	}

	var tables []*dbstructs.TableMetadata
	for _, tableName := range tableNames {
		table := &dbstructs.TableMetadata{TableName: tableName}

		// Get columns
		var columns []*dbstructs.Column
		result := db.Raw(`
				SELECT c.name AS column_name, t.name AS data_type, 
				c.is_nullable = 0 AS not_null,
				CASE WHEN ic.index_column_id IS NOT NULL THEN 1 ELSE 0 END AS is_unique
				FROM sys.columns c
				INNER JOIN sys.types t ON c.user_type_id = t.user_type_id
				LEFT JOIN sys.index_columns ic ON ic.object_id = c.object_id AND ic.column_id = c.column_id
				WHERE c.object_id = OBJECT_ID(?)`, tableName).Scan(&columns)
		if result.Error != nil {
			log.Println("sqlserver.go:[2]", result.Error)
			return nil, result.Error
		}
		table.Columns = columns

		// Get primary keys
		var primaryKeys []string
		rows, err := db.Raw(`
				SELECT c.name AS column_name
				FROM sys.indexes i
				INNER JOIN sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
				INNER JOIN sys.columns c ON ic.object_id = c.object_id AND c.column_id = ic.column_id
				WHERE i.is_primary_key = 1 AND i.object_id = OBJECT_ID(?)
				ORDER BY ic.key_ordinal`, tableName).Rows()
		if err != nil {
			log.Println("sqlserver.go:[3]", err)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var pkColumn string
			if err := rows.Scan(&pkColumn); err != nil {
				log.Println("sqlserver.go:[4]", err)
				return nil, err
			}
			primaryKeys = append(primaryKeys, pkColumn)
		}
		table.PrimaryKey = primaryKeys

		// Get relationships
		var relationships []*dbstructs.RelationshipMetadata
		result = db.Raw(`
				SELECT 
					fk.name AS conname, 
					OBJECT_NAME(fk.referenced_object_id) AS confrelid
				FROM 
					sys.foreign_keys fk
				WHERE 
					fk.parent_object_id = OBJECT_ID(?)`, tableName).Scan(&relationships)
		if result.Error != nil {
			log.Println("sqlserver.go:[5]", result.Error)
			return nil, result.Error
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

func (conn SQLServerConnector) GetIndexes(db *gorm.DB, tableName string) ([]*dbstructs.Index, error) {
	var indexes []*dbstructs.Index
	rows, err := db.Raw(`
			SELECT 
				i.name AS index_name, 
				COLUMN_NAME(ic.object_id, ic.column_id) AS column_name
			FROM 
				sys.indexes i
			INNER JOIN 
				sys.index_columns ic ON i.object_id = ic.object_id AND i.index_id = ic.index_id
			WHERE 
				i.object_id = OBJECT_ID(?) AND i.is_primary_key = 0
			ORDER BY 
				i.name, ic.key_ordinal`, tableName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	currentIndexName := ""
	var index *dbstructs.Index
	for rows.Next() {
		var indexName, columnName string
		if err := rows.Scan(&indexName, &columnName); err != nil {
			return nil, err
		}
		if currentIndexName != indexName {
			if index != nil {
				indexes = append(indexes, index)
			}
			index = &dbstructs.Index{Name: indexName, Columns: []string{}}
			currentIndexName = indexName
		}
		index.Columns = append(index.Columns, columnName)
	}
	if index != nil {
		indexes = append(indexes, index)
	}

	return indexes, nil
}

func (conn SQLServerConnector) GetTableNames(db *gorm.DB) ([]string, error) {
	var tableNames []string
	result := db.Raw(`
			SELECT 
				t.name AS table_name
			FROM 
				sys.tables t
			WHERE 
				t.type = 'U'`).Scan(&tableNames)
	if result.Error != nil {
		log.Println("sqlserver.go:[6]", result.Error)
		return nil, result.Error
	}
	return tableNames, nil
}
