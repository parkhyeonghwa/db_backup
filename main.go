package main

import(
	"fmt"
	"os"
	"golang.org/x/crypto/ssh"
	"log"
	"github.com/pkg/sftp"
	"io/ioutil"
	"time"
	"bufio"
	"strings"
	"github.com/vFlagR/db_backup/unzip"
)

func main(){
	var srcPath string
	var filename string
	var dstPath = "/home/carnyx/Desktop/"
	var gzipSuffix = ".sql.gz"
	var sqlSuffix = ".sql"

	fmt.Println("Please select which database you would like to import:\n" +
				"0 - Old awards (chipshop)\n" +
				"1 - Jobs\n" +
				"2 - Awards Multisite\n" +
				"3 - RAR\n" +
				"4 - The Drum\n")

	reader := bufio.NewReader(os.Stdin)
	c, _ := reader.ReadString('\n')
	c = strings.TrimSpace(c)

	yDate := time.Now().AddDate(0, 0, -1)

	switch c {
		case "0":
			srcPath = "/backup/aws-arctic-fox/"
			filename = "all_arctic_fox_databases-" + yDate.Format("20060102")
		case "1":
			srcPath = "/backup/aws-londinium/"
			filename = "jobs-" + yDate.Format("20060102")
		case "2":
			srcPath = "/backup/aws-londinium/"
			filename = "awards_prod-" + yDate.Format("20060102")
		case "3":
			srcPath = "/backup/aws-robocod-db/"
			filename = yDate.Format("2006_01_02") + "_recommended_agencies"
		case "4":
			srcPath = "/backup/aws-yingyang/"
			filename = "thedrum_prod-" + yDate.Format("20060102")
		default:
			panic("You didn't choose something from the list")
	}


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

	client, err := ssh.Dial("tcp", "backup.carnyx.com:22", config)
	if err != nil {
		panic("Failed to connect: " + err.Error())
	}
	fmt.Println("Successfully connected to the backup server.")

	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	// Open the source file
	srcFile, err := sftp.Open(srcPath + filename + gzipSuffix)
	fmt.Printf("Copying %s to %s \n", filename + gzipSuffix, dstPath)
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

	var finalPath = dstPath + filename

	fmt.Printf("unzipping %s", filename + gzipSuffix)
	unzip.UnpackGzipFile(finalPath + gzipSuffix, finalPath + sqlSuffix)
}

