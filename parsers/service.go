package parsers

import (
	"bufio"
	"errors"
	"os"
	"strings"
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
	Create []Command
	Start  []Command
	Stop   []Command
	Backup []Command
}

func GetService(fileName string) (Service, error) {
	service := Service{}
	var current string

	f, err := os.Open(fileName)
	if err != nil {
		return Service{}, errors.New("Could not find service file")
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		components := strings.Split(strings.TrimLeft(scanner.Text(), "\t \n"), ":")

		switch components[0] {
		case "create", "start", "stop", "backup":
			current = components[0]
		case "exe":
			command := Execute{data: components[1]}
			addToService(current, command, &service)
		case "get":
			body := strings.Split(components[1], ",")
			command := GetFile{source: body[0], destination: body[1]}
			addToService(current, command, &service)
		case "put":
			body := strings.Split(components[1], ",")
			command := PutFile{source: body[0], destination: body[1]}
			addToService(current, command, &service)
		}
	}

	return service, nil
}

func addToService(current string, command Command, service *Service) {
	switch current {
	case "create":
		service.Create = append(service.Create, command)
	case "start":
		service.Start = append(service.Start, command)
	case "stop":
		service.Stop = append(service.Stop, command)
	case "backup":
		service.Backup = append(service.Backup, command)
	}
}
