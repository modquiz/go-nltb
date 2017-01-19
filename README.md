# go-nltoolbox
Go Natural Language Toolbox 

The goal of this project is to provide many natural language modules in one easy to use package similar to NLTK.  

This package will either wrap existing modules, fork existing modules, or include new code

I'm trying to make the API as similiar to thy python NLTK as well.  

### Part of Speech Tagger

The first module I'm creating is a Part of Speech Tagger based of:

https://github.com/EKnapik/goTagger

I have adapted the module to use the Brown Corpus.

You can use this module by:

posTagger = nltb.POSTag{}
posTagger.Init()
taggedWord := posTagger.Do([]byte(str))
