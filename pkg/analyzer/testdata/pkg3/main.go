package pkg3

import "os"

func main() {
	if false {
		os.Exit(1) // is not reported
	}
}
