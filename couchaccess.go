//
// CouchAccess project couchaccess.go
// https://github.com/motaz/couchaccess
// Developed by Motaz, Code 2019
//
package couchaccess

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rhinoman/couchdb-go"
)

type Couchdatabase struct {
	connection *couchdb.Database
	server     string
	username   string
	database   string
	password   string
}

func (adata *Couchdatabase) GetConnection() *couchdb.Database {

	return adata.connection
}

func ConnectToDB(server, username, password, database string) (db *Couchdatabase, err error) {

	db = new(Couchdatabase)
	db.server = server
	db.username = username
	db.password = password
	db.database = database

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

func Insert(db *Couchdatabase, theDoc interface{}, id string) (rev string, insertedID string, err error) {

	if id == "" {
		id = GetNewID()
	}
	insertedID = id
	rev, err = db.connection.Save(theDoc, id, "")

	return

}

func Update(db *Couchdatabase, theDoc interface{}, id string, rev string) (newrev string, err error) {

	newrev, err = db.connection.Save(theDoc, id, rev)

	return

}

func GetOnlyName(filename string) string {

	sep := string(os.PathSeparator)
	println(strings.Count(filename, sep))
	if strings.Count(filename, sep) > 1 {
		filename = filename[strings.LastIndex(filename, sep)+1:]
		println(filename)
	}
	return filename
}

func UploadAttachment(db *Couchdatabase, filename string, fileContents *bufio.Reader, id string, rev string) (arev string, err error) {

	sep := string(os.PathSeparator)

	idx := strings.LastIndex(filename, ".")
	aext := ""
	if idx > 0 {
		aext = filename[idx:]
	}
	if strings.Contains(filename, sep) {
		filename = filename[strings.LastIndex(filename, sep)+1:]

	}
	arev, err = db.connection.SaveAttachment(id, rev, filename, aext, fileContents)
	return
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

// CallView
func CallView(db *Couchdatabase, designname string, viewname string, params string) (result []byte) {

	timeout := time.Duration(30 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	aurl := "http://" + db.server + ":5984/" + db.database + "/_design/" +
		designname + "/_view/" + viewname
	if params != "" {
		aurl += "?" + params
	}

	response, err := client.Get(aurl)
	if err != nil {
		fmt.Println("The HTTP request failed with error :" + err.Error())
	} else {

		result, _ = ioutil.ReadAll(response.Body)

	}
	return
}
