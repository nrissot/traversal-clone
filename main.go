package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	game "github.com/nrissot/traversal-clone/game"
)

type Api struct {
	State        game.BoardState
	RedBotToken  string
	BlueBotToken string
}

var api Api = Api{
	RedBotToken:  "",
	BlueBotToken: "",
}

var registerMutex sync.RWMutex

var stateMutex sync.RWMutex

const DISPLAY_SHOULD_SHOW_HEIGHT = false

func main() {
	// setup the game
	resp, err := game.GetBoard(game.Hell)
	if err != nil {
		log.Fatal(err)
	}
	api.State = game.InitializeBoardState(*resp)
	api.State.DisplayBoard(DISPLAY_SHOULD_SHOW_HEIGHT)

	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleRoot)

	mux.HandleFunc("GET /board", HandleGet_Board)
	mux.HandleFunc("GET /state", HandleGet_State)
	mux.HandleFunc("GET /moves", HandleGet_Moves)
	mux.HandleFunc("GET /moves/struct", HandleGet_Moves_Struct)

	mux.HandleFunc("POST /register", HandlePost_Register)
	mux.HandleFunc("POST /play", HandlePost_Play)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func HandleGet_Board(w http.ResponseWriter, r *http.Request) {

	stateMutex.RLock()
	content, err := json.Marshal(api.State.Board)
	stateMutex.RUnlock()

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleGet_State(w http.ResponseWriter, r *http.Request) {

	stateMutex.RLock()
	content, err := json.Marshal(api.State)
	stateMutex.RUnlock()

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleGet_Moves(w http.ResponseWriter, r *http.Request) {

	stateMutex.RLock()
	moves, err := api.State.GetAllLegalMoves()
	stateMutex.RUnlock()

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	notations := []string{}
	for _, move := range moves {
		notations = append(notations, move.String())
	}
	content, err := json.Marshal(notations)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	api.State.DisplayAllowedMoves(DISPLAY_SHOULD_SHOW_HEIGHT, moves)
}

func HandleGet_Moves_Struct(w http.ResponseWriter, r *http.Request) {

	stateMutex.RLock()
	moves, err := api.State.GetAllLegalMoves()
	stateMutex.RUnlock()

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content, err := json.Marshal(moves)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(content)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type PostRegisterReq struct {
	Token string `json:"token"`
}

type PostRegisterResp struct {
	Color game.PieceColor `json:"color"`
}

func HandlePost_Register(w http.ResponseWriter, r *http.Request) {
	// get the "identitication token" provided by the bot
	var reqStruct PostRegisterReq
	err := json.NewDecoder(r.Body).Decode(&reqStruct)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqStruct.Token == "" {
		log.Println("token is required")
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}

	registerMutex.Lock()
	defer registerMutex.Unlock()
	var selected_color game.PieceColor
	if api.RedBotToken == "" {
		selected_color = game.Red
	} else if api.RedBotToken == reqStruct.Token {
		log.Println("token already registered")
		http.Error(w, "token already registered", http.StatusBadRequest)
		return
	} else if api.BlueBotToken == "" {
		selected_color = game.Blue
	} else {
		log.Println("roster already full")
		http.Error(w, "roster already full", http.StatusBadRequest)
		return
	}
	err = api.Register(selected_color, reqStruct.Token, w)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *Api) Register(color game.PieceColor, token string, w http.ResponseWriter) (err error) {
	switch color {
	case game.Red:
		a.RedBotToken = token
	case game.Blue:
		a.BlueBotToken = token
	}
	log.Printf("%v Bot Registered: %s", color, token)
	respBody, err := json.Marshal(PostRegisterResp{color})
	if err != nil {
		return err
	}
	_, err = w.Write(respBody)
	if err != nil {
		return err
	}
	return nil
}

type PostPlayReq struct {
	Token string `json:"token"`
	Move  string `json:"move"`
}

func HandlePost_Play(w http.ResponseWriter, r *http.Request) {
	var reqStruct PostPlayReq
	err := json.NewDecoder(r.Body).Decode(&reqStruct)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if reqStruct.Token == "" {
		log.Println("token is required")
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}
	if reqStruct.Move == "" {
		log.Println("move is required")
		http.Error(w, "move is required", http.StatusBadRequest)
		return
	}
	stateMutex.Lock()
	defer stateMutex.Unlock()
	is_it_their_turn, err := api.IsItTheirTurn(reqStruct.Token)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !is_it_their_turn {
		log.Println("it's not your turn you silly goose")
		http.Error(w, "it's not your turn you silly goose", http.StatusBadRequest)
		return
	}
	move, err := api.State.LegalMoveFromString(reqStruct.Move)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// we know the move is legal
	err = api.State.ApplyLegalMove(move)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.State.DisplayBoard(DISPLAY_SHOULD_SHOW_HEIGHT)
}

func (a Api) IsItTheirTurn(token string) (is_it_their_turn bool, err error) {
	if token != a.RedBotToken && token != a.BlueBotToken {
		return is_it_their_turn, fmt.Errorf("unknown token")
	}
	if a.State.TurnNumber%2 == 0 {
		return token == a.RedBotToken, nil
	} else {
		return token == a.BlueBotToken, nil
	}
}
