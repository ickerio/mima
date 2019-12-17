package services

import (
	"io"
	"os"
	"strings"

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

	sshClient, err := ssh.Dial("tcp", host, &clientConfig)
	if err != nil {
		return connection, err
	}
	connection.SSHClient = sshClient

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

func (conn Connection) CopyFile(srcPath string, destPath string) error {

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

func (conn Connection) SaveFile(srcPath string, destPath string) error {

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

	srcFile.WriteTo(destFile)

	return nil
}
