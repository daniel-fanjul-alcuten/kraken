package kraken

import (
	"fmt"
	"os"
)

const VERSION = "0.0.1-snapshot"

func ShowVersion() {
	fmt.Fprintf(os.Stdout, "%s version %s\n", os.Args[0], VERSION)
	os.Exit(0)
}
