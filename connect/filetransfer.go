package connect

import (
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

// FileTransfer FTP protocol handler
type FileTransfer struct {
	IP       string
	Username string
	Password string
}

// GetFile retrieves a file from the VPS
func (f FileTransfer) GetFile(conn *ftp.ServerConn, fileName string) {
	r, _ := conn.Retr(fileName)
	outFile, _ := os.Create(fileName)

	defer outFile.Close()
	io.Copy(outFile, r)
}

// Connect Initiates a connection with the host
func (f FileTransfer) Connect() *ftp.ServerConn {
	c, _ := ftp.Dial(f.IP, ftp.DialWithTimeout(5*time.Second))
	c.Login(f.Username, f.Password)
	return c
}
