package databases

import (
	"db_meta/dbstructs"

	"gorm.io/gorm"
)

type DatabaseConnector interface {
	Connect(string, string, string, string, string) (*gorm.DB, error)
	GetTableMetadata(*gorm.DB) ([]*dbstructs.TableMetadata, error)
}
