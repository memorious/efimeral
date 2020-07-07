package contract

import (
	"efimeral/packages/payload"
	"efimeral/packages/resource"
	"fmt"
)

/*EfimeralObject will model our structured contract*/
type EfimeralObject struct {
	Code     string
	Consumer string
	Producer string
	Payload  []payload.EfimeralObject
}

/*Schema to hold object table*/
const Schema = "efimeral"

/*Table to hold object data*/
const Table = "contract"

/*GenerateTableQuery will generate a postgresql table*/
func (structure *EfimeralObject) GenerateTableQuery() string {
	statement := fmt.Sprintf(`-- Drop table
DROP TABLE If EXISTS %s.%s CASCADE;
CREATE TABLE %s.%s (
	code varchar NOT NULL,
	consumer varchar NOT NULL,
	producer varchar NOT NULL,
	payload varchar NOT NULL,
	quantity int4 NOT NULL,
	price float8 NOT NULL,
	CONSTRAINT contract_pk PRIMARY KEY (consumer, producer, payload),
	CONSTRAINT contract_payload_fk FOREIGN KEY (payload, producer) REFERENCES %s.payload(code, party) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT contract_consumer_fk FOREIGN KEY (consumer) REFERENCES %s.party(code) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT contract_producer_fk FOREIGN KEY (producer) REFERENCES %s.party(code) ON UPDATE CASCADE ON DELETE CASCADE
);
`, Schema, Table, Schema, Table, Schema, Schema, Schema)
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
	return efimeralObjectCollection
}

/*GetObjectData to fetch object data*/
func (structure *EfimeralObject) GetObjectData() map[string]interface{} {
	objectData := map[string]interface{}{
		"Code":     structure.Code,
		"Consumer": structure.Consumer,
		"Producer": structure.Producer,
		"Payload":  structure.Payload,
	}
	return objectData
}

/*GetObjectInterface to get interface*/
func (structure *EfimeralObject) GetObjectInterface() resource.EfimeralObjectInterface {
	var efimeralObject resource.EfimeralObjectInterface
	efimeralObject = &EfimeralObject{}
	return efimeralObject
}
