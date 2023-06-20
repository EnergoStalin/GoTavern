package character_test

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/EnergoStalin/GoTavern/pkg/character"

	log "github.com/dsoprea/go-logging"
)

func GetDataPath() string {
	wd, _ := os.Getwd()
	return path.Join(wd, "data")
}

func TestRead(t *testing.T) {

	filepath.Walk(GetDataPath(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			t.FailNow()
			return err
		} else if info.IsDir() {
			return nil
		}

		card, err := character.NewFromFile(path, true)
		log.PanicIf(err)
		// t.Log("Loaded Character...")

		t.Log(card.Character.Name)

		return err
	})
}

func TestWrite(t *testing.T) {
	filepath.Walk(GetDataPath(), func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			t.FailNow()
			return err
		} else if info.IsDir() {
			return nil
		}

		card, err := character.NewFromFile(p, false)
		log.PanicIf(err)
		// t.Log("Loaded Character...")

		op := path.Join(GetDataPath())
		card.WriteToFile(op)

		card, err = character.NewFromFile(card.Path, false)
		log.PanicIf(err)

		t.Log(card.Path)

		return err
	})
}
