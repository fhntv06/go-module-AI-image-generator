package api_kandinsky

import (
	env "AI_image_generator/basic/env_handler"
	"AI_image_generator/basic/fetch"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Params struct {
	Query string `json:"query"`
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

type CheckGeneration struct {
	Uuid             string `json:"uuid"`
	Status           string `json:"status"`
	ErrorDescription string `json:"errorDescription"`
	Result           struct {
		Files    []string `json:"files"`
		Censored bool     `json:"censored"`
	}
}

type ResponsePipelines struct {
	Id            string      `json:"id"`
	Name          string      `json:"name"`
	NameEn        string      `json:"nameEn"`
	Description   string      `json:"description"`
	DescriptionEn string      `json:"descriptionEn"`
	Tags          []string    `json:"tags"`
	Version       float64     `json:"version"`
	ImagePreview  interface{} `json:"imagePreview"`
	Status        string      `json:"status"`
	Type          string      `json:"type"`
	CreatedDate   string      `json:"createdDate"`
	LastModified  string      `json:"lastModified"`
}
type ResponsePiplineRun struct {
	Uuid   string `json:"uuid"`
	Status int    `json:"status"`
}

func getHeadersAuth() []fetch.Headers {
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
	headers := getHeadersAuth()

	b := fetch.Get(u, headers)
	r := ModelList{}
	err := json.Unmarshal(b, &r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in Unmarhal method getModelsList (url: %s): %v\n", u, err)
	}

	return r
}

func generateImageFromPipeline(config RequestGenerate, pipelineId string) (string, error) {
	u := fmt.Sprintf("%s/pipeline/run", URL_API)
	headers := getHeadersAuth()

	data := struct {
		PipelineId string          `json:"pipeline_id"`
		Params     RequestGenerate `json:"params"`
	}{
		pipelineId,
		config,
	}

	body := bytes.NewBuffer([]byte{})
	err := json.NewEncoder(body).Encode(data)

	if err != nil {
		return "", err
	}

	b := fetch.Post(u, body, headers)

	fmt.Printf("Response from generate Image From Pipeline: %s\n", string(b))

	if len(b) == 0 {
		return "", errors.New("empty response from generate Image From Pipeline")
	}

	r := ResponsePiplineRun{}

	err = json.Unmarshal(b, &r)
	if err != nil {
		return "", err
	}

	fmt.Printf("Response generate Image From Pipeline data: %v\n", r)

	return r.Uuid, nil
}

func generateText2image(config RequestGenerate, ModelId int) ResponseText2image {
	u := fmt.Sprintf("%s/text2image/run", URL_API)
	headers := getHeadersAuth()

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

	b := fetch.Post(u, body, headers)
	r := ResponseText2image{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in Unmarhal method getModelsList (url: %s): %v\n", "url", err)
	}

	return r
}

// Получение списка доступный моделей
func getPipeline() (string, error) {
	u := fmt.Sprintf("%s/pipelines", URL_API)
	headers := getHeadersAuth()

	b := fetch.Get(u, headers)

	if len(b) == 0 {
		return "", errors.New("empty response from getPipeline")
	}

	r := []ResponsePipelines{}
	err := json.Unmarshal(b, &r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in Unmarhal method GetPipeline (url: %s): %v\n", u, err)
		return "", err // или panic(err), если нужно прервать выполнение
	}

	if len(r) == 0 {
		fmt.Println("Empty response, no pipelines found")
		return "", errors.New("no pipelines found")
	}

	fmt.Printf("Response pipelines data: %v\n", r)

	return r[0].Id, nil
}

func checkGeneration(uuid string, attempts int, delay int) ([]string, error) {
	u := fmt.Sprintf("%s/pipeline/status/%s", URL_API, uuid)
	headers := getHeadersAuth()

	b := fetch.Get(u, headers)

	if len(b) == 0 {
		fmt.Println("Empty response from checkGeneration")
		return []string{}, errors.New("empty response from check Generation")
	}

	r := CheckGeneration{}

	err := json.Unmarshal(b, &r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in Unmarhal method GetPipeline (url: %s): %v\n", u, err)
	}

	fmt.Printf("Response check generation data: %v\n", r)

	if r.Status == "DONE" {
		return r.Result.Files, nil
	}

	attempts -= 1

	// Pause execution for 2 seconds
	time.Sleep(10 * time.Second)

	if attempts <= 0 {
		return checkGeneration(uuid, attempts, delay)
	}

	return []string{}, errors.New("check generation failed")
}

func ApiKandinsky() {
	prompt := "Придумай логотип для приложения Quata Quick. Оно основано на разделении общего чека, например в ресторане, не нескольких человек."
	config := RequestGenerate{
		Type:      "GENERATE",
		NumImages: 1,
		Width:     1024,
		Height:    1024,
		GenerateParams: Params{
			Query: prompt,
		},
	}

	pipelineId, err := getPipeline()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in getPipeline: %v\n", err)
		return
	}
	if pipelineId == "" {
		fmt.Println("No pipeline found")
		return
	}

	fmt.Printf("Get pipelineId: %v\n", pipelineId)

	uuid, err := generateImageFromPipeline(config, pipelineId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in generate ImageFrom Pipeline: %v\n", err)
		return
	}
	if uuid == "" {
		fmt.Println("No UUID returned from generate Image From Pipeline")
		return
	}

	fmt.Printf("Get uuid: %v\n", uuid)

	files, err := checkGeneration(uuid, 10, 10)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in check Generation: %v\n", err)
		return
	}

	fmt.Printf("Response files: %v\n", files)

	return

	responseModelsList := getModelsList()

	var ModelId int
	if len(responseModelsList.List) > 0 {
		ModelId = responseModelsList.List[0].Id
	}

	fmt.Println(fmt.Sprintf("Response in get models list: %v", responseModelsList))

	if ModelId != 0 {
		responseGenerateText2image := generateText2image(config, ModelId)

		fmt.Println(fmt.Sprintf("Response in Text2image: %v\n", responseGenerateText2image))
	} else {
		fmt.Println(fmt.Sprintf("ModelId is zero!"))
	}

}
