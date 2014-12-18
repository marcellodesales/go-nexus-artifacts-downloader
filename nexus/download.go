package nexus

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Inspired by https://github.com/thbar/golang-playground/blob/master/download-files.go
// http://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go

// The result of a download.
type Result string

// The Download function that takes a result.
type Download func() Result

// Downloads an array of files in parallel http://talks.golang.org/2012/concurrency.slide#47
func fanInDownloads(urls []string) (results []Result) {

	// Retrieve the bars index for each URL
	urlsDownloadIndex := collectBarsIndex(urls)

	// Create a channel to wait for the download results
	c := make(chan Result)

	// Concurrently download the URLs
	for i := 0; i < len(urls); i++ {
		url := urls[i]
		urlDownload := urlsDownloadIndex[url]
		go func() { c <- nexusDownload(urlDownload)() }()
	}

	// The final result
	results = make([]Result, 0, len(urls))

	// Give a total download time for all the resources
	timeout := time.After(20 * time.Minute)

	// Collect the results by listening to the result channel
downloadAllFiles:
	for i := 0; i < len(urls); i++ {
		select {
		case result := <-c:
			results = append(results, result)

		case <-timeout:
			break downloadAllFiles
		}
	}
	// TODO: Compute which files were NOT downloaded due to timeout

	// Sleepting 1s to let the last downloaded file to finish printing the progress update.
	time.Sleep(1 * time.Second)

	return
}

// Downloads a single url http://talks.golang.org/2012/concurrency.slide#47
func nexusDownload(urlDownload *UrlDownload) Download {
	return func() Result {

		// Upon errors, the length will NOT be filled
		if urlDownload.metadata.length == -1 {
			return Result(urlDownload.metadata.err)
		}

		url := urlDownload.metadata.url

		// Download the file
		start := time.Now()

		// HTTP GET the file
		response, err := http.Get(url)
		if err != nil {
			return Result(fmt.Sprintf("Error while downloading %s: %v", url, err))
		}
		// Go's built-in defer statement defers execution of the
		// specified function until the current function returns.
		// https://coderwall.com/p/cp5fya/measuring-execution-time-in-go
		defer response.Body.Close()

		// Verify if the response was ok
		if response.StatusCode != http.StatusOK {
			return Result(fmt.Sprintf("Server return non-200 status: %v\n", response.Status))
		}

		fileName := urlDownload.metadata.name

		// Create a file to store it
		output, err := os.Create(fileName)
		if err != nil {
			log.Println("Error while creating", fileName, "-", err)
			return Result(fmt.Sprintf("Error while creating %s: %v", fileName, err))
		}
		defer output.Close()

		// create multi writer
		writer := io.MultiWriter(output, urlDownload.progressBar)

		// Transfer the bytes to the file, blocking here (both are deferred references)
		// When the output and response.Body are ready, this will continue
		totalBytes, err := io.Copy(writer, response.Body)

		// Stop the bar
		urlDownload.progressBar.Finish()

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

// DownloadAllList downloads all the collected artifacts latest versions in parallel.
func (al *ArtifactsList) DownloadAllList() {
	start := time.Now()

	urls := al.getArtifactsUrlList()

	results := fanInDownloads(urls)
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
