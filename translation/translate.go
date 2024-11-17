package translation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
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

func PdfTrans(path string) {
	vips.Startup(nil)
	defer vips.Shutdown()

	file, err := vips.NewImageFromFile(path)
	checkError(err)
	pages := file.Pages()
	qps := 3
	pageChan := make(chan int, pages)
	for page := range pages {
		pageChan <- page
	}
	trans := ReadYaml()
	group := sync.WaitGroup{}
	group.Add(qps)
	sumdst := make(map[int]string)
	close(pageChan)
	for range qps {
		go func() {
			for {
				if page, ok := <-pageChan; ok {
					pictrans := Pictrans(trans, path, page)
					sumdst[page] = pictrans.Data.SumDst
				} else {
					break
				}
			}
			group.Done()
		}()
	}
	group.Wait()
	fileName := path[strings.LastIndex(path, "/")+1 : strings.LastIndex(path, ".")]
	md, err := os.Create("interpret/" + fileName + ".md")
	checkError(err)
	for k := range len(sumdst) {
		v := "# " + strconv.Itoa(k) + "\n" + sumdst[k] + "\n\n\n"
		_, err = md.Write([]byte(v))
		checkError(err)
	}
	err = md.Close()
	checkError(err)
}

func Pictrans(t Translation, filePath string, page int) PictransRes {
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
		checkError(err)
	}

	params := vips.NewImportParams()
	params.Page.Set(page)
	image, err := vips.LoadImageFromFile(filePath, params)
	checkError(err)
	jpeg, _, err := image.ExportJpeg(&vips.JpegExportParams{})
	checkError(err)
	buffer := bytes.NewBuffer(jpeg)
	index := strings.LastIndex(filePath, "/")
	formFile, err := writer.CreateFormFile("image", filePath[index+1:])
	checkError(err)
	_, err = io.Copy(formFile, buffer)
	checkError(err)
	err = writer.Close()
	checkError(err)

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
		fmt.Println(string(body))
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

func checkError(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
