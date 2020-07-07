package sampledata

import (
	"efimeral/packages/party"
	"efimeral/packages/payload"
	"efimeral/packages/resource"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

var payloadNames = []string{
	"24K Gold", "303 OG Kush", "3 Kings (Three Kings)", "3X Crazy", "501st OG", "707 Headband",
	"818 Headband", "9lb Hammer", "Abracadabra", "Abusive OG", "Acapulco Gold", "AC/DC",
	"Ace of Spades", "A-Dub", "Afghan Big Bud", "Afghan Diesel", "Afghan Haze", "Afghani",
	"Afghan Kush", "Afghan Skunk", "Afgooey", "Afwreck", "Agent Tangie", "AK-47", "AK-48",
	"Alaskan Thunder Fuck", "Albert Walker", "Alexander the Grape", "Alien Cookies", "Alien Dawg",
	"Alien Kush", "Alien OG", "Alien Rock Candy", "Allen Wrench", "Allkush", "Aloha Berry",
	"Alpha Blue", "Alpha OG", "Alpine Blue", "Ambrosia", "American Pie", "Amherst Sour Diesel", "Amnesia",
	"Amnesia Haze", "Ancient OG", "Anesthesia", "Animal Cookies", "Animal Mints", "Apollo 13", "Appalachia",
	"Appalachian Power", "Apple Fritter", "Apple Jack", "Apple Sherbet", "Argyle", "Asian Fantasy", "Atmosphere",
	"Atomical Haze", "Atomic Northern Lights", "Aurora ", "Avalon", "Ayahuasca Purple", "Azure Haze", "Bacio Gelato",
	"Bakerstreet", "Banana Cream OG", "Banana Kush", "Banana Punch", "Banana Sherbet", "Banana Split", "Banjo", "Barba Piss Cat",
}

/*EfimeralObject to hold our sampledata data*/
type EfimeralObject struct {
	Type              string
	PartyCollection   []resource.EfimeralObjectInterface
	PayloadCollection []resource.EfimeralObjectInterface
}

/*SetSampleData will help us generate our calls on randomized data*/
func (structure *EfimeralObject) SetSampleData(center resource.EfimeralObjectInterface, radiusInMiles int, parties int, payloads int) *EfimeralObject {

	// Let's create payloads
	structure.PayloadCollection = structure.GetSamplePayload(payloads)

	// Let's create parties
	structure.PartyCollection = structure.GetSampleParties(center, radiusInMiles, parties)

	return structure
}

/*IsPayloadMappedToParty to prevent duplicate mapping of party to payload*/
func (structure *EfimeralObject) IsPayloadMappedToParty(payloadCollection []resource.EfimeralObjectInterface, payload resource.EfimeralObjectInterface) bool {
	payloadData := payload.GetObjectData()
	for _, payloadInParty := range payloadCollection {
		payloadInPartyData := payloadInParty.GetObjectData()
		if payloadData["Code"] == payloadInPartyData["Code"] {
			return true
		}
	}
	return false
}

/*Radius generator*/
func (structure *EfimeralObject) Radius(center resource.EfimeralObjectInterface, radiusInMiles int) resource.EfimeralObjectInterface {
	rand.Seed(time.Now().UnixNano())
	ramdomMilesAwayFromCenter := rand.Intn(radiusInMiles)
	newParty := new(party.EfimeralObject)
	if ramdomMilesAwayFromCenter == 0 {
		ramdomMilesAwayFromCenter = rand.Intn(radiusInMiles + 1)
	}

	degree := float64(1) / float64(69)
	delta := float64(ramdomMilesAwayFromCenter) * degree
	centerData := center.GetObjectData()
	centerEfimeralObject := new(party.EfimeralObject)
	mapstructure.Decode(centerData, &centerEfimeralObject)
	if (rand.Intn(100-1)+1)%2 == 0 {
		newParty.Lat = centerEfimeralObject.Lat + delta
	} else {
		newParty.Lat = centerEfimeralObject.Lat - delta
	}
	if (rand.Intn(100-1)+1)%2 == 0 {
		newParty.Lng = centerEfimeralObject.Lng + delta
	} else {
		newParty.Lng = centerEfimeralObject.Lng - delta
	}

	newPartyObject := newParty.GetObjectInterface()
	mapstructure.Decode(newParty.GetObjectData, &newParty)

	newPartyObject = newParty

	return newPartyObject
}

/*GetSamplePayload to get a collection of payloads*/
func (structure *EfimeralObject) GetSamplePayload(pageSize int) []resource.EfimeralObjectInterface {
	var samplePayloadCollection []resource.EfimeralObjectInterface
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < pageSize; i++ {
		payload := new(payload.EfimeralObject)
		payload.Name = payloadNames[rand.Intn(len(payloadNames))]

		for {
			goodCode := true
			payload.Code = structure.generateCode()
			for _, item := range samplePayloadCollection {
				if item.GetObjectData()["Code"] == payload.Code {
					goodCode = false
					break
				}
			}
			if goodCode {
				break
			}
		}

		payload.Quantity = rand.Intn(30)
		payload.Price = rand.Float32()
		samplePayloadCollection = append(samplePayloadCollection, payload)
	}
	return samplePayloadCollection
}

/*GetSampleParties to get a collection of parties*/
func (structure *EfimeralObject) GetSampleParties(center resource.EfimeralObjectInterface, radiusInMiles int, pageSize int) []resource.EfimeralObjectInterface {
	var samplePartyCollection []resource.EfimeralObjectInterface

	for i := 0; i < pageSize; i++ {
		var samplePartyPayloadCollection []resource.EfimeralObjectInterface
		newParty := structure.Radius(center, (rand.Intn(radiusInMiles-1) + 1))
		newPartyData := newParty.GetObjectData()
		newPartyEfimeralObject := new(party.EfimeralObject)
		mapstructure.Decode(newPartyData, &newPartyEfimeralObject)

		for {
			goodCode := true
			newPartyEfimeralObject.Code = structure.generateCode()
			for _, item := range samplePartyCollection {
				if item.GetObjectData()["Code"] == newPartyEfimeralObject.Code {
					goodCode = false
					break
				}
			}
			if goodCode {
				break
			}
		}

		for payloadLimit := 0; payloadLimit < 10; payloadLimit++ {
			rand.Seed(time.Now().UnixNano())
			randomPayload := structure.PayloadCollection[rand.Intn(len(structure.PayloadCollection))]
			randomPayloadData := randomPayload.GetObjectData()
			randomPayloadEfimeralObject := new(payload.EfimeralObject)
			mapstructure.Decode(randomPayloadData, &randomPayloadEfimeralObject)
			insertRandomPayload := true

			for _, payloadValue := range samplePartyPayloadCollection {
				if randomPayloadEfimeralObject.Name == payloadValue.GetObjectData()["Name"] {
					insertRandomPayload = false
					break
				}
			}

			randomPayloadEfimeralObject.Party = newPartyEfimeralObject.Code

			if insertRandomPayload {
				samplePartyPayloadCollection = append(samplePartyPayloadCollection, randomPayloadEfimeralObject)
			}
		}

		newPartyEfimeralObject.Payload = samplePartyPayloadCollection
		samplePartyCollection = append(samplePartyCollection, newPartyEfimeralObject)
	}
	return samplePartyCollection
}

func (structure *EfimeralObject) generateCode() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}

