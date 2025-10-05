package scanner

import (
	"errors"
	"github.com/shubhdevelop/YAKVS/token"
	"strconv"
)

type Scanner struct {
	Source  string
	Tokens  []token.Token
	start   int
	current int
}

func (s *Scanner) advance() rune {
	if s.isAtEnd() {
		return '\000'
	}
	runes := []rune(s.Source)
	ch := runes[s.current]
	s.current++
	return ch
}

func (s *Scanner) scanToken() {
	for !s.isAtEnd() {
		s.start = s.current
		ch := s.advance()
		if ch == '\000' {
			break
		}
		switch ch {
		case rune(token.SIMPLE_STRING):
			s.addToken(token.SIMPLE_STRING, nil)
		case rune(token.SIMPLE_ERROR):
			s.addToken(token.SIMPLE_ERROR, nil)
		case rune(token.NUMBER):
			s.addToken(token.NUMBER, nil)
		case rune(token.NULL):
			s.addToken(token.NULL, nil)
		case rune(token.DOUBLE):
			s.addToken(token.DOUBLE, nil)
		case rune(token.BOOLEAN):
			s.addToken(token.BOOLEAN, nil)
		case rune(token.BLOB):
			s.addToken(token.BLOB, nil)
		case rune(token.BLOB_STRING):
			s.addToken(token.BLOB_STRING, nil)
		case rune(token.VERBATIM_STRING):
			s.addToken(token.VERBATIM_STRING, nil)
		case rune(token.ARRAY):
			s.addToken(token.ARRAY, nil)
		case rune(token.SET):
			s.addToken(token.SET, nil)
		case rune(token.PUSH_DATA):
			s.addToken(token.PUSH_DATA, nil)
		case rune(token.CARRIAGE_RETURN):
			s.addToken(token.CARRIAGE_RETURN, nil)
		case rune(token.NEW_LINE):
			s.addToken(token.NEW_LINE, nil)
		case '"':
			s.string()
		default:
			if s.isDigit(ch) {
				s.number()
			} else if s.isAlpha(ch) {
				s.identifier()
			} else {
				// Skip whitespace and other characters
				continue
			}
		}
	}
}

func (s *Scanner) addToken(t token.FirstByte, literal interface{}) {
	if s.start >= len(s.Source) || s.current > len(s.Source) {
		return
	}
	text := s.Source[s.start:s.current] // substring
	tok := token.Token{
		Type:    t,
		Lexeme:  text,
		Literal: literal,
	}
	s.Tokens = append(s.Tokens, tok)
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len([]rune(s.Source)) {
		return '\000'
	}
	runes := []rune(s.Source)
	return runes[s.current+1]
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	runes := []rune(s.Source)
	return runes[s.current]
}

func (s *Scanner) isAlpha(c rune) bool {
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_') {
		return true
	}
	return false
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		s.advance()
	}
	if s.isAtEnd() {
		// Handle unterminated string - could return error or continue
		return
	}
	s.advance()
	if s.start+1 < s.current-1 && s.start+1 < len(s.Source) && s.current-1 <= len(s.Source) {
		value := s.Source[s.start+1 : s.current-1]
		s.addToken(token.STRING_VAL, value)
	} else {
		s.addToken(token.STRING_VAL, "")
	}
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}
	if s.start < s.current && s.start < len(s.Source) && s.current <= len(s.Source) {
		value := s.Source[s.start:s.current]
		valueInFloat, err := strconv.ParseFloat(value, 64)
		if err != nil {
			// Handle parse error - could return error or use string value
			s.addToken(token.STRING_VAL, value)
			return
		}
		s.addToken(token.NUM_VAL, valueInFloat)
	}
}

func (s *Scanner) isDigit(c rune) bool {
	if c >= '0' && c <= '9' {
		return true
	}
	return false
}

func (s *Scanner) ScanTokens() ([]token.Token, error) {
	if len(s.Source) == 0 {
		return nil, errors.New("source is empty")
	}
	s.scanToken()
	return s.Tokens, nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len([]rune(s.Source))
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	if s.start < s.current && s.start < len(s.Source) && s.current <= len(s.Source) {
		text := s.Source[s.start:s.current]
		s.addToken(token.STRING_VAL, text)
	}
}
