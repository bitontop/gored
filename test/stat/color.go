package test

type Color string

const (
	COLOR_RED    Color = "\x1b[31;1m"
	COLOR_GREEN  Color = "\x1b[32;1m"
	COLOR_YELLOW Color = "\x1b[33;1m"
	COLOR_BLUE   Color = "\x1b[34;1m"
	COLOR_PUP    Color = "\x1b[35;1m"
	COLOR_CYAN   Color = "\x1b[36;1m"
	COLOR_B_PINK Color = "\x1b[45;1m"
	COLOR_B_BLUE Color = "\x1b[44;1m"
	COLOR_B_RED  Color = "\x1b[41;1m"
	COLOR_BRU    Color = "\x1b[41;1m"

	COLOR_NO Color = "\x1b[0m" // No Color

)
