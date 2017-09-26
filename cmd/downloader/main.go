package main

import (
	"os"

	gcs "github.com/SoheilSalehian/gcs-image-downloader"
)

func main() {
	gcs.Run(os.Args[1:])
}
