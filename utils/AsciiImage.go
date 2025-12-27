package utils

import (
	"image"
	"log"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/qeesung/image2ascii/convert"
)

type AsciiIamge struct {
	CollectionName string
	Index          int
	Art            string
}

func GenerateAsciiImage(url string, collectionName string, index int, width int, height int) tea.Cmd {

	return func() tea.Msg {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal("error generating image response from get ", collectionName, index, url, err)
			return nil
		}
		defer res.Body.Close()

		img, _, err := image.Decode(res.Body)
		if err != nil {
			log.Fatal("error generating image decode ", collectionName, index, url, err)
			return nil
		}

		convertOptions := convert.DefaultOptions
		convertOptions.FixedWidth = width
		convertOptions.FixedHeight = height
		convertOptions.Colored = true
		convertOptions.Ratio = 0.5 // Best for thumbnails

		converter := convert.NewImageConverter()
		asciiArt := converter.Image2ASCIIString(img, &convertOptions)

		// D. Return with Collection Name
		return AsciiIamge{
			CollectionName: collectionName,
			Index:          index,
			Art:            asciiArt,
		}

	}

}
