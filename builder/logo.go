package builder

import (
	"math/rand"
	"strings"
	"time"
)

func LogoBuild() (logo string) {
	rand.Seed(time.Now().UnixNano())
	color_palette := []string{"\033[38;5;97m", "\033[38;5;98m", "\033[38;5;97m", "\033[38;5;96m"}
	logo_build := []string{
		"\t╔╦╗╔═╗╦╔═   ╔═╗╔═╗",
		"\t ║ ║ ║╠╩╗───║ ╦║ ║",
		"\t ╩ ╚═╝╩ ╩   ╚═╝╚═╝",
		"   # [ Programmed by Z3NTL3 ] #",
		"    ## Discord Token Checker ##",
	}
	var build string
	for i, _ := range logo_build {
		build_essential := strings.Split(logo_build[i], "")

		for index, _ := range build_essential {
			color := color_palette[rand.Intn(len(color_palette))]
			if index == len(build_essential)-1 {
				build += color + build_essential[index] + "\033[0m\n\r"
			} else {
				build += color + build_essential[index] + "\033[0m"
			}
		}
	}
	logo = build
	return
}
