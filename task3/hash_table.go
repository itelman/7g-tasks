package task3

func hash(data []byte, size int) int {
	h := 0
	for i := 0; i < len(data); i++ {
		h = h*31 + int(data[i])
	}
	if h < 0 {
		h = -h
	}
	return h % size
}

type Node struct {
	word  []byte
	count int
	next  *Node
}

type HashTable struct {
	buckets []*Node
	size    int
}

func NewHashTable(size int) *HashTable {
	return &HashTable{
		buckets: make([]*Node, size),
		size:    size,
	}
}

// Insert Time complexity: O(N*1)
// N - number of words in total (in the file)
func (ht *HashTable) Insert(word []byte) {
	index := hash(word, ht.size)
	node := ht.buckets[index]

	// Before adding a word, the code checks if it's already in 'words'
	// If it is, then it's not added, instead the count of its copy in 'words' is incremented by 1

	// If it isn't, the code goes to the next step, which is adding the word in 'words'
	for node != nil {
		if compare(node.word, word) == 0 {
			node.count++
			return
		}
		node = node.next
	}

	// A copy of the word (its byte slice) is added in 'words' instead of the word itself
	// to avoid the "pass by reference" issue
	newWord := make([]byte, len(word))
	copy(newWord, word)

	newNode := &Node{
		word:  newWord,
		count: 1,
		next:  ht.buckets[index],
	}

	ht.buckets[index] = newNode
}
