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

	m["path"] = pathParts[0]

	return m
}

func getSettings(path string, service settingsRetrievalService) map[string]interface{} {
	return service.GetSettings(path)
}
