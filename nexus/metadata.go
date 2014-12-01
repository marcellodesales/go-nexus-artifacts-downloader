// http://tour.golang.org/#4
package nexus

// Refer to the Effective use at https://golang.org/doc/effective_go.html

// http://tour.golang.org/#5
import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	// The Metadata URL where this was downloaded
	Urls *ArtifactResourceUrl
	// The xml Document downloaded from the metadata Url
	xmlDoc string
	// The GroupId of the metadata XML document (PARSED)
	GroupId string `xml:"groupId"`
	// The artifact Id of the metadata XML document (PARSED)
	ArtifactId string `xml:"artifactId"`
	// The versioning of the metadata XML document (PARSED)
	Versioning *Versioning `xml:"versioning"`
}

// parseXml Will bind the values of the xmlDocument to the struct based on the bindings defined.
func (nm *MavenMetadata) parseXml() {
	// For larger XML files, use the events http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/
	// https://github.com/dps/go-xml-parse
	// Pointers http://tour.golang.org/#28, http://tour.golang.org/#29
	if err := xml.Unmarshal([]byte(nm.xmlDoc), nm); err != nil {
		log.Fatalln(err)
	}
}

// loadNexusMetadataXmlDoc Loads the xmlDoc from the given URL given.
func (mm *MavenMetadata) loadNexusMetadataXmlDoc() {
	// short var declaration without "var" http://tour.golang.org/#13
	url := mm.Urls.Metadata
	log.Printf("Received url %s", url)
	// Download a URL http://golang.org/pkg/net/http/
	metadataXmlResponse, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		log.Printf("Error is " + err.Error())
	}

	defer metadataXmlResponse.Body.Close()
	body, err := ioutil.ReadAll(metadataXmlResponse.Body)
	if err != nil {
		fmt.Printf("%s", err)
	}

	mm.xmlDoc = string(body)
}

// NewMavenMetadata creates a new instance based on the give metadataUrl. It will
// fetch the xmlDocument from that URL and parse it to an instance of MavenMetadata.
func NewMavenMetadata(urls *ArtifactResourceUrl) *MavenMetadata {
	metadata := MavenMetadata{
		Urls: urls,
	}

	// Load the XML document
	metadata.loadNexusMetadataXmlDoc()

	// Parse the xml
	metadata.parseXml()
	return &metadata
}
