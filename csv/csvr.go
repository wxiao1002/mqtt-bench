package csv

import (
	"os"

	"github.com/gocarina/gocsv"
)

type Rc struct {
	ClientId string
	Username string
	Password string
}

func Reader(path string) ([]*Rc, error) {
	clientsFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer clientsFile.Close()
	clients := []*Rc{}
	if err := gocsv.UnmarshalFile(clientsFile, &clients); err != nil { // Load clients from file
		panic(err)
	}
	return clients, nil
}
