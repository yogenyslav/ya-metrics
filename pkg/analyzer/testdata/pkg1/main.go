package main

import (
	"os"
)

func main() {
	foo()
	os.Exit(1) // want "os.Exit calls are prohibited in main function of package main"
}
