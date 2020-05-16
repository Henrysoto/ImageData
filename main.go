package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type dimensions struct {
	Width  int
	Height int
}

type uneImage struct {
	Name      string
	Dimension dimensions
	Size      int
	Path      string
}

type bySize []uneImage

func (a bySize) Len() int           { return len(a) }
func (a bySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a bySize) Less(i, j int) bool { return a[i].Size < a[j].Size }

func kbSize(kb int64) int {
	return int(float32(kb) * 0.001)
}

func main() {
	dir := "."
	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		dir = os.Args[1]
	}
	files, _ := ioutil.ReadDir(dir)
	imgArray := []uneImage{}

	for _, imgFile := range files {
		if reader, err := os.Open(filepath.Join(dir, imgFile.Name())); err == nil {
			im, _, err := image.DecodeConfig(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[!] %s -> %v\n", imgFile.Name(), err)
				continue
			}

			imstat, err := reader.Stat()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
				continue
			}

			imsize := imstat.Size()
			imgArray = append(imgArray, uneImage{
				Name: imgFile.Name(),
				Dimension: dimensions{
					Width:  im.Width,
					Height: im.Height,
				},
				Size: kbSize(imsize),
				Path: filepath.Join(dir, imgFile.Name()),
			})
		} else {
			fmt.Println("Can't open file: ", err)
		}
	}
	sort.Sort(bySize(imgArray))
	for _, item := range imgArray {
		fmt.Printf("%s:\n\tSize=%dkB\n\tResolution=%dx%d\n\tPath=%s\n", item.Name, item.Size, item.Dimension.Width, item.Dimension.Height, item.Path)
	}

	fmt.Printf("Number of images: %d\n", len(imgArray))
}
