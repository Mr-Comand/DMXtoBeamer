package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"technikflg.com/dmxToProjector/window"
)

// clientOpts holds the parsed command line flags that affect the output
// sink (the window, or an optional streaming output in builds compiled with
// streaming support).
type clientOpts struct {
	// streamExtra streams the rendered output to an additional sink in
	// builds with that capability.
	streamExtra bool
	// streamOnly renders into the additional sink without a visible window
	// (headless). Rendering still runs with hardware acceleration through
	// raylib/OpenGL. Only recognized by builds with the streaming
	// capability; other builds report it as an unknown option.
	streamOnly bool
	// streamName is the output source name.
	streamName string
	// streamGroups is the comma-separated output discovery group list.
	streamGroups string
	// streamWidth/streamHeight is the extra output resolution. In streamOnly
	// mode it is also the (hidden) render resolution; in windowed mode the
	// rendered framebuffer is scaled to it before being sent. When neither
	// flag is given, streamOnly renders at defaultStreamResolution and
	// windowed mode sends the window resolution unchanged (no scaling).
	streamWidth  int32
	streamHeight int32
	// streamResSet records whether the width/height options were given
	// explicitly, so windowed mode only scales when the user asked for it.
	streamResSet bool
}

// defaultStreamResolution is the internal render resolution used in headless
// streamOnly mode when nothing else is configured.
const (
	defaultStreamWidth  = int32(1280)
	defaultStreamHeight = int32(720)
)

// parseFlags splits args (normally os.Args[1:]) into positional arguments and
// option flags. Positional arguments keep their original order so the existing
// URL/client name handling is unchanged. Unknown -- flags are warned about and
// ignored so existing positional usage keeps working.
//
// parseStreamFlags handles the extra streaming output options and is
// implemented in build-tag specific files: it recognizes and parses them in
// builds with streaming support and reports false in all other builds, so
// those options fall through to the generic unknown-option warning.
func parseFlags(args []string) (positional []string, opts clientOpts, err error) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if parseStreamFlag(args, &i, &opts) {
			continue
		}
		switch {
		case arg == "--help" || arg == "-h":
			printUsage()
			os.Exit(0)
		case strings.HasPrefix(arg, "--"):
			fmt.Fprintf(os.Stderr, "Warning: unknown option %q ignored\n", arg)
		default:
			positional = append(positional, arg)
		}
	}
	if opts.streamWidth == 0 {
		opts.streamWidth = defaultStreamWidth
	}
	if opts.streamHeight == 0 {
		opts.streamHeight = defaultStreamHeight
	}
	return
}

// streamTargetWidth/streamTargetHeight map the parsed flags to the output
// resolution passed to the streaming sink: 0 means "no explicit target" (send
// the window render resolution unchanged in windowed mode). In streamOnly
// mode the defaults always apply (the hidden render window needs a fixed
// size).
func streamTargetWidth(opts clientOpts) int32 {
	if opts.streamOnly || opts.streamResSet {
		return opts.streamWidth
	}
	return 0
}

func streamTargetHeight(opts clientOpts) int32 {
	if opts.streamOnly || opts.streamResSet {
		return opts.streamHeight
	}
	return 0
}

// setupSignalHandler makes Ctrl+C / SIGTERM stop the render loop cleanly so
// output is properly closed (important for headless streamOnly mode, whose
// hidden window never receives a user close event).
func setupSignalHandler() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		log.Println("Shutdown requested (SIGINT/SIGTERM)")
		window.RequestShutdown()
	}()
}

// printUsage is implemented in build-tag specific files so the help text
// matches the build.

