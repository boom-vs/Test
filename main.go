package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
)

type CatBreeds struct {
	Breeds []struct {
		Name        string `json:"name"`
		Origin      string `json:"origin"`
		Temperament string `json:"temperament"`
	} `json:"data"`
}

func main() {
	// Получение данных о породах кошек через API
	response, err := http.Get("https://catfact.ninja/breeds")
	if err != nil {
		fmt.Println("Ошибка при получении данных:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении данных:", err)
		return
	}

	var catBreeds CatBreeds
	err = json.Unmarshal(body, &catBreeds)
	if err != nil {
		fmt.Println("Ошибка при разборе JSON:", err)
		return
	}

	breedsByOrigin := make(map[string][]string)
	for _, breed := range catBreeds.Breeds {
		breedsByOrigin[breed.Origin] = append(breedsByOrigin[breed.Origin], breed.Name)
	}

	for _, breeds := range breedsByOrigin {
		sort.Slice(breeds, func(i, j int) bool {
			return len(breeds[i]) < len(breeds[j])
		})
	}

	outFile, err := os.Create("out.json")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outFile.Close()

	outputData, err := json.Marshal(breedsByOrigin)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON:", err)
		return
	}

	_, err = outFile.Write(outputData)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
		return
	}

	fmt.Println("Данные успешно записаны в out.json")
}
