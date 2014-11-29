// http://tour.golang.org/#4
package nexus

// Refer to the Effective use at https://golang.org/doc/effective_go.html

// http://tour.golang.org/#5
import (
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
	log.Println("Downloading", url, "to", fileName)

	// equivalent to Python's `if os.path.exists(filename)`
	if _, err := os.Stat(fileName); err == nil {
		log.Printf("File %s exists...", fileName)
		return
	}
	// Download the file
	response, err := http.Get(url)
	if err != nil {
		log.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	// Create a file to store it
	output, err := os.Create(fileName)
	if err != nil {
		log.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	// Transfer the bytes to the file.
	n, err := io.Copy(output, response.Body)
	if err != nil {
		log.Println("Error saving the downloaded file", url, "-", err)
		return
	}

	log.Println(n, "bytes downloaded.")
}
