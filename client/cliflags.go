package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"technikflg.com/dmxToProjector/window"
)

// ndiOptions holds the parsed --ndi* command line flags.
type ndiOptions struct {
	// ndi streams to NDI in addition to the window (or instead of it, when
	// ndiOnly is set).
	ndi bool
	// ndiOnly streams to NDI only — no visible window is shown (headless).
	// Rendering still runs with hardware acceleration through raylib/OpenGL.
	ndiOnly bool
	// ndiName is the NDI source name (default: DMXToBeamer-<ClientId>).
	ndiName string
	// ndiGroups is the comma-separated NDI discovery group list ("" = default groups).
	ndiGroups string
	// ndiWidth/ndiHeight is the NDI output resolution. In --ndi-only mode it
	// is also the (hidden) render resolution; in windowed mode the rendered
	// framebuffer is scaled to it before being sent. When neither flag is
	// given, --ndi-only renders at defaultNDIResolution and windowed mode
	// sends the window resolution unchanged (no scaling).
	ndiWidth  int32
	ndiHeight int32
	// ndiResSet records whether --ndi-width/--ndi-height were given
	// explicitly, so windowed mode only scales when the user asked for it.
	ndiResSet bool
}

// defaultNDIResolution is the internal render resolution used in --ndi-only
// mode when nothing else is configured.
const (
	defaultNDIWidth  = int32(1280)
	defaultNDIHeight = int32(720)
)

// parseFlags splits args (normally os.Args[1:]) into positional arguments and
// --ndi* flags. Positional arguments keep their original order so the existing
// URL/client name handling is unchanged. Unknown -- flags are warned about and
// ignored so existing positional usage keeps working.
func parseFlags(args []string) (positional []string, opts ndiOptions, err error) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--ndi":
			opts.ndi = true
		case arg == "--ndi-only":
			opts.ndiOnly = true
			opts.ndi = true
		case strings.HasPrefix(arg, "--ndi-name="):
			opts.ndiName = strings.TrimPrefix(arg, "--ndi-name=")
		case arg == "--ndi-name":
			if i+1 >= len(args) {
				err = fmt.Errorf("--ndi-name requires a value")
				return
			}
			i++
			opts.ndiName = args[i]
		case strings.HasPrefix(arg, "--ndi-groups="):
			opts.ndiGroups = strings.TrimPrefix(arg, "--ndi-groups=")
		case strings.HasPrefix(arg, "--ndi-width="):
			v, e := strconv.ParseInt(strings.TrimPrefix(arg, "--ndi-width="), 10, 32)
			if e != nil || v <= 0 {
				err = fmt.Errorf("invalid value for --ndi-width: %q", arg)
				return
			}
			opts.ndiWidth = int32(v)
			opts.ndiResSet = true
		case strings.HasPrefix(arg, "--ndi-height="):
			v, e := strconv.ParseInt(strings.TrimPrefix(arg, "--ndi-height="), 10, 32)
			if e != nil || v <= 0 {
				err = fmt.Errorf("invalid value for --ndi-height: %q", arg)
				return
			}
			opts.ndiHeight = int32(v)
			opts.ndiResSet = true
		case arg == "--help" || arg == "-h":
			printUsage()
			os.Exit(0)
		case strings.HasPrefix(arg, "--"):
			fmt.Fprintf(os.Stderr, "Warning: unknown option %q ignored\n", arg)
		default:
			positional = append(positional, arg)
		}
	}
	if opts.ndiWidth == 0 {
		opts.ndiWidth = defaultNDIWidth
	}
	if opts.ndiHeight == 0 {
		opts.ndiHeight = defaultNDIHeight
	}
	return
}

// ndiTargetWidth/ndiTargetHeight map the parsed flags to the NDI output
// resolution passed to window.InitNDI: 0 means "no explicit target" (send the
// window render resolution unchanged in windowed mode). In --ndi-only mode the
// defaults always apply (the hidden render window needs a fixed size).
func ndiTargetWidth(opts ndiOptions) int32 {
	if opts.ndiOnly || opts.ndiResSet {
		return opts.ndiWidth
	}
	return 0
}

func ndiTargetHeight(opts ndiOptions) int32 {
	if opts.ndiOnly || opts.ndiResSet {
		return opts.ndiHeight
	}
	return 0
}

// setupSignalHandler makes Ctrl+C / SIGTERM stop the render loop cleanly so
// NDI output is properly closed (important for the headless --ndi-only mode,
// whose hidden window never receives a user close event).
func setupSignalHandler() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("Shutdown requested (SIGINT/SIGTERM)")
		window.RequestShutdown()
	}()
}

// printUsage is implemented in usage_ndi.go (NDI builds) and
// usage_noop.go (non-NDI builds) so help text matches the build.

