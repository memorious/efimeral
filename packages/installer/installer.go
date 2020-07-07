package installer

import (
	"efimeral/packages/contract"
	"efimeral/packages/party"
	"efimeral/packages/payload"
	"efimeral/packages/resource"
	"strings"
)

type EfimeralObject struct {
	packages []resource.EfimeralObjectInterface
}

/*Schema to hold object table*/
const Schema = ""

/*Table to hold object data*/
const Table = ""

/*Prepare to be used for pre-installation setup*/
func (structure *EfimeralObject) Prepare() *EfimeralObject {
	// Insert modules to be installed in this list of Io interface
	structure.packages = []resource.EfimeralObjectInterface{
		&party.EfimeralObject{},
		&payload.EfimeralObject{},
		&contract.EfimeralObject{},
	}
	return structure
}

/*Install to be used for installation of tables/entities*/
func (structure *EfimeralObject) Install() bool {
	var statements []string
	for _, entity := range structure.packages {
		statements = append(statements, entity.GenerateTableQuery())
	}

	statement := strings.Join(statements, "")
	ObjectIo := new(resource.EfimeralObject).GetIoInterface()
	ObjectIo.RunQuery(statement)
	return true
}

/*GenerateTableQuery will generate a postgresql Table*/
func (structure *EfimeralObject) GenerateTableQuery() string {
	statement := ""
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
		"Packages": structure.packages,
	}
	return objectData
}

/*GetObjectInterface to get interface*/
func (structure *EfimeralObject) GetObjectInterface() resource.EfimeralObjectInterface {
	var efimeralObject resource.EfimeralObjectInterface
	efimeralObject = &EfimeralObject{}
	return efimeralObject
}
