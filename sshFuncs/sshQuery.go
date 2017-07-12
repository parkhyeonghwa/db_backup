package sshFuncs

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/sftp"
)

func SshQuery(srcPath string, filename string, gzipSuffix string, dstPath string) {

	//SshConn()
	client := SshConn()
	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	// Open the source file
	srcFile, err := sftp.Open(srcPath + filename + gzipSuffix)
	fmt.Printf("Copying %s to %s \n", filename+gzipSuffix, dstPath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dstPath + filename + gzipSuffix)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// Copy the file
	srcFile.WriteTo(dstFile)

	fmt.Printf("unzipping %s", filename+gzipSuffix)
}
