package cmd

import (
	"github.com/gookit/color"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	RenderYellow = color.FgLightYellow.Render
	RenderGreen  = color.FgGreen.Render
	RenderCyan   = color.FgLightCyan.Render
	RenderRed    = color.FgRed.Render

	FormatNumber = message.NewPrinter(language.BrazilianPortuguese)
)
