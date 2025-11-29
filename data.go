package q2d

import (
	_ "embed"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed fonts/unscii-8-alt.otf
var fontAlternativeData []byte

//go:embed fonts/unscii-8-fantasy.otf
var fontFantasyData []byte

//go:embed fonts/unscii-8-mcr.otf
var fontSciFiData []byte

//go:embed fonts/unscii-8-tall.otf
var fontTallChunkyData []byte

//go:embed fonts/unscii-8-thin.otf
var fontThinData []byte

//go:embed fonts/unscii-8.otf
var fontNormalData []byte

//go:embed fonts/unscii-16.otf
var fontTallData []byte

//go:embed fonts/unscii.txt
var FontLicense []byte

var (
	FontAlternative font.Face
	FontFantasy     font.Face
	FontSciFi       font.Face
	FontTallChunky  font.Face
	FontThin        font.Face
	FontNormal      font.Face
	FontTall        font.Face
)

func init() {
	FontAlternative = loadFont(fontAlternativeData, 8)
	FontFantasy = loadFont(fontFantasyData, 8)
	FontSciFi = loadFont(fontSciFiData, 8)
	FontTallChunky = loadFont(fontTallChunkyData, 8)
	FontThin = loadFont(fontThinData, 8)
	FontNormal = loadFont(fontNormalData, 8)
	FontTall = loadFont(fontTallData, 16)
}

func loadFont(data []byte, size float64) font.Face {
	f, err := opentype.Parse(data)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("failed to create font face: %v", err)
	}
	return face
}
