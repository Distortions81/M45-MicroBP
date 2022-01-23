package main

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func decompress() {
	var input []byte
	var err error

	var decompbp compBPData

	// Open the input file
	input, err = os.ReadFile("micro.txt")
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

	zdata := decode85(string(input))

	r, err := zlib.NewReader(bytes.NewReader([]byte(zdata)))
	if err != nil {
		fmt.Println("ERROR: decompress start failure:", err)
		os.Exit(1)
	}

	dec := gob.NewDecoder(r)

	if err := dec.Decode(&decompbp); err != nil {
		log.Fatal(err)
	}

	var EntNames []string = make([]string, len(decompbp.EntNames))
	var EntRec []string = make([]string, len(decompbp.EntRec))
	var EntType []string = make([]string, len(decompbp.EntType))
	//Tile
	var TileNames []string = make([]string, len(decompbp.TileNames))

	for i, v := range decompbp.EntNames {
		EntNames[i] = string(v)
	}

	for i, v := range decompbp.EntRec {
		EntRec[i] = string(v)
	}

	for i, v := range decompbp.EntType {
		EntType[i] = string(v)
	}

	//Tile
	for i, v := range decompbp.TileNames {
		TileNames[i] = string(v)
	}

	var bp RootData
	EntNumber := 1
	for _, ent := range decompbp.Ents {
		newX := float32(ent.Pos.X)
		newY := float32(ent.Pos.Y)
		if ent.Pos.XYh {
			newX = float32(ent.Pos.X) + 0.5
			newY = float32(ent.Pos.Y) + 0.5
		}
		bp.Blueprint.Entities = append(bp.Blueprint.Entities, Ent{Entity_number: EntNumber, Name: decompbp.EntNames[ent.Name],
			Position: Xy{X: newX, Y: newY}, Direction: ent.Dir, Recipe: EntRec[ent.Rec], Type: EntType[ent.Type], Neighbours: ent.Neighbours})
		EntNumber++
	}

	//Tiles
	for _, tile := range decompbp.Tiles {
		newX := float32(tile.Pos.X)
		newY := float32(tile.Pos.Y)
		if tile.Pos.XYh {
			newX = float32(tile.Pos.X) + 0.5
			newY = float32(tile.Pos.Y) + 0.5
		}
		bp.Blueprint.Tiles = append(bp.Blueprint.Tiles, Tile{Name: decompbp.TileNames[tile.Name], Position: Xy{X: newX, Y: newY}})
		EntNumber++
	}

	bp.Blueprint.Version = decompbp.Version
	bp.Blueprint.Label = decompbp.Label
	bp.Blueprint.Item = decompbp.Item

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	if err := enc.Encode(bp); err != nil {
		fmt.Println("WriteGCfg: enc.Encode failure")
	}
	//fmt.Println("Entities:", outbuf.String())

	fileName := "decomp.txt"
	err = os.WriteFile(fileName, outbuf.Bytes(), 0644)
	if err != nil {
		fmt.Println("ERROR: Failed to write dst:", err)
		os.Exit(1)
	}
	fmt.Println("Wrote dst to", fileName)
}
