package postgresConnector

import (
	"db_meta/dbstructs"
	"fmt"
	"log"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConnector struct{}

func (conn PostgresConnector) Connect(host, port, database, user, password string) (*gorm.DB, error) {
	databaseURL := fmt.Sprintf("host=%s user=%s password=%s database=%s port=%s sslmode=disable", host, user, password, database, port)
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Println("postgres.go:18", err)
		return nil, err
	}
	return db, nil
}

func (conn PostgresConnector) GetTableMetadata(db *gorm.DB) ([]*dbstructs.TableMetadata, error) {
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
				(SELECT count(*) FROM information_schema.table_constraints tc
						JOIN information_schema.constraint_column_usage ccu
						ON ccu.constraint_name = tc.constraint_name
						WHERE tc.table_name = columns.table_name
						AND tc.constraint_type = 'UNIQUE'
						AND ccu.column_name = columns.column_name) > 0 as unique
				FROM information_schema.columns
				WHERE table_name = ?`, tableName).Scan(&columns)
		if result.Error != nil {
			log.Println("postgres.go:42", result.Error)
			return nil, result.Error
		}
		table.Columns = append(table.Columns, columns...)

		// Get primary keys
		var primaryKeys []string
		rows, err := db.Raw(`
				SELECT kcu.column_name
				FROM information_schema.table_constraints tc
				JOIN information_schema.key_column_usage kcu
						ON tc.constraint_name = kcu.constraint_name
						AND tc.table_schema = kcu.table_schema
				WHERE tc.table_name = ? 
						AND tc.constraint_type = 'PRIMARY KEY'
						AND tc.table_schema = 'public'
				ORDER BY kcu.ordinal_position
		`, tableName).Rows()
		if err != nil {
			log.Println("postgres.go:57", err)
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var pkColumn string
			if err := rows.Scan(&pkColumn); err != nil {
				log.Println("postgres.go:72", err)
				return nil, err
			}
			primaryKeys = append(primaryKeys, pkColumn)
		}

		table.PrimaryKey = primaryKeys

		// Get relationships
		var relationships []*dbstructs.RelationshipMetadata
		result = db.Raw(`
				SELECT conname, confrelid::regclass 
				FROM pg_constraint 
				WHERE confrelid = (SELECT oid FROM pg_class WHERE relname = ?)`, table.TableName).Scan(&relationships)
		if result.Error != nil {
			log.Println("postgres.go:87", result.Error)
			return nil, result.Error
		}
		table.Relationships = relationships

		// Get indexes
		indexes, err := conn.GetIndexes(db, table.TableName)
		if err != nil {
			log.Println("Error fetching indexes:", err)
			return nil, err
		}
		table.Indexes = indexes
		tables = append(tables, table)
	}

	return tables, nil
}

func (conn PostgresConnector) GetIndexes(db *gorm.DB, tableName string) ([]*dbstructs.Index, error) {
	var indexes []*dbstructs.Index
	rows, err := db.Raw(`
			SELECT i.relname as indexname, array_agg(a.attname) AS columns
			FROM pg_class t
			INNER JOIN pg_index ix ON t.oid = ix.indrelid
			INNER JOIN pg_class i ON i.oid = ix.indexrelid
			INNER JOIN pg_attribute a ON a.attrelid = t.oid AND a.attnum = ANY(ix.indkey)
			WHERE t.relkind = 'r' AND t.relname = ? AND i.relkind = 'i'
			GROUP BY i.relname
	`, tableName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var index dbstructs.Index
		var colNames pq.StringArray // pq.StringArray to deal with PostgreSQL arrays
		if err := rows.Scan(&index.Name, &colNames); err != nil {
			return nil, err
		}
		index.Columns = colNames // pq.StringArray is a []string alias
		indexes = append(indexes, &index)
	}

	return indexes, nil
}

func (conn PostgresConnector) GetTableNames(db *gorm.DB) ([]string, error) {
	var tableNames []string
	result := db.Raw(`
			SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = 'public'`).Scan(&tableNames)
	if result.Error != nil {
		log.Println("postgres.go:28", result.Error)
		return nil, result.Error
	}
	return tableNames, nil
}
