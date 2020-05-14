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

type Dimension struct {
	Width  int
	Height int
}

type UneImage struct {
	Name      string
	Dimension Dimension
	Size      int
	Path      string
}

type BySize []UneImage

func (a BySize) Len() int           { return len(a) }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySize) Less(i, j int) bool { return a[i].Size < a[j].Size }

func kbSize(kb int64) int {
	return int(float32(kb) * 0.001)
}

func main() {
	dir := "."
	if len(os.Args) > 1 && len(os.Args[1]) > 0 {
		dir = os.Args[1]
	}
	files, _ := ioutil.ReadDir(dir)
	imgArray := []UneImage{}

	for _, imgFile := range files {
		if reader, err := os.Open(filepath.Join(dir, imgFile.Name())); err == nil {
			im, _, err := image.DecodeConfig(reader)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
				continue
			}
			//fmt.Printf("%s: %dx%d\n", imgFile.Name(), im.Width, im.Height)

			// Obtenir les infos fichier
			imstat, err := reader.Stat()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", imgFile.Name(), err)
				continue
			}

			// La taille en octets de l'image
			imsize := imstat.Size()
			// On l'ajoute dans notre tableau de type UneImage
			imgArray = append(imgArray, UneImage{
				Name: imgFile.Name(),
				Dimension: Dimension{
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
	sort.Sort(BySize(imgArray))
	for _, item := range imgArray {
		fmt.Printf("%s:\n\tSize=%dkB\n\tResolution=%dx%d\n\tPath=%s\n", item.Name, item.Size, item.Dimension.Width, item.Dimension.Height, item.Path)
	}
}
