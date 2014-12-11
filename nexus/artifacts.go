package nexus

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

// Closures http://www.golang-book.com/7/index.htm#section4, http://tour.golang.org/#45,
// Closure types http://jordanorelli.com/post/42369331748/function-types-in-go-golang
type ListFilter func(string) bool

// ArtifactsList represents the nexus server properties as the following:
// nexus2
// https://repo1.maven.org/maven2/org/springframework/cloud/(all-artifacts)
type ArtifactsList struct {
	// The hostname with the protocol, up to the nexus: "http://www.yourcompany.net/nexus2"
	serverHost string
	// if it is a repo, proxy, or group, include that too: "/groups/Releases"
	repositoryPath string
	// The groupId of the artifact "org.springframework.cloud"
	groupId string
	// The map of the resources under the package that will be loaded "quadf, taxreturn, etc"
	Index map[string]*Artifact
	// Filter closure for the artifact list.
	filter ListFilter
}

// makeRepoUrl builds a new repo URL based on the given repoName
// [serverHost http://https://repo1.maven.org][repoPath /nexus2]
func (ml *ArtifactsList) makeRepoUrl() string {
	return ml.serverHost + "/" + ml.repositoryPath
}

// makePackageUrl builds the groupId URL based on the given groupId. Note that it must use the "." notation or "/"
// org.springframework.cloud => REPO_URL/org/springframework/cloud
func (ml *ArtifactsList) makeGroupIdUrl() string {
	return ml.makeRepoUrl() + "/" + strings.Replace(ml.groupId, ".", "/", -1)
}

// makeServiceUrl build the artifactId URL based on the artifactId.
// groupID_URL/[artifactId spring-config-server]
func (ml *ArtifactsList) makeArtifactIdUrl(artifactId string) string {
	return ml.makeGroupIdUrl() + "/" + artifactId
}

// GetMetadataUrl builds the URL to retrieve the maven-metadata.xml from a service.
// artifactId/[artifactId spring-config-server/maven-metadata.xml]
func (ml *ArtifactsList) GetArtifactMetadataUrl(artifactId string) string {
	return ml.makeArtifactIdUrl(artifactId) + "/maven-metadata.xml"
}

// GetFileUrl builds the URL to retrieve the binary for the given server version.
func (ml *ArtifactsList) GetArtifactZipUrl(artifactId, version string) string {
	return ml.makeArtifactIdUrl(artifactId) + "/" + version + "/" + artifactId + "-" + version + ".jar"
}

// GetFileUrl builds the URL to retrieve the binary for the given service version.
func (list *ArtifactsList) GetLatestArtifactZipUrl(artifactId string) string {
	mavenMetadata := list.Index[artifactId]
	latestVersion := mavenMetadata.Metadata.Versioning.Latest
	return list.GetArtifactZipUrl(artifactId, latestVersion)
}

// makePackageUrl builds the groupId URL based on the given groupId. Note that it must use the "." notation or "/"
// org.springframework.cloud => REPO_URL/org/springframework/cloud
func (al *ArtifactsList) MakeGroupIdUrl() string {
	return al.makeRepoUrl() + "/" + strings.Replace(al.groupId, ".", "/", -1)
}

// load Will retrieve all the services available in the packages Url using screen-scrapping.
func (al *ArtifactsList) fetch() {
	groupIdUrl := al.MakeGroupIdUrl()
	log.Println("Loading the package list " + groupIdUrl)

	// Screen-scrape the nexus packages list https://github.com/PuerkitoBio/goquery#examples
	htmlDoc, err := goquery.NewDocument(groupIdUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Scrapping the links from the package URL
	htmlDoc.Find("a").Each(func(i int, s *goquery.Selection) {
		// Retrieve the text from the link
		artifactId := s.Text()

		// Skip the directory parent directory and the user's filter
		if strings.Contains(artifactId, " ") || !al.filter(artifactId) {
			return
		}

		// Remove the trailing "/"
		artifactId = artifactId[:len(artifactId)-1]

		// Index them for the artifactId
		al.Index[artifactId] = NewArtifact(al, artifactId)
	})
}

// NewArtifactsList is a factory method that creates a new instalce of ArtifactsList based on the
// given host (full path to nexus), repoPath (group|proxy/repoName) and a package (java package path).
// The filter is a closure to select which elements from the list to return.
func NewArtifactsList(nexusHost, repoPath, groupId string, filter ListFilter) *ArtifactsList {
	ml := ArtifactsList{
		serverHost:     nexusHost,
		repositoryPath: repoPath,
		groupId:        groupId,
		// arrays and slices http://tour.golang.org/#32
		// Maps http://tour.golang.org/#39, literals http://tour.golang.org/#40, http://tour.golang.org/#42,
		// http://tour.golang.org/#43
		Index:  make(map[string]*Artifact),
		filter: filter,
	}
	ml.fetch()
	return &ml
}
