package main

import (
	"strings"
)

type settingsRetrievalService interface {
	GetSettings(path string) map[string]interface{}
}

type compositeSettingsRetrievalService struct {
}

func (s *compositeSettingsRetrievalService) GetSettings(path string) map[string]interface{} {
	pathParts := strings.Split(path, "/")

	m := make(map[string]interface{})

	database, err := getMongoDatabase()
	if err != nil {
		log.Fatal(err)
	}

	docid, defaultData, err := getDefaultSettings(database, path)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(docid)
	fmt.Println(defaultData)

	userid := "7b803e2d-ee0e-4213-a025-9db732bcbb2e"
	// userid := "ad2ea197-310a-4832-940c-2935bd6fa511"

	userData, err := getUserSettings(database, docid, user1id)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(userData)

	m, err := mergeSettings(defaultData, user1Data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(m)

	m["path"] = pathParts[0]

	return m
}

func getSettings(path string, service settingsRetrievalService) map[string]interface{} {
	return service.GetSettings(path)
}
