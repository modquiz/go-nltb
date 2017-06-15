/*
Copyright (c) 2015 Eric Knapik, All Rights Reserved

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

  1. Redistributions of source code must retain the above copyright
     notice, this list of conditions and the following disclaimer.

  2. Redistributions in binary form must reproduce the above copyright
     notice, this list of conditions and the following disclaimer in the
     documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/

// This file is about the creation of tables that should
// be able to become import files. And using made tables to tag the parts
// of speech of possible copyright laced text. Once the text is tagged
// it can be called on by the file copyright.go to extract a possible
// notice, using a DFA.

package tagger

import (
	"strings"
)

// TagBytes returns a slice of TaggedWord objects
// representing that word in the sentence and the part of speech for
// that word
func (copyrightTagger *Tagger) TagBytes(rawBytes []byte) []TaggedWord {
	// ERROR AND SANITIZATION CHECKS
	var wrdArry = make([]TaggedWord, 0)
	if len(rawBytes) < 1 { // do I even need to do any work
		return wrdArry
	}

	// perform several regular expression subs and other so that the string is in a desired
	// form
	rawBytes = formatSent(rawBytes)
	// split the sentence propperly
	wrdArry = mkWrdArray(rawBytes)

	sentLength := len(wrdArry) + 1             // I need 1 more for the start of the sentence
	sentMatrix := make([][]float32, numOfTags) // Create the sentence Matrix
	for row := range sentMatrix {
		sentMatrix[row] = make([]float32, sentLength)
	}

	// initialize the first column
	var lastBestProb float32 = 1.0
	var lastBestTag = TagStrToInt["."]
	var currBestProb float32 = 0.0
	var currBestTag int = 0

	sentMatrix[TagStrToInt["."]][0] = 1.0 // the max probability something can be
	for wrdIndex := 0; wrdIndex < len(wrdArry); wrdIndex++ {
		for tagIndex := 0; tagIndex < numOfTags; tagIndex++ {
			var lowerCaseWord = strings.ToLower(wrdArry[wrdIndex].Word)

			var currTrans float32 = copyrightTagger.TransMatrix[lastBestTag][tagIndex]
			var currProb float32 = lastBestProb * currTrans

			if len(copyrightTagger.Dictionary[wrdArry[wrdIndex].Word]) != 0 { // has the word been seen before?

				if wrdArry[wrdIndex].Word == "." || wrdArry[wrdIndex].Word == "?" || wrdArry[wrdIndex].Word == "!" {
					sentMatrix[TagStrToInt["."]][wrdIndex+1] = 1.0
				} else {
					for _, tagObject := range copyrightTagger.Dictionary[wrdArry[wrdIndex].Word] {
						if TagIntToStr[tagIndex] == tagObject.tag {
							sentMatrix[tagIndex][wrdIndex+1] = currProb * tagObject.freq
						}
					}
				}
				// check for the word not caring about capitalization
			} else if len(copyrightTagger.Dictionary[lowerCaseWord]) != 0 {
				for _, tagObject := range copyrightTagger.Dictionary[lowerCaseWord] {
					if TagIntToStr[tagIndex] == tagObject.tag {
						sentMatrix[tagIndex][wrdIndex+1] = currProb * tagObject.freq
					}
				}
			} else { // Try to determine tag based on transitional probability and word itself
				if currTrans >= 0.7 {
					sentMatrix[tagIndex][wrdIndex+1] = currProb
				} else {
					likelyTag := tagUnkown(wrdArry[wrdIndex].Word)
					sentMatrix[TagStrToInt[likelyTag]][wrdIndex+1] = currProb * 0.95
				}
			}
			// see if this is the best transition for next column
			if currProb > currBestProb {
				currBestProb = currProb
				currBestTag = tagIndex
			}
		}
		lastBestProb = currBestProb
		lastBestTag = currBestTag
	}
	// Sentence Matrix Created.
	// Now walk through the matrix assigning the best tag to each word

	for wrdIndex := 0; wrdIndex < len(wrdArry); wrdIndex++ {
		var tagProb float32 = 0.0
		for tagIndex := 0; tagIndex < numOfTags; tagIndex++ {
			if sentMatrix[tagIndex][wrdIndex+1] > tagProb {
				tagProb = sentMatrix[tagIndex][wrdIndex+1]
				wrdArry[wrdIndex].Tag = TagIntToStr[tagIndex]
			}
		}
	}

	// compress numbers and propper nouns that might have been split
	//  wrdArry = compressNumInString(wrdArry)
	//	wrdArry = compressNP(wrdArry)

	return wrdArry
}

/*
// The DFA required for number compression
// into one number not number period number
// START is start state
// INTERM is intermediate state
// REJECT is reject state
// ACCEPT is accept state
func mkNumCompressDFA() map[Tri]int {
	// symbols := "cd."
	dfa := make(map[Tri]int)

	dfa[Tri{state: START, word: "X", pos: "cd"}] = START
	dfa[Tri{state: INTERM, word: "X", pos: "cd"}] = ACCEPT
	dfa[Tri{state: REJECT, word: "X", pos: "cd"}] = START
	dfa[Tri{state: ACCEPT, word: "X", pos: "cd"}] = START

	dfa[Tri{state: START, word: ".", pos: "."}] = INTERM
	dfa[Tri{state: INTERM, word: ".", pos: "."}] = REJECT
	dfa[Tri{state: REJECT, word: ".", pos: "."}] = REJECT
	dfa[Tri{state: ACCEPT, word: ".", pos: "."}] = REJECT

	dfa[Tri{state: START, word: "?", pos: "."}] = REJECT
	dfa[Tri{state: INTERM, word: "?", pos: "."}] = REJECT
	dfa[Tri{state: REJECT, word: "?", pos: "."}] = REJECT
	dfa[Tri{state: ACCEPT, word: "?", pos: "."}] = REJECT

	dfa[Tri{state: START, word: "!", pos: "."}] = REJECT
	dfa[Tri{state: INTERM, word: "!", pos: "."}] = REJECT
	dfa[Tri{state: REJECT, word: "!", pos: "."}] = REJECT
	dfa[Tri{state: ACCEPT, word: "!", pos: "."}] = REJECT

	dfa[Tri{state: START, word: "X", pos: "X"}] = REJECT
	dfa[Tri{state: INTERM, word: "X", pos: "X"}] = REJECT
	dfa[Tri{state: REJECT, word: "X", pos: "X"}] = REJECT
	dfa[Tri{state: ACCEPT, word: "X", pos: "X"}] = REJECT

	return dfa
}

// Given a slice of TaggedWord objects this will
// compress floating point and numbers containing periods
// that might have been split up by the tagger's formatting
func compressNumInString(inSent []TaggedWord) []TaggedWord {
	// dfa := mkNumCompressDFA()

	var finalSent []TaggedWord = make([]TaggedWord, 0)

	currentState := REJECT // the dead state
	var compNum []string = make([]string, 0)
	var saveNum []TaggedWord = make([]TaggedWord, 0)
	var saveStartByte int

	for _, taggedWord := range inSent {
		// Make the transition to the next state based on the input
		if taggedWord.tag == "." {
			currentState = dfa[Tri{state: currentState, word: taggedWord.word, pos: taggedWord.tag}]
		} else if taggedWord.tag == "cd" {
			currentState = dfa[Tri{state: currentState, word: "X", pos: taggedWord.tag}]
		} else {
			currentState = dfa[Tri{state: currentState, word: "X", pos: "X"}]
		}

		// Based on the current input decide how to save information
		if currentState == START {
			finalSent = append(finalSent, saveNum...)
			compNum = nil
			saveNum = nil
			compNum = append(compNum, taggedWord.word)
			saveStartByte = taggedWord.byteStart
			saveNum = append(saveNum, taggedWord)
		} else if currentState == INTERM {
			compNum = append(compNum, ".")
			saveNum = append(saveNum, taggedWord)
		} else if currentState == REJECT {
			finalSent = append(finalSent, saveNum...)
			saveNum = nil
			finalSent = append(finalSent, taggedWord)
		} else if currentState == ACCEPT {
			compNum = append(compNum, taggedWord.word)
			saveNum = nil
			saveNum = append(saveNum, TaggedWord{word: strings.Join(compNum, ""), tag: "cd", byteStart: saveStartByte})
			currentState = START
		}
	}
	if currentState != REJECT {
		finalSent = append(finalSent, saveNum...)
	}

	// returns the joined string with 'floating point' numbers combined
	return finalSent
}

// Similar to the compressNumInString this recompresses
// propper nouns that the tagger possibly separated to generalize
// tagging and account for words it has not seen before.
func compressNP(inSent []TaggedWord) []TaggedWord {
	var finalSent []TaggedWord = make([]TaggedWord, 0)

	prevTag := ""
	var saveWord []string = make([]string, 0)
	var saveByteStart int
	for _, taggedWord := range inSent {

		if prevTag == "np" && taggedWord.word == "." {
			saveWord = append(saveWord, ".")
			finalSent = append(finalSent, TaggedWord{word: strings.Join(saveWord, ""), tag: "np", byteStart: saveByteStart})
			saveWord = nil
		} else if prevTag == "np" && taggedWord.tag == "np" {
			finalSent = append(finalSent, TaggedWord{word: strings.Join(saveWord, ""), tag: "np", byteStart: saveByteStart})
			saveWord = nil
			saveWord = append(saveWord, taggedWord.word)
		} else if prevTag == "np" && taggedWord.word != "." {
			finalSent = append(finalSent, TaggedWord{word: strings.Join(saveWord, ""), tag: "np", byteStart: saveByteStart}, taggedWord)
			saveWord = nil
		} else if taggedWord.tag == "np" {
			saveWord = append(saveWord, taggedWord.word)
		} else {
			finalSent = append(finalSent, taggedWord)
		}

		saveByteStart = taggedWord.byteStart
		prevTag = taggedWord.tag

	}
	if prevTag == "np" {
		finalSent = append(finalSent, TaggedWord{word: strings.Join(saveWord, ""), tag: "np", byteStart: saveByteStart})
	}

	return finalSent
}

func toString(inSent []TaggedWord) string {
	var finalSent = make([]string, 0)
	for _, taggedWord := range inSent {
		finalSent = append(finalSent, taggedWord.word)
	}
	return strings.Join(finalSent, " ")
}

*/
