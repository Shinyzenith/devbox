package main

import "errors"

type TokenKind int

const (
	TK_WORD TokenKind = iota + 1
	TK_EQUALS
	TK_LBRACE
	TK_RBRACE
	TK_COMMA
	TK_DOUBLE_QUOTE
	TK_SINGLE_QUOTE
	TK_INTEGER
	TK_SEMICOLON
)

type SourceFile struct {
	FileName string
	contents string
	Char     byte
	Position Position
	IsEmpty  bool
}

func (f SourceFile) NextChar() (byte, error) {
	if f.Position.index >= len(f.contents) {
		f.IsEmpty = true
		return 0, errors.New("Reached end of file")
	}

	f.Position.index += 1

	f.Char = f.contents[f.Position.index]

	return f.Char, nil
}

func NewSourceFile(contents string, fileName string) SourceFile {
	return SourceFile{
		fileName,
		contents,
		contents[0],
		Position{1, 1, 0},
		false,
	}
}

type Token struct {
	lexeme string
	file   string
	pos    PositionRange
	kind   TokenKind
}

type Position struct {
	lineNo int
	column int
	index  int
}

func (p Position) Next() {
	p.column += 1
	p.index += 1
}

func (p Position) NextLine() {
	p.lineNo += 1
	p.column = 1
	p.index += 1
}

type PositionRange struct {
	start Position
	end   Position
}

func (p PositionRange) Next() {
	p.end.Next()
}

func (p PositionRange) NextLine() {
	p.start.NextLine()
	p.end.NextLine()
}

type Lexer struct{}

func (l Lexer) Lex(file SourceFile) []Token {
	tokens := []Token{}

	for !file.IsEmpty {
		position := PositionRange{file.Position, file.Position}

		switch {
		case file.Char == '{':
			tokens = append(tokens, Token{"{", file.FileName, position, TK_LBRACE})
			file.NextChar()

		case file.Char == '}':
			tokens = append(tokens, Token{"}", file.FileName, position, TK_RBRACE})
			file.NextChar()

		case file.Char == ',':
			tokens = append(tokens, Token{",", file.FileName, position, TK_COMMA})
			file.NextChar()

		default:
			file.NextChar()
		}
	}

	return tokens
}
