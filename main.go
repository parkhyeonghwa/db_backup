package main

import(
	"fmt"
	"os"
	"golang.org/x/crypto/ssh"
	"log"
	"github.com/pkg/sftp"
	"io/ioutil"
)

func main(){

	key, err := ioutil.ReadFile("/home/carnyx/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User:            "carnyx",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	client, err := ssh.Dial("tcp", "172.16.0.224:22", config)
	if err != nil {
		panic("Failed to connect: " + err.Error())
	}
	fmt.Println("Successfully connected to ssh server.")

	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	srcPath := "/backup/aws-londinium/"
	dstPath := "/home/carnyx/Desktop/"
	filename := "awards_prod-20170510.sql.gz"

	// Open the source file
	srcFile, err := sftp.Open(srcPath + filename)
	fmt.Printf("Copying %s to %s", filename, dstPath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dstPath + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()

	// Copy the file
	srcFile.WriteTo(dstFile)
}
