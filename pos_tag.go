package nltb

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jinzhu/copier"
	tagger "github.com/modquiz/go-nltb/lib/tagger"
)

type TaggedWord struct {
	tagger.TaggedWord
}

/* Does Parts of Speech Tagging */
func POSTag(byteString []byte) []TaggedWord {
	_, file, _, _ := runtime.Caller(0)
	filename := filepath.Dir(file) + string(os.PathSeparator) + "brown" + string(os.PathSeparator)
	fmt.Println(filename)
	goTagger := tagger.New(filename)
	taggedWord := goTagger.TagBytes(byteString)

	var returnTaggedWord []TaggedWord
	copier.Copy(&returnTaggedWord, &taggedWord)
	fmt.Println(returnTaggedWord)
	return returnTaggedWord
}
