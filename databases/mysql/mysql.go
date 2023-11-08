package mysqlConnector

import (
	"db_meta/dbstructs"
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLConnector struct{}

func (conn MySQLConnector) Connect(host, port, database, user, password string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("mysql.go:[1]", err)
		return nil, err
	}
	return db, nil
}

func (conn MySQLConnector) GetTableMetadata(db *gorm.DB) ([]*dbstructs.TableMetadata, error) {
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
				SELECT column_name, data_type, is_nullable = 'NO' as not_null,
				(SELECT count(*) FROM information_schema.statistics
						WHERE table_name = ? AND non_unique = 0
						AND column_name = columns.column_name) > 0 as unique
				FROM information_schema.columns
				WHERE table_name = ?`, tableName, tableName).Scan(&columns)
		if result.Error != nil {
			log.Println("mysql.go:[2]", result.Error)
			return nil, result.Error
		}
		table.Columns = append(table.Columns, columns...)

		// Get primary keys
		var primaryKeys []string
		rows, err := db.Raw(`
				SELECT column_name
				FROM information_schema.columns
				WHERE table_name = ? AND column_key = 'PRI'
		`, tableName).Rows()
		if err != nil {
			log.Println("mysql.go:[3]", err)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var pkColumn string
			if err := rows.Scan(&pkColumn); err != nil {
				log.Println("mysql.go:[4]", err)
				return nil, err
			}
			primaryKeys = append(primaryKeys, pkColumn)
		}

		table.PrimaryKey = primaryKeys

		// Get relationships
		var relationships []*dbstructs.RelationshipMetadata
		result = db.Raw(`
				SELECT constraint_name, referenced_table_name
				FROM information_schema.key_column_usage
				WHERE table_name = ? AND referenced_table_name IS NOT NULL`, tableName).Scan(&relationships)
		if result.Error != nil {
			log.Println("mysql.go:[5]", result.Error)
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

func (conn MySQLConnector) GetIndexes(db *gorm.DB, tableName string) ([]*dbstructs.Index, error) {
	var indexes []*dbstructs.Index
	rows, err := db.Raw(`
			SELECT index_name, GROUP_CONCAT(column_name ORDER BY seq_in_index) AS columns
			FROM information_schema.statistics
			WHERE table_name = ? AND table_schema = (SELECT DATABASE())
			GROUP BY index_name
	`, tableName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var index dbstructs.Index
		var colNames string
		if err := rows.Scan(&index.Name, &colNames); err != nil {
			return nil, err
		}
		index.Columns = strings.Split(colNames, ",") // Split the string into a slice
		indexes = append(indexes, &index)
	}

	return indexes, nil
}

func (conn MySQLConnector) GetTableNames(db *gorm.DB) ([]string, error) {
	var tableNames []string
	result := db.Raw(`
			SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = (SELECT DATABASE()) AND table_type = 'BASE TABLE'`).Scan(&tableNames)
	if result.Error != nil {
		log.Println("mysql.go:[6]", result.Error)
		return nil, result.Error
	}
	return tableNames, nil
}
