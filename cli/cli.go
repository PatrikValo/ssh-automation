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

func (auth *Auth) GetAuthMethod() (ssh.AuthMethod, error) {
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

func fileExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func getFlags() (*string, *string, *string, *bool) {
	playbookFilename := flag.String("f", Filename, "source file of the tasks")
	keyFilename := flag.String("k", SSHKey, "ssh private key for machines")
	user := flag.String("u", User, "user login")
	usePassword := flag.Bool("p", false, "usage of password")
	flag.Parse()
	return playbookFilename, keyFilename, user, usePassword
}

type Cli struct {
	getPassword func() (string, error)
	fileExist   func(string) bool
	getFlags    func() (*string, *string, *string, *bool)
}

func (cli *Cli) GetUserInput() (UserFlags, error) {
	pf, kf, user, usePwd := cli.getFlags()

	if !cli.fileExist(*pf) {
		return UserFlags{}, errors.New(*pf + " file doesn't exist")
	}

	if !*usePwd {
		if !cli.fileExist(*kf) {
			return UserFlags{}, errors.New(*kf + " file doesn't exist")
		}
		return UserFlags{Filename: *pf, Auth: Auth{User: *user, Password: "", PrivateKeyFilename: *kf}}, nil
	}

	pwd, err := cli.getPassword()
	if err != nil {
		return UserFlags{}, err
	}

	return UserFlags{Filename: *pf, Auth: Auth{User: *user, Password: pwd, PrivateKeyFilename: ""}}, nil
}

func CreateCli() Cli {
	return Cli{getPassword: getPassword, fileExist: fileExist, getFlags: getFlags}
}
