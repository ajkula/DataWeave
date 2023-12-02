package databases

import (
	"context"
	"db_meta/dbstructs"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"gorm.io/gorm"
)

const expectedPostgresJSON = `[{"tableName":"table1","columns":[{"columnName":"id","data_type":"integer","not_null":true,"unique":false},{"columnName":"name","data_type":"character varying","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"table1_pkey","columns":["id"]}],"relationships":[{"Conname":"table2_table1_id_fkey","RelatedTableName":"table1"}]},{"tableName":"table2","columns":[{"columnName":"id","data_type":"integer","not_null":true,"unique":false},{"columnName":"table1_id","data_type":"integer","not_null":false,"unique":false},{"columnName":"description","data_type":"character varying","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"table2_pkey","columns":["id"]}],"relationships":null},{"tableName":"table3","columns":[{"columnName":"id","data_type":"integer","not_null":true,"unique":false},{"columnName":"info","data_type":"character varying","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"table3_pkey","columns":["id"]}],"relationships":null}]`
const expectedMySQLJSON = `[{"tableName":"table1","columns":[{"columnName":"id","data_type":"bigint","not_null":true,"unique":true},{"columnName":"name","data_type":"varchar","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"PRIMARY","columns":["id"]}],"relationships":null},{"tableName":"table2","columns":[{"columnName":"description","data_type":"varchar","not_null":false,"unique":false},{"columnName":"id","data_type":"bigint","not_null":true,"unique":true},{"columnName":"table1_id","data_type":"bigint","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"PRIMARY","columns":["id"]},{"name":"table1_id","columns":["table1_id"]}],"relationships":[{"Conname":"table2_ibfk_1","RelatedTableName":"table1"}]},{"tableName":"table3","columns":[{"columnName":"id","data_type":"bigint","not_null":true,"unique":true},{"columnName":"info","data_type":"varchar","not_null":false,"unique":false}],"primary_key":["id"],"indexes":[{"name":"PRIMARY","columns":["id"]}],"relationships":null}]`
const expectedSQLServerJSON = `[{"tableName":"spt_fallback_db","columns":[{"columnName":"xserver_name","data_type":"varchar","not_null":false,"unique":false},{"columnName":"xdttm_ins","data_type":"datetime","not_null":false,"unique":false},{"columnName":"xdttm_last_ins_upd","data_type":"datetime","not_null":false,"unique":false},{"columnName":"xfallback_dbid","data_type":"smallint","not_null":false,"unique":false},{"columnName":"name","data_type":"varchar","not_null":false,"unique":false},{"columnName":"dbid","data_type":"smallint","not_null":false,"unique":false},{"columnName":"status","data_type":"smallint","not_null":false,"unique":false},{"columnName":"version","data_type":"smallint","not_null":false,"unique":false}],"primary_key":null,"indexes":null,"relationships":null},{"tableName":"spt_fallback_dev","columns":[{"columnName":"xserver_name","data_type":"varchar","not_null":false,"unique":false},{"columnName":"xdttm_ins","data_type":"datetime","not_null":false,"unique":false},{"columnName":"xdttm_last_ins_upd","data_type":"datetime","not_null":false,"unique":false},{"columnName":"xfallback_low","data_type":"int","not_null":false,"unique":false},{"columnName":"xfallback_drive","data_type":"char","not_null":false,"unique":false},{"columnName":"low","data_type":"int","not_null":false,"unique":false},{"columnName":"high","data_type":"int","not_null":false,"unique":false},{"columnName":"status","data_type":"smallint","not_null":false,"unique":false},{"columnName":"name","data_type":"varchar","not_null":false,"unique":false},{"columnName":"phyname","data_type":"varchar","not_null":false,"unique":false}],"primary_key":null,"indexes":null,"relationships":null},{"tableName":"spt_fallback_usg","columns":[{"columnName":"xserver_name","data_type":"varchar","not_null":false,"unique":false},{"columnName":"xdttm_ins","data_type":"datetime","not_null":false,"unique":false},{"columnName":"xdttm_last_ins_upd","data_type":"datetime","not_null":false,"unique":false},{"columnName":"xfallback_vstart","data_type":"int","not_null":false,"unique":false},{"columnName":"dbid","data_type":"smallint","not_null":false,"unique":false},{"columnName":"segmap","data_type":"int","not_null":false,"unique":false},{"columnName":"lstart","data_type":"int","not_null":false,"unique":false},{"columnName":"sizepg","data_type":"int","not_null":false,"unique":false},{"columnName":"vstart","data_type":"int","not_null":false,"unique":false}],"primary_key":null,"indexes":null,"relationships":null},{"tableName":"table1","columns":[{"columnName":"id","data_type":"int","not_null":false,"unique":true},{"columnName":"name","data_type":"varchar","not_null":false,"unique":false}],"primary_key":["id"],"indexes":null,"relationships":null},{"tableName":"table2","columns":[{"columnName":"id","data_type":"int","not_null":false,"unique":true},{"columnName":"description","data_type":"varchar","not_null":false,"unique":false},{"columnName":"table1_id","data_type":"int","not_null":false,"unique":false}],"primary_key":["id"],"indexes":null,"relationships":[{"Conname":"FK__table2__table1_i__22CA2527","RelatedTableName":"table1"}]},{"tableName":"table3","columns":[{"columnName":"id","data_type":"int","not_null":false,"unique":true},{"columnName":"info","data_type":"varchar","not_null":false,"unique":false}],"primary_key":["id"],"indexes":null,"relationships":null},{"tableName":"spt_monitor","columns":[{"columnName":"lastrun","data_type":"datetime","not_null":false,"unique":false},{"columnName":"cpu_busy","data_type":"int","not_null":false,"unique":false},{"columnName":"io_busy","data_type":"int","not_null":false,"unique":false},{"columnName":"idle","data_type":"int","not_null":false,"unique":false},{"columnName":"pack_received","data_type":"int","not_null":false,"unique":false},{"columnName":"pack_sent","data_type":"int","not_null":false,"unique":false},{"columnName":"connections","data_type":"int","not_null":false,"unique":false},{"columnName":"pack_errors","data_type":"int","not_null":false,"unique":false},{"columnName":"total_read","data_type":"int","not_null":false,"unique":false},{"columnName":"total_write","data_type":"int","not_null":false,"unique":false},{"columnName":"total_errors","data_type":"int","not_null":false,"unique":false}],"primary_key":null,"indexes":null,"relationships":null},{"tableName":"MSreplication_options","columns":[{"columnName":"optname","data_type":"sysname","not_null":false,"unique":false},{"columnName":"value","data_type":"bit","not_null":false,"unique":false},{"columnName":"major_version","data_type":"int","not_null":false,"unique":false},{"columnName":"minor_version","data_type":"int","not_null":false,"unique":false},{"columnName":"revision","data_type":"int","not_null":false,"unique":false},{"columnName":"install_failures","data_type":"int","not_null":false,"unique":false}],"primary_key":null,"indexes":null,"relationships":null}]`

