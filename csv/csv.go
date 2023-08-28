package csv

import (
	"encoding/csv"
	"io"
	"os"
)

type Rc struct {
	ID       string `json:"id"`
	Username string `json:"user"`
	Password string `json:"pass"`
	Type     string `json:"type"`
	Ex1      string `json:"ex1"`
	Ex2      string `json:"ex2"`
	Ex3      string `json:"ex3"`
}

func ReaderCSV(path string) ([]*Rc, error) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	clients := []*Rc{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		clients = append(clients, &Rc{Username: record[0], Password: record[1]})
	}
	return clients, nil
}
