package api

import (
	env "AI_image_generator/basic/env_handler"
	"fmt"
)

func Api() string {
	fmt.Println("Get env api key: ", env.GetEnvParam("APP_API_KEY_GETIMG"))

	return "Api"
}
