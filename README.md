# API Documentation

this is a clone of the game [Traversal](https://www.traversalgame.com/) available as an API. This project is meant to be used in a bot making competion that a bunch of friend and I are planning to do.

This project uses only go 1.26.9 and the standard library, no external packages.

You can simply git clone the project and run `go run .` to execute it, or run `go buid` to get an executable.

## struct definition:

### piece
```application/json
{
  "color":      string, // ∈ {"Red", "Blue"}
  "type":       string,	// ∈ {"S", "V", "E"}
  "direction":  string,	// ∈ {"W", "E"}
  "x":			    int,	  // ∈ [0..9]
  "y":			    int,	  // ∈ [0..9]
  "alive":		  bool,	  // ∈ {true, false}
}
```
```example
ex: Red Eastern Vanguard, located on g1, currently alive
{
  "color": "Red",
  "type": "V",
  "direction": "E",
  "x": 6,
  "y": 0,
  "alive": true
},
```

move
```application/json
{
  "from_x":	int,	  // ∈ [0..9]
  "from_y":	int,	  // ∈ [0..9]
  "to_x":		int,	  // ∈ [0..9]
  "to_y":		int,	  // ∈ [0..9]
  "piece":	piece,  // cf piece
  "takes":	bool,	  // ∈ {true, false}
}
```
```example
ex: first move, EWe3
{
  "from_x": 4,
  "from_y": 0,
  "to_x": 4,
  "to_y": 2,
  "piece": {
    "color": "Red",
    "type": "E",
    "direction": "W",
    "x": 4,
    "y": 0,
    "alive": true
  },
  "takes": false
},
```

## `GET /board`

Get the board as list of list of int.

### Responses:

#### 200 OK
```application/json
[10][10]int
```
```example
[
  [1,3,2,1,1,2,3,3,3,3],
  [3,2,2,1,0,2,3,2,2,3],
  [4,3,3,2,1,0,1,1,1,3],
  [4,3,3,4,2,2,2,0,1,1],
  [3,3,2,2,2,2,2,1,1,2],
  [2,2,2,2,2,2,0,0,0,2],
  [2,2,2,1,2,2,1,1,1,2],
  [3,2,0,2,2,2,2,3,2,1],
  [3,2,2,3,3,3,3,3,1,1],
  [3,3,4,3,3,3,3,3,2,2]
]
```

## `GET /state`

Get the current state of the game.

### Responses:

#### 200 OK
```application/json
{
  "board": [10][10]int,
  "pieces": [12]piece,
  "turn_number":int,
  "score":[2]int,
  "has_traversed":int,
  "history":[]move,
}
```
```example
{
  "board": [
    [1,3,2,1,1,2,3,3,3,3],
    [3,2,2,1,0,2,3,2,2,3],
    [4,3,3,2,1,0,1,1,1,3],
    [4,3,3,4,2,2,2,0,1,1],
    [3,3,2,2,2,2,2,1,1,2],
    [2,2,2,2,2,2,0,0,0,2],
    [2,2,2,1,2,2,1,1,1,2],
    [3,2,0,2,2,2,2,3,2,1],
    [3,2,2,3,3,3,3,3,1,1],
    [3,3,4,3,3,3,3,3,2,2]
  ],
  "pieces": [
    {
      "color": "Red",
      "type": "S",
      "direction": "W",
      "x": 2,
      "y": 0,
      "alive": true
    },
    {
      "color": "Red",
      "type": "V",
      "direction": "W",
      "x": 3,
      "y": 0,
      "alive": true
    },
    {
      "color": "Red",
      "type": "E",
      "direction": "W",
      "x": 4,
      "y": 2,
      "alive": true
    },
    {
      "color": "Red",
      "type": "E",
      "direction": "E",
      "x": 5,
      "y": 0,
      "alive": true
    },
    {
      "color": "Red",
      "type": "V",
      "direction": "E",
      "x": 6,
      "y": 0,
      "alive": true
    },
    {
      "color": "Red",
      "type": "S",
      "direction": "E",
      "x": 7,
      "y": 0,
      "alive": true
    },
    {
      "color": "Blue",
      "type": "S",
      "direction": "W",
      "x": 2,
      "y": 9,
      "alive": true
    },
    {
      "color": "Blue",
      "type": "V",
      "direction": "W",
      "x": 3,
      "y": 9,
      "alive": true
    },
    {
      "color": "Blue",
      "type": "E",
      "direction": "W",
      "x": 4,
      "y": 9,
      "alive": true
    },
    {
      "color": "Blue",
      "type": "E",
      "direction": "E",
      "x": 5,
      "y": 7,
      "alive": true
    },
    {
      "color": "Blue",
      "type": "V",
      "direction": "E",
      "x": 6,
      "y": 9,
      "alive": true
    },
    {
      "color": "Blue",
      "type": "S",
      "direction": "E",
      "x": 7,
      "y": 9,
      "alive": true
    }
  ],
  "turn_number": 2,
  "score": [
    0,
    0
  ],
  "has_traversed": 0,
  "history": [
    {
      "from_x": 4,
      "from_y": 0,
      "to_x": 4,
      "to_y": 2,
      "piece": {
        "color": "Red",
        "type": "E",
        "direction": "W",
        "x": 4,
        "y": 0,
        "alive": true
      },
      "takes": false
    },
    {
      "from_x": 5,
      "from_y": 9,
      "to_x": 5,
      "to_y": 7,
      "piece": {
        "color": "Blue",
        "type": "E",
        "direction": "E",
        "x": 5,
        "y": 9,
        "alive": true
      },
      "takes": false
    }
  ]
}
```

## `GET /moves`

Get the list of legal moves in our custom notation format

``` 
id (2) | x (1?) | col (1) | row (1-2) | 

id:   ["E"|"V"|"S"]["W"|"E"]  required (piece identifier)
x:    "x"                     optionnal (capture ?)
col:  [a..j]                  required (destination col)
row:  [1..10]|"*"             required (destination row, if traversal: "*" instead)
```

### Response:

#### 200 OK
```application/json
[]string
```
```example
["SWd2","SWb2","SWc2","SWb1","SWc3","VWd3","VWe2","VWc2","VWb2","VWd2","VWe1","EWd2","EWf2","EWe4","EWe2","EWd3","EWg3","EWf3","EWe6","EWg4","EWc4","EWe7","EWe1","EWd4","EWd5","EWb3","EWf6","EWc5","EWe5","EWf5","EWh3","EWd6","EWg5","EWc3","EWf4","EEe1","EEh2","EEf3","EEg2","EEe2","EEf2","VEg2","VEg3","VEh2","VEf2","VEi2","SEh2","SEi1","SEi2","SEj1","SEj2"]
```

## `GET /moves/struct`

Get the list of legal moves in a structured format (cf struct definition section)

### Response:

#### 200 OK
```application/json
[]move
```
```example
[
  {
    "from_x": 2,
    "from_y": 0,
    "to_x": 2,
    "to_y": 1,
    "piece": {
      "color": "Red",
      "type": "S",
      "direction": "W",
      "x": 2,
      "y": 0,
      "alive": true
    },
    "takes": false
  },
  {
    "from_x": 2,
    "from_y": 0,
    "to_x": 1,
    "to_y": 0,
    "piece": {
      "color": "Red",
      "type": "S",
      "direction": "W",
      "x": 2,
      "y": 0,
      "alive": true
    },
    "takes": false
  },
  // [...],
  {
    "from_x": 7,
    "from_y": 0,
    "to_x": 9,
    "to_y": 1,
    "piece": {
      "color": "Red",
      "type": "S",
      "direction": "E",
      "x": 7,
      "y": 0,
      "alive": true
    },
    "takes": false
  }
]
```

## `GET /is-move-legal/{move}`

Check the legality of a move.

### Response:

#### 200 OK

If the move is legal
```application/json
true
```

If the move is correctly formed, but illegal (ie, you tried teleportation, necromancy, or violence without the ability to do so...)
```application/json
false
```

#### 400 Bad Request

You tried to verify a malformed move
```
examples: 
[LegalMoveFromString] EZe4 invalid piece side identifier (expected 'W' or 'E', got 'Z')
[LegalMoveFromString] EWez unable to process coordinates: [NotationToCoords] invalid number (expected [1..10], got 'z')
[LegalMoveFromString] EWz1 unable to process coordinates: [NotationToCoords] invalid letter (expected [a..j], got 'z')
...
```

## `POST /register`

Register your bot as a player at the start of the game, provide a "token" string to be used as an identifier to register plays (see below)

### Body

```application/json
{
  "token":<string>
}
```

### Responses

#### 200 OK

```application/json
{
  "color":string // ∈ {"Red", "Blue"}
}
```
```example
{
  "color":"Red"
}
```

#### 400 Bad Request

You tried to register twice, or attempted to register using the same token as your opponent.
```application/json
token already registered
```

You tried to register when two bots were already registered.
```application/json
roster already full
```

You tried to register without sending a token, or with an empty string as a token.
```
token is required
```



## `POST /play`

### Responses

#### 200 OK

no response body

#### 400 Bad Request

You tried to play when it was NOT your turn
```
it's not your turn you silly goose
```

You tried to play an Illegal move
```
[LegalMoveFromString] (Illegal Move) (expected [...], got 'EWe4')
```

You tried to play a malformed move
```
examples: 
[LegalMoveFromString] EZe4 invalid piece side identifier (expected 'W' or 'E', got 'Z')
[LegalMoveFromString] EWez unable to process coordinates: [NotationToCoords] invalid number (expected [1..10], got 'z')
[LegalMoveFromString] EWz1 unable to process coordinates: [NotationToCoords] invalid letter (expected [a..j], got 'z')
...
```

You tried to become a necromancer and play a dead piece
```
[LegalMoveFromString] unable to find an alive piece matching the provided identifier (expected [SW VW EW EE VE SE SW VW EW EE VE SE], got 'EW')
```

You tried to play with an unknown token
```
unknown token
```

You forgot to include your token
```
token is required
```

You forgot to include your move
```
move is required
```