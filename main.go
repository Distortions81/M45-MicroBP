package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
)

//Max sizes
const maxInput = 10 * 1024 * 1024 //10MB
const maxJson = 100 * 1024 * 1024 //100MB

type Xy struct {
	X float32 `json:"y"`
	Y float32 `json:"x"`
}

type Tile struct {
	Position Xy     `json:"position"`
	Name     string `json:"name"`
}

type Ent struct {
	Entity_number int      `json:"entity_number"`
	Name          string   `json:"name"`
	Position      Xy       `json:"position"`
	Direction     uint8    `json:"direction"`
	Recipe        string   `json:"recipe"`
	Neighbours    []uint16 `json:"neighbours"`
	Type          string   `json:"type"`
}

type SignalData struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Icn struct {
	Signal SignalData `json:"signaldata"`
	Index  uint16     `json:"index"`
}

type BluePrintData struct {
	Entities []Ent  `json:"entities"`
	Tiles    []Tile `json:"tiles"`
	Icons    []Icn  `json:"icons"`
	Item     string `json:"item"`
	Label    string `json:"label"`
	Version  uint64 `json:"version"`
}
type BluePrintsData struct {
	Blueprint BluePrintData `json:"blueprint"`
}
type BluePrintBookData struct {
	Blueprints []BluePrintsData `json:"blueprints"`
}
type RootData struct {
	Blueprintbook BluePrintBookData `json:"blueprint_book"`
	Blueprint     BluePrintData     `json:"blueprint"`
}

//bool pointer to bool
func bAddr(b *bool) bool {
	boolVar := *b

	if boolVar {
		return boolVar
	}

	return false

}

var doCompress bool = false
var doDecompress bool = false

func main() {

	doCompressP := flag.Bool("doCompress", false, "convert blueprint to uBP format")
	doDecompressP := flag.Bool("doDecompress", true, "convert uBP to blueprint format")
	flag.Parse()

	doCompress = bAddr(doCompressP)
	doDecompress = bAddr(doDecompressP)

	if doCompress {
		fmt.Println("Compressing to uBP format")
		compress()
		return
	} else if doDecompress {
		fmt.Println("Decompressing from uBP format")
		decompress()
		return
	} else {
		fmt.Println("No action specified, try -h for help")
		return
	}

}

type compXy struct {
	X   int16
	Y   int16
	XYh bool
}

type compBPData struct {
	Ents     []compEntity
	EntNames []string
	EntRec   []string
	EntType  []string
	Items    []string

	Tiles     []compTile
	TileNames []string

	Label   string
	Item    string
	Version uint64
}

type compTile struct {
	Pos  compXy
	Name uint16
}

type compItems struct {
	Name  int
	Count int
}

type compEntity struct {
	Name       uint16
	Pos        compXy
	Dir        uint8
	Type       uint16
	Rec        uint16
	Items      []compItems
	Neighbours []uint16
}

func compressGzip(data []byte) []byte {
	var b bytes.Buffer
	w, err := zlib.NewWriterLevel(&b, zlib.BestCompression)
	if err != nil {
		fmt.Println("ERROR: Gzip writer failure:", err)
	}
	w.Write(data)
	w.Close()
	return b.Bytes()
}
