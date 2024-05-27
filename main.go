package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/cloudfunctionsv2"
	"github.com/pulumi/pulumi/sdk/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type goldResponse struct {
	Timestamp      uint64  `json:"timestamp"`
	Price          float64 `json:"price"`
	Price_gram_24k float64 `json:"price_gram_24k"`
	Price_gram_22k float64 `json:"price_gram_22k"`
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		function, err := cloudfunctionsv2.NewFunction(ctx, "gold-alert", &cloudfunctionsv2.FuntionArgs{
			Name:        pulumi.String("gold-alert"),
			Location:    pulumi.String("eu-west1"),
			Description: pulumi.String("Gold alerts on a polling api call"),
			BuildConfig: &cloudfunctionsv2.FunctionBuildConfigArgs{
				Runtime: pulumi.String("go"),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})

	godotenv.Load()
	response, err := getPrice()
	if err != nil {
		panic(err)
	}
	g := goldResponse{}
	if err := json.Unmarshal(response, &g); err != nil {
		panic(err)
	}

	if g.Price_gram_24k < 55 {
		client := twilio.NewRestClient()

		params := &api.CreateMessageParams{}
		params.SetBody(fmt.Sprintf("Gold Price per gram: %f.4", g.Price_gram_24k))
		params.SetFrom("+447446970386")
		params.SetTo("+447391605016")

		resp, err := client.Api.CreateMessage(params)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if resp.Sid != nil {
				fmt.Println(*resp.Sid)
			} else {
				fmt.Println(resp.Sid)
			}
		}
	}

}

func getPrice() ([]byte, error) {
	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.goldapi.io/api/XAU/GBP/", nil)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("x-access-token", os.Getenv("ACCESS_TOKEN"))
	req.Header.Set("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