func createTestPostgresSchema(db *gorm.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS table1 (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255)
	);

	CREATE TABLE IF NOT EXISTS table2 (
			id SERIAL PRIMARY KEY,
			description VARCHAR(255),
			table1_id INT,
			FOREIGN KEY (table1_id) REFERENCES table1(id)
	);

	CREATE TABLE IF NOT EXISTS table3 (
			id SERIAL PRIMARY KEY,
			info VARCHAR(255)
	);`

	return db.Exec(schema).Error
}

func createTestMySQLSchema(db *gorm.DB) error {
	statements := []string{`
	CREATE TABLE IF NOT EXISTS table1 (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255)
	);`,

		`CREATE TABLE IF NOT EXISTS table2 (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			description VARCHAR(255),
			table1_id BIGINT
	);`,

		`ALTER TABLE table2
		ADD FOREIGN KEY (table1_id) REFERENCES table1(id);`,

		`CREATE TABLE IF NOT EXISTS table3 (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			info VARCHAR(255)
	);`}

	for _, stmt := range statements {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}
	return nil
}

func createTestSQLServerSchema(db *gorm.DB) error {
	schema := `
	CREATE TABLE table1 (
		id INT IDENTITY PRIMARY KEY,
		name VARCHAR(255)
	);

	CREATE TABLE table2 (
		id INT IDENTITY PRIMARY KEY,
		description VARCHAR(255),
		table1_id INT,
		FOREIGN KEY (table1_id) REFERENCES table1(id)
	);

	CREATE TABLE table3 (
		id INT IDENTITY PRIMARY KEY,
		info VARCHAR(255)
	);`

	return db.Exec(schema).Error
}

type DBContainerConfig struct {
	Image       string
	ExposedPort string
	Env         map[string]string
	InternalDB  string
	WaitingFor  func(exposedPort nat.Port) wait.Strategy
	Mounts      []testcontainers.ContainerMount
	Cmd         []string
}

func startDBContainer(t *testing.T, config DBContainerConfig) (*testcontainers.Container, string, string, string, string, string, string) {
	ctx := context.Background()

	split := strings.Split(config.ExposedPort, "/")
	portNumber, protocol := split[0], split[1]

	exposedPort, err := nat.NewPort(protocol, portNumber)
	if err != nil {
		t.Fatalf("Invalid exposed port %s for %s: %s", config.ExposedPort, config.InternalDB, err)
	}

	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		ExposedPorts: []string{config.ExposedPort},
		Env:          config.Env,
		WaitingFor:   config.WaitingFor(exposedPort),
	}

	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start %s container: %s", config.InternalDB, err)
	}

	// Get mapped port
	host, err := dbContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host for %s: %s", config.InternalDB, err)
	}

	// Get mapped port without protocol
	mappedPort, err := dbContainer.MappedPort(ctx, exposedPort)
	if err != nil {
		t.Fatalf("Failed to get container port for %s: %s", config.InternalDB, err)
	}
	var pwd string
	if val, ok := config.Env["DB_PASSWORD"]; ok {
		pwd = val
	} else {
		pwd = config.Env["SA_PASSWORD"]
	}
	return &dbContainer, config.InternalDB, host, mappedPort.Port(), config.Env["DB_NAME"], config.Env["DB_USER"], pwd
}

func startPostgresContainer(t *testing.T) (*testcontainers.Container, string, string, string, string, string, string) {
	postgresConfig := DBContainerConfig{
		Image:       "postgres:latest",
		ExposedPort: "5432/tcp",
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
			"DB_NAME":           "testdb",
			"DB_USER":           "user",
			"DB_PASSWORD":       "password",
		},
		WaitingFor: func(exposedPort nat.Port) wait.Strategy {
			return wait.ForListeningPort(exposedPort).WithStartupTimeout(60 * time.Second)
		},
		InternalDB: "postgres",
	}
	return startDBContainer(t, postgresConfig)
}

func startMySQLContainer(t *testing.T) (*testcontainers.Container, string, string, string, string, string, string) {
	mysqlConfig := DBContainerConfig{
		Image:       "mysql:latest",
		ExposedPort: "3306/tcp",
		Env: map[string]string{
			"MYSQL_DATABASE":      "testdb",
			"MYSQL_USER":          "user",
			"MYSQL_PASSWORD":      "password",
			"MYSQL_ROOT_PASSWORD": "rootpassword",
			"DB_NAME":             "testdb",
			"DB_USER":             "user",
			"DB_PASSWORD":         "password",
		},
		WaitingFor: func(exposedPort nat.Port) wait.Strategy {
			return wait.ForListeningPort(exposedPort).WithStartupTimeout(60 * time.Second)
		},
		InternalDB: "mysql",
	}
	return startDBContainer(t, mysqlConfig)
}

func startSQLServerContainer(t *testing.T) (*testcontainers.Container, string, string, string, string, string, string) {
	sqlServerConfig := DBContainerConfig{
		Image:       "mcr.microsoft.com/azure-sql-edge",
		ExposedPort: "1433/tcp",
		Env: map[string]string{
			"ACCEPT_EULA": "1",
			"DB_NAME":     "master",
			"DB_USER":     "sa",
			"SA_PASSWORD": "StrongP@ssw0rd!",
		},
		WaitingFor: func(exposedPort nat.Port) wait.Strategy {
			return wait.ForLog("SQL Server is now ready for client connections").WithStartupTimeout(60 * time.Second)
		},
		InternalDB: "sqlserver",
		Mounts: []testcontainers.ContainerMount{
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: "../init_script.sh",
				},
				Target: "/var/opt/mssql/scripts/init-script.sh",
			},
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: "../init_script.sql",
				},
				Target: "/var/opt/mssql/scripts/mssql_init.sql",
			},
		},
		Cmd: []string{"/bin/bash", "/var/opt/mssql/scripts/init-script.sh"},
	}

	return startDBContainer(t, sqlServerConfig)
}

func TestDatabaseManager_postgres_GetTableMetadata(t *testing.T) {
	postgresContainer, postgresDB, host, port, database, user, password := startPostgresContainer(t)
	defer (*postgresContainer).Terminate(context.Background())

	// Connexion to PostgreSQL container
	dbm := DatabaseManager{}
	db, err := dbm.Connect(postgresDB, host, port, database, user, password)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Create DB schema
	err = createTestPostgresSchema(db)
	assert.NoError(t, err)

	// Testing the method
	metaData, err := dbm.GetTableMetadata()
	assert.NoError(t, err)
	metaDataJSON, err := json.Marshal(metaData)
	assert.NoError(t, err)
	assert.Equal(t, expectedPostgresJSON, string(metaDataJSON))

	// Close Db after test
	SQLDB, err := db.DB()
	assert.NoError(t, err)
	defer SQLDB.Close()
}

func TestDatabaseManager_MySQL_GetTableMetadata(t *testing.T) {
	mysqlContainer, mysqldb, host, port, database, user, password := startMySQLContainer(t)
	defer (*mysqlContainer).Terminate(context.Background())

	dbm := DatabaseManager{}
	db, err := dbm.Connect(mysqldb, host, port, database, user, password)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = createTestMySQLSchema(db)
	assert.NoError(t, err)

	actualData, err := dbm.GetTableMetadata()
	assert.NoError(t, err)

	var expectedData []*dbstructs.TableMetadata
	err = json.Unmarshal([]byte(expectedMySQLJSON), &expectedData)
	assert.NoError(t, err)

	if !reflect.DeepEqual(expectedData, actualData) {
		t.Errorf("Data mismatch: expected %+v, got %+v", expectedData, actualData)
	}

	SQLDB, err := db.DB()
	assert.NoError(t, err)
	defer SQLDB.Close()
}

func TestDatabaseManager_SQLServer_GetTableMetadata(t *testing.T) {
	sqlserverContainer, sqlserverdb, host, port, database, user, password := startSQLServerContainer(t)
	defer (*sqlserverContainer).Terminate(context.Background())

	dbm := DatabaseManager{}
	db, err := dbm.Connect(sqlserverdb, host, port, database, user, password)
	if err != nil {
		t.Fatalf("Failed to connect: %s", err)
	}
	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = createTestSQLServerSchema(db)
	assert.NoError(t, err)

	metaData, err := dbm.GetTableMetadata()
	assert.NoError(t, err)

	metaDataJSON, err := json.Marshal(metaData)
	assert.NoError(t, err)
	assert.Equal(t, expectedSQLServerJSON, string(metaDataJSON))

	SQLDB, err := db.DB()
	assert.NoError(t, err)
	defer SQLDB.Close()
}
