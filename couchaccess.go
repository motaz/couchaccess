// CouchAccess project couchaccess.go
package couchaccess

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/rhinoman/couchdb-go"
)

type Couchdatabase struct {
	connection *couchdb.Database
}

func ConnectToDB(server, username, password, database string) (db *Couchdatabase, err error) {

	db = new(Couchdatabase)
	var timeout = time.Duration(500 * time.Millisecond)
	conn, err := couchdb.NewConnection(server, 5984, timeout)
	if err != nil {
		println("Error: ", err.Error())
	} else {
		auth := couchdb.BasicAuth{Username: username, Password: password}
		db.connection = conn.SelectDB(database, &auth)
		err = db.connection.DbExists()
		if err != nil {
			err = conn.CreateDB(database, &auth)
			if err != nil {
				fmt.Println("Error creating database: " + err.Error())
			}
		}
	}
	return
}

func Insert(db *Couchdatabase, theDoc interface{}, id string) (rev string, err error) {

	if id == "" {
		id = GetNewID()
	}

	rev, err = db.connection.Save(theDoc, id, "")

	return

}

func Update(db *Couchdatabase, theDoc interface{}, id string, rev string) {

	rev, err := db.connection.Save(theDoc, id, rev)

	if err != nil {
		println("Error in .Save: ", err.Error())
	} else {
		println("Updated: ", rev)
	}

}

func Search(db *Couchdatabase, selector string, sort interface{}, result interface{}) (err error) {

	var selectorObj interface{}
	err = json.Unmarshal([]byte(selector), &selectorObj)
	params := couchdb.FindQueryParams{Selector: &selectorObj}
	if sort != "" {

		params.Sort = sort

	}
	err = db.connection.Find(&result, &params)
	if err != nil {
		println("Error in .Find: ", err.Error())
	}
	return
}

func GetNewID() string {

	source := rand.NewSource(time.Now().UnixNano())
	rand := rand.New(source)
	randPart := strconv.Itoa(rand.Intn(10000000))
	randID := time.Now().Format("06102150405") + "-" + randPart
	return randID

}
