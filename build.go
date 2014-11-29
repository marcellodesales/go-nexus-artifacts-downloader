// http://tour.golang.org/#4
package main

// http://tour.golang.org/#5
import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	//	"runtime"
	//	"time"
)

// http://tour.golang.org/#16
const (
	NEXUS_REPO           = "http://pprdnexusas301.corp.intuit.net/nexus/content/groups/ENG.CTG-Releases/"
	CFP_PACKAGE_METADATA = "com/intuit/cfp/sp/service/_SERVICE_NAME_/maven-metadata.xml"
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
	Name     string        `The name of the service`
	Metadata MavenMetadata `The metadata of the service`
}

func (sm *ServiceMetadata) getUrl() string {
	// short var declaration without "var" http://tour.golang.org/#13
	url := strings.Replace(NEXUS_URL, "_SERVICE_NAME_", sm.Name, -1)
	log.Printf("Will return %s", url)
	return url
}

func (nm *MavenMetadata) parseXml() {
	// Pointers http://tour.golang.org/#28, http://tour.golang.org/#29
	if err := xml.Unmarshal([]byte(nm.xmlDoc), nm); err != nil {
		log.Fatalln(err)
	}
}

func (sm *ServiceMetadata) load() (MavenMetadata, string) {
	url := sm.getUrl()
	log.Printf("Received url %s", url)
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

	metadata := MavenMetadata{
		// Type conversion http://tour.golang.org/#15, Types http://tour.golang.org/#14
		xmlDoc: string(body),
	}
	metadata.parseXml()

	// Multiple returns http://tour.golang.org/#9
	return metadata, metadata.xmlDoc
}

func main() {
	// Struct literals http://tour.golang.org/#29, https://golang.org/doc/effective_go.html#composite_literals
	// Less boilerplate code that uses the new (type) http://tour.golang.org/#30
	instance := ServiceMetadata{
		Name: "quadf",
	}
	log.Printf("Going to load Service '" + instance.Name + "' from " + instance.getUrl())

	// Multiple returns http://tour.golang.org/#9
	updatedInstance, xml := instance.load()
	log.Printf("Updated the service %#v", updatedInstance.Versioning)
	log.Printf("Downloaded xml %s", xml)
}
