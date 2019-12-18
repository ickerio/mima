package services

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type Minecraft struct {
	SavesDir string
	Name     string
	Host     string
	Username string
	Password string
}

func (m Minecraft) Start() error {

	conn, err := Connect(m.Host, m.Username, m.Password)
	if err != nil {
		return err
	}
	defer conn.Close()

	savePath := path.Join(m.SavesDir, m.Name)
	if fileExists(savePath) {
		fmt.Print(savePath + " exists!\n")
		var files []string

		err := filepath.Walk(savePath, func(currentPath string, info os.FileInfo, err error) error {
			files = append(files, currentPath)
			return nil
		})
		if err != nil {
			return err
		}

		if len(files) == 0 {

			if err := m.createServer(&conn); err != nil {
				return err
			}
		} else {

			if err := m.uploadServer(&conn); err != nil {
				return err
			}
		}
	} else {
		fmt.Print(savePath + " does not exists! Creating directory...\n")
		if err := os.MkdirAll(savePath, os.ModePerm); err != nil {
			return err
		}
		if err := m.createServer(&conn); err != nil {
			return err
		}
	}

	return nil
}

func (m Minecraft) Stop() error {

	conn, err := Connect(m.Host, m.Username, m.Password)
	if err != nil {
		return err
	}

	if err := m.stopServerJar(&conn); err != nil {
		return err
	}

	if err := m.backupServer(&conn); err != nil {
		return err
	}

	conn.Close()
	return nil
}

func (m Minecraft) uploadServer(conn *Connection) error {
	return nil
}

func (m Minecraft) createServer(conn *Connection) error {

	conn.Execute(`mkdir server&&cd server&&sudo apt-get update&&sudo apt-get install default-jre&&wget https://launcher.mojang.com/v1/objects/3737db93722a9e39eeada7c27e7aca28b144ffa7/server.jar&&screen&&sudo java -Xmx3584M -jar server.jar nogui`)

	// conn.Execute("mkdir server")
	// conn.Execute("cd server")

	// conn.Execute("sudo apt-get update")
	// conn.Execute("sudo apt-get install default-jre")
	// conn.Execute("wget https://launcher.mojang.com/v1/objects/3737db93722a9e39eeada7c27e7aca28b144ffa7/server.jar")
	// conn.Execute("screen")
	// conn.Execute("sudo java -Xmx3584M -jar server.jar nogui")

	return nil
}

func (m Minecraft) startServerJar(conn *Connection) error {
	return nil
}

func (m Minecraft) backupServer(conn *Connection) error {
	return nil
}

func (m Minecraft) stopServerJar(conn *Connection) error {
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
