package program

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Parser struct {
	Filename string
}

func (parser Parser) Parse() (*Program, error) {
	yamlFile, err := ioutil.ReadFile(parser.Filename)

	if err != nil {
		return nil, err
	}

	var program Program
	err = yaml.Unmarshal(yamlFile, &program)

	return &program, err
}
