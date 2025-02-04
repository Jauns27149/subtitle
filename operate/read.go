package operate

import (
	"github.com/Jauns27149/subtitle/model"
	"github.com/Jauns27149/subtitle/tools"
	"os"
	"strings"
)

func ReadSrt(path string) model.Srt {
	file, err := os.ReadFile(path)
	tools.CheckErr(err)
	blocks := strings.Split(string(file), "\n\n")
	srt := model.Srt{make([]model.Block, 0, len(blocks))}
	for i, b := range blocks {
		b = strings.TrimSpace(b)
		lines := strings.Split(b, "\n")
		if len(lines) < 2 {
			continue
		}
		subtitle := lines[2:]
		tools.CheckErr(err)

		block := model.Block{
			Sequence: i + 1,
			Time:     lines[1],
			Subtitle: strings.Join(subtitle, ""),
		}
		srt.Blocks = append(srt.Blocks, block)
	}
	return srt
}
