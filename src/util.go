package main

import "strings"

func (p *Prettier) lastIs(kind string) bool {
	if p.cursor == nil || p.cursor.Node() == nil || p.cursor.Node().PrevSibling() == nil {
		return false
	}
	return p.cursor.Node().PrevSibling().Kind() == kind
}

func (p *Prettier) nextIs(kind string) bool {
	if p.cursor == nil || p.cursor.Node() == nil || p.cursor.Node().NextSibling() == nil {
		return false
	}
	return p.cursor.Node().NextSibling().Kind() == kind
}

func (p *Prettier) curText() string {
	return p.cursor.Node().Utf8Text(p.content)
}

func (p *Prettier) sepLine() {
	if strings.HasSuffix(string(p.buf.Bytes()), "\n\n") {
		return
	}
	p.write([]byte{'\n'})
}

func (p *Prettier) sepSpace() {
	c := p.buf.Bytes()[p.buf.Len()-1]
	if c == '\n' || c == '\r' || c == '\t' || c == ' ' {
		return
	}
	p.buf.WriteByte(' ')
}

func (p *Prettier) writeIndent() {
	for i := 0; i < p.indent; i++ {
		p.write([]byte{'\t'})
	}
}

func (p *Prettier) writeNewLine() {
	p.buf.WriteByte('\n')
}

func (p *Prettier) write(data []byte) {
	p.buf.Write(data)
}

func (p *Prettier) writeString(s string) {
	p.buf.WriteString(s)
}
