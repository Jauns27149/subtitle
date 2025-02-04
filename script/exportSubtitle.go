package main

import (
	"fmt"
	"github.com/Jauns27149/subtitle/tools"
	"os"
	"os/exec"
	"strings"
)

func main() {
	open, err := os.Open("./")
	tools.CheckErr(err)
	defer open.Close()
	names, err := open.Readdirnames(-1)
	tools.CheckErr(err)
	for _, name := range names {
		if ok := strings.HasSuffix(name, "mkv"); ok {
			fmt.Println(name)
			cmd := exec.Command("/usr/bin/mkvextract", "tracks", name, "11:"+name[:len(name)-len(".mkv")]+"srt")
			_, err = cmd.CombinedOutput()
			tools.CheckErr(err)
		}
	}
	fmt.Println("done")
}
