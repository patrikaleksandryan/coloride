package syntax

const (
	// Classes of symbols
	CNone = iota
	CComment
	CIdent
	CKeyword
	CString
	CNumber
	CProcCall
)

var (
	keywords = map[string]bool{
		"break":       true,
		"default":     true,
		"func":        true,
		"interface":   true,
		"select":      true,
		"case":        true,
		"defer":       true,
		"go":          true,
		"map":         true,
		"struct":      true,
		"chan":        true,
		"else":        true,
		"goto":        true,
		"package":     true,
		"switch":      true,
		"const":       true,
		"fallthrough": true,
		"if":          true,
		"range":       true,
		"type":        true,
		"continue":    true,
		"for":         true,
		"import":      true,
		"return":      true,
		"var":         true,
	}
)

func isWhitespace(c rune) bool {
	return c <= ' '
}

func isNumeric(c rune) bool {
	return '0' <= c && c <= '9'
}

func isAlpha(c rune) bool {
	return c == '_' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

func isAlphaNumeric(c rune) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}

// isEndOfMLComment report whether s starts with the end of a multi-line comment ("*/").
func isEndOfMLComment(s []rune) bool {
	return len(s) >= 2 && s[0] == '*' && s[1] == '/'
}

func isKeyword(s string) bool {
	return keywords[s]
}

// startsWithLPareen rports whether s starts with arbitary amount of whitespace and a left parenthesis '('.
func startsWithLParen(s []rune) bool {
	i := 0
	for i != len(s) && isWhitespace(s[i]) {
		i++
	}
	return i != len(s) && s[i] == '('
}

// searchEnd returns length of continuing nestedClass symbol and the changed newNestingLevel.
func searchEnd(s []rune, nestingLevel, nestedClass int) (length, newNestingLevel int) {
	length = 0
	newNestingLevel = nestingLevel
	if nestedClass == CString { // `-string
		for length != len(s) && s[length] != '`' {
			length++
		}
		if length != len(s) {
			length++
			newNestingLevel--
		}
	} else /* nestedClass == CComment */ {
		for length != len(s) && !isEndOfMLComment(s[length:]) {
			length++
		}
		if length != len(s) {
			length += 2
			newNestingLevel--
		}
	}
	return
}

// Scan returns class (one of c- constants, see above) and length of symbol in the beginning of s.
func Scan(s []rune, nestingLevel, nestedClass int) (class, length, newNestingLevel int) {
	newNestingLevel = nestingLevel
	c := s[0]
	if nestingLevel != 0 {
		class = nestedClass
		length, newNestingLevel = searchEnd(s, nestingLevel, nestedClass)
	} else if isWhitespace(c) {
		class = CNone
		length = 1
	} else if isAlpha(c) {
		length = 1
		for length != len(s) && isAlphaNumeric(s[length]) {
			length++
		}
		if isKeyword(string(s[:length])) {
			class = CKeyword
		} else if startsWithLParen(s[length:]) {
			class = CProcCall
		} else {
			class = CIdent
		}
	} else if isNumeric(c) {
		class = CNumber
		length = 1
		for length != len(s) && isNumeric(s[length]) {
			length++
		}
	} else if c == '/' { // Start of single-line comment ("//") or multi-line comment ("/*")
		if len(s) > 1 && s[1] == '/' {
			class = CComment
			length = len(s)
		} else if len(s) > 1 && s[1] == '*' {
			class = CComment
			length = 2
			for length != len(s) && !isEndOfMLComment(s[length:]) {
				length++
			}
			if length != len(s) {
				length += 2
			} else {
				newNestingLevel++
			}
		} else {
			class = CNone
			length = 1
		}
	} else if c == '"' || c == '\'' || c == '`' {
		class = CString
		length = 1
		for length != len(s) && s[length] != c {
			length++
		}
		if length != len(s) {
			length++
		} else if c == '`' {
			newNestingLevel++
		}
	} else {
		class = CNone
		length = 1
	}
	return
}
