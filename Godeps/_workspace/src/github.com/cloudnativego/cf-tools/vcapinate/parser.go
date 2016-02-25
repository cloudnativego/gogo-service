package main

import "gopkg.in/yaml.v2"

//SpaceConfiguration holds the parsed configurations for user-provided
//and brokered services to be applied to a given space.
type SpaceConfiguration struct {
	UserProvided []UserProvidedService
	Brokered     []BrokeredService
}

//UserProvidedService holds a parsed configuration for a single
//user-provided service.
type UserProvidedService struct {
	Name        string
	SyslogURL   string
	Credentials map[interface{}]interface{}
}

//BrokeredService holds a parsed configuration for a single brokered
//service.
type BrokeredService struct {
	Name        string
	Label       string
	Plan        string
	Tags        string
	Credentials map[interface{}]interface{}
}

func parseSpaceConfiguration(config string) (spaceConfig *SpaceConfiguration, err error) {
	spaceConfig = &SpaceConfiguration{}
	err = yaml.Unmarshal([]byte(config), spaceConfig)
	return spaceConfig, err
}
