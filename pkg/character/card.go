package character

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"path"

	png "github.com/dsoprea/go-png-image-structure"
)

type Card struct {
	Path      string
	Image     []byte
	Character Character
}

func newTEXtChunk(bytes []byte) *png.Chunk {
	c := &png.Chunk{
		Length: uint32(len(bytes)),
		Type:   "tEXt",
		Data:   bytes,
	}

	c.UpdateCrc32()

	return c
}

func NewFromFile(path string, metaOnly bool) (char *Card, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	char, err = NewFromReader(f, metaOnly)
	if err != nil {
		return nil, err
	}

	char.Path = path

	return
}

func NewFromReader(r io.Reader, metaOnly bool) (char *Card, err error) {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	char, err = NewFromBytes(bytes, metaOnly)
	if err != nil {
		return nil, err
	}

	if !metaOnly {
		char.Image = bytes
	}
	return
}

func NewFromBytes(b []byte, metaOnly bool) (char *Card, err error) {
	metaBytes, err := extractMetadata(b)
	if err != nil {
		return nil, err
	}

	metaBytes, err = decodeMetadata(metaBytes)
	if err != nil {
		return nil, err
	}

	meta, err := UnmarshalMetadata(metaBytes)
	if err != nil {
		return nil, err
	}

	char = &Card{
		Character: *meta,
	}

	return
}

func (c *Card) read() error {
	f, err := os.Open(c.Path)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	c.Image = b
	return nil
}

func (c *Card) WriteToWriter(w io.Writer) error {
	if c.Image == nil {
		err := c.read()
		if err != nil {
			return err
		}
	}

	slice, err := c.buildChunkSlice()
	if err != nil {
		return err
	}

	return slice.WriteTo(w)
}

func (c *Card) WriteToFile(root string) (err error) {
	p := path.Join(root, c.Character.Name+".png")
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	err = c.WriteToWriter(f)
	if err != nil {
		return err
	}

	c.Path = p

	return
}

func extractMetadata(bytes []byte) (meta []byte, err error) {
	p := png.NewPngMediaParser()
	intfc, err := p.ParseBytes(bytes)
	if err != nil {
		return nil, err
	}

	cs := intfc.(*png.ChunkSlice)
	ch := cs.Index()["tEXt"][0]

	return ch.Data, err
}

func decodeMetadata(metaBytes []byte) (meta []byte, err error) {
	in := bytes.NewBuffer(metaBytes[6:])
	out := bytes.NewBuffer(nil)

	io.Copy(out, base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(in.Bytes())))

	return out.Bytes(), nil
}

func (c *Card) encodeMetadata() (meta []byte, err error) {
	chara, err := c.Character.Marshal()
	if err != nil {
		return nil, err
	}

	out := bytes.NewBuffer(nil)
	enc := base64.NewEncoder(base64.StdEncoding, out)
	defer enc.Close()

	_, err = enc.Write(chara)
	if err != nil {
		return nil, err
	}
	enc.Close()

	res := bytes.NewBuffer([]byte("chara"))
	res.WriteByte(0)
	res.Write(out.Bytes())

	return res.Bytes(), nil
}

func (c *Card) buildChunkSlice() (*png.ChunkSlice, error) {
	metaBytes, err := c.encodeMetadata()
	if err != nil {
		return nil, err
	}

	p := png.NewPngMediaParser()
	img, _ := p.ParseBytes(c.Image)

	ch := img.(*png.ChunkSlice)

	chunks := make([]*png.Chunk, 0)
	hasText := false
	for _, c := range ch.Chunks() {
		if c.Type == "tEXt" {
			chunks = append(chunks, newTEXtChunk(metaBytes))
			hasText = true
		} else {
			chunks = append(chunks, c)
		}
	}

	if !hasText {
		iend := chunks[len(chunks)-1]
		chunks[len(chunks)-1] = newTEXtChunk(metaBytes)
		chunks = append(chunks, iend)
	}

	return png.NewChunkSlice(chunks), nil
}
