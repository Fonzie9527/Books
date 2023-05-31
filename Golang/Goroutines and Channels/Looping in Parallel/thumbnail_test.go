package thumbnail_test

import (
	"code-example/bookcode/gopl.io/ch8/thumbnail"
	"log"
	"os"
	"sync"
)

func makeThumbnails3(filenames []string) {
	ch := make(chan struct{})

	for _, f := range filenames {
		go func(f string) {
			thumbnail.ImageFile(f) // NOTE: ignoring errors
			ch <- struct{}{}
		}(f)
	}

	// Wait for goroutines to complete.
	// 事先已知 goroutine 的数量，不实用
	for range filenames {
		<-ch
	}
}

func makeThumbnails4(filenames []string) error {
	errors := make(chan error)

	for _, f := range filenames {
		go func(f string) {
			_, err := thumbnail.ImageFile(f)
			errors <- err
		}(f)
	}

	// 存在的问题：有可能导致上面创建的 goroutine 无法正常退出
	for range filenames {
		if err := <-errors; err != nil {
			return err // NOTE: incorrect: goroutine leak!
		}
	}

	return nil
}

func makeThumbnails5(filenames []string) (thumbfiles []string, err error) {
	type item struct {
		thumbfile string
		err       error
	}

	// 创建一个 buffered channel
	ch := make(chan item, len(filenames))

	for _, f := range filenames {
		go func(f string) {
			var it item
			it.thumbfile, it.err = thumbnail.ImageFile(f)
			ch <- it
		}(f)
	}

	for range filenames {
		it := <-ch

		// 存在的问题：如果发生错误，虽然不会有 goroutine leak，但是 ch 永远不会关闭
		if it.err != nil {
			return nil, it.err
		}

		thumbfiles = append(thumbfiles, it.thumbfile)
	}

	return thumbfiles, nil
}

// looping in parallel 的常用方式：使用 sync.WaitGroup 计数器
func makeThumbnails6(filenames <-chan string) int64 {
	sizes := make(chan int64)

	var wg sync.WaitGroup

	for f := range filenames {
		wg.Add(1)
		// worker
		go func(f string) {
			defer wg.Done()

			thumb, err := thumbnail.ImageFile(f)
			if err != nil {
				log.Println(err)
				return
			}

			info, _ := os.Stat(thumb) // OK to ignore error
			sizes <- info.Size()
		}(f)
	}

	// closer
	go func() {
		wg.Wait()
		close(sizes)
	}()

	// main goroutine 中统计 size
	var total int64
	for size := range sizes {
		total += size
	}
	return total
}
