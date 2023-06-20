package internal

import (
	"github.com/EnergoStalin/GoTavern/pkg/character"
)

type CharacterStore interface {
	GetCharacters() (chars []*character.Card, err error)
	SaveCharacter(card *character.Card, overwrite bool) (err error)
}
