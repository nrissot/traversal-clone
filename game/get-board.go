package game

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const ENDPOINT = "https://api.traversalgame.com/api/boards/generate-assembled"
const APP_TYPE = "application/json"

type BoardType string

const (
	Random           BoardType = "random"
	Mountains        BoardType = "mountains"
	Valleys          BoardType = "valleys"
	Highlands        BoardType = "highlands"
	Woodlands        BoardType = "woodlands"
	Desert           BoardType = "desert"
	Dunes            BoardType = "dunes"
	Jungle           BoardType = "jungle"
	Meadows          BoardType = "meadows"
	Canyons          BoardType = "canyons"
	Rivers           BoardType = "rivers"
	Ridges           BoardType = "ridges"
	Petrified_Forest BoardType = "petrified"
	Mount_Everest    BoardType = "everest"
	Mariana_Trench   BoardType = "mariana"
	Hell             BoardType = "hell"
)

type RespStruct struct {
	Board           [][]int `json:"board"`
	Rows            int     `json:"rows"`
	Columns         int     `json:"columns"`
	Terrain         string  `json:"terrain"`
	TerrainSelected string  `json:"terrain_selected"`
	// AssemblyPlan    [][]int  `json:"assembly_plan"`
}

func GetBoard(boardType BoardType) (*RespStruct, error) {
	request_body := fmt.Appendf(nil, "{\"terrain\":\"%s\"}", boardType)
	resp, err := http.Post(ENDPOINT, APP_TYPE, bytes.NewReader(request_body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resp_struct *RespStruct = &RespStruct{}

	err = json.Unmarshal(content, resp_struct)
	if err != nil {
		return nil, err
	}
	return resp_struct, nil
}
