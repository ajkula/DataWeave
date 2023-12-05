package sqlServerConnector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresConnector_Connect_error(t *testing.T) {
	connector := SQLServerConnector{}

	// Test with invalid params
	_, err := connector.Connect("invalidHost", "invalidPort", "invalidDatabase", "invalidUser", "invalidPassword")
	assert.Error(t, err, "Invalid parameters should fail to connect")
}
