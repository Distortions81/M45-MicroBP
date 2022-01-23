package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func compress() {
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

	newbook := BluePrintsData{}

	err = json.Unmarshal(enflated, &newbook)
	if err != nil {
		fmt.Println("ERROR: JSON Unmarshal failure:", err)
		os.Exit(1)
	}

	fileName := "start.json"
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
	var typeMap = make(map[string]uint16)
	var recipeMap = make(map[string]uint16)

	var entNumber uint16
	entNumber = 1
	for _, ent := range newbook.Blueprint.Entities {
		if recipeMap[ent.Recipe] == 0 && ent.Recipe != "" {
			recipeMap[ent.Recipe] = entNumber
			entNumber++
		}
	}
	entNumber = 1
	for _, ent := range newbook.Blueprint.Entities {
		if entNameMap[ent.Name] == 0 && ent.Name != "" {
			entNameMap[ent.Name] = entNumber
			entNumber++
		}
	}
	entNumber = 1
	for _, tile := range newbook.Blueprint.Entities {
		if typeMap[tile.Type] == 0 && tile.Type != "" {
			typeMap[tile.Type] = entNumber
			entNumber++
		}
	}
	//Tiles
	entNumber = 1
	for _, tile := range newbook.Blueprint.Tiles {
		if tileNameMap[tile.Name] == 0 && tile.Name != "" {
			tileNameMap[tile.Name] = entNumber
			entNumber++
		}
	}

	compbp.EntNames = make([]string, len(entNameMap))
	for key, nameNum := range entNameMap {
		compbp.EntNames[nameNum-1] = key
	}
	compbp.EntRec = make([]string, len(recipeMap))
	for key, nameNum := range recipeMap {
		compbp.EntRec[nameNum-1] = key
	}
	compbp.EntType = make([]string, len(typeMap))
	for key, nameNum := range typeMap {
		compbp.EntType[nameNum-1] = key
	}
	//Tiles
	compbp.TileNames = make([]string, len(tileNameMap))
	for key, nameNum := range tileNameMap {
		compbp.TileNames[nameNum-1] = key
	}

	for _, ent := range newbook.Blueprint.Entities {
		xyh := false
		xhf := float32(ent.Position.X - float32(int(ent.Position.X)))
		yhf := float32(ent.Position.Y - float32(int(ent.Position.Y)))

		if xhf > 0.0 || yhf > 0.0 {
			xyh = true
		}
		compbp.Ents = append(compbp.Ents, compEntity{
			Pos: compXy{X: int16(ent.Position.X), Y: int16(ent.Position.Y), XYh: xyh},
			Dir: ent.Direction, Name: entNameMap[ent.Name], Type: typeMap[ent.Type], Rec: recipeMap[ent.Recipe]},
		)
	}

	//Tiles
	for _, tile := range newbook.Blueprint.Tiles {
		xyh := false
		xhf := float32(tile.Position.X - float32(int(tile.Position.X)))
		yhf := float32(tile.Position.Y - float32(int(tile.Position.Y)))

		if xhf > 0.0 || yhf > 0.0 {
			xyh = true
		}
		compbp.Tiles = append(compbp.Tiles, compTile{
			Pos:  compXy{X: int16(tile.Position.X), Y: int16(tile.Position.Y), XYh: xyh},
			Name: tileNameMap[tile.Name]},
		)
	}

	compbp.Item = newbook.Blueprint.Item
	compbp.Label = newbook.Blueprint.Label
	compbp.Version = newbook.Blueprint.Version

	//JSON OUT DEBUG
	outbuf := new(bytes.Buffer)
	enca := json.NewEncoder(outbuf)
	enca.SetIndent("", "\t")
	if err := enca.Encode(compbp); err != nil {
		fmt.Println("WriteGCfg: enc.Encode failure")
	}
	fmt.Println(outbuf.String())

	//GOB ENCODE
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

	dst := encode85(string(gz))
	fmt.Printf("asci85 length: %v\n", len(dst))

	fileName = "micro.txt"
	err = os.WriteFile(fileName, []byte(dst), 0644)
	if err != nil {
		fmt.Println("ERROR: Failed to write dst:", err)
		os.Exit(1)
	}
	fmt.Println("Wrote dst to", fileName)
}
