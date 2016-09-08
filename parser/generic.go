package parser

import (
	"encoding/json"
	"io"
)

type Generic struct {
	File
}

func (g *Generic) Build(r io.Reader, t interface{}) error {
	err := json.NewDecoder(r).Decode(&t)
	if err != nil {
		return parseErrorFactory(JSON_ERROR, err.Error())
	}

	return nil
}

func (g *Generic) ParserFactory(filePath string, t interface{}) error {
	relationFile, err := g.Parse(filePath)
	if err != nil {
		return err
	}
	if err := g.Build(relationFile, &t); err != nil {
		return err
	}

	return nil
}
