package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		UserName string `yaml:"username"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	}
	Cookie struct {
		Key   string `yaml:"key"`
		Block string `yaml:"block"`
	}
	Sendgrid struct {
		ApiKey string `yaml:"api_key"`
	}
	AWS struct {
		AWSRegion          string `yaml:"aws_region"`
		AWSAccessKey       string `yaml:"aws_access_key"`
		AWSSecretAccessKey string `yaml:"aws_secret_access_key"`
		BucketName         string `yaml:"bucket_name"`
	}
}

func DecodeFile(fileName string) (Config, error) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error opening config connection")
		fmt.Println(err.Error())
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Println("Error Marshling yaml")
		fmt.Println(err.Error())
	}
	return config, err
}
