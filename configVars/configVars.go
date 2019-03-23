package configVars

import (
	"encoding/json"
	"os"
)

// Add any config variables here
type Configuration struct {
	GITHUB_PASS string
}

var Config Configuration

func InitConfigVars() {

	file, err := os.Open("./configVars/config.json")
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(file).Decode(&Config)
	if err != nil {
		panic(err)
	}

}
