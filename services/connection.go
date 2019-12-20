package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/pkg/sftp"
)

type Connection struct {
	SSHClient  *ssh.Client
	SFTPClient *sftp.Client
}

func Connect(host string, username string, password string) (Connection, error) {

	var connection Connection

	if !strings.Contains(host, ":") {
		host += ":22"
	}

	clientConfig, err := auth.PasswordKey(username, password, ssh.InsecureIgnoreHostKey())
	if err != nil {
		return connection, err
	}

	fmt.Printf("[MIMA] Dialing VPS...")
	for i := 0; i < 36; i++ {
		sshClient, err := ssh.Dial("tcp", host, &clientConfig)
		if err == nil {
			connection.SSHClient = sshClient
			break
		}
		fmt.Printf("\r[MIMA] Dialing VPS - %v tries...", i+1)
		time.Sleep(time.Second * 10)
	}
	fmt.Println()

	if connection.SSHClient == nil {
		return connection, errors.New("Failed to establish SSH connection to VPS.")
	}

	sftpClient, err := sftp.NewClient(connection.SSHClient)
	if err != nil {
		return connection, err
	}
	connection.SFTPClient = sftpClient

	return connection, nil
}

func (conn Connection) Close() {
	conn.SFTPClient.Close()

	conn.SSHClient.Close()
}

func (conn Connection) PutFile(srcPath string, destPath string) error {

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := conn.SFTPClient.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func (conn Connection) GetFile(srcPath string, destPath string) error {

	srcFile, err := conn.SFTPClient.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = srcFile.WriteTo(destFile)
	if err != nil {
		return err
	}

	return nil
}

func (conn Connection) Execute(cmd string) error {
	session, err := conn.SSHClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var buff bytes.Buffer
	session.Stdout = &buff
	fmt.Print("[MIMA] Executing $ " + cmd + "\n")
	if err := session.Run(cmd); err != nil {
		return errors.New(err.Error() + "\n" + buff.String())
	}

	return nil
}
