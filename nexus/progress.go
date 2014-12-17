package nexus

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Inspired by https://github.com/thbar/golang-playground/blob/master/download-files.go
// http://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go

// The Url Metadata collected during the HTTP Head.
type UrlMetadata struct {
	url    string
	name   string
	length int64
	err    string
}

// The UrlDownload carries the UrlMetadata and the instance of the Progress Bar.
type UrlDownload struct {
	metadata    *UrlMetadata
	progressBar *pb.ProgressBar
}

// The download type that's a function that returns the UrlMetadata.
type ResourceSizeDownload func() UrlMetadata

// collectBarsIndex will return the index of the url and the URL metadata.
func collectBarsIndex(urls []string) map[string]*UrlDownload {
	// Channel for the results
	c := make(chan UrlMetadata)

	// The resulting index of the resources fileNames and their progress bars
	barsIndex := make(map[string]*UrlDownload)

	// Retrieve the resources Length in parallel
	for i := 0; i < len(urls); i++ {
		url := urls[i]
		go func() { c <- retrieveResourceLength(url)() }()
	}

	timeout := time.After(1 * time.Minute)

collectResourceSize:
	for i := 0; i < len(urls); i++ {
		select {
		case urlMetadata := <-c:
			if urlMetadata.length == -1 {
				continue
			}

			// create bar
			bar := pb.New64(urlMetadata.length).SetUnits(pb.U_BYTES)
			bar.SetRefreshRate(time.Millisecond * 10).Prefix(urlMetadata.name)
			bar.ShowSpeed = true

			// Index the bar for the filename of the resource
			barsIndex[urlMetadata.url] = &UrlDownload{
				metadata:    &urlMetadata,
				progressBar: bar,
			}

		case <-timeout:
			break collectResourceSize
		}
	}
	return barsIndex
}

// retrieveResourceLength will execute an HTTP HEAD request and collect the Content-Length of the given
// url. It will return the UrlMetadata instance as a result. In addition, it will verify if the file
// exists as well, so that there's no need to download the file.
func retrieveResourceLength(url string) ResourceSizeDownload {
	return func() UrlMetadata {
		// Parse the URL
		tokens := strings.Split(url, "/")
		fileName := tokens[len(tokens)-1]
		log.Println("Processing", url, "to", fileName)

		// equivalent to Python's `if os.path.exists(filename)`
		if _, err := os.Stat(fileName); err == nil {
			return UrlMetadata{
				url:    url,
				name:   fileName,
				err:    "File " + fileName + " already exists",
				length: -1,
			}
		}

		// Make an HTTP HEAD for the given URL
		response, err := http.Head(url)

		// If any problem occurs, just return the error.
		if err != nil {
			log.Println("Error while downloading", url, ":", err)
			return UrlMetadata{
				url:    url,
				name:   fileName,
				err:    fmt.Sprintf("Error while downloading", url, ": ", err),
				length: -1,
			}
		}

		// Verify if the response was ok
		if response.StatusCode != http.StatusOK {
			log.Println("Server return non-200 status: %v\n", response.Status)
			return UrlMetadata{
				url:    url,
				name:   fileName,
				err:    "Server return non-200 status: " + response.Status,
				length: -1,
			}
		}

		// Retrieve the value of the length of the resource by looking at the HTTP Header below
		length, _ := strconv.Atoi(response.Header.Get("Content-Length"))
		sourceSize := int64(length)
		return UrlMetadata{
			url:    url,
			name:   fileName,
			length: sourceSize,
			err:    "",
		}
	}
}
