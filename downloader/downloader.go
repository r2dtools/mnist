package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const httpClientTimeout = 60 * time.Second

func DownloadFiles(directory string, fileUrlMap map[string]string) error {
	if directory == "" {
		return errors.New("download directory is not specified")
	}

	httpClient := &http.Client{
		Timeout: httpClientTimeout,
	}
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)

	wg := sync.WaitGroup{}
	wg.Add(len(fileUrlMap))

	for fileName, fileUrl := range fileUrlMap {
		go func(name, url string) {
			defer wg.Done()
			err := downloadFile(httpClient, ctx, name, url, directory)
			if err != nil {
				cancel(err)
				return
			}
		}(fileName, fileUrl)
	}

	wg.Wait()

	return nil
}

func downloadFile(
	httpClient *http.Client,
	ctx context.Context,
	name string,
	url string,
	directory string,
) error {
	filePath := filepath.Join(directory, name)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code: %d", response.StatusCode)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
