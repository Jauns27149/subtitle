package main

import (
	"IELTS/translation"
	"fmt"
)

func main() {
	yaml := translation.ReadYaml()
	//token := translation.GetAccessToken(yaml)
	pictrans := translation.Pictrans(yaml, "translation/english.png")
	fmt.Println(pictrans.Data.SumDst)
}
