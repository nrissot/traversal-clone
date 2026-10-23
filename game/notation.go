package game

import (
	"fmt"
	"strconv"
)

func NotationToCoords(pos string) (x, y int, err error) {
	letter := pos[0]
	if letter < 97 || letter > 106 {
		//      ^ 'a'          ^ 'j'
		return x, y, fmt.Errorf("[NotationToCoords] invalid letter (expected [a..j], got '%s')", string(letter))
	}
	number, err := strconv.Atoi(pos[1:])
	if err != nil || number < 1 || number > 10 {
		// thanks lazy eval :)
		if pos[1:] == "*" {
			number = 0
		} else {
			return x, y, fmt.Errorf("[NotationToCoords] invalid number (expected [1..10], got '%s')", string(pos[1:]))
		}
	}
	x = int(letter) - 97
	y = number - 1
	return x, y, nil
}

func CoordsToNotation(x, y int) (pos string, err error) {
	if x < 0 || x > 9 {
		return "", fmt.Errorf("[CoordsToNotation] invalid x value (expected [0..9], got '%d')", x)
	}
	if y < 0 || y > 9 {
		if y == -1 || y == 10 {
			return fmt.Sprintf("%s*", string(byte(97+x))), nil
		}
		return "", fmt.Errorf("[CoordsToNotation] invalid y value (expected [0..9], got '%d')", y)
	}
	return fmt.Sprintf("%s%d", string(byte(97+x)), y+1), nil
}
