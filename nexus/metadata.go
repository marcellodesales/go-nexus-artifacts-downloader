// http://tour.golang.org/#4
package nexus

// Refer to the Effective use at https://golang.org/doc/effective_go.html

// http://tour.golang.org/#5
import (
	"encoding/xml"
	"log"
)

const (
	NEXUS_REPO = "https://repo1.maven.org/maven2"
)

// Parse the array http://play.golang.org/p/7lQnQOCh0I
type Versioning struct {
	Latest      string   `xml:"latest"`
	Release     string   `xml:"release"`
	LastUpdated string   `xml:"lastUpdated"`
	Versions    []string `xml:"versions>version"`
}

// http://tour.golang.org/#26
type MavenMetadata struct {
	GroupId    string      `xml:"groupId"`
	ArtifactId string      `xml:"artifactId"`
	Versioning *Versioning `xml:"versioning"`
	xmlDoc     string      `The downloaded xml`
}

func (nm *MavenMetadata) parseXml() {
	// For larger XML files, use the events http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
	// https://github.com/dps/go-xml-parse
	// Pointers http://tour.golang.org/#28, http://tour.golang.org/#29
	if err := xml.Unmarshal([]byte(nm.xmlDoc), nm); err != nil {
		log.Fatalln(err)
	}
}

func NewMavenMetadata(xmlDocument string) *MavenMetadata {
	metadata := MavenMetadata{
		// Type conversion http://tour.golang.org/#15, Types http://tour.golang.org/#14
		xmlDoc: xmlDocument,
	}
	// Parse the xml
	metadata.parseXml()
	return &metadata
}
