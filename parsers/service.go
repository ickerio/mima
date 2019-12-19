package parsers

import (
	"errors"
	"fmt"
	"os"
)

type Command interface {
	Run() error
}

type Execute struct {
	data string
}

func (e Execute) Run() error {
	return nil
}

type GetFile struct {
	source      string
	destination string
}

func (g GetFile) Run() error {
	return nil
}

type PutFile struct {
	source      string
	destination string
}

func (p PutFile) Run() error {
	return nil
}

type Service struct {
	commands []Command
}

func GetService(fileName string) (Service, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return Service{}, errors.New("Could not find config file")
	}

	fmt.Println(f)

	return Service{}, nil
}
