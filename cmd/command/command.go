package command

import (
	"bufio"
	"efimeral/packages/installer"
	"efimeral/packages/party"
	"efimeral/packages/payload"
	"efimeral/packages/resource"
	"efimeral/packages/sampledata"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/mitchellh/mapstructure"
)

/*EfimeralObject to hold payload values*/
type EfimeralObject struct {
	Radius, Parties, Payloads int
	Install, Sample, Test     bool
	Region, Item              string
}

/*Entry for I/O*/
func (structure *EfimeralObject) Entry() {

	/*Let's get our cli params*/
	structure.parseParams()

	/*Let's regenerate entity tables if "install" flag is true*/
	if structure.Install {
		fmt.Println("Rebuilding Schema Tables in DB")
		installation := new(installer.EfimeralObject)
		installation.Prepare().Install()
	}

	/*Let's generate and upsert sample data if "sample" flag is true*/
	if structure.Sample {
		fmt.Println("Generating/Inserting Sample Data")
		/*Let's generate our request location*/
		center, err := new(party.EfimeralObject).PostalToCoords(structure.Region)
		if err != nil {
			log.Fatal(err)
		}

		sampledata := new(sampledata.EfimeralObject)
		centerObject := center.GetObjectInterface()
		centerObject = &center
		sampledata.SetSampleData(centerObject, structure.Radius, structure.Parties, structure.Payloads).WriteToDb()
	}

	if structure.Test {
		fmt.Println(fmt.Sprintf("\n\tLooking for: %s\n\tAround: %s\n\tDistance: %d\n", structure.Item, structure.Region, structure.Radius))
		searchResult := structure.search(structure.Region, structure.Item, structure.Radius)
		structure.outputSearchResult(searchResult)
	}

	if !structure.Test {
		for {
			consoleReader := bufio.NewReader(os.Stdin)
			color.Green("\nWhat Do You Want (Strain)?")
			structure.Item, _ = consoleReader.ReadString('\n')
			structure.Item = strings.ToLower(structure.Item)
			structure.Item = strings.TrimSuffix(strings.ToLower(structure.Item), "\n")

			if strings.HasPrefix(structure.Item, "bye") {
				fmt.Println("\nGood bye!")
				os.Exit(0)
			}

			color.Green("Where Do You Want It (Zip)?")
			structure.Region, _ = consoleReader.ReadString('\n')
			structure.Region = strings.TrimSuffix(strings.ToLower(structure.Region), "\n")

			color.Green("Distance Around You (Miles)?")
			distance, _ := consoleReader.ReadString('\n')
			distance = strings.TrimSuffix(strings.ToLower(distance), "\n")
			distanceInt, _ := strconv.Atoi(distance)

			color.Yellow(fmt.Sprintf("\n\tLooking for: %s\n\tAround: %s\n\tDistance: %d\n", structure.Item, structure.Region, distanceInt))
			searchResult := structure.search(structure.Region, structure.Item, distanceInt)
			structure.outputSearchResult(searchResult)
		}
	}
}

func (structure *EfimeralObject) outputSearchResult(searchResult []resource.EfimeralObjectInterface) {
	if len(searchResult) > 0 {

		for _, partyValue := range searchResult {
			partyCodeValueData := partyValue.GetObjectData()
			partyCodeValueEfimeralObject := new(party.EfimeralObject)
			mapstructure.Decode(partyCodeValueData, &partyCodeValueEfimeralObject)

			for _, payloadValue := range partyCodeValueEfimeralObject.Payload {
				payloadValueData := payloadValue.GetObjectData()
				payloadValueEfimeralObject := new(payload.EfimeralObject)
				mapstructure.Decode(payloadValueData, &payloadValueEfimeralObject)
				color.Green(fmt.Sprintf("\n\tSeller: %s\n\t\tDistance: %s\n\t\tProduct: %s\n",
					payloadValueEfimeralObject.Party,
					partyCodeValueEfimeralObject.Dist,
					payloadValueEfimeralObject.Name))
			}
		}
	}
}

func (structure *EfimeralObject) search(region string, payloadQuery string, distance int) []resource.EfimeralObjectInterface {

	/*Let's generate our request location*/
	center, err := new(party.EfimeralObject).PostalToCoords(region)
	if err != nil {
		log.Fatal(err)
	}

	/*Let's find locations around request location*/
	centerLocations := center.FindLocations(distance, "mi", 10)
	PayloadObject := new(payload.EfimeralObject).GetObjectInterface()
	var payloads []resource.EfimeralObjectInterface
	var parties []resource.EfimeralObjectInterface

	for _, partyCodeValue := range centerLocations {
		partyCodeValueData := partyCodeValue.GetObjectData()
		partyCodeValueEfimeralObject := new(party.EfimeralObject)
		mapstructure.Decode(partyCodeValueData, &partyCodeValueEfimeralObject)

		filters := map[string][]string{
			"party":    {"=", "'" + fmt.Sprintf("%v", partyCodeValueData["Code"]) + "'"},
			"\"name\"": {"ilike", "'%" + payloadQuery + "%'"},
		}
		ploads := PayloadObject.GetCollection(filters)
		for _, p := range ploads {
			payloads = append(payloads, p)
		}
		partyCodeValueEfimeralObject.Payload = payloads
		parties = append(parties, partyCodeValueEfimeralObject)
	}
	return parties
}

func (structure *EfimeralObject) parseParams() {

	radiusPtr := flag.Int("radius", 1, "An Integer")
	partyPtr := flag.Int("parties", 1, "An integer")
	payloadPtr := flag.Int("payload", 1, "An integer")
	installPtr := flag.String("install", "no", "A String")
	samplePtr := flag.String("sample", "no", "A String")
	regionPtr := flag.String("region", "10001", "A String")
	testPtr := flag.String("test", "no", "A String")
	itemPtr := flag.String("item", "Afghan", "A String")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	structure.Region = *regionPtr
	structure.Item = *itemPtr

	structure.Radius = *radiusPtr
	if structure.Radius == 0 {
		structure.Radius = rand.Intn(*radiusPtr + 1)
	}

	structure.Parties = *partyPtr
	if structure.Parties == 0 {
		structure.Parties = rand.Intn(*partyPtr + 1)
	}

	structure.Payloads = *payloadPtr
	if structure.Payloads == 0 {
		structure.Payloads = rand.Intn(*payloadPtr + 1)
	}

	structure.Install = false
	if *installPtr == "yes" {
		structure.Install = true
	}

	structure.Sample = false
	if *samplePtr == "yes" {
		structure.Sample = true
	}

	structure.Test = false
	if *testPtr == "yes" {
		structure.Test = true
	}
}
