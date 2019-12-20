package parsers

import (
	"bufio"
	"errors"
	"github.com/ickerio/mima/services"
	"os"
	"strings"
)

type Service struct {
	SavesDir string
	Name     string
	Host     string
	Username string
	Password string
	Commands struct {
		Create []Command
		Start  []Command
		Stop   []Command
		Backup []Command
	}
}

func (s Service) Create() error {
	conn, _ := services.Connect(s.Host, s.Username, s.Password)

	for _, element := range s.Commands.Create {
		err := element.Run(conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Service) Start() error {
	conn, _ := services.Connect(s.Host, s.Username, s.Password)

	for _, element := range s.Commands.Start {
		err := element.Run(conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Service) Stop() error {
	conn, _ := services.Connect(s.Host, s.Username, s.Password)

	for _, element := range s.Commands.Stop {
		err := element.Run(conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Service) Backup() error {
	conn, _ := services.Connect(s.Host, s.Username, s.Password)

	for _, element := range s.Commands.Backup {
		err := element.Run(conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetService(fileName string, savesDir string, ip string, password string) (Service, error) {
	service := Service{
		SavesDir: savesDir,
		Host:     ip,
		Username: "root",
		Password: password,
	}
	var current string

	f, err := os.Open(fileName)
	if err != nil {
		return service, errors.New("Could not find service file")
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		components := strings.Split(strings.TrimLeft(scanner.Text(), "\t \n"), ":")

		switch components[0] {
		case "create", "start", "stop", "backup":
			current = components[0]
		case "exe":
			command := Execute{data: strings.TrimLeft(components[1], "\t \n")}
			addToService(current, command, &service)
		case "get":
			body := strings.Split(strings.TrimLeft(components[1], "\t \n"), ",")
			command := GetFile{source: body[0], destination: body[1]}
			addToService(current, command, &service)
		case "put":
			body := strings.Split(strings.TrimLeft(components[1], "\t \n"), ",")
			command := PutFile{source: body[0], destination: body[1]}
			addToService(current, command, &service)
		}
	}

	return service, nil
}

type Command interface {
	Run(services.Connection) error
}

type Execute struct {
	data string
}

func (e Execute) Run(conn services.Connection) error {
	conn.Execute(e.data)
	return nil
}

type GetFile struct {
	source      string
	destination string
}

func (g GetFile) Run(conn services.Connection) error {
	conn.GetFile(g.source, g.destination)
	return nil
}

type PutFile struct {
	source      string
	destination string
}

func (p PutFile) Run(conn services.Connection) error {
	conn.PutFile(p.source, p.destination)
	return nil
}

func addToService(current string, command Command, service *Service) {
	switch current {
	case "create":
		service.Commands.Create = append(service.Commands.Create, command)
	case "start":
		service.Commands.Start = append(service.Commands.Start, command)
	case "stop":
		service.Commands.Stop = append(service.Commands.Stop, command)
	case "backup":
		service.Commands.Backup = append(service.Commands.Backup, command)
	}
}
