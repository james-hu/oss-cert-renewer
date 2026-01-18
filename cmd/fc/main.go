package main

import (
	"fmt"
	"log"
	"os"
	"osscert"
	"strings"
	"sync"

	"github.com/aliyun/fc-runtime-go-sdk/fc"
)

func main() {
	fc.Start(handler)
}

func handler() (string, error) {
	// Parse buckets from environment variables
	bucketSet := make(map[string]struct{})

	// Multiple buckets support
	if bs := os.Getenv("OSS_BUCKETS"); bs != "" {
		parts := strings.Split(bs, ",")
		for _, p := range parts {
			if trimmed := strings.TrimSpace(p); trimmed != "" {
				bucketSet[trimmed] = struct{}{}
			}
		}
	}

	if len(bucketSet) == 0 {
		return "No buckets configured for renewal.", nil
	}

	var (
		sb    strings.Builder
		mutex sync.Mutex
		wg    sync.WaitGroup
	)

	count := 0
	total := len(bucketSet)

	for bucket := range bucketSet {
		wg.Add(1)
		count++
		
		go func(b string, c int) {
			defer wg.Done()
			
			log.Printf("[%d/%d] Starting processing for bucket: %s", c, total, b)

			// Run the renewal logic
			rst, err := osscert.Run(b)
			
			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				// Log error but continue with other buckets
				log.Printf("[%s] Error: %v", b, err)
				sb.WriteString(fmt.Sprintf("[%s] Error: %v\n", b, err))
			} else {
				log.Printf("[%s] Success: %s", b, strings.TrimSpace(rst))
				sb.WriteString(fmt.Sprintf("[%s] %s\n", b, rst))
			}
		}(bucket, count)
	}

	wg.Wait()
	return sb.String(), nil
}
