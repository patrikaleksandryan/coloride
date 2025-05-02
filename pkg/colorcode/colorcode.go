package colorcode

const (
	// Symbols

	Number         = iota // i.e. "12"
	Letter                // i.e. "R"
	NumberedLetter        // i.e. "12R"
	EOC                   // End of code
)

type Scanner struct {
	Sym    int // One of symbol constants
	Number int
	Letter rune

	code []rune
	ch   rune // invariant: ch = code[0], or 0
}

func NewScanner(code []rune) *Scanner {
	s := &Scanner{
		code: code,
	}
	if len(s.code) != 0 {
		s.ch = code[0]
	} else {
		s.ch = 0
	}
	return s
}

func isLetter(ch rune) bool {
	return 'A' <= ch && ch <= 'Z' || 'a' <= ch && ch <= 'z'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Scanner) Scan() {
	s.skipWhitespace()
	if len(s.code) == 0 {
		s.Sym = EOC
	} else {
		if isDigit(s.ch) {
			s.Number = 0
			for isDigit(s.ch) {
				s.Number = s.Number*10 + int(s.ch-'0')
				s.read()
			}
			//s.Number may be negative because of possible integer overflow
			if s.Number < 0 {
				s.Number = 0
			}
			// Check for case of a numbered letter
			if isLetter(s.ch) {
				s.Letter = s.ch
				s.read()
				s.Sym = NumberedLetter
			} else {
				s.Sym = Number
			}
		} else if isLetter(s.ch) {
			s.Sym = Letter
			s.Letter = s.ch
			s.read()
		} else { // Undefined character
			s.Sym = EOC
		}
	}
}

// read removes first character from s.code.
func (s *Scanner) read() {
	if len(s.code) != 0 {
		s.code = s.code[1:]
	}
	// Restore invariant of ch
	if len(s.code) != 0 {
		s.ch = s.code[0]
	} else {
		s.ch = 0
	}
}

// skipWhitespace removes leading whitespaces from s.code.
func (s *Scanner) skipWhitespace() {
	for len(s.code) != 0 && s.ch <= ' ' {
		s.read()
	}
}

// Color returns s.Letter as a color number.
func (s *Scanner) Color() int {
	switch s.Letter {
	case 'r':
		return 1
	case 'g':
		return 2
	case 'b':
		return 3
	case 'y':
		return 4
	case 'R':
		return 5
	case 'G':
		return 6
	case 'B':
		return 7
	case 'Y':
		return 8
	default:
		return 0
	}
}

// ToLetter returns the letter of the given color number.
// color must not be 0.
func ToLetter(color int) rune {
	switch color {
	case 1:
		return 'r'
	case 2:
		return 'g'
	case 3:
		return 'b'
	case 4:
		return 'y'
	case 5:
		return 'R'
	case 6:
		return 'G'
	case 7:
		return 'B'
	case 8:
		return 'Y'
	default:
		panic("impossible: color number")
	}
}
