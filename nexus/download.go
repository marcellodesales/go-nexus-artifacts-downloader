package nexus

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Result string
type Download func() Result

// Downloads an array of files in parallel http://talks.golang.org/2012/concurrency.slide#47
func fanInDownloads(files []string) (results []Result) {
	c := make(chan Result)

	for i := 0; i < len(files); i++ {
		fileName := files[i]
		go func() { c <- nexusDownload(fileName)() }()
	}

	timeout := time.After(1 * time.Minute)
	for i := 0; i < len(files); i++ {
		select {
		case result := <-c:
			fmt.Println("FINISHED: " + result)
			results = append(results, result)

		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}

// Downloads a single url http://talks.golang.org/2012/concurrency.slide#47
func nexusDownload(url string) Download {
	return func() Result {
		tokens := strings.Split(url, "/")
		fileName := tokens[len(tokens)-1]
		log.Println("Downloading", url, "to", fileName)

		// equivalent to Python's `if os.path.exists(filename)`
		if _, err := os.Stat(fileName); err == nil {
			log.Printf("File %s exists...", fileName)
			return Result("File " + fileName + " already exists")
		}

		// Download the file
		start := time.Now()
		response, err := http.Get(url)
		elapsed := time.Since(start)
		if err != nil {
			log.Println("Error while downloading", url, "-", err)
			return Result("Error while downloading: " + url)
		}
		defer response.Body.Close()

		// Create a file to store it
		output, err := os.Create(fileName)
		if err != nil {
			log.Println("Error while creating", fileName, "-", err)
			return Result("Error while creating " + fileName)
		}
		defer output.Close()

		// Transfer the bytes to the file.
		n, err := io.Copy(output, response.Body)
		if err != nil {
			log.Println("Error saving the downloaded file", url, "-", err)
			return Result("Error while saving the file " + fileName)
		}

		log.Println(n, "bytes downloaded from ", url, " in ", elapsed)
		return Result(fmt.Sprintf("%s downloaded in %d", fileName, elapsed))
	}
}

// getArtifactsUrlList Builds the list of URLs to be downloaded based on the
// Lavest version of the artifacts.
func (al *ArtifactsList) getArtifactsUrlList() []string {
	urls := make([]string, 0, len(al.Index))
	for _, artifact := range al.Index {
		// TODO: Filter the filtered list of artifacts already loaded.
		version := artifact.Metadata.Versioning.Latest
		extension := ".zip"
		urls = append(urls, artifact.GetArtifactUrl(version, extension))
	}
	return urls
}

// DownloadAllList downloads all the collected artifacts latest versions in parallel.
func (al *ArtifactsList) DownloadAllList() {
	start := time.Now()

	files := al.getArtifactsUrlList()

	results := fanInDownloads(files)
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
