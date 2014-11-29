// http://tour.golang.org/#4
package main

// http://tour.golang.org/#5
import (
	"./nexus"
	"log"
)

func main() {
	// Struct literals http://tour.golang.org/#29, https://golang.org/doc/effective_go.html#composite_literals
	// Less boilerplate code that uses the new (type) http://tour.golang.org/#30
	service := nexus.NewServiceMetadata("spring-cloud-config-server")
	log.Printf("Going to load Service '" + service.Name + "' from " + service.GetMetadataUrl())
	version := service.Metadata.Versioning.Latest
	log.Printf("Downloading selected version %s", version)

	// Multiple returns http://tour.golang.org/#9
	log.Printf("Going to load latest Version '" + service.Name + "' from " + service.GetFileUrl(version))

	// Downloads the file to the current directory
	service.Download(version)
}
