package main

import (
	"fmt"
	"github.com/ricochet2200/go-disk-usage/du"
)

var KB = uint64(1024)

func main() {
	usage := du.NewDiskUsage("C:\\")
	fmt.Println("Free:", usage.Free()/(KB*KB))
	fmt.Println("Available:", usage.Available()/(KB*KB))
	fmt.Println("Size:", usage.Size()/(KB*KB))
	fmt.Println("Used:", usage.Used()/(KB*KB))
	fmt.Println("Usage:", usage.Usage()*100, "%")
}
