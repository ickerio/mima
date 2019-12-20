package parsers

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ickerio/mima/services"
)

type Service struct {
	SavesDir string
	Name     string
	Host     string
	Username string
	Password string
	Commands struct {
		InstallTools     []Command
		CreateExecutable []Command
		StartExecutable  []Command
		StopExecutable   []Command
	}
}

func (s Service) Create() error {
	conn, err := services.Connect(s.Host, s.Username, s.Password)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := s.installTools(&conn); err != nil {
		return err
	}

	if err := s.createExecutable(&conn); err != nil {
		return err
	}

	if err := s.startExecutable(&conn); err != nil {
		return err
	}

	return nil
}

func (s Service) Start() error {
	conn, err := services.Connect(s.Host, s.Username, s.Password)
	if err != nil {
		return err
	}
	defer conn.Close()

	savePath := path.Join(s.SavesDir, s.Name)
	if fileExists(savePath) {
		var files []string

		err := filepath.Walk(savePath, func(currentPath string, info os.FileInfo, err error) error {
			if currentPath[len(currentPath)-7:] == ".tar.gz" {
				files = append(files, currentPath)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if len(files) == 0 {
			return errors.New("No backups found in " + savePath + ". Use the 'init' command to create a new server.")
		} else {
			if err := s.installTools(&conn); err != nil {
				return err
			}
			if err := s.uploadBackup(&conn, files); err != nil {
				return err
			}
			if err := s.startExecutable(&conn); err != nil {
				return err
			}
		}
	} else {
		return errors.New("No backups found in " + savePath + ". Use the 'init' command to create a new server.")
	}

	return nil
}

func (s Service) Stop() error {
	conn, err := services.Connect(s.Host, s.Username, s.Password)
	if err != nil {
		return err
	}
	defer conn.Close()

	s.stopExecutable(&conn)

	if err := s.saveBackup(&conn); err != nil {
		return err
	}

	return nil
}

func GetService(fileName string, savesDir string, name string, ip string, password string) (Service, error) {
	service := Service{
		SavesDir: savesDir,
		Host:     ip,
		Name:     name,
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
	Run(*services.Connection) error
}

type Execute struct {
	data string
}

func (e Execute) Run(conn *services.Connection) error {
	conn.Execute(e.data)
	return nil
}

type GetFile struct {
	source      string
	destination string
}

func (g GetFile) Run(conn *services.Connection) error {
	conn.GetFile(g.source, g.destination)
	return nil
}

type PutFile struct {
	source      string
	destination string
}

func (p PutFile) Run(conn *services.Connection) error {
	conn.PutFile(p.source, p.destination)
	return nil
}

func addToService(current string, command Command, service *Service) {
	switch current {
	case "install_tools":
		service.Commands.InstallTools = append(service.Commands.InstallTools, command)
	case "create_executable":
		service.Commands.CreateExecutable = append(service.Commands.CreateExecutable, command)
	case "start_executable":
		service.Commands.StartExecutable = append(service.Commands.StartExecutable, command)
	case "stop_executable":
		service.Commands.StopExecutable = append(service.Commands.StopExecutable, command)
	}
}

func (s Service) installTools(conn *services.Connection) error {
	for _, element := range s.Commands.InstallTools {
		err := element.Run(conn)
		if err != nil {
			return err
		}
	}
	return nil //conn.Execute(`sudo apt-get -y update && sudo apt-get -y install default-jre`)
}

func (s Service) createExecutable(conn *services.Connection) error {
	for _, element := range s.Commands.CreateExecutable {
		err := element.Run(conn)
		if err != nil {
			return err
		}
	}
	return nil //conn.Execute(`cd ~/ && mkdir server && cd server && wget https://launcher.mojang.com/v1/objects/3737db93722a9e39eeada7c27e7aca28b144ffa7/server.jar && echo "eula=true" > eula.txt`)
}

func (s Service) startExecutable(conn *services.Connection) error {
	for _, element := range s.Commands.StartExecutable {
		err := element.Run(conn)
		if err != nil {
			return err
		}
	}
	return nil //conn.Execute(`cd ~/server && screen -S mima -dm sudo java -Xmx3584M -jar server.jar nogui`)
}

func (s Service) stopExecutable(conn *services.Connection) error {
	for _, element := range s.Commands.StopExecutable {
		err := element.Run(conn)
		if err != nil {
			return err
		}
	}
	return nil //conn.Execute(`screen -S mima -p 0 -X stuff "stop^M"`)
}

func (s Service) saveBackup(conn *services.Connection) error {

	if err := conn.Execute(`cd ~/ && tar -cz ./server > server.tar.gz`); err != nil {
		return err
	}

	os.MkdirAll(path.Join(s.SavesDir, s.Name), os.ModeDir)

	fileName := strconv.Itoa(int(time.Now().Unix())) + ".tar.gz"
	backupPath := "./" + path.Join(s.SavesDir, s.Name, fileName)
	fmt.Println("[MIMA] Saving backup to " + backupPath + "...")
	if err := conn.GetFile("/root/server.tar.gz", backupPath); err != nil {
		return err
	}

	return nil
}

func (s Service) uploadBackup(conn *services.Connection, files []string) error {

	// Sort files from highest to lowest.
	sort.Slice(files, func(i, j int) bool {
		iname := filepath.Base(files[i])
		inum, err := strconv.ParseInt(iname[:len(iname)-7], 10, 64)
		if err != nil {
			fmt.Println("[MIMA] Invalid file name " + iname + "\n" + err.Error())
			return false
		}

		jname := filepath.Base(files[j])
		jnum, err := strconv.ParseInt(jname[:len(jname)-7], 10, 64)
		if err != nil {
			fmt.Println("[MIMA] Invalid file name " + jname + "\n" + err.Error())
			return true
		}

		return inum > jnum
	})

	fmt.Println("[MIMA] Uploading " + files[0] + "...")
	if err := conn.PutFile(files[0], "/root/server.tar.gz"); err != nil {
		return err
	}

	if err := conn.Execute("tar xvzf /root/server.tar.gz"); err != nil {
		return err
	}

	return nil
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
