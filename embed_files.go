package main

import (
	_ "embed"
)

//go:embed resource/images/logo.png
var logopng []byte

//go:embed resource/fonts/SourceHanSans-Bold.ttf
var fyfont []byte

// var resourceFontsTtf = &fyne.StaticResource{
// 	StaticName:    "SourceHanSans-Bold.ttf",
// 	StaticContent: fyfont,
// }
