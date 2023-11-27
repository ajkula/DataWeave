package postgresConnector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
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

	// Récupération de l'adresse IP et du port
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
	// Démarrage du conteneur PostgreSQL et récupération des informations de connexion
	postgresContainer, host, port, database, user, password := startPostgresContainer(t)
	defer (*postgresContainer).Terminate(context.Background())

	// Connexion au conteneur PostgreSQL
	connector := PostgresConnector{}
	db, err := connector.Connect(host, port, database, user, password)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// Création du schéma de test dans la base de données
	err = createTestSchema(db)
	assert.NoError(t, err)

	// Fermer la connexion à la fin du test
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	defer sqlDB.Close()
}

func TestPostgresConnector_GetTableNames(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Erreur lors de la création de l'instance db, mock: %s", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Erreur lors de la création de l'instance GORM: %s", err)
	}

	// Configurer le mock pour simuler une réponse pour GetTableNames
	mock.ExpectQuery("^SELECT table_name").WillReturnRows(sqlmock.NewRows([]string{"table_name"}).AddRow("table1").AddRow("table2"))

	// Test de GetTableNames
	connector := PostgresConnector{}
	tableNames, err := connector.GetTableNames(gormDB)
	assert.NoError(t, err)
	assert.Equal(t, []string{"table1", "table2"}, tableNames)
}

func TestPostgresConnector_Connect_error(t *testing.T) {
	connector := PostgresConnector{}

	// Test avec des paramètres invalides
	_, err := connector.Connect("invalidHost", "invalidPort", "invalidDatabase", "invalidUser", "invalidPassword")
	assert.Error(t, err, "Une connexion avec des paramètres invalides devrait échouer")

}

// Autres tests pour GetTableMetadata, GetIndexes, etc.