package tagger

import (
	"bytes"
	"regexp"
	"strings"
)

// global const:
const numOfTags int = 38

// The character used to split the string from the part of speech tag in the Corpus
const SPLITCHARS string = "/"

// global regex
var copyright = regexp.MustCompile("(\\\\[(]co)")

// A struct/pair for the dictionary value
// The dictionary actually stores an array of these.
type TagFrequency struct {
	tag  string
	freq float32
}

// The Tagger Object
type Tagger struct {
	Dictionary  map[string][]TagFrequency
	TransMatrix [][]float32
	// for the copyright extraction
	CopyrightDFA  map[Tri]int
	CopyrightSyms string
}

type TaggedWord struct {
	Word      string
	Tag       string
	byteStart int
}

// three variable structure used in DFA translation
type Tri struct {
	state int
	word  string
	pos   string
}

// This is the counter of tag transitions. Moving from one part of speech tag
// to the other. When reading the input corpus this function is called to
// increment/make note of every part of speech tag transition.
// transitionOccurances = transMatrix[prev POS tag][current POS tag]
func incrementTransMatrix(transMatrix *[][]float32, prevTagIndex int, currTagIndex int) {
	(*transMatrix)[prevTagIndex][currTagIndex]++
}

// Given the unigram word dictionary, a word and the given part of speech
// tag for the word this will increment if the word already existed in the dictionary
// if the word did not this will create a new entry and set the times seen to 1
func incrementUnigramWrd(dictionary map[string][]TagFrequency, word string, tag string) {
	// dictionary is the map used for unigram word count/frequency
	// it is a key->slice of TagFrequency objects
	if tag == "nil" {
		return
	}
	if dictionary[word] != nil {
		for i := 0; i < len(dictionary[word]); i++ {
			if tag == dictionary[word][i].tag {
				dictionary[word][i].freq++
				return
			}
		}
		dictionary[word] = append(dictionary[word], TagFrequency{tag, 1})
		return
	} else {
		dictionary[word] = append(dictionary[word], TagFrequency{tag, 1})
		return
	}
}

// This will convert the dictionary which was in the form of
// counted occurances into a dictionary of probability for each part of speech
// tag given a specific word
func convertDictToProb(dictionary map[string][]TagFrequency) {
	// dictionary is a global variable
	var total float32
	for key := range dictionary {
		total = 0
		for i := 0; i < len(dictionary[key]); i++ {
			total = total + dictionary[key][i].freq
		}
		for i := 0; i < len(dictionary[key]); i++ {
			dictionary[key][i].freq = dictionary[key][i].freq / total
		}
	}
}

// This will convert the Transition Matrix to the probability
// Transition matrix the likelyhood of a given part of speech tag transition.
// Moving from tag A to tag B will result in what probility.
// transMatrix[FromTagA][ToTagB] = Probability X
// This is where the smoothing will be implemented
// I am using Laplace Smoothing across the transitional probability
// This means that every transition has a small probability of happeing
func convertTransMatrixToProb(transMatrix *[][]float32) {
	// transMatrix is a global variable
	var total float32

	for row := 0; row < numOfTags; row++ {
		total = float32(numOfTags)
		for col := 0; col < numOfTags; col++ {
			total += (*transMatrix)[row][col]
		}

		for col := 0; col < numOfTags; col++ {
			(*transMatrix)[row][col] = ((*transMatrix)[row][col] + 1) / total
		}
	}
}

