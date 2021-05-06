package checkstyle

import (
	"encoding/xml"
	"io"

	"golang.org/x/net/html/charset"
)

// Defined severity levels from http://checkstyle.sourceforge.net/property_types.html#severity
type SeverityLevel string

const (
	SeverityIgnore  SeverityLevel = "ignore"
	SeverityInfo    SeverityLevel = "info"
	SeverityWarning SeverityLevel = "warning"
	SeverityError   SeverityLevel = "error"
)

type Checkstyle struct {
	XMLName xml.Name        `xml:"checkstyle"`
	Version string          `xml:"version,attr"`
	File    CheckstyleFiles `xml:"file"`
}

type CheckstyleFile struct {
	Name  string                `xml:"name,attr"`
	Error []CheckstyleFileError `xml:"error"`
}

type CheckstyleFileError struct {
	Line     int           `xml:"line,attr"`
	Severity SeverityLevel `xml:"severity,attr"`
	Message  string        `xml:"message,attr"`
	Source   string        `xml:"source,attr"`
}

type CheckstyleFiles []CheckstyleFile

func (chk CheckstyleFiles) FromName(name string) CheckstyleFile {
	var out CheckstyleFile
	for _, c := range chk {
		if c.Name == name {
			out.Error = append(out.Error, c.Error...)
		}
	}

	return out
}

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
