// Copyright (c) 2019 suquiya
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package tools

import (
	"bufio"
	"bytes"
)

const (
	slash    = '/'
	asterisk = '*'
	cr       = '\r'
	lf       = '\n'
)

var (
	startLineComment []byte
	startWrapComment []byte
	endWrapComment   []byte
	crlf []byte
)

func init() {
	arr := func(a byte, b byte) []byte {
		barr := make([]byte, 2, 2)
		barr[0] = a
		barr[1] = b
		return barr
	}
	startLineComment = arr(slash, slash)
	startWrapComment = arr(slash, asterisk)
	endWrapComment = arr(asterisk, slash)
	crlf = arr(cr,lf)
}

//Lexer is lexer for sethead. Lexer tokenizes comment part.
type Lexer struct {
	s         *bufio.Reader
	t         Token
	nextBytes []byte
	next      TokenType
	err       error
}

//NewLexer returns Lexer setted br as scanner. if reader's stream empty, it returns nil.
func NewLexer(br *bufio.Reader) *Lexer {
	b, err := br.ReadByte()
	if err != nil {
		return nil
	}
	return &Lexer{br, nil, getSingleByteArray(b), UnknownToken, nil}
}

//Next scans the next token and return its Token struct
func (l *Lexer) Next() bool {
	switch l.next {
	case UnknownToken:
		nextByte := l.nextBytes[0]
		switch nextByte {
		case cr:
			nb, canR := l.readByte()
			if canR {
				if nb == lf {
					l.t = &BlankLine{[]byte{cr, lf}}
					l.nextBytes[0], canR = l.readByte()
					if canR {
						l.next = UnknownToken
						return true
					}
					l.next = EOFToken
					return true
				}
				l.t = &BlankLine{[]byte{cr}}
				l.next = UnknownToken
				l.nextBytes[0] = nb
				return true
			}
			l.next = EOFToken
			l.t = &BlankLine{[]byte{cr}}
		case lf:
			l.t = &BlankLine{[]byte{lf}}
			nb, canR := l.readByte()
			l.nextBytes[0] = nb
			if !canR {
				l.next = EOFToken
			}
		case slash:
			nb, canR := l.readByte()
			if canR {
				switch nb {
				case slash:
					ct := l.GetLinesBlock()
					l.t = ct
				case asterisk:
					ct := l.GetWrapCommentBlock()
					l.t = ct
				}
			} else {
				l.next = EOFToken
				r := getSingleByteArray('/')
				l.t = &NotCommentToken{OtherToken, r, r}
			}

		}
	}

	return false
}

//GetLinesBlock detect and return multi line comment block.
func (l *Lexer) GetLinesBlock() *CommentBlock {
	ct := &CommentBlock{Lines, []byte{}, []byte{}}
	ct.r = append(ct.r, startLineComment...)
	line, CanR := l.readBytes(lf)
	if CanR {
		ct.r = append(ct.r, line...)
		ct.content = append(ct.content, bytes.TrimSpace(line)...)
		detect := true
		for detect{
			line, CanR = l.readBytes(lf)
			if CanR{
				if bytes.HasPrefix(line, startLineComment){
					ct.r = append(ct.r, line...)
					ct.content = append(ct.content, bytes.TrimSpace(line)...)
					//detect = true
				}else{
					switch {
					case bytes.HasPrefix(line, startWrapComment):
						l.next = CommentWrapBlockToken
					case bytes.Equal(line, crlf) || line[0] == lf:
						l.next = BlankLineToken
					case bytes.HasPrefix("package"):
						l.next = PackageDeclarationToken
					default:
						l.next = NormalStringToken
					}
					l.nextBytes = line
					detect = false
				}
			}else{
				l.next = EOFToken
				detect = canR
			}
		}
		return ct
	}
	l.next = EOFToken
	return ct

}

