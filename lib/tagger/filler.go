package tagger

import (
	"fmt"
	"strings"
)

// create maps for converting tag string to integer and vice versa
var TagStrToInt = make(map[string]int)
var TagIntToStr = make(map[int]string)

// intitalizes the map used to convert the integer and string
// representation of part of speech tags
func initTagConversionMap() {

	TagStrToInt["bos"] = 0
	TagStrToInt["$"] = 1
	TagStrToInt["\""] = 2
	TagStrToInt["("] = 3
	TagStrToInt[")"] = 4
	TagStrToInt[","] = 5
	TagStrToInt["--"] = 6
	TagStrToInt["."] = 7
	TagStrToInt[":"] = 8
	TagStrToInt["cc"] = 9
	TagStrToInt["cd"] = 10
	TagStrToInt["dt"] = 11
	TagStrToInt["fw"] = 12
	TagStrToInt["jj"] = 13
	TagStrToInt["jj-tl"] = 14
	TagStrToInt["ls"] = 15
	TagStrToInt["nn"] = 16
	TagStrToInt["nn-tl"] = 17
	TagStrToInt["nns"] = 18
	TagStrToInt["nnp"] = 19
	TagStrToInt["np"] = 20
	TagStrToInt["nps"] = 21
	TagStrToInt["pos"] = 22
	TagStrToInt["pr"] = 23
	TagStrToInt["rb"] = 24
	TagStrToInt["sym"] = 25
	TagStrToInt["to"] = 26
	TagStrToInt["uh"] = 27
	TagStrToInt["vb"] = 28
	TagStrToInt["vbn"] = 29
	TagStrToInt["vbd"] = 30
	TagStrToInt["vbz"] = 31
	TagStrToInt["md"] = 32
	TagStrToInt["in"] = 33
	TagStrToInt["ap"] = 34
	TagStrToInt["at"] = 35
	TagStrToInt["bez"] = 36
	TagStrToInt["ppss"] = 37

	TagIntToStr[0] = "bos"
	TagIntToStr[1] = "$"
	TagIntToStr[2] = "\""
	TagIntToStr[3] = "("
	TagIntToStr[4] = ")"
	TagIntToStr[5] = ","
	TagIntToStr[6] = "--"
	TagIntToStr[7] = "."
	TagIntToStr[8] = ":"
	TagIntToStr[9] = "cc"
	TagIntToStr[10] = "cd"
	TagIntToStr[11] = "dt"
	TagIntToStr[12] = "fw"
	TagIntToStr[13] = "jj"
	TagIntToStr[14] = "jj-tl"
	TagIntToStr[15] = "ls"
	TagIntToStr[16] = "nn"
	TagIntToStr[17] = "nn-tl"
	TagIntToStr[18] = "nns"
	TagIntToStr[19] = "nnp"
	TagIntToStr[20] = "np"
	TagIntToStr[21] = "nps"
	TagIntToStr[22] = "pos"
	TagIntToStr[23] = "pr"
	TagIntToStr[24] = "rb"
	TagIntToStr[25] = "sym"
	TagIntToStr[26] = "to"
	TagIntToStr[27] = "uh"
	TagIntToStr[28] = "vb"
	TagIntToStr[29] = "vbn"
	TagIntToStr[30] = "vbd"
	TagIntToStr[31] = "vbz"
	TagIntToStr[32] = "md"
	TagIntToStr[33] = "in"
	TagIntToStr[34] = "ap"
	TagIntToStr[35] = "at"
	TagIntToStr[36] = "bez"
	TagIntToStr[37] = "ppss"
}

func addToDictionary(dictionary map[string][]TagFrequency, transMatrix [][]float32, path string) (map[string][]TagFrequency, [][]float32) {
	// read through the corpus file to populate the dictionary and transMatrix

	raw, err := Asset(path)

	//raw, err := ioutil.ReadFile(path)
	if err != nil && !strings.HasSuffix(path, "\\") {
		fmt.Println("-")
		fmt.Printf("could not read the file %v for tagging\n", path)
		fmt.Println(err)
		return dictionary, transMatrix
	}
	rawString := string(raw[:])

	prevTag := "."
	currTag := ""
	// I need to use the split feature on the coprus. So the input Corpus must have three
	// spaces between each word|~|tag pair. Once I have each word|~|tag pair I can
	// then split on the delimeter. Assumptions are made but I am assuming a safe input
	// file which I feel is acceptable, since if you save a file you know what it will
	// look like.
	textArry := strings.Split(rawString, " ")
	// textArry = textArry[:len(textArry)-1]

	for _, word := range textArry {
		wrdArry := strings.Split(word, SPLITCHARS)
		if len(wrdArry) > 1 {
			currTag = wrdArry[1]
			incrementUnigramWrd(dictionary, wrdArry[0], currTag)
			incrementTransMatrix(&transMatrix, TagStrToInt[prevTag], TagStrToInt[currTag])
			prevTag = currTag
		}
	}
	// everything is counted now convert the dictionary and TransMatrix to probabilistic
	convertDictToProb(dictionary)
	convertTransMatrixToProb(&transMatrix)

	return dictionary, transMatrix

}

// Initialization for the Tagger object
// Takes a file path and will create the unigram dictionary and transition
// matrix required for sentence tagging and NLP processing
func New(searchDir string) *Tagger {

	// initialize my TagStrToInt and TagIntToStr
	initTagConversionMap()

	// initialize the dictionary
	var dictionary = make(map[string][]TagFrequency)

	// Initialize the transition Matrix,
	var transMatrix = make([][]float32, numOfTags)
	for row := range transMatrix {
		transMatrix[row] = make([]float32, numOfTags)
	}

	// for every fileint he brown corpus do
	//err := filepath.Walk(searchDir, func(searchDir string, f os.FileInfo, err error) error {
	//	if f.Name() != "brown" {

	for _, asset := range AssetNames() {
		dictionary, transMatrix = addToDictionary(dictionary, transMatrix, asset)
	}

	//	}
	//	return nil
	//})
	/*
		if err != nil {
			fmt.Println(err)
		}
	*/
	// SETUP THE COPYRIGHT DFA
	// symbols, dfa := mkNoticeDFA()
	return &Tagger{Dictionary: dictionary, TransMatrix: transMatrix}
}
