package window

import (
	"fmt"
	"log"
	"math"

	"github.com/gordonklaus/portaudio"
	"github.com/mjibson/go-dsp/fft"
)

var audioProcessor *AudioProcessor

// Initialize PortAudio and create a new audio stream
func InitAudioProcessor(bufferSize int) {
	// Initialize PortAudio
	err := portaudio.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize PortAudio: %v", err)
	}

	// List all available devices
	devices, err := portaudio.Devices()
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	// Print out available devices and their index
	for i, device := range devices {
		fmt.Printf("%d: %s \t%s  \t%s\n", i, device.Name, device.HostApi.Name, device.DefaultSampleRate)
	}

	// Create a new AudioProcessor instance
	audioProcessor = NewAudioProcessor(bufferSize, 1)
	audioProcessor.Start()
}

// AudioProcessor struct to manage audio capture and FFT processing
type AudioProcessor struct {
	Stream     *portaudio.Stream
	BufferSize int
	SampleRate float64
	Data       []float32
}

// NewAudioProcessor initializes the PortAudio stream and prepares the data buffer
func NewAudioProcessor(bufferSize int, deviceIndex int) *AudioProcessor {

	// Set device index (default is 0)
	devices, err := portaudio.Devices()
	if err != nil {
		log.Fatalf("Failed to get the device: %v", err)
	}
	device := devices[deviceIndex]
	fmt.Println("-------------------------------------------")
	fmt.Println("Using:                                    "+device.Name, device.DefaultSampleRate)
	// Data buffer to store audio samples
	data := make([]float32, bufferSize)
	sampleRate := device.DefaultSampleRate
	// Open the default audio stream (1 input channel, 0 output channels)
	stream, err := portaudio.OpenStream(portaudio.StreamParameters{
		Input:      portaudio.StreamDeviceParameters{Device: device, Channels: 1, Latency: device.DefaultLowInputLatency},
		SampleRate: sampleRate,
	}, data)

	if err != nil {
		log.Fatalf("Failed to open PortAudio stream: %v", err)
	}

	// Return the audio processor instance
	return &AudioProcessor{
		Stream:     stream,
		BufferSize: bufferSize,
		SampleRate: sampleRate,
		Data:       data,
	}
}

// Start the audio stream to begin capturing microphone input
func (ap *AudioProcessor) Start() {
	err := ap.Stream.Start()
	if err != nil {
		log.Fatalf("Failed to start audio stream: %v", err)
	}
	fmt.Println("Microphone stream started...")
}

// Stop the audio stream
func (ap *AudioProcessor) Stop() {
	err := ap.Stream.Stop()
	if err != nil {
		log.Fatalf("Failed to stop audio stream: %v", err)
	}
	portaudio.Terminate()
	fmt.Println("Microphone stream stopped.")
}

// Process reads audio samples and performs FFT to get frequency data
const fftSize = 4096 // Set fftSize to 4096 to get 2048 frequency bins

func (ap *AudioProcessor) Process() []float64 {
	// Read data from the microphone stream
	err := ap.Stream.Read()
	if err != nil {
		log.Fatalf("Failed to read from audio stream: %v", err)
	}

	// Check if data is populated
	if len(ap.Data) == 0 {
		log.Fatalf("No audio data received from stream")
	}

	// Convert float32 data to float64
	data := toFloat64Slice(ap.Data)

	// Ensure fftSize samples are available
	if len(data) < fftSize {
		// Zero-pad the data if it's smaller than fftSize
		paddedData := make([]float64, fftSize)
		copy(paddedData, data) // Copy the available data
		data = paddedData
	} else {
		// Slice the data to exactly fftSize
		data = data[:fftSize]
	}

	// Calculate the maximum amplitude (peak) from raw audio data
	maxAmplitude := 0.0
	for _, sample := range data {
		absSample := math.Abs(sample)
		if absSample > maxAmplitude {
			maxAmplitude = absSample
		}
	}

	// Convert peak amplitude to decibels (dBFS)
	// if maxAmplitude > 0 {
	// 	peakLevel := 20 * math.Log10(maxAmplitude)
	// 	fmt.Printf("Peak Amplitude Level: %.2f dBFS\n", peakLevel)
	// } else {
	// 	fmt.Println("Peak Amplitude Level: -Infinity dBFS (silence)")
	// }

	// Calculate FFT with exactly fftSize samples
	fftData := fft.FFTReal(data)

	// Calculate magnitude of first 2048 frequency bins
	// fftData will have fftSize / 2 complex values due to symmetry (Nyquist theorem)
	binCount := 2048
	frequencies := make([]float64, binCount)
	for i := 0; i < binCount; i++ {
		realPart := real(fftData[i])
		imagPart := imag(fftData[i])
		frequencies[i] = math.Sqrt(realPart*realPart + imagPart*imagPart)
	}

	return frequencies
}

// Helper function to convert float32 slice to float64 slice
func toFloat64Slice(data []float32) []float64 {
	result := make([]float64, len(data))
	for i, v := range data {
		result[i] = float64(v)
	}
	return result
}

// Process audio in real-time and get frequency data (called from your main loop)
func processAudio() []float64 {
	if audioProcessor == nil {
		return nil
	}

	// Get the FFT frequency data from the microphone
	frequencies := audioProcessor.Process()
	return frequencies
}
