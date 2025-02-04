package model

type Srt struct {
	Blocks []Block
}
type Block struct {
	Sequence int
	Time     string
	Subtitle string
}
