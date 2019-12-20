package services

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"
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
			if err := m.installTools(&conn); err != nil {
				return err
			}
			if err := m.uploadServer(&conn, files); err != nil {
				return err
			}
			if err := m.startServerJar(&conn); err != nil {
				return err
			}
		}
	} else {
		return errors.New("No backups found in " + savePath + ". Use the 'init' command to create a new server.")
	}

	return nil
}

func (m Minecraft) Stop() error {

	conn, err := Connect(m.Host, m.Username, m.Password)
	if err != nil {
		return err
	}
	defer conn.Close()

	m.stopServerJar(&conn)

	if err := m.backupServer(&conn); err != nil {
		return err
	}

	return nil
}

func (m Minecraft) Create() error {
	conn, err := Connect(m.Host, m.Username, m.Password)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := m.installTools(&conn); err != nil {
		return err
	}

	if err := m.createServer(&conn); err != nil {
		return err
	}

	if err := m.startServerJar(&conn); err != nil {
		return err
	}

	return nil
}

func (m Minecraft) uploadServer(conn *Connection, files []string) error {

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

func (m Minecraft) createServer(conn *Connection) error {
	return conn.Execute(`cd ~/ && mkdir server && cd server && wget https://launcher.mojang.com/v1/objects/3737db93722a9e39eeada7c27e7aca28b144ffa7/server.jar && echo "eula=true" > eula.txt`)
}

func (m Minecraft) installTools(conn *Connection) error {
	return conn.Execute(`sudo apt-get -y update && sudo apt-get -y install default-jre`)
}

func (m Minecraft) startServerJar(conn *Connection) error {
	return conn.Execute(`cd ~/server && screen -S mima -dm sudo java -Xmx3584M -jar server.jar nogui`)
}

func (m Minecraft) stopServerJar(conn *Connection) error {
	return conn.Execute(`screen -S mima -p 0 -X stuff "stop^M"`)
}

func (m Minecraft) backupServer(conn *Connection) error {

	if err := conn.Execute(`cd ~/ && tar -cz ./server > server.tar.gz`); err != nil {
		return err
	}

	os.MkdirAll(path.Join(m.SavesDir, m.Name), os.ModeDir)

	fileName := strconv.Itoa(int(time.Now().Unix())) + ".tar.gz"
	backupPath := "./" + path.Join(m.SavesDir, m.Name, fileName)
	fmt.Println("[MIMA] Saving backup to " + backupPath + "...")
	if err := conn.GetFile("/root/server.tar.gz", backupPath); err != nil {
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
