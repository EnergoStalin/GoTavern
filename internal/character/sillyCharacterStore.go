package internal

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/EnergoStalin/GoTavern/pkg/character"
)

type SillyCharacterStore struct {
	CharacterStore
	Path string
}

func (s *SillyCharacterStore) GetCharacters(metaOnly bool) (chars []*character.Card, err error) {
	chars = make([]*character.Card, 0)
	filepath.Walk(s.Path, func(p string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		char, err := character.NewFromFile(p, metaOnly)
		if err != nil {
			return err
		}
		chars = append(chars, char)

		return err
	})

	return
}

func (s *SillyCharacterStore) SaveCharacter(card *character.Card, overwrite bool) (err error) {
	name := path.Join(s.Path, card.Character.Name+".png")
	if _, err := os.Stat(name); !os.IsNotExist(err) && !overwrite {
		return errors.New("character file exists")
	}

	f, err := os.Create(name)
	if err != nil {
		return
	}
	defer f.Close()

	card.WriteToFile(s.Path)

	return
}
