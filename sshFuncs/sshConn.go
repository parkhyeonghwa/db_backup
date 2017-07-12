package sshFuncs

import (
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"
)

func SshConn() (client *ssh.Client) {
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

	client, err = ssh.Dial("tcp", "backup.carnyx.com:22", config)
	if err != nil {
		panic("Failed to connect: " + err.Error())
	}
	fmt.Println("Successfully connected to the backup server.")

	return
}
