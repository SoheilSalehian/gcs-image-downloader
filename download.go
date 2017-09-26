package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"github.com/SoheiLSalehian/gcs-image-downloader/lib/entity"
	"github.com/SoheiLSalehian/gcs-image-downloader/lib/file"
)

func bulkDownload(results []entity.Link) int64 {
	var wg sync.WaitGroup
	var cnt int
	var errCount int64

	for _, v := range results {
		if cnt == num {
			break
		}
		cnt++

		wg.Add(1)
		go func(link string) {
			defer wg.Done()

			err := download(link)
			if err != nil {
				atomic.AddInt64(&errCount, 1)
				log.Println(err)
			}
		}(v.Link)
		time.Sleep(time.Millisecond * 100)
	}
	wg.Wait()

	return errCount
}

func download(link string) error {
	resp, err := http.Get(link)
	if err != nil || resp.StatusCode != 200 {
		return err
	}
	defer resp.Body.Close()

	n, err := file.Name(link)
	if err != nil {
		return err
	}
	f, err := os.Create(path.Join(outputsPath, n))
	if err != nil {
		return err
	}

	_, err = io.Copy(f, resp.Body)
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	return err
}
