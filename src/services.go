package main

import (
	"log"
	"strings"
)

type SettingsRetrievalService interface {
	GetSettings(path string) map[string]interface{}
}

type CompositeSettingsRetrievalService struct {
}

func (s *CompositeSettingsRetrievalService) GetSettings(path string) map[string]interface{} {
	pathParts := strings.Split(path, "/")

	database, err := getMongoDatabase()
	if err != nil {
		log.Fatal(err)
	}

	docid, defaultData, err := getDefaultSettings(database, pathParts)

	if err != nil {
		log.Fatal(err)
	}

	userid := "7b803e2d-ee0e-4213-a025-9db732bcbb2e"
	// userid := "ad2ea197-310a-4832-940c-2935bd6fa511"

	userData, err := getUserSettings(database, docid, userid)

	if err != nil {
		log.Fatal(err)
	}

	m, err := mergeSettings(defaultData, userData)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

func getSettings(path string, service SettingsRetrievalService) map[string]interface{} {
	return service.GetSettings(path)
}
