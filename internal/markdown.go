package internal

import (
	"fmt"
	"io"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func parseString(data []byte) ([]byte, ast.Node, error) {
	markdown := goldmark.New(
		goldmark.WithExtensions(),
	)

	doc := markdown.Parser().Parse(text.NewReader(data))
	return data, doc, nil
}

func parseFile(path string) ([]byte, ast.Node, error) {

	file, err := os.Open(path)
	if err != nil {
		Logger.Log(fmt.Sprintf("Failed reading \"%s\"", path))
		return nil, nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		Logger.Log(fmt.Sprintf("Failed reading \"%s\"", path))
		return nil, nil, err
	}

	return parseString(data)

}
