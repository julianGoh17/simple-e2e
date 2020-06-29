package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	models "github.com/julianGoh17/simple-e2e/framework/models"
	"gopkg.in/yaml.v2"
)

func main() {
	testProcedure := models.TestProcedure{}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadFile(fmt.Sprintf("%s/test.yaml", dir))
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	fmt.Println(string(body))

	err = yaml.Unmarshal([]byte(body), &testProcedure)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- test procedure:\n%v\n\n", testProcedure)
}
