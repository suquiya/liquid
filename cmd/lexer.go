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

package cmd

import "bufio"

//Lexer is lexer for sethead. Lexer tokenizes comment part.
type Lexer struct {
	s *bufio.Reader
	t *Token
}

//TokenType is type of Token
type TokenType int

const (
	//CommentBlockToken represents Block of comments
	CommentBlockToken TokenType = iota
	//BlankLineToken represents empty line, only contains "" until new line code.
	BlankLineToken
	//PackageDeclarationToken is declaration line of package
	PackageDeclarationToken
	//OtherToken is other Token.
	OtherToken
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