//GetWrapCommentBlock return wrap comments
func GetWrapCommentBlock() *CommentBlock{
	cb := &CommentBlock{Wrap, []byte{},[]byte{}}
	cb.r = append(cb.r, startWrapComment)

	seek := true
	for seek{
		block, CanR := l.readBytes(slash)
		if CanR{
			blen := len(block)
			if blen > 1{
				cb.r = append(cb.r,block)
				if block[blen-2] == asterisk{
					seek = false
				}
			}
		}else{
			seek = CanR
			l.next = EOFToken
		}
	}
	
}

func getBytesArray(b byte, n int) []byte {
	bArray := make([]byte, n, n)
	for i := 0; i < n; i++ {
		bArray[i] = b
	}
	return bArray
}

func getSingleByteArray(b byte) (bArray []byte) {
	bArray = make([]byte, 0, 1)
	bArray = append(bArray, b)
	return
}

func (l *Lexer) readBytes(b byte) ([]byte, bool) {
	if l.err != nil {
		return nil, false
	}
	barr, err := l.s.ReadBytes(b)
	if len(barr) <1　&& err != nil{
		return nil, false
	}
	l.err = err
	return barr, true
}

func (l *Lexer) readByte() (byte, bool) {
	if l.err != nil {
		return 0, false
	}
	b, err := l.s.ReadByte()
	l.err = err
	if err != nil && b == 0{
		return b,false
	}
	return b, true
}

//TokenType is type of Token
type TokenType int

const (
	//CommentLineBlockToken represents Block of line comments
	CommentLineBlockToken TokenType = iota
	//CommentWrapBlockToken represent Comment block of /* and */ wrap.
	CommentWrapBlockToken
	//BlankLineToken represents empty line, only contains "" until new line code.
	BlankLineToken
	//PackageDeclarationToken is declaration line of package
	PackageDeclarationToken
	//NormalStringToken represents a normal string block - not comment, blankLine, EOF
	NormalStringToken
	//EOFToken represents EOF
	EOFToken
	//OtherToken is other Token.
	OtherToken
	//UnknownToken is Unknown. This type is used in lexing.
	UnknownToken
)

//Token is interface for Lexer
type Token interface {
	Type() TokenType
	Raw() []byte
	Content() []byte
	ContentString() string
}

//NotCommentToken represents tokens that is not comment block.
type NotCommentToken struct {
	tt TokenType
	r  []byte
	c  []byte
}

//Type returns n's TokenType
func (n *NotCommentToken) Type() TokenType {
	return n.tt
}

//Raw returns n's RawData
func (n *NotCommentToken) Raw() []byte {
	return n.r
}

//Content returns n's content with byte array form
func (n *NotCommentToken) Content() []byte {
	return n.c
}

//ContentString returns n's content with string form
func (n *NotCommentToken) ContentString() string {
	return string(n.Content())
}

//BlankLine represents blankLineBlock
type BlankLine struct {
	NLCode []byte
}

//Type returns blankLineToken
func (b *BlankLine) Type() TokenType {
	return BlankLineToken
}

//Raw returns raw of block
func (b *BlankLine) Raw() []byte {
	return b.NLCode
}

//Content returns content (nil)
func (b *BlankLine) Content() []byte {
	return nil
}

//ContentString returns empty string
func (b *BlankLine) ContentString() string {
	return ""
}

//CommentType represents type of comment
type CommentType int

//CommentBlock is data for comment block of source code.
type CommentBlock struct {
	ct      CommentType
	r       []byte
	content []byte
}

//Type method returns c's TokenType
func (c *CommentBlock) Type() TokenType {
	return CommentBlockToken
}

//CommentMethod returns c's CommentType
func (c *CommentBlock) CommentMethod() CommentType {
	return c.ct
}

//Raw returns c's raw byte array.
func (c *CommentBlock) Raw() []byte {
	return c.r
}

//Content returns contents of comment block in the form of byte array.
func (c *CommentBlock) Content() []byte {
	return c.content
}

//ContentString return contents of comment block in the form of string.
func (c *CommentBlock) ContentString() string {
	return string(c.content)
}

const (
	//Lines means type of "//"
	Lines CommentType = iota
	//Wrap means type of "/*" ~ "*/"
	Wrap
)
