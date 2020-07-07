package resource

import (
	"database/sql"
	"efimeral/packages/connection/postgresql"
	"fmt"
	"strings"
)

/*EfimeralObjectInterface interface for generation of initial resources*/
type EfimeralObjectInterface interface {
	GenerateTableQuery() string
	GetSchema() string
	GetTable() string
	GetObjectInterface() EfimeralObjectInterface
	GetCollection(map[string][]string) []EfimeralObjectInterface
	GetObjectData() map[string]interface{}
}

/*EfimeralDbIo interface for communication with persistent data*/
type EfimeralDbIo interface {
	RunQuery(string) (*sql.Rows, *sql.DB)
	GetCollection(EfimeralObjectInterface, map[string][]string) (*sql.Rows, *sql.DB)
}

/*EfimeralObject to hold resource structure info*/
type EfimeralObject struct {
	// connection type
	postgresql.EfimeralObject
}

/*GetCollection to fetch data from db*/
func (structure *EfimeralObject) GetCollection(efobj EfimeralObjectInterface, filters map[string][]string) (*sql.Rows, *sql.DB) {
	statement := fmt.Sprintf("SELECT * FROM %s.%s", efobj.GetSchema(), efobj.GetTable())
	var wheres []string
	currentClause := " WHERE"
	for attributeOpp, attributeValue := range filters {
		clause := fmt.Sprintf("%s %s %s %s", currentClause, attributeOpp, attributeValue[0], attributeValue[1])
		wheres = append(wheres, clause)
		currentClause = " AND"
	}
	statement = statement + strings.Join(wheres, "") + ";"
	return structure.RunQuery(statement)
}

/*GetIoInterface to get interface*/
func (structure *EfimeralObject) GetIoInterface() EfimeralDbIo {
	var ObjectIo EfimeralDbIo
	ObjectIo = &EfimeralObject{}
	return ObjectIo
}
