// http://tour.golang.org/#4
package main

// http://tour.golang.org/#5
import (
	//	"./cfp"
	"./nexus"
	"log"
	"strings"
)

const (
	NEXUS_SERVER   = "https://repo1.maven.org"
	NEXUS_REPO     = "maven2"
	NEXUS_GROUP_ID = "org.springframework.cloud"
)

func main() {
	// Struct literals http://tour.golang.org/#29, https://golang.org/doc/effective_go.html#composite_literals
	// Less boilerplate code that uses the new (type) http://tour.golang.org/#30
	log.Printf("Loading services from Nexus at %s", NEXUS_SERVER)

	// Closures http://www.golang-book.com/7/index.htm#section4,
	// http://tour.golang.org/#45, http://jordanorelli.com/post/42369331748/function-types-in-go-golang
	filter := func(artifactId string) bool {
		return !strings.Contains(artifactId, "server")
	}
	artifactsList := nexus.NewArtifactsList(NEXUS_SERVER, NEXUS_REPO, NEXUS_GROUP_ID, filter)

	log.Printf("The maven list of artifacts is %#v \n", artifactsList.Index)

	/*	service := cfp.NewServiceMetadata("spring-cloud-config-server")
		log.Printf("Going to load Service '" + service.Name + "' from " + service.GetMetadataUrl())

		version := service.Metadata.Versioning.Latest
		log.Printf("Downloading selected version %s", version)

		// Multiple returns http://tour.golang.org/#9
		log.Printf("Going to load latest Version '" + service.Name + "' from " + service.GetFileUrl(version))

		// Downloads the file to the current directory
		service.Download(version)
	*/
}