// Given a word with an unknown part of speech. Using a model based from the
// Brill tagger, Krymolowski and Roth 1998 research (http://www.aclweb.org/anthology/P98-2186)
// This returns a guessed part of speech for unknown words
func tagUnkown(word string) string {

	// perform an N for loop checking for integer ascii value
	var i int
	for i = 0; i < len(word); i++ {
		if word[i] > 47 && word[i] < 58 {
			return "cd"
		}
	}

	loWord := strings.ToLower(word)

	switch {
	case strings.HasSuffix(loWord, "able"):
		return "jj"
	case strings.HasSuffix(loWord, "ible"):
		return "jj"
	case strings.HasSuffix(loWord, "ic"):
		return "jj"
	case strings.HasSuffix(loWord, "ous"):
		return "jj"
	case strings.HasSuffix(loWord, "al"):
		return "jj"
	case strings.HasSuffix(loWord, "ful"):
		return "jj"
	case strings.HasSuffix(loWord, "less"):
		return "jj"
	case strings.HasSuffix(loWord, "ly"):
		return "rb"
	case strings.HasSuffix(loWord, "ate"):
		return "vb"
	case strings.HasSuffix(loWord, "fy"):
		return "vb"
	case strings.HasSuffix(loWord, "ize"):
		return "vb"
	}

	// perform an N for loop checking for capital letter
	for i = 0; i < len(word); i++ {
		if word[i] > 64 && word[i] < 91 {
			return "np"
		}
	}

	switch {
	case strings.HasSuffix(loWord, "ion"):
		return "nn"
	case strings.HasSuffix(loWord, "ess"):
		return "nn"
	case strings.HasSuffix(loWord, "ment"):
		return "nn"
	case strings.HasSuffix(loWord, "er"):
		return "nn"
	case strings.HasSuffix(loWord, "or"):
		return "nn"
	case strings.HasSuffix(loWord, "ist"):
		return "nn"
	case strings.HasSuffix(loWord, "ism"):
		return "nn"
	case strings.HasSuffix(loWord, "ship"):
		return "nn"
	case strings.HasSuffix(loWord, "hood"):
		return "nn"
	case strings.HasSuffix(loWord, "ology"):
		return "nn"
	case strings.HasSuffix(loWord, "ty"):
		return "nn"
	case strings.HasSuffix(loWord, "y"):
		return "nn"
	default:
		return "fw"
	}
}

// Performs several string substitutions so that the tagger has an easier job
// These calls are to substitute parts of the string for other parts
// Once the sentence is formatted correctly it returns the string
func formatSent(rawBytes []byte) []byte {
	// to ensure a propper formatting.
	// replace weird copyright symbols
	// replaces \(co with (c)
	rawBytes = copyright.ReplaceAll(rawBytes, []byte("(c) ")) // added extra space to preserve byte offset

	// replace contractions
	// for byte preservation can not do these, but for more accurate tagging
	// replacing contractions can be useful
	/*
		rawBytes = bytes.Replace(rawBytes, []byte("ain't"), []byte("are not"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("won't"), []byte("will not"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("can't"), []byte("cannot"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("n't"), []byte(" not"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("'re"), []byte(" are"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("'m"), []byte(" am"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("'ll"), []byte(" will"), -1)
		rawBytes = bytes.Replace(rawBytes, []byte("'ve"), []byte(" have"), -1)
	*/
	return rawBytes
}

// returns true if the given byte is a white space character
func isSpace(b ...byte) bool {

	return bytes.Contains([]byte(" \n\r\t"), b)
}

// returns true if the given byte is a ASCII symbolic character
func isSymbol(b ...byte) bool {
	return bytes.Contains([]byte("~!`@#$%^&*()[]_+-=|}{:;'\"/\\.?><,"), b)
}

// Given a slice of raw bytes will convert this into a slice of
// TaggedWord objects with no tag set. This slice of TaggedWord objects will
// then be given to the tagger for determining the part of speech tag
func mkWrdArray(rawBytes []byte) []TaggedWord {

	currByte := 0
	wordStart := currByte
	var taggedWords []TaggedWord = make([]TaggedWord, 0)

	for currByte < len(rawBytes) {
		if isSpace(rawBytes[currByte]) {
			if wordStart != currByte { // add the word if I can
				taggedWords = append(taggedWords, TaggedWord{Word: string(rawBytes[wordStart:currByte]), Tag: "", byteStart: wordStart})
			}
			currByte++
			wordStart = currByte
		} else if isSymbol(rawBytes[currByte]) {
			if wordStart != currByte { // add the word if I can
				taggedWords = append(taggedWords, TaggedWord{Word: string(rawBytes[wordStart:currByte]), Tag: "", byteStart: wordStart})
			}
			wordStart = currByte
			currByte++
			taggedWords = append(taggedWords, TaggedWord{Word: string(rawBytes[wordStart:currByte]), Tag: "", byteStart: wordStart})
			wordStart = currByte
		} else {
			currByte++
		}
	}
	taggedWords = append(taggedWords, TaggedWord{Word: string(rawBytes[wordStart:currByte]), Tag: "", byteStart: wordStart})
	return taggedWords
}
