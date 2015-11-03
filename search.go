package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Tag struct {
	Name string `json:"name"`
}

func search(name string) error {
	request := fmt.Sprintf("https://api.github.com/repos/%s/tags", name)
	resp, err := http.Get(request)
	if err != nil {
		return err
	}
	var tags []Tag
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tags)
	if err != nil {
		return fmt.Errorf("Can't find %s package\n", name)
	}
	fmt.Println("latest")
	for _, tag := range tags {
		fmt.Println(tag.Name)
	}
	return nil
}
