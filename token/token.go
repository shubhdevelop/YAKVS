package token

import (
	"fmt"
)

// firstByte represents all possible token types
type FirstByte rune

const (
	SIMPLE_STRING   FirstByte = '+'
	SIMPLE_ERROR    FirstByte = '-'
	NUMBER          FirstByte = ':'
	NULL            FirstByte = '_'
	DOUBLE          FirstByte = ','
	BOOLEAN         FirstByte = '#'
	BLOB            FirstByte = '!'
	VERBATIM_STRING FirstByte = '='

	// Aggregate Types
	BLOB_STRING     FirstByte = '$'
	ARRAY           FirstByte = '*'
	MAP             FirstByte = '%'
	SET             FirstByte = '~'
	PUSH_DATA       FirstByte = '>'
	CARRIAGE_RETURN FirstByte = '\r'
	NEW_LINE        FirstByte = '\n'
	NUM_VAL         FirstByte = 'n'
	STRING_VAL      FirstByte = 's'
)

func (t FirstByte) String() string {
	switch t {
	case SIMPLE_STRING:
		return "SIMPLE_STRING"
	case SIMPLE_ERROR:
		return "SIMPLE_ERROR"
	case NUMBER:
		return "NUMBER"
	case NULL:
		return "NULL"
	case DOUBLE:
		return "DOUBLE"
	case BOOLEAN:
		return "BOOLEAN"
	case BLOB:
		return "BLOB"
	case VERBATIM_STRING:
		return "VERBATIM_STRING"
	case BLOB_STRING:
		return "BLOB_STRING"
	case ARRAY:
		return "ARRAY"
	case MAP:
		return "MAP"
	case SET:
		return "SET"
	case PUSH_DATA:
		return "PUSH_DATA"
	case CARRIAGE_RETURN:
		return "CARRIAGE_RETURN"
	case NEW_LINE:
		return "NEW_LINE"
	case NUM_VAL:
		return "NUM_VAL"
	case STRING_VAL:
		return "STRING_VAL"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(t))
	}
}

type Token struct {
	Type    FirstByte   
	Lexeme  string      
	Literal interface{} 
}

func (t Token) String() string {
	if t.Literal == nil {
		return fmt.Sprintf("%v %s", t.Type, t.Lexeme)
	}
	return fmt.Sprintf("%v %s %v", t.Type, t.Lexeme, t.Literal)
}
