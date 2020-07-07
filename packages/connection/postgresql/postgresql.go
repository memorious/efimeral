package postgresql

import (
	"database/sql"
	"fmt"

	/*Importing pq library*/
	_ "github.com/lib/pq"
)

/*Default Connection Settings*/
const (
	host     = ""
	port     = 
	user     = ""
	password = ""
	dbname   = ""
)

/*EfimeralObject to hold connection data */
type EfimeralObject struct {
	Connection      *sql.DB
	ConnectionError error
	Statement       *sql.Stmt
	Result          *sql.Rows
	Host            string
	Port            int
	User            string
	Password        string
	Dbname          string
}

/*Conn to establish a connection */
func (structure *EfimeralObject) Conn() bool {

	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		structure.Host,
		structure.Port,
		structure.User,
		structure.Password,
		structure.Dbname,
	)
	structure.Connection, structure.ConnectionError = sql.Open("postgres", connectionString)
	if structure.ConnectionError != nil {
		return true
	}
	return false
}

/*SetDefaultCredentials to generate a default connection*/
func (structure *EfimeralObject) SetDefaultCredentials() bool {
	structure.Host = host
	structure.Port = port
	structure.User = user
	structure.Password = password
	structure.Dbname = dbname
	return true
}

/*RunQuery to run Queries in the db*/
func (structure *EfimeralObject) RunQuery(statement string) (*sql.Rows, *sql.DB) {
	structure.SetDefaultCredentials()
	structure.Conn()
	structure.Result, structure.ConnectionError = structure.Connection.Query(statement)
	if structure.ConnectionError != nil {
		fmt.Println(structure.ConnectionError.Error())
		return structure.Result, structure.Connection
	}
	return structure.Result, structure.Connection
}
