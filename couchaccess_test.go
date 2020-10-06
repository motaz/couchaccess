package couchaccess

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

type rec struct {
	Filename string
}

type DocumentSection struct {
	SectionID int
	TypeName  string
	Type      string
}

type Documents struct {
	Docs []DocumentSection
}

func TestHi(t *testing.T) {
	println("Testing..")
	db, err := ConnectToDB("localhost", "tester", "999", "codedocuments")
	if err == nil {
		sections := Documents{}
		println("Connected")
		err = Search(db, `{"Type": "sections"}`, "", &sections)
		if err != nil {
			println(err.Error())
		} else {
			fmt.Printf("%+v", sections)
		}
		var file rec

		afile, _ := os.Open("/home/motaz/Pictures/2019.jpg")
		filer := bufio.NewReader(afile)
		onlyname := GetOnlyName(afile.Name())
		file.Filename = onlyname
		rev, id, err := Insert(db, file, "")

		arev, err := UploadAttachment(db, afile.Name(), filer, id, rev)
		if err != nil {
			println("Error in Insert: " + err.Error())
		} else {
			println("Revision: " + rev)
			println("ARev: " + arev)
		}

	} else {
		println("Error: " + err.Error())
	}
	db.GetConnection()
}
