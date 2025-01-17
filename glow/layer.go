package glow

import (
	"fmt"
	"image"
	"image/color"
)

type Layer struct {
	Length    uint16 `yaml:"length" json:"length"`
	Rows      uint16 `yaml:"rows" json:"rows"`
	Grid      Grid   `yaml:"grid" json:"grid"`
	Chroma    Chroma `yaml:"chroma" json:"chroma"`
	HueShift  int16  `yaml:"hue_shift" json:"hue_shift"`
	Scan      uint16 `yaml:"scan" json:"scan"`
	Begin     uint16 `yaml:"begin" json:"begin"`
	End       uint16 `yaml:"end" json:"end"`
	Rate      uint32 `yaml:"rate" json:"rate"`
	ImageName string `yaml:"image_name" json:"image_name"`

	position uint16
	first    uint16
	last     uint16
	picture  image.Image
}

func NewLayer() *Layer {
	var layer Layer
	layer.Chroma.Colors = append(layer.Chroma.Colors,
		HSV{Hue: 180, Saturation: 100, Value: 100})
	return &layer
}

func (layer *Layer) Setup(length, rows uint16,
	grid *Grid, chroma *Chroma, hueShift int16,
	scan uint16, begin uint16, end uint16) error {

	layer.Length = length
	layer.Rows = rows
	layer.Grid = *grid
	layer.Chroma = *chroma
	layer.HueShift = hueShift
	layer.Scan = scan
	layer.Begin = begin
	layer.End = end

	return layer.Validate()
}

func (layer *Layer) SetupLength(length, rows uint16) error {
	layer.Length = length
	layer.Rows = rows
	layer.Grid.SetupLength(length, rows)
	layer.Chroma.SetupLength(length, layer.HueShift)
	return layer.Validate()
}

func (layer *Layer) SetRate(rate uint32) {
	layer.Rate = rate
}

func (layer *Layer) Validate() error {
	if layer.Length == 0 {
		return fmt.Errorf("Layer.Setup zero length")
	}
	if layer.Rows == 0 {
		return fmt.Errorf("Layer.Setup zero rows")
	}
	if err := layer.Grid.SetupLength(layer.Length, layer.Rows); err != nil {
		return err
	}
	if err := layer.Chroma.SetupLength(layer.Length, layer.HueShift); err != nil {
		return err
	}
	if layer.Scan > layer.Length {
		layer.Scan = layer.Length
	}
	if layer.End == 0 {
		layer.End = 100
	}
	if layer.Rate == 0 {
		layer.Rate = 48
	}
	layer.setBounds()

	return nil
}

func (layer *Layer) setBounds() {
	ratio := func(offset, length uint16) float32 {
		if offset > 100 {
			offset %= 100
		}
		return float32(offset) / 100.0 * float32(length)
	}

	layer.first = layer.Grid.AdjustBounds(ratio(layer.Begin, layer.Length))
	layer.last = layer.Grid.AdjustBounds(ratio(layer.End, layer.Length))

	if layer.last < layer.first {
		layer.first, layer.last = layer.last, layer.first
	}
}

func (layer *Layer) Spin(light Light) {
	if layer.picture != nil {
		layer.spinImage(light)
		return
	}

	startAt := layer.first
	endAt := layer.last
	if layer.Scan > 0 {
		startAt, endAt = layer.updateScanPosition()
	}

	for i := startAt; i < endAt; i++ {
		x := layer.first + (i % (layer.last - layer.first))
		light.Set(layer.Grid.Map(x), layer.Chroma.Map(x))
	}

	layer.Chroma.UpdateColors()
}

func (layer *Layer) updateScanPosition() (startAt, endAt uint16) {
	startAt = layer.position
	endAt = layer.position + layer.Scan
	layer.position++
	if layer.position >= layer.last {
		layer.position = layer.first
	}
	return startAt, endAt
}

func (layer *Layer) MakeCode() string {
	s := fmt.Sprintf("{%d,%d,%s,%s,%d,%d,%d,%d},",
		layer.Length,
		layer.Rows,
		layer.Grid.MakeCode(),
		layer.Chroma.MakeCode(),
		layer.HueShift, layer.Scan, layer.Begin, layer.End)
	return s
}

func (layer *Layer) LoadImage(rows, cols int) (err error) {
	layer.picture, err = LoadPicPath(layer.ImageName, rows, cols)
	return
}

func (layer *Layer) spinImage(light Light) {
	if layer.picture == nil {
		return
	}
	pic := layer.picture
	b := pic.Bounds()
	for x := b.Min.X; x < b.Max.X && x < int(layer.Grid.columns); x++ {
		for y := b.Min.Y; y < b.Max.Y && y < int(layer.Rows); y++ {
			c := pic.At(x, y)
			r, g, b, a := c.RGBA()
			if a != 0 {
				light.Set(uint16(y)*layer.Rows+uint16(x),
					color.NRGBA{uint8(r), uint8(g), uint8(b), 255})
			}
		}
	}
}
