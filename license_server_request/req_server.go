package license_server_request

import (
	"encoding/json"
	"fmt"
	"github.com/GoHippo/license_checker/pkg/uuid"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

var INVALID_AUTHORIZATION = "invalid authorization"

type Request struct {
	Token    string `json:"token" validate:"required"`
	UUID     string `json:"uuid" validate:"required"`
	SoftName string `json:"soft_name" validate:"required"`
	Payload  string `json:"payload"`
}

type Response struct {
	Data string `json:"data"`

	Status string `json:"status"` // error, ok
	Error  string `json:"error"`
}

type DataServerOptions struct {
	Url      string
	Token    string
	SoftName string
	Payload  string
}

func GetDataFromServer(opt DataServerOptions) (string, error) {

	if opt.Token == "" {
		return "", fmt.Errorf("Key License is not set!")
	}

	if opt.SoftName == "" {
		return "", fmt.Errorf("Softname is not set!")
	}
	if opt.Url == "" {
		return "", fmt.Errorf("Url is not set!")
	}

	uuid, err := uuid.GetUUID()
	if err != nil {
		return "", fmt.Errorf("get uuid: %w", err)
	}

	data := Request{
		Token:    opt.Token,
		UUID:     uuid,
		SoftName: opt.SoftName,
		Payload:  opt.Payload,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("json marshal request: %w", err)
	}

	//url := `http://127.0.0.1:8045/check`

	client := &fasthttp.Client{
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	req.SetRequestURI(opt.Url)
	req.SetBody(body)
	req.Header.SetMethod("POST")

	if err = client.DoTimeout(req, res, time.Second*30); err != nil {
		errStr := strings.ReplaceAll(err.Error(), string(req.URI().Host()), "SERVER")
		return "", fmt.Errorf("server license request: %s", errStr)
	}

	//s, b, err := fasthttp.Post(body, url, nil)
	//if err != nil {
	//	return "", fmt.Errorf("server license request: %w", err)
	//}

	if res.StatusCode() != 200 {
		return "", fmt.Errorf("server license request code: %v", res.StatusCode())
	}

	var resp = Response{}
	if err = json.Unmarshal(res.Body(), &resp); err != nil {
		return "", fmt.Errorf("json unmarshal license: %w", err)
	}

	if resp.Status != "ok" {
		if resp.Error == INVALID_AUTHORIZATION {
			return "", fmt.Errorf(resp.Error)
		}

		return "", fmt.Errorf("server license err in status: %v", resp.Error)
	}

	return resp.Data, nil
}
