package main

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"

	"github.com/dsnet/compress/bzip2"
)

//Max sizes
const maxInput = 10 * 1024 * 1024 //10MB
const maxJson = 100 * 1024 * 1024 //100MB

type Xy struct {
	X float32 `json:"y"`
	Y float32 `json:"x"`
}

type Tile struct {
	Position Xy     `json:"position,omitempty"`
	Name     string `json:"name,omitempty"`
}

type Ent struct {
	Entity_number int      `json:"entity_number,omitempty"`
	Name          string   `json:"name"`
	Position      Xy       `json:"position,omitempty"`
	Direction     uint8    `json:"direction,omitempty"`
	Recipe        string   `json:"recipe,omitempty"`
	Neighbours    []uint16 `json:"neighbours,omitempty"`
	Type          string   `json:"type,omitempty"`
}

type SignalData struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

type Icn struct {
	Signal SignalData `json:"signal,omitempty"`
	Index  uint16     `json:"index,omitempty"`
}

type BluePrintData struct {
	Entities []Ent  `json:"entities,omitempty"`
	Tiles    []Tile `json:"tiles,omitempty"`
	Icons    []Icn  `json:"icons,omitempty"`
	Item     string `json:"item,omitempty"`
	Label    string `json:"label,omitempty"`
	Version  uint64 `json:"version,omitempty"`
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
	X   int16 `json:"x,omitempty"`
	Y   int16 `json:"y,omitempty"`
	XYh bool  `json:"z,omitempty"`
}

type compBPData struct {
	Ents     []compEntity `json:"e,omitempty"`
	EntNames []string     `json:"n,omitempty"`
	EntRec   []string     `json:"r,omitempty"`
	EntType  []string     `json:"t,omitempty"`

	Tiles     []compTile `json:"i,omitempty"`
	TileNames []string   `json:"m,omitempty"`

	Icons []Icn `json:"c,omitempty"`

	Label   string `json:"l,omitempty"`
	Item    string `json:"f,omitempty"`
	Version uint64 `json:"v,omitempty"`
}

type compTile struct {
	Pos  compXy `json:"p,omitempty"`
	Name uint16 `json:"n,omitempty"`
}
type compEntity struct {
	Name       uint16   `json:"n,omitempty"`
	Pos        compXy   `json:"p,omitempty"`
	Dir        uint8    `json:"d,omitempty"`
	Type       uint16   `json:"t,omitempty"`
	Rec        uint16   `json:"r,omitempty"`
	Neighbours []uint16 `json:"n,omitempty"`
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

func compressBzip2(data []byte) []byte {
	var b bytes.Buffer
	w, err := bzip2.NewWriter(&b, &bzip2.WriterConfig{Level: 9})
	if err != nil {
		fmt.Println("ERROR: Bzip writer failure:", err)
	}
	w.Write(data)
	w.Close()
	return b.Bytes()
}
