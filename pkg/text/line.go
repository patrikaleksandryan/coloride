package text

import (
	"github.com/patrikaleksandryan/coloride/pkg/colorcode"
	"github.com/patrikaleksandryan/coloride/pkg/scanner"
)

// Run holds color attributes of a run of characters in a line.
type Run struct {
	length int
	color  int
	next   *Run
}

type Line struct {
	chars       []rune // Does not include Line.sapces, color marker or color code
	spaces      []rune // All successive whitespace characters right before the color marker ("///")
	colorCode   []rune
	NewLineType int // One of New Line Type constants in scanner.go
	runs        *Run
	prev, next  *Line
}

// Line

func NewLine() *Line {
	return &Line{
		runs:        &Run{length: 1},
		NewLineType: scanner.LF,
	}
}

func (l *Line) Chars() []rune {
	return l.chars
}

func (l *Line) Prev() *Line {
	return l.prev
}

func (l *Line) Next() *Line {
	return l.next
}

func (l *Line) DeleteChar(pos int) {
	l.chars = append(l.chars[:pos], l.chars[pos+1:]...)
	r, _ := l.FindRun(pos)
	r.length--
	l.NormalizeRuns()
}

func (l *Line) InsertChar(pos int, ch rune) {
	l.chars = append(l.chars, 0)
	copy(l.chars[pos+1:], l.chars[pos:])
	l.chars[pos] = ch

	r, _ := l.FindRun(pos - 1)
	r.length++
}

// FindRun returns the run, to which character at the given position belongs to,
// and the position of that character within the run. If given -1, returns the first run (if any).
func (l *Line) FindRun(pos int) (*Run, int) {
	r := l.runs
	if r != nil {
		for r != nil && pos >= r.length {
			pos -= r.length
			r = r.next
		}
	}
	return r, pos
}

// RemoveEmptyRuns insures there are no runs with length = 0.
func (l *Line) RemoveEmptyRuns() {
	r := l.runs
	var rOld *Run
	for r != nil {
		if r.length == 0 {
			if rOld == nil {
				l.runs = r.next
			} else {
				rOld.next = r.next
			}
		} else {
			rOld = r
		}
		r = r.next
	}
}

// MergeSameRuns insures there are no two successive runs with the same color.
func (l *Line) MergeSameRuns() {
	r := l.runs
	if r != nil { // Line is not empty
		for r.next != nil {
			if r.IsSameColor(r.next) {
				r.length += r.next.length
				r.next = r.next.next
			} else {
				r = r.next
			}
		}
	}
}

// NormalizeRuns restored the invariant of runs.
// 1. There are no runs with length = 0.
// 2. There are no two successive runs with the same color.
// 3. Any line has at least one run (at least the new line character).
// 4. There are no characters in the line that do not belong to a run.
// 5. There are no characters in the line that belong to more than one run.
// 6. The sum of lengths of all runs in a line equals len(chars) + 1.
func (l *Line) NormalizeRuns() {
	l.RemoveEmptyRuns()
	l.MergeSameRuns()
}

// Split splits the line in two parts at the given position, placing the new line after the current line.
func (l *Line) Split(pos int) {
	newLine := NewLine()
	newLine.next = l.next
	newLine.prev = l
	if l.next != nil {
		l.next.prev = newLine
	}
	l.next = newLine

	newLine.chars = append(newLine.chars, l.chars[pos:]...)
	l.chars = l.chars[:pos]

	l.SplitRuns(pos)
}

// SplitRuns splits the list of runs of the given line at the given position, placing the second part of
// runs to the next line. It also increases the length of the last run of the first line by one.
func (l *Line) SplitRuns(pos int) {
	if pos == 0 {
		color := 0
		if l.prev != nil {
			prevColor := l.prev.LastRun().color
			if prevColor == l.runs.color {
				color = prevColor
			}
		}
		l.next.runs = l.runs
		l.runs = &Run{length: 1, color: color}
	} else {
		// (pos - 1) because we want to find the previous run if the character is on the 0-th index of the run
		r, pos := l.FindRun(pos - 1)
		pos++
		if r.length != pos {
			r.Split(pos)
		}
		l.next.runs = r.next
		r.next = nil
		r.length++
	}
}

// LastRun returns the last (right-most) run of the line.
func (l *Line) LastRun() *Run {
	r := l.runs
	for r.next != nil {
		r = r.next
	}
	return r
}

// CutRun cuts the run at the given position into two parts. If the position is between two runs, no cut is made.
func (l *Line) CutRun(pos int) {
	run, pos := l.FindRun(pos)
	if pos != 0 {
		run.Split(pos)
	}
}

// IsColorized reports whether the line has at least one non-standard run.
func (l *Line) IsColorized() bool {
	// If there is more than one run, then at least one of them must be non-standard.
	// If there is only one run, then it can be standard (color = 0) or non-standard.
	return l.runs.next != nil || l.runs.color != 0
}

// Colorize sets the color of the given range of characters [from; to).
func (l *Line) Colorize(color, from, to int) {
	if 0 <= from && from < to && to <= len(l.chars)+1 {
		l.CutRun(from)
		l.CutRun(to)
		run, _ := l.FindRun(from)
		length := to - from
		for length != 0 {
			run.color = color
			length -= run.length
			run = run.next
		}
		l.NormalizeRuns()
	}
}

// StringRange returns characters [from; to) of the line as string.
func (l *Line) StringRange(from, to int) string {
	return string(l.chars[from:to])
}

// DeleteRange deletes characters in the given range [from; to).
func (l *Line) DeleteRange(from, to int) {
	if 0 <= from && from < to && to <= len(l.chars)+1 {
		l.CutRun(from)
		l.CutRun(to)
		fromRun, _ := l.FindRun(from - 1) // Find previous run, to handle run.next
		var next *Run
		if to != len(l.chars)+1 {
			next, _ = l.FindRun(to)
		}

		if from == 0 {
			l.runs = next
		} else {
			fromRun.next = next
			if next == nil {
				fromRun.length++
			}
		}
		l.NormalizeRuns()

		if to == len(l.chars)+1 {
			to--
		}
		l.chars = append(l.chars[:from], l.chars[to:]...)
	}
}

// ApplyColorCode parses l.colorCode and applies the instructions as Colorize commands.
func (l *Line) ApplyColorCode() {
	s := colorcode.NewScanner(l.colorCode)
	column := 0
	s.Scan()
	for s.Sym != colorcode.EOC {
		switch s.Sym {
		case colorcode.Number:
			column += s.Number
		case colorcode.NumberedLetter:
			l.Colorize(s.Color(), column, column+s.Number)
			column += s.Number
		case colorcode.Letter:
			l.Colorize(s.Color(), column, len(l.chars)+1)
		}
		s.Scan()
	}
}

// Run

func (r *Run) Next() *Run {
	return r.next
}

func (r *Run) IsSameColor(other *Run) bool {
	return r.color == other.color
}

// Split splits the given run in two parts. Does not check if pos is in the correct range.
func (r *Run) Split(pos int) {
	newR := &Run{
		length: r.length - pos,
		color:  r.color,
		next:   r.next,
	}
	r.length = pos
	r.next = newR
}
