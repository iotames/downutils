package fyne

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type myTheme struct {
	defaultFont fyne.Resource
}

func (t *myTheme) Font(s fyne.TextStyle) fyne.Resource {
	// if s.Monospace {
	// 	return theme.DefaultTheme().Font(s)
	// }
	// if s.Bold {
	// 	if s.Italic {
	// 		return theme.DefaultTheme().Font(s)
	// 	}
	// 	return t.defaultFont
	// }
	// if s.Italic {
	// 	return theme.DefaultTheme().Font(s)
	// }
	return t.defaultFont
}

func (*myTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (*myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (*myTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}

func NewFyneTheme(font fyne.Resource) *myTheme {
	return &myTheme{defaultFont: font}
}
