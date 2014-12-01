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
	log.Printf("Loading services from Nexus at %s/%s", NEXUS_SERVER, NEXUS_REPO)

	// Closures http://www.golang-book.com/7/index.htm#section4,
	// http://tour.golang.org/#45, http://jordanorelli.com/post/42369331748/function-types-in-go-golang
	serverFilter := func(artifactId string) bool {
		// runnable servers from spring https://repo1.maven.org/maven2/org/springframework/cloud/
		return strings.Contains(artifactId, "server")
	}

	// TODO: If the filter does not resolve to anything, it is generating the error
	//2014/11/29 13:56:21 Loading services from Nexus at https://repo1.maven.org/maven2
	//2014/11/29 13:56:21 Loading the package list https://repo1.maven.org/maven2/org/springframework/cloud
	//2014/11/29 13:56:21 Received url https://repo1.maven.org/maven2/org/springframework/cloud/../maven-metadata.xml
	//2014/11/29 13:56:21 XML syntax error on line 6: element <hr> closed by </body>

	artifactsList := nexus.NewArtifactsList(NEXUS_SERVER, NEXUS_REPO, NEXUS_GROUP_ID, serverFilter)

	log.Printf("The maven list of artifacts is %#v \n", artifactsList.Index)
	log.Printf("DOWNLOADING NOW... \n\n")

	artifactsList.DownloadAllList()
}
