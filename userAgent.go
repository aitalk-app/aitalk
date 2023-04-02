package main

import "fmt"

const NAME = "aitalk"

var VERSION = "0.1.0"

func userAgent() string {
	return fmt.Sprintf("%s/%s", NAME, VERSION)
}
