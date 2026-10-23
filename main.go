package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

const SHOULD_DISPLAY = true
const DISPLAY_SHOULD_SHOW_HEIGHT = false

func main() {
	// setup the game
	resp, err := game.GetBoard(game.Random)
	if err != nil {
		log.Fatal(err)
	}
	api.State = game.InitializeBoardState(*resp)
	api.State.DisplayBoard(DISPLAY_SHOULD_SHOW_HEIGHT)

	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleNotFound)

	// get data
	mux.HandleFunc("GET /board", HandleGet_Board)
	mux.HandleFunc("GET /state", HandleGet_State)

	// move related
	mux.HandleFunc("GET /moves", HandleGet_Moves)
	mux.HandleFunc("GET /moves/struct", HandleGet_Moves_Struct)
	mux.HandleFunc("GET /is-move-legal/{move}", HandleGet_IsMoveLegal)

	// move translations
	// struct -> notation
	mux.HandleFunc("POST /translate", HandlePost_Translate)
	// notation -> struct
	mux.HandleFunc("GET /translate/{move}", HandleGet_Translate)

	mux.HandleFunc("POST /register", HandlePost_Register)
	mux.HandleFunc("POST /play", HandlePost_Play)

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", http.StatusNotFound)
}

func HandleGet_Board(w http.ResponseWriter, r *http.Request) {
	stateMutex.RLock()
	defer stateMutex.RUnlock()

	err := json.NewEncoder(w).Encode(api.State.Board)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleGet_State(w http.ResponseWriter, r *http.Request) {
	stateMutex.RLock()
	defer stateMutex.RUnlock()

	err := json.NewEncoder(w).Encode(api.State)
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
	err = json.NewEncoder(w).Encode(notations)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if SHOULD_DISPLAY {
		api.State.DisplayAllowedMoves(DISPLAY_SHOULD_SHOW_HEIGHT, moves)
	}
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

	err = json.NewEncoder(w).Encode(moves)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	err = json.NewEncoder(w).Encode(PostRegisterResp{color})
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

func HandleGet_IsMoveLegal(w http.ResponseWriter, r *http.Request) {
	nota := r.PathValue("move")
	stateMutex.RLock()
	_, err := api.State.LegalMoveFromString(nota)
	stateMutex.RUnlock()
	isMoveLegal := true
	if err != nil {
		if strings.Contains(err.Error(), "unable to check move legality") {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else if strings.Contains(err.Error(), "Illegal Move") {
			isMoveLegal = false
		} else {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	err = json.NewEncoder(w).Encode(isMoveLegal)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func HandlePost_Translate(w http.ResponseWriter, r *http.Request) {
	var move game.Move
	err := json.NewDecoder(r.Body).Decode(move)
	if err != nil || move.String() == "####" {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(move.String())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleGet_Translate(w http.ResponseWriter, r *http.Request) {
	nota := r.PathValue("move")
	stateMutex.RLock()
	move, err := api.State.LegalMoveFromString(nota)
	stateMutex.RUnlock()
	if err != nil {
		if strings.Contains(err.Error(), "unable to check move legality") {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	err = json.NewEncoder(w).Encode(move)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
