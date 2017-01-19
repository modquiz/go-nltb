package nltb

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/jinzhu/copier"
	tagger "github.com/modquiz/go-nltb/lib/tagger"
)

type TaggedWord struct {
	Word      string
	Tag       string
	byteStart int
}

type POSTag struct {
	goTagger *tagger.Tagger
}

/* Init parts of speech Tagging */
func (p *POSTag) Init() {
	_, file, _, _ := runtime.Caller(0)
	filename := filepath.Dir(file) + string(os.PathSeparator) + "brown" + string(os.PathSeparator)
	p.goTagger = tagger.New(filename)
}

/* Does Parts of Speech Tagging */
func (p *POSTag) Do(byteString []byte) []TaggedWord {
	taggedWord := p.goTagger.TagBytes(byteString)

	var returnTaggedWord []TaggedWord
	copier.Copy(&returnTaggedWord, &taggedWord)

	return returnTaggedWord
}
