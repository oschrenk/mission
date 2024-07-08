package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/itchyny/gojq"
)

func parseFocusFile(path string) string {
	query, err := gojq.Parse("try .data[].storeAssertionRecords[].assertionDetails.assertionDetailsModeIdentifier")
	if err != nil {
		log.Println(err)
	}
	bytes, err := os.ReadFile(path) // just pass the file name
	if err != nil {
		log.Println(err)
	}
	dynamic := make(map[string]interface{})
	err = json.Unmarshal(bytes, &dynamic)
	if err != nil {
		log.Println(err)
	}
	iter := query.Run(dynamic) // or query.RunWithContext
	var value string
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			log.Println(err)
		}
		value = fmt.Sprintf("%s", v)
	}
	return value
}

func GetFocus(path string) string {
	return parseFocusFile(path)
}
