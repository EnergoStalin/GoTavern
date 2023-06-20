package main

import (
	"os"

	internal "github.com/EnergoStalin/GoTavern/internal/character"
)

func main() {
	store := &internal.SillyCharacterStore{
		Path: os.Args[1],
	}

	chars, err := store.GetCharacters(true)
	if err != nil {
		panic(err)
	}

	err = chars[0].WriteToFile("./")
	if err != nil {
		panic(err)
	}
}
