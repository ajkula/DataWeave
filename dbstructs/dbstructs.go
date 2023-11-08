package dbstructs

type Column struct {
	ColumnName string `gorm:"column:column_name" json:"columnName"`
	DataType   string `json:"data_type"`
	NotNull    bool   `json:"not_null"`
	Unique     bool   `json:"unique"`
}

type RelationshipMetadata struct {
	Conname          string `gorm:"column:conname"`
	RelatedTableName string `gorm:"column:confrelid"`
}

type Index struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}

type TableMetadata struct {
	TableName     string                  `json:"tableName"`
	Columns       []*Column               `json:"columns"`
	PrimaryKey    []string                `json:"primary_key"`
	Indexes       []*Index                `json:"indexes"`
	Relationships []*RelationshipMetadata `json:"relationships"`
}

type NodeElement struct {
	Data *NodeData `json:"data"`
}

type NodeData struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Columns    []*Column `json:"columns"`
	PrimaryKey []string  `json:"primary_key"`
	Indexes    []*Index  `json:"indexes"`
}

type RelationshipEdge struct {
	Data *EdgeData `json:"data"`
}

type EdgeData struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type GraphResponse struct {
	Edges []*RelationshipEdge `json:"edges"`
	Nodes []*NodeElement      `json:"nodes"`
}

type PrimaryKeyIssue struct {
	TableName        string `json:"tableName"`
	IssueDescription string `json:"issueDescription"`
}

type UniqueIndexIssue struct {
	TableName        string `json:"tableName"`
	ColumnName       string `json:"columnName"`
	IssueDescription string `json:"issueDescription"`
}

type RedundantIndexIssue struct {
	TableName        string `json:"tableName"`
	IndexName        string `json:"indexName"`
	RedundantWith    string `json:"redundantWith"`
	IssueDescription string `json:"issueDescription"`
}

type ForeignKeyIssue struct {
	TableName        string `json:"tableName"`
	ColumnName       string `json:"columnName"`
	RelatedTableName string `json:"relatedTableName"`
	IssueDescription string `json:"issueDescription"`
}

type NullableColumnIssue struct {
	TableName        string `json:"tableName"`
	ColumnName       string `json:"columnName"`
	IssueDescription string `json:"issueDescription"`
}

type SchemaVerificationResults struct {
	MissingPrimaryKeys   []*PrimaryKeyIssue     `json:"missingPrimaryKeys"`
	NullableColumns      []*NullableColumnIssue `json:"nullableColumns"`
	MissingUniqueIndexes []*UniqueIndexIssue    `json:"missingUniqueIndexes"`
	ForeignKeyIssues     []*ForeignKeyIssue     `json:"foreignKeyIssues"`
	RedundantIndexes     []*RedundantIndexIssue `json:"redundantIndexes"`
}
