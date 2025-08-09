package logging

import (
	"os"
	"strings"
)

type Color string
type Role string

type SystemInfo struct {
	Role  Role
	Color Color
}

func (SystemInfo) ParseColor(color string) Color {
	color = strings.ToLower(color)
	color = strings.ReplaceAll(color, " ", "")
	switch color {
	case "red":
		return Red
	case "green":
		return Green
	case "yellow":
		return Yellow
	case "blue":
		return Blue
	case "purple":
		return Purple
	case "cyan":
		return Cyan
	case "gray":
		return Gray
	default:
		return White
	}
}

var config SystemInfo

func ReadConfig() {
	config.Role = Role(os.Getenv("ROLE"))
	config.Color = config.ParseColor(os.Getenv("COLOR"))
}
