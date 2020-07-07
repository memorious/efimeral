package efimeralhttp

import (
	"io/ioutil"
	"net/http"
)

/*EfimeralObject to hold our structure data*/
type EfimeralObject struct {
	Connection      string
	ConnectionError error
}

/*Get function to "GET over http*/
func (structure *EfimeralObject) Get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return url, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return url, err
	}

	return string(body), nil
}
