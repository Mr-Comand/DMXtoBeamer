package animation_MH

import (
	"fmt"
	"image/color"

	"technikflg.com/dmxToProjector/animations/animation_helpers"
	"technikflg.com/dmxToProjector/animations/preset_animation"
	"technikflg.com/dmxToProjector/artnet"
)

const (
	addressesNeeded = 10
	shapeCount      = 4
)

type MH struct {
	preset_animation.Animation
	DynamicConfig *DynamicConfig

	artnetnode *artnet.Artnet
	particle   animation_helpers.Particle
}

type DynamicConfig struct {
	Universe int  `parameter:"Universe,default=0,min=0"`
	Address  int  `parameter:"Address,default=1,min=1,max=502"`
	XYMode   bool `parameter:"XYMode,default=false"`
}

type MhState struct {
	Intensity uint8
	Red       uint8
	Green     uint8
	Blue      uint8
	Size      uint8
	Pan       uint16
	Tilt      uint16
	Shape     animation_helpers.Shape
	Hollow    bool
	LineWidth uint8
}

func (data *MhState) FromDMX(dmx [512]byte, startAddr int) {
	if startAddr < 1 {
		startAddr = 1
	} else if startAddr > 512-addressesNeeded {
		startAddr = 512 - addressesNeeded
	}
	startAddr -= 1
	data.Intensity = dmx[startAddr]
	data.Red = dmx[startAddr+1]
	data.Green = dmx[startAddr+2]
	data.Blue = dmx[startAddr+3]
	data.Size = dmx[startAddr+4]
	data.Pan = uint16(dmx[startAddr+5])<<8 | uint16(dmx[startAddr+6])
	data.Tilt = uint16(dmx[startAddr+7])<<8 | uint16(dmx[startAddr+8])

	shapeVal := dmx[startAddr+9]

	if shapeVal <= shapeCount {
		data.Hollow = false
	} else {
		data.Hollow = true
		shapeVal -= shapeCount
		data.LineWidth = shapeVal % 51
		shapeVal /= 51
	}
	switch shapeVal {
	case 0:
		data.Shape = animation_helpers.Circle
	case 1:
		data.Shape = animation_helpers.Square
	case 2:
		data.Shape = animation_helpers.Rectangle
	case 3:
		data.Shape = animation_helpers.Triangle
	case 4:
		data.Shape = animation_helpers.Snowflake
	default:
		data.Shape = animation_helpers.Circle
	}
}

type MHGenerator struct {
}

func NewMHGenerator() *MHGenerator {

	return &MHGenerator{}
}

func (g *MHGenerator) Unload() {
}

func (g *MHGenerator) Create(config preset_animation.AnimationParameters) preset_animation.AnimationInterface {
	MH := MH{
		DynamicConfig: &DynamicConfig{},
		artnetnode:    artnet.GetArtnet("10.8.8.1/24", "MH-animation-node"),
		particle:      animation_helpers.Particle{Shape: animation_helpers.Circle, Size: 50, Position: animation_helpers.Position{X: 500, Y: 500}},
	}
	MH.Configure(config)
	return &MH
}

func (a *MH) Render(data *[]float64, dt float64) {
	dmxData := a.artnetnode.GetDMXData(a.DynamicConfig.Universe)

	if len(dmxData) < a.DynamicConfig.Address {
		fmt.Printf("MH Render: Not enough DMX data for address %d\n", a.DynamicConfig.Address)
		return
	}

	mhState := MhState{}
	mhState.FromDMX(dmxData, a.DynamicConfig.Address)

	a.particle.Color = color.RGBA{R: mhState.Red, G: mhState.Green, B: mhState.Blue, A: mhState.Intensity}
	a.particle.Shape = mhState.Shape
	a.particle.Hollow = mhState.Hollow
	a.particle.LineWidth = mhState.LineWidth
	a.particle.Size = float64(mhState.Size) * 2.0

	if a.DynamicConfig.XYMode {
		x := (-float64(mhState.Pan)/65535.0+1)*1500.0 - 250.0
		y := (-float64(mhState.Tilt)/65535.0+1)*1500.0 - 250.0
		a.particle.Position = animation_helpers.Position{X: x, Y: y}
	} else {
		angle := float64(mhState.Pan)/65535.0*360.0 - 180.0
		dmxY := (-float64(mhState.Tilt)/65535.0+1)*1500.0 - 250.0
		x, y := animation_helpers.RotatePoint(500, dmxY, 500, 500, angle)
		a.particle.Position = animation_helpers.Position{X: float64(x), Y: float64(y)}

	}
	a.particle.Draw()
}

// Configure allows updating the configuration for the MH (e.g., new font path or text)
func (a *MH) Configure(config preset_animation.AnimationParameters) {
	preset_animation.Parse(config, a.DynamicConfig)
	if a.DynamicConfig.Address < 1 {
		a.DynamicConfig.Address = 1
	} else if a.DynamicConfig.Address > 512-addressesNeeded {
		a.DynamicConfig.Address = 512 - addressesNeeded
	}

}
func (a *MH) Unload() {}
