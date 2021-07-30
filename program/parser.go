package program

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Parser struct {
	Filename string
}

func (parser *Parser) addDefaultPort(program *Program) error {
	for i, host := range program.Hosts {
		split := strings.Split(host, ":")

		if len(split) > 2 {
			return errors.New("bad format of host ip address or domain")
		}

		if len(split) == 2 && split[1] == "" {
			program.Hosts[i] = host + "22"
		}

		if len(split) == 1 {
			program.Hosts[i] = host + ":22"
		}
	}
	return nil
}

func (parser *Parser) Parse() (*Program, error) {
	yamlFile, err := ioutil.ReadFile(parser.Filename)

	if err != nil {
		return nil, err
	}

	var program Program
	err = yaml.Unmarshal(yamlFile, &program)

	if err == nil {
		err = parser.addDefaultPort(&program)
	}

	return &program, err
}
