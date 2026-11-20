//go:build ndi

package main

import "fmt"

// printUsage prints the client usage including the NDI flags. This variant is
// only compiled into NDI-enabled builds so non-NDI builds carry no NDI help
// text or trademark attribution.
func printUsage() {
	fmt.Print(`Usage: client [options] [ws://host:port] [ClientName]

NDI® options (NDI® is a registered trademark of Vizrt NDI AB, https://ndi.video/):
  --ndi                 Additionally stream the rendered output to NDI.
  --ndi-only            Stream to NDI only — no window is opened (headless).
                        Rendering still uses the GPU (hardware accelerated).
  --ndi-name=NAME       NDI source name (default: DMXToBeamer-<ClientName>).
  --ndi-groups=LIST     Comma separated NDI discovery groups.
  --ndi-width=N         NDI output width. --ndi-only: render resolution
                        (default 1280). Windowed mode: the framebuffer is
                        scaled to N before being sent (unset = window size).
  --ndi-height=N        NDI output height (same rules, default 720).

The NDI runtime must be installed separately; it is not bundled with this
application. Get it at https://ndi.video/tools/.
`)
}
