package translation

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type PictransRes struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
	Data      struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Content []struct {
			Src       string `json:"src"`
			Dst       string `json:"dst"`
			Rect      string `json:"rect"`
			LineCount int    `json:"lineCount"`
			Points    []struct {
				X int `json:"x"`
				Y int `json:"y"`
			} `json:"points"`
			PasteImg string `json:"pasteImg"`
		} `json:"content"`
		SumSrc   string `json:"sumSrc"`
		SumDst   string `json:"sumDst"`
		PasteImg string `json:"pasteImg"`
	} `json:"data"`
}

func Pictrans(t Translation, filePath string) PictransRes {
	u, err := url.Parse(t.Api.Pictrans)
	if err != nil {
		log.Fatal(err)
	}
	u.RawQuery = url.Values{
		"access_token": {GetAccessToken(t)},
	}.Encode()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	param := map[string]string{
		"v":    "3",
		"from": "en",
		"to":   "zh",
	}
	for k, v := range param {
		err = writer.WriteField(k, v)
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	if filePath == "" {
		log.Fatal("File path is empty")
	}
	index := strings.LastIndex(filePath, "/")
	formFile, err := writer.CreateFormFile("image", filePath[index+1:])
	if err != nil {
		log.Fatalf("CreateFormFile failed: %s", err)
	}
	_, err = io.Copy(formFile, file)
	if err != nil {
		log.Fatalf("io.Copy failed: %s", err)
	}
	err = writer.Close()
	if err != nil {
		log.Fatalf("writer.Close failed: %s", err)
	}

	req, err := http.NewRequest("POST", u.String(), &b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do failed: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll failed: %s", err)
	}
	response := PictransRes{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatalf("json.Unmarshal failed: %s", err)
	}
	return response
}

func GetAccessToken(t Translation) string {
	u, err := url.Parse(t.Api.AccessToken)
	if err != nil {
		log.Fatalf("url.Parse failed: %s", err)
	}
	v := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {t.AK},
		"client_secret": {t.SK},
	}
	u.RawQuery = v.Encode()
	ur := u.String()
	client := &http.Client{}
	req, err := http.NewRequest("POST", ur, nil)
	if err != nil {
		log.Fatalf("NewRequest failed: %s", err)
	}
	req.Header = http.Header{
		"Content-Type": {"application/x-www-form-urlencoded"},
		"Accept":       {"application/json"},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do failed: %s", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("ReadAll failed: %s", err)
	}
	s := string(body)
	i := strings.Index(s, "access_token") + 15
	return s[i : i+71]
}
