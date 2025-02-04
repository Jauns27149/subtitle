package main

import (
	"fmt"
	"github.com/Jauns27149/subtitle/model"
	"github.com/Jauns27149/subtitle/operate"
	"github.com/Jauns27149/subtitle/tools"
	"github.com/Jauns27149/subtitle/translation"
	"os"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	path := "translation/raw"
	open, err := os.Open(path)
	tools.CheckErr(err)
	names, err := open.Readdirnames(-1)

	for _, name := range names {
		tran(path + "/" + name)
	}

	//m := sync.Map{}
	//
	////m := make(map[int]*model.Block, len(srt.Blocks))
	//jobs := make(chan int, len(srt.Blocks))
	//for i, b := range srt.Blocks {
	//	m.Store(i+1, &b)
	//	//m[i+1] = &b
	//	jobs <- i + 1
	//}
	//close(jobs)
	//var wait sync.WaitGroup
	//wait.Add(8)
	//for i := 0; i < 8; i++ {
	//	go func() {
	//		for {
	//			key := <-jobs
	//			if key == 0 {
	//				wait.Done()
	//				break
	//			}
	//			temp, _ := m.Load(key)
	//			b := temp.(*model.Block)
	//			//b := m[key]
	//			des := translation.Texttrans(b.Subtitle)
	//			m[key].Subtitle = m[key].Subtitle + "\n" + des
	//		}
	//	}()
	//}
	//wait.Wait()
	//operate.WriteSrt(m)

	//file, err := os.ReadFile("translation/config.yaml")
	//checkError(err)
	//t := translation.Translation{}
	//err = yaml.Unmarshal(file, &t)
	//checkError(err)
	//
	//srt := t.Texttrans("hello")
	//fmt.Println(srt)
	//path, err := filepath.Abs("translation/raw/en.srt")
	//checkError(err)
	//file, err := os.ReadFile(path)
	//s := strings.TrimPrefix(string(file), "\ufeff")
	//checkError(err)
	//srt, err := subtitles.NewFromSRT(s)
	//checkError(err)
	//captions := srt.Captions
	//fmt.Println(captions)

	//srt := "hello,\n\nelephant!"
	//ch := translation.Texttrans(srt)
	//fmt.Println(ch)

	elapsed := time.Since(start)
	fmt.Printf("Elapsed: %v\n", elapsed.Seconds())
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func tran(path string) {
	srt := operate.ReadSrt(path)
	jobs := make(chan model.Block, len(srt.Blocks))
	for _, b := range srt.Blocks {
		jobs <- b
	}
	close(jobs)

	results := make(chan model.Block, len(srt.Blocks))
	mun := 3
	var wg sync.WaitGroup
	wg.Add(mun)
	for range mun {
		go func() {
			for {
				if b, ok := <-jobs; ok {
					des := translation.Texttrans(b.Subtitle)
					b.Subtitle = b.Subtitle + "\n" + des
					results <- b
				} else {
					wg.Done()
					break
				}
			}
		}()
	}
	wg.Wait()
	close(results)

	blocks := make([]model.Block, len(srt.Blocks))
	for b := range results {
		blocks[b.Sequence-1] = b
	}
	srt.Blocks = blocks

	operate.WriteSrt(srt, path)
}
