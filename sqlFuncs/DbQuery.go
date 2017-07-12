package sqlFuncs

import "os/exec"

func DbQuery(dbName string, finalPath string, sqlSuffix string) {

	sqlConn()

	db.Query("CREATE DATABASE IF NOT EXISTS " + dbName)

	cmd := exec.Command("docker", "exec", "-i", "thedrum_mysql-yingyang_1", "mysql", "-uroot", "-pcarnyx007", dbName)
	catCmd := exec.Command("cat", finalPath+sqlSuffix)

	cmd.Stdin, _ = catCmd.StdoutPipe()

	_ = catCmd.Start()
	_ = cmd.Run()
	_ = catCmd.Wait()
}
