package dtsfmt

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/joelspadin/tree-sitter-devicetree/bindings/go"
	"github.com/tree-sitter/go-tree-sitter"
)

type Prettier struct {
	isDebug bool

	integerCellSize int

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

	for {
		p.traverse()
		if !p.cursor.GotoNextSibling() {
			break
		}
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
	if p.isDebug {
		fmt.Println("node_kind:", p.cursor.Node().Kind())
	}

	// todo 应该尽可能地减少token枚举类型, 只处理关键类型
	switch p.cursor.Node().Kind() {
	case NodeKindComment:
		p.WriteComment()
	case NodeKindFileVersion:
		p.WriteFileVersion()
	case NodeKindPreprocInclude:
		p.WritePreprocInclude()
	case NodeKindNode:
		p.WriteNode()
	case NodeKindLeftBracket:
		p.WriteLeftBracket()
	case NodeKindRightBracket:
		p.WriteRightBracket()
	case NodeKindProperty:
		p.WriteProperty()
	case NodeKindStringLiteral:
		p.WriteStringLiteral()
	case NodeKindIntegerCells:
		p.WriteIntegerCells()
	default:
		p.WriteDefault()
	}
}

func (p *Prettier) WriteComment() {
	if !p.lastIs(NodeKindComment) {
		p.sepLine()
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
	p.sepLine()
	p.writeString(p.curText())
	p.writeNewLine()
	p.sepLine()
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

func (p *Prettier) WriteNode() {
	p.sepLine()
	if !p.lastIs(NodeKindColon) {
		p.writeIndent()
	}

	p.cursor.GotoFirstChild()
	p.writeString(p.curText())
	for p.cursor.GotoNextSibling() {
		p.traverse()
	}
	p.writeIndent()
	p.writeNewLine()
	p.cursor.GotoParent()
}

func (p *Prettier) WriteLeftBracket() {
	p.sepSpace()
	p.writeString(p.curText())
	if !p.nextIs(NodeKindRightBracket) {
		p.writeNewLine()
		p.indent += 1
	}
}

func (p *Prettier) WriteRightBracket() {
	if !p.lastIs(NodeKindLeftBracket) {
		p.indent -= 1
		p.writeIndent()
	}
	p.writeString(p.curText())
}

func (p *Prettier) WriteIntegerCells() {
	p.writeString("<")
	p.cursor.GotoFirstChild()

	cellSize, tmpCursor := 0, p.cursor.Copy()
	for tmpCursor.GotoNextSibling() {
		if tmpCursor.Node().Utf8Text(p.content) != ">" {
			cellSize += 1
		}
	}

	if cellSize >= p.integerCellSize {
		p.indent += 1
		p.writeNewLine()
		p.writeIndent()
	}

	cnt := 0
	for p.cursor.GotoNextSibling() {
		if p.curText() == ">" {
			break
		}
		if cnt >= p.integerCellSize {
			p.writeNewLine()
			p.writeIndent()
			cnt = 0
		}
		if !p.lastIs("<") {
			p.sepSpace()
		}

		p.writeString(p.curText())
		cnt += 1
	}

	if cellSize >= p.integerCellSize {
		p.indent -= 1
		p.writeNewLine()
		p.writeIndent()
	}
	p.writeString(">")

	p.cursor.GotoParent()
}

func (p *Prettier) WriteProperty() {
	p.writeIndent()
	p.cursor.GotoFirstChild()
	p.writeString(p.curText())
	for p.cursor.GotoNextSibling() {
		switch p.cursor.Node().Kind() {
		case NodeKindComma:
			p.writeString(", ")
		case NodeKindEq:
			p.writeString(" = ")
		case NodeKindSemiColon:
		default:
			p.traverse()
		}
	}
	p.writeString(";")
	p.writeNewLine()
	p.cursor.GotoParent()
	if p.nextIs(NodeKindNode) {
		p.writeNewLine()
	}
}

func (p *Prettier) WriteStringLiteral() {
	p.writeString(p.curText())
}

func (p *Prettier) WriteLT() {
	p.writeString("<")
	tmpCursor := p.cursor.Copy()
	cellSize := 0
	for tmpCursor.GotoNextSibling() {
		if tmpCursor.Node().Utf8Text(p.content) == ">" {
			break
		}
		cellSize += 1
	}

	if cellSize > p.integerCellSize {
		p.indent += 1
		p.writeNewLine()
		p.writeIndent()
	}
	cnt := 0
	for p.cursor.GotoNextSibling() {
		if cnt >= p.integerCellSize {
			p.writeNewLine()
			p.writeIndent()
			cnt = 0
		}
		p.writeString(p.curText())
		if !p.nextIs(">") && cnt != p.integerCellSize-1 {
			p.writeString(" ")
		}
		if p.nextIs(">") {
			if cnt >= p.integerCellSize {
				p.writeNewLine()
			}
			break
		}
		cnt += 1
	}
	if cellSize > 1 {
		p.indent -= 1
		p.writeNewLine()
		p.writeIndent()
	}
	p.writeString(">")
	p.cursor.GotoParent()
}

func (p *Prettier) WriteDefault() {
	if !p.cursor.GotoFirstChild() {
		p.writeString(p.curText())
		return
	}
	for {
		p.traverse()
		if !p.cursor.GotoNextSibling() {
			break
		}
	}
	p.cursor.GotoParent()
}
