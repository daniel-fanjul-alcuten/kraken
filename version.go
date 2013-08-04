package kraken

import (
	"fmt"
	"os"
)

const VERSION = "0.0.1-snapshot"

func ShowVersion() {
	fmt.Fprintf(os.Stderr, "%s version %s\n", os.Args[0], VERSION)
	os.Exit(0)
}
