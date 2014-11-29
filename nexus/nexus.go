// http://tour.golang.org/#4
package nexus

// Refer to the Effective use at https://golang.org/doc/effective_go.html

// http://tour.golang.org/#5
import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// http://tour.golang.org/#16
const (
	NEXUS_REPO                  = "https://repo1.maven.org/maven2"
	NEXUS_PACKAGE               = "/org/springframework/cloud/_SERVICE_NAME_"
	NEXUS_ARTIFACT_METADATA_URL = NEXUS_REPO + NEXUS_PACKAGE + "/maven-metadata.xml"
	NEXUS_ARTIFACT_FILE_URL     = NEXUS_REPO + NEXUS_PACKAGE + "/_VERSION_/_SERVICE_NAME_-_VERSION_-exec.jar"
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

func (sm *ServiceMetadata) GetMetadataUrl() string {
	return strings.Replace(NEXUS_ARTIFACT_METADATA_URL, "_SERVICE_NAME_", sm.Name, -1)
}

func (sm *ServiceMetadata) GetFileUrl(version string) string {
	serviceUrl := strings.Replace(NEXUS_ARTIFACT_FILE_URL, "_SERVICE_NAME_", sm.Name, -1)
	return strings.Replace(serviceUrl, "_VERSION_", version, -1)
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
	// short var declaration without "var" http://tour.golang.org/#13
	url := sm.GetMetadataUrl()
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

func (sm *ServiceMetadata) Download(version string) {
	url := sm.GetFileUrl(version)
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	// TODO: check if the file to be downloaded ETag is different
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}
