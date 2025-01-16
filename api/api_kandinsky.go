package api_kandinsky

import (
	env "AI_image_generator/basic/env_handler"
	"AI_image_generator/basic/fetch"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type Params struct {
	Query string
}

type Model struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Version float64 `json:"version"`
	Type    string  `json:"type"`
}

type ModelList struct {
	List []Model
}

type ResponseText2image struct {
	Uuid   string `json:"uuid"`
	Status string `json:"status"`
}

var URL_API string = "https://api-key.fusionbrain.ai/key/api/v1"

type RequestGenerate struct {
	Type           string `json:"type"`
	NumImages      int    `json:"numImages"`
	Width          int    `json:"width"`
	Height         int    `json:"height"`
	GenerateParams Params `json:"generateParams"`
}

func GetHeadersAuth() []fetch.Headers {
	return []fetch.Headers{
		{
			Key:   "X-Key",
			Value: "Key " + env.GetEnvParam("APP_API_KEY_FUSIONBRAIN_PUBLIC"),
		},
		{
			Key:   "X-Secret",
			Value: "Secret " + env.GetEnvParam("APP_API_KEY_FUSIONBRAIN_SECRET"),
		},
	}
}

func getModelsList() ModelList {
	u := fmt.Sprintf("%s/models", URL_API)
	headers := GetHeadersAuth()

	b := fetch.Get(u, headers)
	r := ModelList{}
	err := json.Unmarshal(b, &r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in Unmarhal method getModelsList (url: %s): %v\n", u, err)
	}

	return r
}

func GenerateText2image(config RequestGenerate, ModelId int) ResponseText2image {
	url := fmt.Sprintf("%s/text2image/run", URL_API)

	data := struct {
		ModelId int             `json:"model_id"`
		Params  RequestGenerate `json:"params"`
	}{
		ModelId,
		config,
	}

	body := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(body).Encode(data)

	if err != nil {
		fmt.Println("Error in Encode payload in Text2image method!")
	}

	b := fetch.Post(url, "application/json", body)
	r := ResponseText2image{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in Unmarhal method getModelsList (url: %s): %v\n", "url", err)
	}

	return r
}

func ApiKandinsky() {
	responseModelsList := getModelsList()

	var ModelId int
	if len(responseModelsList.List) > 0 {
		ModelId = responseModelsList.List[0].Id
	}

	fmt.Println(fmt.Sprintf("Response in get models list: %v", responseModelsList))

	prompt := "Придумай логотип для приложения Quato Quick. Оно основано на разделении общего чека, например в ресторане, не нескольких человек."
	config := RequestGenerate{
		Type:      "GENERATE",
		NumImages: 1,
		Width:     256,
		Height:    256,
		GenerateParams: Params{
			Query: prompt,
		},
	}

	if ModelId != 0 {
		responseGenerateText2image := GenerateText2image(config, ModelId)

		fmt.Println(fmt.Sprintf("Response in Text2image: %v\n", responseGenerateText2image))
	} else {
		fmt.Println(fmt.Sprintf("ModelId is zero!"))
	}

}
