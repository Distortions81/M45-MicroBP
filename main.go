package main

import (
	"bytes"
	"compress/zlib"
	"encoding/ascii85"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	Entity_number int    `json:"entity_number"`
	Name          string `json:"name"`
	Position      Xy     `json:"position"`
	Direction     uint8  `json:"direction"`
}

type SignalData struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Icn struct {
	Signal SignalData `json:"signaldata"`
	Index  int16      `json:"index"`
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

func main() {

	var input []byte
	var err error

	// Open the input file
	input, err = os.ReadFile("bp.txt")
	if err != nil {
		fmt.Println("ERROR: Opening input file", err)
		os.Exit(1)
	}

	fmt.Println("Input file size:", len(input))

	//Max input
	if len(input) > maxInput {
		fmt.Println("ERROR: Input data too large.")
		os.Exit(1)
	}

	data, err := base64.StdEncoding.DecodeString(string(input[1:]))
	if err != nil {
		fmt.Println("ERROR: decoding input:", err)
		os.Exit(1)
	}

	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		fmt.Println("ERROR: decompress start failure:", err)
		os.Exit(1)
	}
	enflated, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("ERROR: Decompress read failure:", err)
		os.Exit(1)
	}

	//Max decompressed size
	if len(enflated) > maxJson {
		fmt.Println("ERROR: Input data too large.")
		os.Exit(1)
	}

	newbook := RootData{}

	err = json.Unmarshal(enflated, &newbook)
	if err != nil {
		fmt.Println("ERROR: JSON Unmarshal failure:", err)
		os.Exit(1)
	}

	fileName := "shrunk.json"
	err = os.WriteFile(fileName, enflated, 0644)
	if err != nil {
		fmt.Println("ERROR: Failed to write json:", err)
		os.Exit(1)
	}
	fmt.Println("Wrote json to", fileName)

	//Make compact format bp
	var compbp compBPData
	var entNameMap = make(map[string]uint16)
	var tileNameMap = make(map[string]uint16)

	var entNumber uint16 = 1
	for _, ent := range newbook.Blueprint.Entities {
		if entNameMap[ent.Name] == 0 {
			entNameMap[ent.Name] = entNumber
			entNumber++
		}
	}

	for _, ent := range newbook.Blueprint.Entities {
		xyh := false
		xhf := float32(ent.Position.X - float32(int(ent.Position.X)))
		yhf := float32(ent.Position.Y - float32(int(ent.Position.Y)))

		if xhf > 0.0 || yhf > 0.0 {
			xyh = true
		}
		compbp.Ents = append(compbp.Ents, compEntity{
			X: int32(ent.Position.X), Y: int32(ent.Position.Y),
			XYh: xyh,
			Dir: ent.Direction, Name: entNameMap[ent.Name],
		})
		compbp.EntNames = make([]string, len(entNameMap))
		for key, nameNum := range entNameMap {
			compbp.EntNames[nameNum-1] = key
		}
	}

	entNumber = 1
	for _, tile := range newbook.Blueprint.Tiles {
		if tileNameMap[tile.Name] == 0 {
			tileNameMap[tile.Name] = entNumber
			entNumber++
		}
	}

	for _, tile := range newbook.Blueprint.Tiles {
		xyh := false
		xhf := float32(tile.Position.X - float32(int(tile.Position.X)))
		yhf := float32(tile.Position.Y - float32(int(tile.Position.Y)))

		if xhf > 0.0 || yhf > 0.0 {
			xyh = true
		}
		compbp.Tiles = append(compbp.Tiles, compTile{
			X: int32(tile.Position.X), Y: int32(tile.Position.Y),
			XYh: xyh, Name: tileNameMap[tile.Name],
		})
		compbp.TileNames = make([]string, len(tileNameMap))
		for key, nameNum := range tileNameMap {
			compbp.TileNames[nameNum-1] = key
		}
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(compbp)
	if err != nil {
		fmt.Println("ERROR: Gob encode failure:", err)
		os.Exit(1)
	}
	fmt.Printf("gob length: %v\n", len(buf.Bytes()))

	gz := compressGzip(buf.Bytes())
	fmt.Printf("gz length: %v\n", len(gz))

	dst := make([]byte, ascii85.MaxEncodedLen(4+len(gz)))
	ascii85.Encode(dst, gz)
	fmt.Printf("asci85 length: %v\n", len(dst))
	//fmt.Println(string(dst))
}

type compBPData struct {
	Ents     []compEntity
	EntNames []string

	Tiles     []compTile
	TileNames []string

	Label   string
	Version uint64
}

type compTile struct {
	X    int32
	Y    int32
	XYh  bool
	Name uint16
}

type compEntity struct {
	X   int32
	Y   int32
	XYh bool

	Dir uint8

	Name uint16
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
