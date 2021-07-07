package cli

import (
	"errors"
	"testing"
)

func createGetFlagsMock(pf string, kf string, user string, usePwd bool) func() (*string, *string, *string, *bool) {
	return func() (*string, *string, *string, *bool) {
		return &pf, &kf, &user, &usePwd
	}
}


func createGetPassword(pwd string, err error) func() (string, error)  {
	return func() (string, error) {
		return pwd, err
	}
}

func createFileExist(names ...string) func(string) bool {
	return func(name string) bool {
		for _, n := range names {
			if n == name {
				return true
			}
		}
		return false
	}
}


func TestGetUserInputPfNotExist(t *testing.T) {
	playbookFilename := "playbook.yaml"

	cli := Cli{getFlags: createGetFlagsMock(playbookFilename, "", "root", false), fileExist: createFileExist()}

	_, err := cli.GetUserInput()

	if err == nil {
		t.Fatalf("GetUserInput() should return error")
	}
}

func TestGetUserInputKfNotExist(t *testing.T) {
	pf := "playbook.yaml"
	kf := "id_rsa"

	cli := Cli{getFlags: createGetFlagsMock(pf, kf, "root", false), fileExist: createFileExist(pf)}

	_, err := cli.GetUserInput()

	if err == nil {
		t.Fatalf("GetUserInput() should return error")
	}
}

func TestGetUserInputKfExist(t *testing.T) {
	pf := "playbook.yaml"
	kf := "id_rsa"
	user := "pat"

	cli := Cli{getFlags: createGetFlagsMock(pf, kf, user, false), fileExist: createFileExist(pf, kf)}

	res, err := cli.GetUserInput()

	if err != nil {
		t.Fatalf("GetUserInput() shouldn't return error")
	}

	if res.Filename != pf {
		t.Fatalf("%s was return, but %s was expected", res.Filename, pf)
	}

	if res.Auth.User != user {
		t.Fatalf("%s was return, but %s was expected", res.Auth.User, user)
	}

	if res.Auth.Password != "" {
		t.Fatalf("%s should be empty", res.Auth.Password)
	}

	if res.Auth.PrivateKeyFilename != kf {
		t.Fatalf("%s was return, but %s was expected", res.Auth.PrivateKeyFilename, kf)
	}
}

func TestGetUserInputPwdError(t *testing.T) {
	pf := "playbook_2.yaml"
	user := "pat"

	cli := Cli{
		getFlags: createGetFlagsMock(pf, "", user, true),
		fileExist: createFileExist(pf),
		getPassword: createGetPassword("", errors.New("error")),
	}

	_, err := cli.GetUserInput()

	if err == nil {
		t.Fatalf("GetUserInput() should return error")
	}
}

func TestGetUserInputPwdValid(t *testing.T) {
	pf := "playbook.yaml"
	user := "pat"
	pwd := "pwd123"

	cli := Cli{
		getFlags: createGetFlagsMock(pf, "", user, true),
		fileExist: createFileExist(pf),
		getPassword: createGetPassword(pwd, nil),
	}

	res, err := cli.GetUserInput()

	if err != nil {
		t.Fatalf("GetUserInput() shouldn't return error")
	}

	if res.Filename != pf {
		t.Fatalf("%s was return, but %s was expected", res.Filename, pf)
	}

	if res.Auth.User != user {
		t.Fatalf("%s was return, but %s was expected", res.Auth.User, user)
	}

	if res.Auth.PrivateKeyFilename != "" {
		t.Fatalf("%s should be empty", res.Auth.PrivateKeyFilename)
	}

	if res.Auth.Password != pwd {
		t.Fatalf("%s was return, but %s was expected", res.Auth.Password, pwd)
	}

}