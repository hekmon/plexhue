package huecolors

import (
	"strings"

	"github.com/EdlinOrg/prominentcolor"
	colorful "github.com/lucasb-eyer/go-colorful"
)

/*
	ColorSets (sortable)
*/

// ColorSets is a collection of ColorSet sortable on their distance
type ColorSets []ColorSet

func (cs ColorSets) Len() int {
	return len(cs)
}

func (cs ColorSets) Less(i, j int) bool {
	return cs[i].TotalDistance() < cs[j].TotalDistance()
}

func (cs ColorSets) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

/*
	ColorSet
*/

// NewColorSet create a metadata rich colorset from a prominentcolor set
func NewColorSet(set []prominentcolor.ColorItem, params int) ColorSet {
	// Convert colors betweeb libs
	list := make([]colorful.Color, len(set))
	for index, pcolor := range set {
		list[index] = colorful.Color{
			R: float64(pcolor.Color.R) / 255.0,
			G: float64(pcolor.Color.G) / 255.0,
			B: float64(pcolor.Color.B) / 255.0,
		}
	}
	// Compute the total distance
	var (
		i, j     int
		distance float64
	)
	for i = 0; i < len(list)-1; i++ {
		for j = i + 1; j < len(list); j++ {
			distance += list[i].DistanceCIE94(list[j])
		}
	}
	// Return the object
	return ColorSet{
		params:        params,
		paramsStr:     genPColorParamsString(params),
		totalDistance: distance,
		set:           list,
	}
}

// ColorSet is a collection of colors holding metadata like:
// prominentcolor parameters used to create the set in int and str format,
// total distance between each colors (diversity)
// and the colors theirself in colorful format allowing all type of conversions.
type ColorSet struct {
	params        int
	paramsStr     string
	totalDistance float64
	set           []colorful.Color
}

// GetColorfullSet returns the set in the colorful format
func (cs ColorSet) GetColorfullSet() []colorful.Color {
	return cs.set
}

// GetPColorParams returns the prominentcolor params used to generate the set
func (cs ColorSet) GetPColorParams() int {
	return cs.params
}

// GetPColorParamsString returns the string format of prominentcolor params used to generate the set
func (cs ColorSet) GetPColorParamsString() string {
	return cs.paramsStr
}

// GetPColorSet rebuild a prominentcolor compatible set (minus the cnt metadata which has been lost)
func (cs ColorSet) GetPColorSet() (pcolorSet []prominentcolor.ColorItem) {
	pcolorSet = make([]prominentcolor.ColorItem, len(cs.set))
	for index, color := range cs.set {
		pcolorSet[index] = prominentcolor.ColorItem{
			Color: prominentcolor.ColorRGB{
				R: uint32(color.R * 255),
				G: uint32(color.G * 255),
				B: uint32(color.B * 255),
			},
		}
	}
	return
}

// TotalDistance returns the total distance between every colors within the set
func (cs ColorSet) TotalDistance() float64 {
	return cs.totalDistance
}

func genPColorParamsString(params int) string {
	list := make([]string, 4)
	// Seed random
	if prominentcolor.IsBitSet(params, prominentcolor.ArgumentSeedRandom) {
		list[0] = "Random seed"
	} else {
		list[0] = "Kmeans++"
	}
	// Average mean
	if prominentcolor.IsBitSet(params, prominentcolor.ArgumentAverageMean) {
		list[1] = "Mean"
	} else {
		list[1] = "Median"
	}
	// LAB or RGB
	if prominentcolor.IsBitSet(params, prominentcolor.ArgumentLAB) {
		list[2] = "LAB"
	} else {
		list[2] = "RGB"
	}
	// Cropping or no cropping
	if prominentcolor.IsBitSet(params, prominentcolor.ArgumentNoCropping) {
		list[3] = "No cropping"
	} else {
		list[3] = "Cropping center"
	}
	// build str
	return strings.Join(list, ", ")
}
