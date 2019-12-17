package services

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/pkg/sftp"
)

func CopyFile(srcPath string, destPath string, host string, username string, password string) error {

	if !strings.Contains(host, ":") {
		host += ":22"
	}

	clientConfig, err := auth.PasswordKey(username, password, ssh.InsecureIgnoreHostKey())
	if err != nil {
		fmt.Printf("Error creating client config.\n")
		fmt.Printf(err.Error())
		return err
	}

	conn, err := ssh.Dial("tcp", host, &clientConfig)
	if err != nil {
		fmt.Printf("Error creating ssh connection.\n")
		fmt.Printf(err.Error())
		return err
	}

	sftp, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Printf("Error creating sftp client.\n")
		fmt.Printf(err.Error())
		return err
	}
	defer sftp.Close()

	srcFile, err := os.Open(srcPath)
	if err != nil {
		fmt.Printf("Error opening source file " + srcPath + "\n")
		fmt.Printf(err.Error())
		return err
	}
	defer srcFile.Close()

	destFile, err := sftp.Create(destPath)
	if err != nil {
		fmt.Printf("Error creating destination file " + destPath + "\n")
		fmt.Printf(err.Error())
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		fmt.Printf("Error copying bytes " + destPath + "\n")
		fmt.Printf(err.Error())
		return err
	}

	return nil
}

func SaveFile(srcPath string, destPath string, host string, username string, password string) error {

	if !strings.Contains(host, ":") {
		host += ":22"
	}

	clientConfig, err := auth.PasswordKey(username, password, ssh.InsecureIgnoreHostKey())
	if err != nil {
		fmt.Printf("Error creating client config.\n")
		fmt.Printf(err.Error())
		return err
	}

	conn, err := ssh.Dial("tcp", host, &clientConfig)
	if err != nil {
		fmt.Printf("Error creating ssh connection.\n")
		fmt.Printf(err.Error())
		return err
	}

	sftp, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Printf("Error creating sftp client.\n")
		fmt.Printf(err.Error())
		return err
	}
	defer sftp.Close()

	srcFile, err := sftp.Open(srcPath)
	if err != nil {
		fmt.Printf("Error opening source file " + srcPath + "\n")
		fmt.Printf(err.Error())
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		fmt.Printf("Error creating destination file " + destPath + "\n")
		fmt.Printf(err.Error())
		return err
	}
	defer destFile.Close()

	srcFile.WriteTo(destFile)

	return nil
}
