module technikflg.com/dmxToProjector

go 1.23

toolchain go1.23.6

require (
	github.com/gen2brain/raylib-go/raylib v0.0.0-20241207114308-a9ad86d5018c
	github.com/gordonklaus/portaudio v0.0.0-20230709114228-aafa478834f5
	github.com/gorilla/websocket v1.5.3
	github.com/mjibson/go-dsp v0.0.0-20180508042940-11479a337f12
)

require github.com/sirupsen/logrus v1.9.3 // indirect

require (
	github.com/ebitengine/purego v0.7.1 // indirect
	github.com/jsimonetti/go-artnet v0.0.0-20251001161948-ff57fcafff73
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/sys v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/gordonklaus/portaudio => ./portaudio-local

replace github.com/jsimonetti/go-artnet => ./artnet/go-artnet
