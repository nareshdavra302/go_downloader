package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "time"
	"crypto/rsa"
	"net/url"
	"strings"    
)

type ProgressReader struct {
    Reader io.Reader
    Size   int64
    Pos    int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
    n, err := pr.Reader.Read(p)
    if err == nil {
        pr.Pos += int64(n)
        fmt.Printf("\rDownloading... %.2f%%", float64(pr.Pos)/float64(pr.Size)*100)
    }
    return n, err
}

func main() {
    start := time.Now().UnixMilli()
    tempPath := ".tmp"
    outPath := "200MB.zip"
    req, _ := http.NewRequest("GET", "https://speed.hetzner.de/100MB.bin", nil)
    resp, _ := http.DefaultClient.Do(req)
    if resp.StatusCode != 200 {
        log.Fatalf("Error while downloading: %v", resp.StatusCode)
    }
    defer resp.Body.Close()

    f, _ := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY, 0644)
    defer f.Close()

    progressReader := &ProgressReader{
        Reader: resp.Body,
        Size:   resp.ContentLength,
    }

    if _, err := io.Copy(f, progressReader); err != nil {
        log.Fatalf("Error while downloading: %v", err)
    }

    os.Rename(tempPath, outPath)
    fmt.Println(" - Download completed!")

    fmt.Printf("Took: %.2fs\n", float64(time.Now().UnixMilli()-start)/1000)
}

