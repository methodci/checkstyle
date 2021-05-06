package checkstyle

import (
	"encoding/xml"
	"io"

	"golang.org/x/net/html/charset"
)

// SeverityLevel as defined by severity levels from https://checkstyle.sourceforge.io/property_types.html#SeverityLevel
type SeverityLevel string

// SeverityLevels as defined by the checkstyle standard
const (
	SeverityIgnore  SeverityLevel = "ignore"
	SeverityInfo    SeverityLevel = "info"
	SeverityWarning SeverityLevel = "warning"
	SeverityError   SeverityLevel = "error"
)

// Checkstyle <checkstyle /> XML struct as defined by the checkstyle standard
type Checkstyle struct {
	XMLName xml.Name `xml:"checkstyle"`
	Version string   `xml:"version,attr"`
	File    Files    `xml:"file"`
}

// File <file /> XML struct as defined by the checkstyle standard
type File struct {
	Name  string      `xml:"name,attr"`
	Error []FileError `xml:"error"`
}

// FileError <error /> XML struct as defined by the checkstyle standard
type FileError struct {
	Line     int           `xml:"line,attr"`
	Severity SeverityLevel `xml:"severity,attr"`
	Message  string        `xml:"message,attr"`
	Source   string        `xml:"source,attr"`
}

// Files is a collection of <file /> XML structs with added helpers
type Files []File

// FromName finds a File matching the given filename. If multiple files match a
// given name their errors will be merged
func (chk Files) FromName(name string) File {
	var out File
	for _, c := range chk {
		if c.Name == name {
			out.Name = name
			out.Error = append(out.Error, c.Error...)
		}
	}

	return out
}

// Decode reads a checkstyle file from a reader
func Decode(r io.Reader) (*Checkstyle, error) {
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReaderLabel

	chk := Checkstyle{}
	err := dec.Decode(&chk)
	if err != nil {
		return nil, err
	}

	return &chk, nil
}
