package program

import (
	"testing"
)

func createProgramWithHosts(hosts []string) *Program {
	return &Program{Hosts: hosts}
}

func TestAddDefaultPortIPBasic(t *testing.T) {
	playbookFilename := "playbook.yaml"
	host := "192.168.2.1"
	hosts := []string{host}

	program := createProgramWithHosts(hosts)

	parser := Parser{Filename: playbookFilename}
	err := parser.addDefaultPort(program)

	if err != nil {
		t.Fatalf("addDefaultPort(program) should not return error")
	}

	if program.Hosts[0] != (host + ":22") {
		t.Fatalf("host should be extended with port section")
	}
}

func TestAddDefaultPortDomainBasic(t *testing.T) {
	playbookFilename := "playbook.yaml"
	host := "google.com"
	hosts := []string{host}

	program := createProgramWithHosts(hosts)

	parser := Parser{Filename: playbookFilename}
	err := parser.addDefaultPort(program)

	if err != nil {
		t.Fatalf("addDefaultPort(program) should not return error")
	}

	if program.Hosts[0] != (host + ":22") {
		t.Fatalf("host should be extended with port section")
	}
}

func TestAddDefaultPortBadFormat1(t *testing.T) {
	playbookFilename := "playbook.yaml"
	host := "http://192.1.168.25"
	hosts := []string{host}

	program := createProgramWithHosts(hosts)

	parser := Parser{Filename: playbookFilename}
	err := parser.addDefaultPort(program)

	if err == nil {
		t.Fatalf("addDefaultPort(program) should return error")
	}
}

func TestAddDefaultPortBadFormat2(t *testing.T) {
	playbookFilename := "playbook.yaml"
	host := "http://192.1.168.25:22"
	hosts := []string{host}

	program := createProgramWithHosts(hosts)

	parser := Parser{Filename: playbookFilename}
	err := parser.addDefaultPort(program)

	if err == nil {
		t.Fatalf("addDefaultPort(program) should return error")
	}
}

func TestAddDefaultPortBadFormat3(t *testing.T) {
	playbookFilename := "playbook.yaml"
	host := "192.1.168.25:"
	hosts := []string{host}

	program := createProgramWithHosts(hosts)

	parser := Parser{Filename: playbookFilename}
	err := parser.addDefaultPort(program)

	if err == nil {
		t.Fatalf("addDefaultPort(program) should return error")
	}
}

func TestAddDefaultPortCustomPort(t *testing.T) {
	playbookFilename := "playbook.yaml"
	host := "192.1.168.25:8080"
	hosts := []string{host}

	program := createProgramWithHosts(hosts)

	parser := Parser{Filename: playbookFilename}
	err := parser.addDefaultPort(program)

	if err != nil {
		t.Fatalf("addDefaultPort(program) should not return error")
	}

	if program.Hosts[0] != host {
		t.Fatalf("host should not be either extended or replaced")
	}
}

func TestAddDefaultPortMultiple(t *testing.T) {
	playbookFilename := "playbook.yaml"
	host1 := "192.1.168.25"
	host2 := "192.1.168.15"
	host3 := "192.1.168.250:8080"
	hosts := []string{host1, host2, host3}

	program := createProgramWithHosts(hosts)

	parser := Parser{Filename: playbookFilename}
	err := parser.addDefaultPort(program)

	if err != nil {
		t.Fatalf("addDefaultPort(program) should not return error")
	}

	if program.Hosts[0] != (host1 + ":22") {
		t.Fatalf("host should be extended with port section")
	}

	if program.Hosts[1] != (host2 + ":22") {
		t.Fatalf("host should be extended with port section")
	}

	if program.Hosts[2] != host3 {
		t.Fatalf("host should not be extended with port section")
	}
}
