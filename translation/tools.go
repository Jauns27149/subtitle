package translation

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Translation struct {
	Api struct {
		Pictrans    string `yaml:"pictrans"`
		AccessToken string `yaml:"access_token"`
	} `yaml:"api"`
	AK string `yaml:"ak"`
	SK string `yaml:"sk"`
}

func ReadYaml() Translation {
	file, err := os.ReadFile("translation/config.yaml")
	if err != nil {
		panic(err)
	}
	transfer := Translation{}
	err = yaml.Unmarshal(file, &transfer)
	if err != nil {
		log.Fatal(err)
	}
	return transfer
}
