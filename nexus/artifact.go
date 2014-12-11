package nexus

// Refer to the Effective use at https://golang.org/doc/effective_go.html
type ArtifactResourceUrl struct {
	Resource string
	Metadata string
}

// Artifact represents a given Nexus artifact that has its own metadata:
type Artifact struct {
	list *ArtifactsList
	// The metadata of the artifact
	Metadata *MavenMetadata
}

// makeArtifactIdUrl retuns the URL for a given artifact in Nexus
func makeArtifactIdUrl(al *ArtifactsList, artifactId string) string {
	return al.MakeGroupIdUrl() + "/" + artifactId
}

// getArtifactMetadataUrl builds the URL to retrieve the maven-metadata.xml for an artifact
// https://repo1.maven.org/nexus2/org/springframework/cloud/[...server...]
// artifactId = quadf
func getArtifactMetadataUrl(al *ArtifactsList, artifactId string) string {
	return makeArtifactIdUrl(al, artifactId) + "/maven-metadata.xml"
}

// makeServiceUrl build the artifactId URL based on the artifactId.
func (art *Artifact) makeArtifactIdUrl() string {
	return makeArtifactIdUrl(art.list, art.Metadata.ArtifactId)
}

// GetMetadataUrl builds the URL to retrieve the maven-metadata.xml from a service.
func (art *Artifact) GetArtifactMetadataUrl() string {
	return getArtifactMetadataUrl(art.list, art.Metadata.ArtifactId)
}

// GetArtifactUrl builds the URL to retrieve the binary for the given service version.
func (art *Artifact) GetArtifactUrl(version, extension string) string {
	return art.makeArtifactIdUrl() + "/" + version + "/" + art.Metadata.ArtifactId + "-" + version + extension
}

// NewArtifact creates a new Artifact based on the given serviceName
func NewArtifact(al *ArtifactsList, artifactId string) *Artifact {
	// Build the meadata urls
	urls := ArtifactResourceUrl{
		Resource: makeArtifactIdUrl(al, artifactId),
		Metadata: getArtifactMetadataUrl(al, artifactId),
	}

	// Build a new Maven Metadata, extracting the XML from the urls
	artifactMetadata := NewMavenMetadata(&urls)

	// Build the artifact instance
	art := Artifact{
		list:     al,
		Metadata: artifactMetadata,
	}

	return &art
}
