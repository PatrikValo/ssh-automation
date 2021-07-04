package cli

import (
	"errors"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"os"
)

const (
	Filename string = "playbook.yaml"
	SSHKey   string = "~/.ssh/id_rsa"
	User     string = "root"
)

type UserFlags struct {
	Filename string
	Auth     Auth
}

type Auth struct {
	User               string
	Password           string
	PrivateKeyFilename string
}

func (auth Auth) GetAuthMethod() (ssh.AuthMethod, error) {
	if auth.PrivateKeyFilename == "" {
		return ssh.Password(auth.Password), nil
	}

	key, err := ioutil.ReadFile(auth.PrivateKeyFilename)
	if err != nil {
		return ssh.Password(""), errors.New("unable to read private key")
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return ssh.Password(""), errors.New("unable to parse private key")
	}

	return ssh.PublicKeys(signer), nil
}

func getPassword() (string, error) {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(0)
	fmt.Println()
	return string(password), err
}

func GetUserInput() (UserFlags, error) {
	filename := flag.String("f", Filename, "source file of the tasks")
	pkf := flag.String("k", SSHKey, "ssh private key for machines")
	user := flag.String("u", User, "user login")
	usePassword := flag.Bool("p", false, "usage of password")
	flag.Parse()

	if _, err := os.Stat(*filename); os.IsNotExist(err) {
		return UserFlags{}, err
	}

	if !*usePassword {
		if _, err := os.Stat(*pkf); os.IsNotExist(err) {
			return UserFlags{}, err
		}

		return UserFlags{Filename: *filename, Auth: Auth{User: *user, Password: "", PrivateKeyFilename: *pkf}}, nil
	}

	password, err := getPassword()
	if err != nil {
		return UserFlags{}, err
	}

	return UserFlags{Filename: *filename, Auth: Auth{User: *user, Password: password, PrivateKeyFilename: ""}}, nil
}
