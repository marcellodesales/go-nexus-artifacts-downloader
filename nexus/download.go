package nexus

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// The result of a download.
type Result string

// The Download function that takes a result.
type Download func() Result

// Download Progress
type Progress struct {
	log      string
	finished bool
}

// Downloads an array of files in parallel http://talks.golang.org/2012/concurrency.slide#47
func fanInDownloads(files []string, progressChannel chan Progress) (results []Result) {
	c := make(chan Result)

	for i := 0; i < len(files); i++ {
		fileName := files[i]
		go func() { c <- nexusDownload(fileName)() }()
	}

	go monitor(progressChannel)

	timeout := time.After(20 * time.Minute)

downloadAllFiles:
	for i := 0; i < len(files); i++ {
		select {
		case result := <-c:
			progressChannel <- Progress{
				log: "FINISHED[" + strconv.Itoa(i) + "]: " + string(result),
			}
			results = append(results, result)

		case <-timeout:
			progressChannel <- Progress{
				log: "1m TIMEOUT[" + strconv.Itoa(i) + "]: Missing " + strconv.Itoa(len(files)-i),
			}
			break downloadAllFiles
		}
	}

	progressChannel <- Progress{
		log:      "FINISHED ALL",
		finished: true,
	}
	close(progressChannel)

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
			return Result("File " + fileName + " already exists")
		}

		// Download the file
		start := time.Now()

		// HTTP GET the file
		response, err := http.Get(url)
		if err != nil {
			log.Println("Error while downloading", url, "-", err)
			return Result(fmt.Sprintf("Error while downloading", url, ":", err))
		}
		// Go's built-in defer statement defers execution of the
		// specified function until the current function returns.
		// https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
		defer response.Body.Close()

		// Create a file to store it
		output, err := os.Create(fileName)
		if err != nil {
			log.Println("Error while creating", fileName, "-", err)
			return Result("Error while creating " + fileName)
		}
		defer output.Close()

		// Transfer the bytes to the file, blocking here (both are deferred references)
		// When the output and response.Body are ready, this will continue
		totalBytes, err := io.Copy(output, response.Body)

		// https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
		elapsed := time.Since(start)

		if err != nil {
			log.Println("Error saving the downloaded file", url, "-", err)
			return Result("Error while saving the file " + fileName)
		}

		return Result(fmt.Sprintf("%d bytes downloaded and saved as %s in %s", totalBytes, fileName, elapsed))
	}
}

// getArtifactsUrlList Builds the list of URLs to be downloaded based on the
// Lavest version of the artifacts.
func (al *ArtifactsList) getArtifactsUrlList() []string {
	urls := make([]string, 0, len(al.Index))
	for _, artifact := range al.Index {
		// TODO: Filter the filtered list of artifacts already loaded.
		version := artifact.Metadata.Versioning.Latest
		extension := ".jar"
		urls = append(urls, artifact.GetArtifactUrl(version, extension))
	}
	return urls
}

// Monitors the progress of a Progress channel.
func monitor(progress chan Progress) {
	for {
		select {
		case event := <-progress:
			fmt.Println("MONITORING: " + event.log)
			if event.finished {
				return
			}
		}
	}

}

// DownloadAllList downloads all the collected artifacts latest versions in parallel.
func (al *ArtifactsList) DownloadAllList() {
	start := time.Now()

	progressChannel := make(chan Progress)

	files := al.getArtifactsUrlList()

	results := fanInDownloads(files, progressChannel)
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
