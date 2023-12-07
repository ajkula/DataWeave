package main

import (
	"db_meta/databases"
	"db_meta/dbstructs"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
)

const postgresJSON = `[{"tableName":"table1","columns":[{"columnName":"id","data_type":"integer","not_null":true,"unique":false},{"columnName":"name","data_type":"character varying","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"table1_pkey","columns":["id"]}],"relationships":[{"Conname":"table2_table1_id_fkey","RelatedTableName":"table1"}]},{"tableName":"table2","columns":[{"columnName":"id","data_type":"integer","not_null":true,"unique":false},{"columnName":"table1_id","data_type":"integer","not_null":false,"unique":false},{"columnName":"description","data_type":"character varying","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"table2_pkey","columns":["id"]}],"relationships":null},{"tableName":"table3","columns":[{"columnName":"id","data_type":"integer","not_null":true,"unique":false},{"columnName":"info","data_type":"character varying","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"table3_pkey","columns":["id"]}],"relationships":null}]`
const ExpectedTablesListJSON = `["table1","table2","table3"]`
const ExpectedGraphJSON = `{"edges":[],"nodes":[{"data":{"id":"0","name":"table1","columns":[{"columnName":"id","data_type":"integer","not_null":true,"unique":false},{"columnName":"name","data_type":"character varying","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"table1_pkey","columns":["id"]}]}},{"data":{"id":"1","name":"table2","columns":[{"columnName":"id","data_type":"integer","not_null":true,"unique":false},{"columnName":"table1_id","data_type":"integer","not_null":false,"unique":false},{"columnName":"description","data_type":"character varying","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"table2_pkey","columns":["id"]}]}},{"data":{"id":"2","name":"table3","columns":[{"columnName":"id","data_type":"integer","not_null":true,"unique":false},{"columnName":"info","data_type":"character varying","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"table3_pkey","columns":["id"]}]}}]}`

type DatabaseConnectorMock struct {
	ConnectFunc          func(string, string, string, string, string) (*gorm.DB, error)
	GetTableMetadataFunc func(*gorm.DB) ([]*dbstructs.TableMetadata, error)
}

func (m *DatabaseConnectorMock) Connect(host, port, database, user, password string) (*gorm.DB, error) {
	if m.ConnectFunc != nil {
		return m.ConnectFunc(host, port, database, user, password)
	}
	return nil, nil
}

func (m *DatabaseConnectorMock) GetTableMetadata(db *gorm.DB) ([]*dbstructs.TableMetadata, error) {
	if m.GetTableMetadataFunc != nil {
		return m.GetTableMetadataFunc(db)
	}
	return nil, nil
}

func TestApp_GetTablesList(t *testing.T) {
	var tablesMetadataMock []*dbstructs.TableMetadata
	err := json.Unmarshal([]byte(postgresJSON), &tablesMetadataMock)
	assert.NoError(t, err)

	mockConnector := &DatabaseConnectorMock{
		GetTableMetadataFunc: func(*gorm.DB) ([]*dbstructs.TableMetadata, error) {
			tables := tablesMetadataMock
			return tables, nil
		},
	}

	mockDBManager := &databases.DatabaseManager{}
	originalDBM := databases.GetDatabaseManagerInstance()
	defer databases.SetInstance(originalDBM)
	mockDBManager.DB = new(gorm.DB)
	databases.SetInstance(mockDBManager)
	mockDBManager.SetConnector(mockConnector)

	_, err = mockDBManager.GetTableMetadata()
	assert.NoError(t, err)

	app := NewApp()
	tables := app.GetTablesList()

	assert.Len(t, tables, len(tablesMetadataMock))
	assert.Equal(t, "table1", tables[0].TableName)
	assert.Equal(t, "table2", tables[1].TableName)
	assert.Equal(t, "table3", tables[2].TableName)
}
