package payload

import (
	"efimeral/packages/resource"
	"fmt"
	"log"
)

/*EfimeralObject to hold payload values*/
type EfimeralObject struct {
	Code     string
	Name     string
	Party    string
	Quantity int
	Price    float32
}

/*Schema to hold object table*/
const Schema = "efimeral"

/*Table to hold object data*/
const Table = "payload"

/*GenerateTableQuery will generate a postgresql Table*/
func (structure *EfimeralObject) GenerateTableQuery() string {
	statement := fmt.Sprintf(`-- Drop Table
DROP Table If EXISTS %s.%s CASCADE;
CREATE Table  %s.%s (
	code varchar NOT NULL,
	"name" varchar NOT NULL,
	party varchar NOT NULL,
	quantity int4 NOT NULL,
	price float8 NOT NULL,
	CONSTRAINT payload_pk PRIMARY KEY (code, party),
	CONSTRAINT payload_fk FOREIGN KEY (party) REFERENCES %s.party(code) ON UPDATE CASCADE ON DELETE CASCADE
);
`, Schema, Table, Schema, Table, Schema)
	return statement
}

/*GetSchema ...*/
func (structure *EfimeralObject) GetSchema() string {
	return Schema
}

/*GetTable ...*/
func (structure *EfimeralObject) GetTable() string {
	return Table
}

/*GetCollection retuns (un)filtered collections of objects*/
func (structure *EfimeralObject) GetCollection(filters map[string][]string) []resource.EfimeralObjectInterface {
	var efimeralObjectCollection []resource.EfimeralObjectInterface
	ObjectIo := new(resource.EfimeralObject).GetIoInterface()
	efimeralObjects, queryResultConnection := ObjectIo.GetCollection(structure.GetObjectInterface(), filters)
	for efimeralObjects.Next() {
		objectEfimeralObject := new(EfimeralObject)
		err := efimeralObjects.Scan(&objectEfimeralObject.Code, &objectEfimeralObject.Name, &objectEfimeralObject.Party, &objectEfimeralObject.Quantity, &objectEfimeralObject.Price)
		if err != nil {
			log.Fatal(err)
		}
		efimeralObjectCollection = append(efimeralObjectCollection, objectEfimeralObject)
	}
	efimeralObjects.Close()
	queryResultConnection.Close()
	return efimeralObjectCollection
}

/*GetObjectData to fetch object data*/
func (structure *EfimeralObject) GetObjectData() map[string]interface{} {
	objectData := map[string]interface{}{
		"Code":     structure.Code,
		"Name":     structure.Name,
		"Party":    structure.Party,
		"Quantity": structure.Quantity,
		"Price":    structure.Price,
	}
	return objectData
}

/*GetObjectInterface to get interface*/
func (structure *EfimeralObject) GetObjectInterface() resource.EfimeralObjectInterface {
	var efimeralObject resource.EfimeralObjectInterface
	efimeralObject = &EfimeralObject{}
	return efimeralObject
}
