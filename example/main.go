package main

import (
	"fmt"
	"io/ioutil"

	"github.com/q231950/sputnik/sputnik"
)

func main() {
	fmt.Println("hello sputnik")

	payload := `{
    "operations": [
        {
            "operationType": "create",
            "record": {
                "recordType": "Shelve",
                "fields": {
                    "title": {
                        "value": "panda panda 🐯🐯"
                    }
                }
            }
        }
    ]
}`

	response, _ := sputnik.Post("records/modify", payload)

	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}
