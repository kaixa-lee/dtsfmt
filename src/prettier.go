package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/joelspadin/tree-sitter-devicetree/bindings/go"
	"github.com/tree-sitter/go-tree-sitter"
)

type Prettier struct {
	buf *bytes.Buffer

	content []byte

	indent int

	cursor *tree_sitter.TreeCursor
}

func (p *Prettier) Prettify(fn string) ([]byte, error) {
	if err := p.loadFile(fn); err != nil {
		return nil, err
	}

	parser := tree_sitter.NewParser()
	if err := parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_devicetree.Language())); err != nil {
		return nil, err
	}

	p.cursor = parser.Parse(p.content, nil).Walk()
	p.cursor.GotoFirstChild()
	for ; p.hasNext(); p.cursor.GotoNextSibling() {
		p.traverse()
	}

	return p.buf.Bytes(), nil
}

func (p *Prettier) hasNext() bool {
	return p.cursor != nil && p.cursor.Node() != nil && p.cursor.Node().NextSibling() != nil
}

func (p *Prettier) loadFile(fn string) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		return err
	}
	p.content = data
	p.buf = bytes.NewBuffer(nil)
	return nil
}

func (p *Prettier) traverse() {
	fmt.Println("node_kind:", p.cursor.Node().Kind())

	switch p.cursor.Node().Kind() {
	case NodeKindComment:
		p.WriteComment()
	case NodeKindFileVersion:
		p.WriteFileVersion()
	case NodeKindPreprocInclude:
		p.WritePreprocInclude()
	case NodeKindPreprocDef:
	case NodeKindPreprocFunctionDef:
	case NodeKindPreprocIfdef:
	case NodeKindLabeledItem:
	case NodeKindNode:
	case NodeKindProperty:
	case NodeKindStringLiteral:
	case NodeKindIntegerCells:
	default:
		// todo 如果没有具体的类型, 是否可以直接把str写下来
		p.WriteDefault()
	}
}

func (p *Prettier) WriteComment() {
	if !p.lastIs(NodeKindComment) {
		p.sep()
	}
	p.writeIndent()
	commentText := p.curText()
	if strings.HasPrefix(commentText, "//") {
		p.writeString("// ")
		p.writeString(strings.TrimSpace(strings.TrimPrefix(commentText, "//")))
	} else {
		p.writeString(commentText)
	}
	p.writeNewLine()
}

func (p *Prettier) WriteFileVersion() {
	p.sep()
	p.writeString(p.curText())
	p.writeNewLine()
	p.sep()
}

func (p *Prettier) WritePreprocInclude() {
	p.writeIndent()
	p.writeString("#include ")
	p.cursor.GotoFirstChild()
	p.cursor.GotoNextSibling()
	p.writeString(p.curText())
	p.writeNewLine()
	p.cursor.GotoParent()
}

func (p *Prettier) WriteDefault() {
	if !p.cursor.GotoFirstChild() {
		return
	}
	p.indent += 1
	for ; p.hasNext(); p.cursor.GotoNextSibling() {
		p.traverse()
	}
	p.indent -= 1
	p.cursor.GotoParent()
}
