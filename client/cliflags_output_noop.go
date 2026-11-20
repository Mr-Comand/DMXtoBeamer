//go:build !ndi

package main

// parseStreamFlag handles the streaming output options. Builds without the
// streaming capability do not recognize them, so the flags fall through to
// the generic unknown-option warning in parseFlags.
func parseStreamFlag(args []string, i *int, opts *clientOpts) bool {
	return false
}

// startOutput is a no-op in builds without streaming support; it always
// reports success so the caller continues normally.
func startOutput(clientId string, opts clientOpts) bool {
	return true
}

// streamStartFailed is always false in builds without streaming support.
func streamStartFailed() bool {
	return false
}
