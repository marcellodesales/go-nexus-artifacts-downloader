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
	"strings"
)

// http://tour.golang.org/#16
const (
	NEXUS_REPO           = "https://repo1.maven.org/maven2/"
	CFP_PACKAGE_METADATA = "org/springframework/cloud/_SERVICE_NAME_/maven-metadata.xml"
	NEXUS_URL            = NEXUS_REPO + CFP_PACKAGE_METADATA
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

type ServiceMetadata struct {
	Name     string         `The name of the service`
	Metadata *MavenMetadata `The metadata of the service`
}

func (sm *ServiceMetadata) GetUrl() string {
	// short var declaration without "var" http://tour.golang.org/#13
	url := strings.Replace(NEXUS_URL, "_SERVICE_NAME_", sm.Name, -1)
	log.Printf("Will return %s", url)
	return url
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
	log.Println("The value is %#v", metadata)
	return &metadata
}

func (sm *ServiceMetadata) load() {
	url := sm.GetUrl()
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

	sm.Metadata = NewMavenMetadata(string(body))
}

func NewServiceMetadata(serviceName string) *ServiceMetadata {
	sm := ServiceMetadata{
		Name: serviceName,
	}
	log.Printf("Going to load Service '" + serviceName)

	sm.load()
	return &sm
}
