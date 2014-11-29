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
	instance := nexus.NewServiceMetadata("spring-cloud-config-server")
	log.Printf("Going to load Service '" + instance.Name + "' from " + instance.GetUrl())

	// Multiple returns http://tour.golang.org/#9
	log.Printf("Updated the service %#v", instance.Metadata.Versioning)
}
