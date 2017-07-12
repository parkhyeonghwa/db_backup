package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vFlagR/db_backup/sqlFuncs"
	"github.com/vFlagR/db_backup/sshFuncs"
	"github.com/vFlagR/db_backup/unzip"
)

func main() {
	var srcPath string
	var filename string
	var dbName string
	var dstPath = "/home/carnyx/Desktop/"
	var finalPath = dstPath + filename
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
		dbName = "award_chip_shop"
	case "1":
		srcPath = "/backup/aws-londinium/"
		filename = "jobs-" + yDate.Format("20060102")
		dbName = "jobs"
	case "2":
		srcPath = "/backup/aws-londinium/"
		filename = "awards_prod-" + yDate.Format("20060102")
		dbName = "awards_prod"
	case "3":
		srcPath = "/backup/aws-robocod-db/"
		filename = yDate.Format("2006_01_02") + "_recommended_agencies"
	case "4":
		srcPath = "/backup/aws-yingyang/"
		filename = "thedrum_prod-" + yDate.Format("20060102")
		dbName = "thedrum_prod"
	default:
		panic("You didn't choose something from the list")
	}

	// Creates a connection to the backup database
	sshFuncs.SshConn()

	// Opens a session then copies the file from the server to the host machine
	sshFuncs.SshQuery(srcPath, filename, gzipSuffix, dstPath)

	// Unzips the copied .gz file then deletes the .gz version
	unzip.UnpackGzipFile(finalPath+gzipSuffix, finalPath+sqlSuffix)

	// Executes the sql import using the relevant database from the switch statement.
	sqlFuncs.DbQuery(dbName, finalPath, sqlSuffix)

	fmt.Println("I'm all done chief")

}
