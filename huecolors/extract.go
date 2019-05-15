package huecolors

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // will register the jpg decode for image.Decode()
	_ "image/png"  // will register the png decode for image.Decode()
	"sort"

	"github.com/EdlinOrg/prominentcolor"
)

var (
	prominentcolorsParams = []int{
		prominentcolor.ArgumentAverageMean | prominentcolor.ArgumentNoCropping,
		prominentcolor.ArgumentAverageMean,
		prominentcolor.ArgumentNoCropping,
		prominentcolor.ArgumentDefault,
	}
)

// CIExyY represents a color in the CIE xyY colorspace
type CIExyY struct {
	X         float64
	Y         float64
	Luminance float64
}

// GetHueColors returns the x (nbColors) main colors from imgData (jpg or png)
func GetHueColors(nbColors int, imgData []byte) (colors []CIExyY, params string, err error) {
	// Decode data as image
	img, _, err := image.Decode(bytes.NewBuffer(imgData))
	if err != nil {
		err = fmt.Errorf("can't decode data as image: %v", err)
		return
	}
	// Create each set
	var RGBcolors []prominentcolor.ColorItem
	possibilities := make(ColorSets, len(prominentcolorsParams))
	for index, param := range prominentcolorsParams {
		// Extract main colors with current params
		RGBcolors, err = prominentcolor.KmeansWithAll(nbColors, img, param,
			prominentcolor.DefaultSize, prominentcolor.GetDefaultMasks())
		if err != nil {
			err = fmt.Errorf("prominent colors extraction failed with params '%s': %v", genPColorParamsString(param), err)
			return
		}
		// Add them to the list
		possibilities[index] = NewColorSet(RGBcolors, param)
	}
	// Get the set with the most differents colors
	sort.Sort(sort.Reverse(possibilities))
	// Build the answer
	params = possibilities[0].GetPColorParamsString()
	colors = make([]CIExyY, nbColors)
	for index, color := range possibilities[0].GetColorfullSet() {
		colors[index].X, colors[index].Y, colors[index].Luminance = color.Xyy()
	}
	return
}
