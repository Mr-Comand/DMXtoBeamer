//go:build ndi

package main

import (
	"strconv"
	"strings"

	"technikflg.com/dmxToProjector/window"
)

// parseStreamFlag parses the streaming output options in NDI builds and
// returns true if it consumed arg. It may consume additional args (value
// taking options) by advancing *i. Malformed --ndi-width/--ndi-height values
// produce the specific error messages here.
func parseStreamFlag(args []string, i *int, opts *clientOpts) bool {
	arg := args[*i]
	switch {
	case arg == "--ndi":
		opts.streamExtra = true
	case arg == "--ndi-only":
		opts.streamOnly = true
		opts.streamExtra = true
	case strings.HasPrefix(arg, "--ndi-name="):
		opts.streamName = strings.TrimPrefix(arg, "--ndi-name=")
	case arg == "--ndi-name":
		if *i+1 >= len(args) {
			// parseFlags reports the generic unknown-option warning; the
			// more specific error is handled by the caller afterwards.
			return false
		}
		*i++
		opts.streamName = args[*i]
	case strings.HasPrefix(arg, "--ndi-groups="):
		opts.streamGroups = strings.TrimPrefix(arg, "--ndi-groups=")
	case strings.HasPrefix(arg, "--ndi-width="):
		v, e := strconv.ParseInt(strings.TrimPrefix(arg, "--ndi-width="), 10, 32)
		if e != nil || v <= 0 {
			return false
		}
		opts.streamWidth = int32(v)
		opts.streamResSet = true
	case strings.HasPrefix(arg, "--ndi-height="):
		v, e := strconv.ParseInt(strings.TrimPrefix(arg, "--ndi-height="), 10, 32)
		if e != nil || v <= 0 {
			return false
		}
		opts.streamHeight = int32(v)
		opts.streamResSet = true
	default:
		return false
	}
	return true
}

// startOutput starts the streaming output sink when requested. It returns
// false if the user asked for streaming output but it could not be started.
// In NDI builds the default source name is derived from the client id.
func startOutput(clientId string, opts clientOpts) bool {
	if !opts.streamExtra {
		return true
	}
	streamName := opts.streamName
	if streamName == "" {
		streamName = "DMXToBeamer-" + clientId
	}
	window.InitOutput(streamName, opts.streamGroups, opts.streamOnly,
		streamTargetWidth(opts), streamTargetHeight(opts))
	return window.OutputManager() != nil
}

// streamStartFailed reports whether streaming output was requested but could
// not be started (used to abort in headless mode).
func streamStartFailed() bool {
	return window.OutputManager() == nil
}
