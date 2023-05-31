package thumbnail_test

import (
	"code-example/bookcode/gopl.io/ch8/thumbnail"
	"log"
)

func makeThumbnails(filenames []string) {
	for _, f := range filenames {
		if _, err := thumbnail.ImageFile(f); err != nil {
			log.Println(err)
		}
	}
}
