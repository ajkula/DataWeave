package postgresConnector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"gorm.io/gorm"
)

func createTestSchema(db *gorm.DB) error {
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

func startPostgresContainer(t *testing.T) (*testcontainers.Container, string, string, string, string, string) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start postgres container: %s", err)
	}

	// Get host and port
	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %s", err)
	}
	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %s", err)
	}

	return &postgresContainer, host, port.Port(), "testdb", "user", "password"
}

func TestPostgresConnector_Connect(t *testing.T) {
	// Starting PostgreSQL container and get connection infos
	postgresContainer, host, port, database, user, password := startPostgresContainer(t)
	defer (*postgresContainer).Terminate(context.Background())

	// Connection to PostgreSQL
	connector := PostgresConnector{}
	db, err := connector.Connect(host, port, database, user, password)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Test schema creation in DB
	err = createTestSchema(db)
	assert.NoError(t, err)

	// Close DB at the end
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	defer sqlDB.Close()
}

func TestPostgresConnector_GetTableNames(t *testing.T) {
	postgresContainer, host, port, database, user, password := startPostgresContainer(t)
	defer (*postgresContainer).Terminate(context.Background())

	connector := PostgresConnector{}
	db, err := connector.Connect(host, port, database, user, password)
	assert.NoError(t, err)

	err = createTestSchema(db)
	assert.NoError(t, err)

	tableNames, err := connector.GetTableNames(db)
	assert.NoError(t, err)
	assert.Equal(t, []string{"table1", "table2", "table3"}, tableNames)
}

func TestPostgresConnector_Connect_error(t *testing.T) {
	connector := PostgresConnector{}

	// Test with invalid params
	_, err := connector.Connect("invalidHost", "invalidPort", "invalidDatabase", "invalidUser", "invalidPassword")
	assert.Error(t, err, "Invalid parameters should fail to connect")

}

// GetTableMetadata, GetIndexes, etc.
