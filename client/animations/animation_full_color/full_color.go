package animation_full_color

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
	"technikflg.com/dmxToProjector/artnet"
)

type DMXColor struct {
	preset_animation.Animation
	DynamicConfig *DynamicConfig

	artnetnode *artnet.Artnet
}
type DynamicConfig struct {
	UseDMX   bool                    `parameter:"UseDMX,default=true"`
	Color    animation_helpers.Color `parameter:"Color,default=#ffffff"`
	Universe int                     `parameter:"Universe,default=0,min=0"`
	Address  int                     `parameter:"Address,default=1,min=1,max=502"`
}
type DMXColorGenerator struct {
}

type DMXState struct {
	Intensity uint8
	Red       uint8
	Green     uint8
	Blue      uint8
}

func (data *DMXState) FromDMX(dmx [512]byte, startAddr int) {
	if startAddr < 1 {
		startAddr = 1
	} else if startAddr > 512-5 {
		startAddr = 512 - 5
	}
	startAddr -= 1
	data.Intensity = dmx[startAddr]
	data.Red = dmx[startAddr+1]
	data.Green = dmx[startAddr+2]
	data.Blue = dmx[startAddr+3]
}

func NewDMXColorGenerator() *DMXColorGenerator {
	return &DMXColorGenerator{}
}
func (g *DMXColorGenerator) Unload() {
}
func (g *DMXColorGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {

	DMXColor := DMXColor{
		DynamicConfig: &DynamicConfig{},
		artnetnode:    artnet.GetArtnet("10.144.92.3/24", "MH-animation-node"),
	}

	(&DMXColor).Configure(config)
	return &DMXColor
}
func (a *DMXColor) Render(data *[]float64, dt float64) {
	if a.DynamicConfig == nil {
		return
	}
	if !a.DynamicConfig.UseDMX {
		c := a.DynamicConfig.Color.RGBA
		c.A = 255
		rl.ClearBackground(c)
		return
	}
	dmxData := a.artnetnode.GetDMXData(a.DynamicConfig.Universe)
	dmxState := DMXState{}
	dmxState.FromDMX(dmxData, a.DynamicConfig.Address)

	//BUG: For some reason ClearBackground expects values 0-100 for A
	rl.ClearBackground(color.RGBA{R: dmxState.Red, G: dmxState.Green, B: dmxState.Blue, A: uint8(float64(dmxState.Intensity) / 2.55)})
}
func (a *DMXColor) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
}
