package statusline

import "github.com/fatih/color"

const barWidth = 10

var (
	colorBold     = color.New(color.Bold)
	colorDim      = color.New(color.Faint)
	colorBoldCyan = color.New(color.FgCyan, color.Bold)
	colorYellow   = color.New(color.FgYellow)
	colorGreen    = color.New(color.FgGreen)
	colorRed      = color.New(color.FgRed)
	colorMagenta  = color.New(color.FgMagenta)
)
