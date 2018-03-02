package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	defaultWidth  = 128
	defaultHeight = 128
)

type options struct {
	width  int
	height int
}

func parseArgs(args []string) options {
	options := options{
		defaultWidth,
		defaultHeight,
	}

	option := ""

	for _, arg := range args {
		if len(option) > 0 {
			switch option {
			case "width":
				options.width, _ = strconv.Atoi(arg)
				if options.width <= 0 {
					fmt.Printf("Width must be positive '" + string(options.width) + "'\n")
				}
				break
			case "height":
				options.height, _ = strconv.Atoi(arg)
				if options.height <= 0 {
					fmt.Printf("Height must be positive '" + string(options.width) + "'\n")
				}
				break
			default:
				fmt.Printf("Unsupported option '" + option + "'\n")
			}
			option = ""
			continue
		} else if strings.HasPrefix(arg, "-") {
			switch strings.ToLower(arg) {
			case "-w", "--width":
				option = "width"
				continue
			case "-h", "--height":
				option = "height"
				continue
			case "-v", "--verbose":
				continue
			default:
			}
		}
		fmt.Printf("Unsupported argument '" + arg + "'\n")
	}

	return options
}
