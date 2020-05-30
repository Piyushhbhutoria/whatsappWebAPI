package main

import (
	"encoding/csv"
	"os"
	"strings"
)

func sendBulk(file string) string {
	csvFile, err := os.Open(dir + file)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, each := range csvData {
		each[0] = strings.Replace(each[0], " ", "", -1)
		if each[0] != "" {
			v := sendText{
				Receiver: each[0],
				Message:  each[1],
			}
			textChannel <- v
		}
	}

	return "Done"
}

func sendBulkImg(file string) string {
	csvFile, err := os.Open(dir + file)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = -1

	csvData, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, each := range csvData {
		if each[0] != "" {
			each[0] = strings.Replace(each[0], " ", "", -1)
			v := sendImage{
				Receiver: each[0],
				Message:  each[1],
				Image:    each[2],
			}
			imageChannel <- v
		}
	}

	return "Done"
}
