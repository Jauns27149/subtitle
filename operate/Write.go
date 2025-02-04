package operate

import (
	"github.com/Jauns27149/subtitle/model"
	"github.com/Jauns27149/subtitle/tools"
	"os"
	"strconv"
	"strings"
)

func WriteSrt(srt model.Srt, path string) {
	s := strings.Split(path, "/")
	file, err := os.Create("tran" + s[len(s)-1])
	tools.CheckErr(err)
	defer file.Close()
	for i, b := range srt.Blocks {
		temp := []string{strconv.Itoa(i), b.Time, b.Subtitle, "\n"}
		block := strings.Join(temp, "\n")
		_, err = file.WriteString(block)
		tools.CheckErr(err)
	}
}
