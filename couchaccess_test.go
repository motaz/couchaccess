package couchaccess

import (
	"testing"
)

type rec struct {
	Filename string
}

func TestHi(t *testing.T) {
	println("Testing..")
	db, err := ConnectToDB("localhost", "tester", "999", "my")
	if err == nil {
		println("Connected")
		var file rec
		file.Filename = "Attari"
		rev, id, err := Insert(db, file, "")
		arev, err := UploadAttachment(db, "/home/motaz/Pictures/1920px-Atari-2600-Jr-FL.jpg", id, rev)
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
