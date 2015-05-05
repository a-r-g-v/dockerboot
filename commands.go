package main

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/cli"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/exec"
)

var Commands = []cli.Command{
	commandEnable,
	commandDisable,
	commandList,
	commandAwake,
}
var commandEnable = cli.Command{
	Name: "enable",
	Usage: `awake を実行した時に起動するコンテナを追加する. 
	使用法:dockerboot enable < container-name(or continer-id) > `,
	Description: `
`,
	Action: doEnable,
}

var commandDisable = cli.Command{
	Name: "disable",
	Usage: `awake を実行した時に起動するコンテナを削除する. 
	使用法:dockerboot disable < container-name(or continer-id) > `,
	Description: `
`,
	Action: doDisable,
}

var commandList = cli.Command{
	Name: "list",
	Usage: `awake を実行した時に起動するコンテナを一覧表示する. 
	使用法:dockerboot list`,
	Description: `
`,
	Action: doList,
}
var commandAwake = cli.Command{
	Name:  "awake",
	Usage: "登録されているコンテナを起動する",
	Description: `
`,
	Action: doAwake,
}

func debug(v ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		log.Println(v...)
	}
}

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getContainerName(c *cli.Context) string {
	args := c.Args()
	if len(args) <= 0 {
		log.Fatal("引数としてコンテナIDが必須です")
		os.Exit(1)
	}
	return args[0]
}

func doEnable(c *cli.Context) {

	container := getContainerName(c)
	db, err := sql.Open("sqlite3", "/var/lib/dockerboot.db")
	assert(err)
	stmt, err := db.Prepare("INSERT INTO containers(id) values(?)")
	assert(err)
	res, err := stmt.Exec(container)
	assert(err)
	debug(res)

	fmt.Printf("コンテナ %s を登録しました \n", container)

}

func doList(c *cli.Context) {

	db, err := sql.Open("sqlite3", "/var/lib/dockerboot.db")
	assert(err)
	rows, err := db.Query("SELECT * FROM containers")
	assert(err)
	for rows.Next() {
		var container string
		err = rows.Scan(&container)
		assert(err)
		fmt.Println(container)
	}
}

func doDisable(c *cli.Context) {
	container := getContainerName(c)

	db, err := sql.Open("sqlite3", "/var/lib/dockerboot.db")
	assert(err)
	stmt, err := db.Prepare("DELETE FROM containers where id = ?")
	assert(err)
	res, err := stmt.Exec(container)
	assert(err)
	debug(res)

	fmt.Printf("コンテナ %s を削除しました \n", container)
}

func doAwake(c *cli.Context) {
	db, err := sql.Open("sqlite3", "/var/lib/dockerboot.db")
	assert(err)
	rows, err := db.Query("SELECT * FROM containers")
	assert(err)
	for rows.Next() {
		var container string
		err = rows.Scan(&container)
		assert(err)
		out, err := exec.Command("docker", "start", container).Output()
		assert(err)
		debug(string(out))
		fmt.Printf("コンテナ %s を起動しました \n", container)
	}
}
