package party

import (
	"efimeral/packages/connection/efimeralhttp"
	"efimeral/packages/errors"
	"efimeral/packages/resource"
	"encoding/json"
	"fmt"
	"log"
)

/*EfimeralObject to hold party values*/
type EfimeralObject struct {
	Dist     string
	Lat, Lng float64
	Code     string
	Payload  []resource.EfimeralObjectInterface
}

/*Schema to hold object table*/
const Schema = "efimeral"

/*Table to hold object data*/
const Table = "party"

const coordsQuery = "https://public.opendatasoft.com/api/records/1.0/search/?dataset=us-zip-code-latitude-and-longitude&q=%s"

const harvesine = "SELECT lat, lng, code, distance " +
	"FROM ( " +
	"SELECT " +
	"z.lat, z.lng, z.code, p.radius, p.distance_unit " +
	"* DEGREES(ACOS(LEAST(1.0, COS(RADIANS(p.lat)) " +
	"* COS(RADIANS(z.lat)) " +
	"* COS(RADIANS(p.lng - z.lng)) " +
	"+ SIN(RADIANS(p.lat)) " +
	"* SIN(RADIANS(z.lat))))) AS distance " +
	"FROM %s.%s AS z " +
	"JOIN (SELECT %f AS lat, %f AS lng, %d AS radius, %f AS distance_unit) AS p ON 1=1 " +
	"WHERE z.lat " +
	"BETWEEN p.lat  - (p.radius / p.distance_unit) AND p.lat  + (p.radius / p.distance_unit) AND z.lng " +
	"BETWEEN p.lng - (p.radius / (p.distance_unit * COS(RADIANS(p.lat)))) AND p.lng + (p.radius / (p.distance_unit * COS(RADIANS(p.lat)))) " +
	") AS d " +
	"WHERE distance <= radius " +
	"ORDER BY distance ASC " +
	"LIMIT %d ;"

var distanceUnits = map[string]float32{
	"km": 111.045,
	"mi": 69.0,
}

/*FindLocations to locate nearby party*/
func (structure *EfimeralObject) FindLocations(radius int, distanceUnit string, limit int) []resource.EfimeralObjectInterface {

	statement := fmt.Sprintf(harvesine, Schema, Table, structure.Lat, structure.Lng, radius, distanceUnits[distanceUnit], limit)
	ObjectIo := new(resource.EfimeralObject).GetIoInterface()
	efimeralObjects, queryResultConnection := ObjectIo.RunQuery(statement)

	var parties []resource.EfimeralObjectInterface
	for efimeralObjects.Next() {
		currentParty := new(EfimeralObject)
		err := efimeralObjects.Scan(&currentParty.Lat, &currentParty.Lng, &currentParty.Code, &currentParty.Dist)
		if err != nil {
			log.Fatal(err)
		}

		parties = append(parties, currentParty)
	}
	efimeralObjects.Close()
	queryResultConnection.Close()
	return parties
}

/*PostalToCoords converts a zip/postal code into coordinates*/
func (structure *EfimeralObject) PostalToCoords(postalcode string) (EfimeralObject, error) {
	statement := fmt.Sprintf(coordsQuery, postalcode)
	client := new(efimeralhttp.EfimeralObject)
	geolocation, err := client.Get(statement)
	if err != nil {
		return *structure, err
	}

	type Coordinates struct {
		Latitude  float64 `json:"Latitude"`
		Longitude float64 `json:"Longitude"`
	}

	type Fields struct {
		Fields Coordinates `json:"fields"`
	}

	type Records struct {
		Records []Fields `json:"records"`
	}

	var cont Records
	json.Unmarshal([]byte(geolocation), &cont)
	if len(cont.Records) == 0 {
		return *structure, &errors.EmptyPostalToCordsError{}
	}

	structure.Lat = float64(cont.Records[0].Fields.Latitude)
	structure.Lng = float64(cont.Records[0].Fields.Longitude)
	return *structure, err
}

/*GenerateTableQuery will generate a postgresql table*/
func (structure *EfimeralObject) GenerateTableQuery() string {
	statement := fmt.Sprintf(`-- Drop table
DROP TABLE IF EXISTS %s.%s CASCADE;
CREATE TABLE IF NOT EXISTS %s.%s (
	code varchar NOT NULL,
	lat float8 NOT NULL,
	lng float8 NOT NULL,
	CONSTRAINT party_pk PRIMARY KEY (code)
);
`, Schema, Table, Schema, Table)
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
		err := efimeralObjects.Scan(&objectEfimeralObject.Code, &objectEfimeralObject.Lat, &objectEfimeralObject.Lng)
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
		"Code":    structure.Code,
		"Payload": structure.Payload,
		"Lat":     structure.Lat,
		"Lng":     structure.Lng,
		"Dist":    structure.Dist,
	}
	return objectData
}

/*GetObjectInterface to get interface*/
func (structure *EfimeralObject) GetObjectInterface() resource.EfimeralObjectInterface {
	var efimeralObject resource.EfimeralObjectInterface
	efimeralObject = &EfimeralObject{}
	return efimeralObject
}
