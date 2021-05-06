package checkstyle

import (
	"regexp"
)

// DiffOptions are the options for the Diff command
type DiffOptions struct {
	MaxLineDiff int
}

// Diff runs a diff operation on two checkstyle structs and returns two
// new checkstyle structs of fixed and created errors
func Diff(left, right *Checkstyle, opt DiffOptions) (*Checkstyle, *Checkstyle) {
	lout := &Checkstyle{}
	rout := &Checkstyle{}

	names := allNames(left, right)
	for n, s := range names {

		if s == sideLeft {
			lout.File = append(lout.File, left.File.FromName(n))
			continue
		}

		if s == sideRight {
			rout.File = append(rout.File, right.File.FromName(n))
			continue
		}

		lf := left.File.FromName(n)
		rf := right.File.FromName(n)

		lfe := lf.Error
		rfe := rf.Error

		// lfe := lf.Error
		// rfe := rf.Error

		// running through zero first prevernts accidental lose equality where
		// actual equality exists several lines below
		lfe, rfe = fileErrorWoExactEq(lfe, rfe, 0)
		if opt.MaxLineDiff > 0 {
			lfe, rfe = fileErrorWoExactEq(lfe, rfe, opt.MaxLineDiff)
		}

		if len(lfe) > 0 {
			lout.File = append(lout.File, File{
				Name:  n,
				Error: lfe,
			})
		}

		if len(rfe) > 0 {
			rout.File = append(rout.File, File{
				Name:  n,
				Error: rfe,
			})
		}
	}

	return lout, rout
}

type side int

const (
	sideLeft side = iota + 1
	sideRight
	sideBoth
)

func fileErrorWoExactEq(left, right []FileError, maxLineDiff int) ([]FileError, []FileError) {
	lout := []FileError{}
	rout := []FileError{}

leftloop:
	for _, l := range left {
		for _, r := range right {
			if fileErrorEq(l, r, maxLineDiff) {
				continue leftloop
			}
		}

		lout = append(lout, l)
	}

rightloop:
	for _, r := range right {
		for _, l := range left {
			if fileErrorEq(l, r, maxLineDiff) {
				continue rightloop
			}
		}

		rout = append(rout, r)
	}

	return lout, rout
}

func fileErrorEq(lfile, rfile FileError, maxLineDiff int) bool {
	if lfile.Severity != rfile.Severity || lfile.Source != rfile.Source || !msgEq(lfile.Message, rfile.Message) {
		return false
	}

	if abs(lfile.Line-rfile.Line) <= maxLineDiff {
		return true
	}

	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

var msgClean = regexp.MustCompile(`[^\p{L}?:().]+`)

// This is an initial attempt at this and should be reworked later to consider
// things like the line number difference in parts of this that perhaps include
// line numbers
func msgEq(left, right string) bool {
	left = msgClean.ReplaceAllString(left, "|")
	right = msgClean.ReplaceAllString(right, "|")

	return left == right
}

func allNames(left, right *Checkstyle) map[string]side {
	ret := map[string]side{}

	for _, c := range left.File {
		ret[c.Name] = sideLeft
	}

	for _, c := range right.File {
		if _, ok := ret[c.Name]; ok {
			ret[c.Name] = sideBoth
		} else {
			ret[c.Name] = sideRight
		}
	}

	return ret
}
