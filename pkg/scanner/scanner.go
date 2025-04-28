package scanner

import (
	"bufio"
	"fmt"
	"io"
)

const (
	// Symbols

	String      = iota
	ColorMarker // Three slashes ("///")
	NewLine     // "\r", "\n" or "\r\n"
	EOT         // End of text
)

type Scanner struct {
	Ch                  rune
	file                *bufio.Reader
	Error               error
	eof                 bool
	colorMarkerDetected bool

	Sym    int
	String []rune // Actual data of the last scanned symbol if sym = String
}

func NewScanner(file *bufio.Reader) *Scanner {
	s := &Scanner{
		file: file,
	}
	s.read()
	return s
}

// Scan scans the opened file for the next symbol and returns it.
// It saves the result it Scanner.Sym and Scanner.String.
func (s *Scanner) Scan() {
	if s.eof {
		s.Sym = EOT
	} else if s.colorMarkerDetected {
		s.Sym = ColorMarker
		s.colorMarkerDetected = false
	} else if s.Ch == '\n' {
		s.read()
		s.Sym = NewLine
	} else if s.Ch == '\r' {
		s.read()
		if s.Ch == '\n' {
			s.read()
		}
		s.Sym = NewLine
	} else { // String or ColorMarker
		s.String = make([]rune, 0, 20)
		slashes := 0
		for s.Ch != 0 && s.Ch != '\r' && s.Ch != '\n' && slashes != 3 {
			s.String = append(s.String, s.Ch)
			if s.Ch == '/' {
				slashes++
			} else {
				slashes = 0
			}
			s.read()
		}
		if slashes == 3 {
			s.colorMarkerDetected = true
			s.String = s.String[:len(s.String)-3]
		}
		s.Sym = String
	}
}

func (s *Scanner) read() {
	var err error
	s.Ch, _, err = s.file.ReadRune()
	if err != nil {
		if err == io.EOF {
			s.Ch = 0
			s.eof = true
		}
		s.Error = fmt.Errorf("read rune: %w", err)
	}
}
