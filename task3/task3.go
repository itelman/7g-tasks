package task3

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Word entity for working with words, because by task requirement, strings are NOT allowed
type Word struct {
	data  []byte
	count int
}

func isLetter(b byte) bool {
	if b >= 'a' && b <= 'z' {
		return true
	}
	if b >= 'A' && b <= 'Z' {
		return true
	}
	return false
}

func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}

// Compare two words lexicographically
// If they're the same, the function returns 0, otherwise it's either 1 or -1
func compare(a, b []byte) int {
	// Get the length of the shortest word among the two words to be compared
	// It's needed to avoid the "index out of range" error in the next step,
	// where the words are compared letter by letter
	min := len(a)
	if len(b) < min {
		min = len(b)
	}

	// Compare the words letter by letter
	for i := 0; i < min; i++ {
		if a[i] < b[i] {
			return 1
		}
		if a[i] > b[i] {
			return -1
		}
	}

	// Compare the words by length
	if len(a) < len(b) {
		return 1
	}
	if len(a) > len(b) {
		return -1
	}

	return 0
}

// Time complexity: O(N*K)
// N - number of words in total (in the file)
// K - number of unique words
// Worst case: O(N^2) - all words in the file are unique
func addWord(words *[]Word, w []byte) {
	// Before adding a word, the code checks if it's already in 'words'
	// If it is, then it's not added, instead the count of its copy in 'words' is incremented by 1

	// If it isn't, the code goes to the next step, which is adding the word in 'words'
	for i := range *words {
		if compare((*words)[i].data, w) == 0 {
			(*words)[i].count++
			return
		}
	}

	// A copy of the word (its byte slice) is added in 'words' instead of the word itself
	// to avoid the "pass by reference" issue
	newWord := make([]byte, len(w))
	copy(newWord, w)

	*words = append(*words, Word{
		data:  newWord,
		count: 1,
	})
}

func Run() {
	if len(os.Args) < 2 {
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		return
	}
	defer file.Close()

	// var words []Word
	hashTable := NewHashTable(100003)

	buffer := make([]byte, 4096)
	// 'currentWord' is used to store the letters of the current word the code is working with
	// It's needed, because the file content is read character by character, and
	// words are put together that way
	var currentWord []byte

	for {
		// The file content is read in bytes, because by task requirement, strings are NOT allowed
		n, err := file.Read(buffer)
		if n > 0 {
			for i := 0; i < n; i++ {
				b := buffer[i]

				// If a character (its byte) happens to be a letter, its "lowered" version
				// is appended to 'currentWord' (because by task requirement,
				// the code has to be case-insensitive)

				// If it isn't a letter, the "putting the word together" process stops,
				// and the assembled word is added to 'words', before going next
				if isLetter(b) {
					currentWord = append(currentWord, toLower(b))
				} else {
					if len(currentWord) > 0 {
						// addWord(&words, currentWord)
						hashTable.Insert(currentWord)
						currentWord = nil
					}
				}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}

	// Add the last word that was assembled before the code reached the end of the file
	if len(currentWord) > 0 {
		// addWord(&words, currentWord)
		hashTable.Insert(currentWord)
	}

	var words []Word

	for i := 0; i < hashTable.size; i++ {
		node := hashTable.buckets[i]
		for node != nil {
			words = append(words, Word{
				data:  node.word,
				count: node.count,
			})
			node = node.next
		}
	}

	// All collected words will be sorted by:
	// 1. count descending
	// 2. lexicographically descending
	sort.Slice(words, func(i, j int) bool {
		if words[i].count == words[j].count {
			return compare(words[i].data, words[j].data) < 0
		}
		return words[i].count > words[j].count
	})

	limit := 20
	if len(words) < 20 {
		limit = len(words)
	}

	for i := 0; i < limit; i++ {
		fmt.Printf("%7d %s\n", words[i].count, words[i].data)
	}
}