/*WriteToDb to insert sample data into db*/
func (structure *EfimeralObject) WriteToDb() *EfimeralObject {

	var statements []string

	partyStatement := "insert into efimeral.party (lat, lng, code) values "
	var partyValues []string

	payloadStatement := "insert into efimeral.payload (\"name\", code, party, quantity, price) values "
	var payloadValues []string

	for _, pty := range structure.PartyCollection {
		partyEfimeralObject := new(party.EfimeralObject)
		mapstructure.Decode(pty.GetObjectData(), &partyEfimeralObject)
		fmt.Printf("|%8s|%12f|%12f|\n", partyEfimeralObject.Code, partyEfimeralObject.Lat, partyEfimeralObject.Lng)
		partyValue := fmt.Sprintf("(%f, %f, '%s')", partyEfimeralObject.Lat, partyEfimeralObject.Lng, partyEfimeralObject.Code)
		partyValues = append(partyValues, partyValue)
		for _, pload := range partyEfimeralObject.Payload {
			payloadEfimeralObject := new(payload.EfimeralObject)
			mapstructure.Decode(pload.GetObjectData(), &payloadEfimeralObject)
			fmt.Printf("\t|%30s|%10s|%10s|%4d|%10f\n", payloadEfimeralObject.Name, payloadEfimeralObject.Code, payloadEfimeralObject.Party, payloadEfimeralObject.Quantity, payloadEfimeralObject.Price)
			payloadValue := fmt.Sprintf("('%s', '%s', '%s', %d, %f)", payloadEfimeralObject.Name, payloadEfimeralObject.Code, payloadEfimeralObject.Party, payloadEfimeralObject.Quantity, payloadEfimeralObject.Price)
			payloadValues = append(payloadValues, payloadValue)
		}
	}

	partyStatement += strings.Join(partyValues, ",")
	statements = append(statements, partyStatement)
	payloadStatement += strings.Join(payloadValues, ",")
	statements = append(statements, payloadStatement)
	statement := strings.Join(statements, "; ")

	ObjectIo := new(resource.EfimeralObject).GetIoInterface()
	ObjectIo.RunQuery(statement)

	return structure
}
