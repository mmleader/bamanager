//go:build !windows

package main

import _ "embed"

//go:embed build/bsmanager.png
var iconData []byte
