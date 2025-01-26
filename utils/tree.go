package utils

import (
	"fmt"
	"strings"
	"time"
)

// RadixNode represents a node in the Radix Tree
type RadixNode struct {
	prefix         string
	expirationTime time.Time
	isLeaf         bool                // Indicates if this node is a leaf (i.e., the end of a word)
	nodes          map[byte]*RadixNode // Child nodes mapped by the first character of their prefix
}

// NewRadixNode creates a new RadixNode
func NewRadixNode(prefix string, expirationTime time.Time, isLeaf bool) *RadixNode {
	return &RadixNode{
		prefix:         prefix,
		expirationTime: expirationTime,
		isLeaf:         isLeaf,
		nodes:          make(map[byte]*RadixNode),
	}
}

// Match computes the common substring between the current node's prefix and a given word.
// It returns three parts:
// - The common substring between the prefix and the word
// - The remaining part of the prefix
// - The remaining part of the word
func (node *RadixNode) Match(word string) (string, string, string) {
	x := 0
	// Compare characters of the prefix and word one by one
	for i := 0; i < len(node.prefix) && i < len(word); i++ {
		if node.prefix[i] != word[i] {
			break
		}
		x++
	}
	return node.prefix[:x], node.prefix[x:], word[x:]
}

// InsertMany inserts multiple words into the Radix Tree by calling the Insert method for each word.
// func (node *RadixNode) InsertMany(words []string) {
// 	for _, word := range words {
// 		node.Insert(word)
// 	}
// }

// Insert inserts a single word into the Radix Tree, handling multiple cases:
// - If the word matches the node's prefix exactly, mark the node as a leaf.
// - If no child node matches the first character of the word, create a new child node.
// - If a partial match exists with a child, adjust the tree structure to accommodate the new word.
func (node *RadixNode) Insert(word string, expirationTime time.Time) {
	if node.prefix == word {
		node.isLeaf = true // Case 1: The word matches the current node's prefix
		return
	}

	// Check if a child node starts with the first character of the word
	if child, exists := node.nodes[word[0]]; !exists {
		// Case 2: No matching child node, so create a new one
		node.nodes[word[0]] = NewRadixNode(word, expirationTime, true)
	} else {
		// Case 3 and 4: A partial or full match exists with a child node
		matching, remainingPrefix, remainingWord := child.Match(word)

		if remainingPrefix == "" {
			// Case 3: The word continues beyond the current child node's prefix
			child.Insert(remainingWord, expirationTime)
		} else {
			// Case 4: Partial match, so create intermediate nodes
			child.prefix = remainingPrefix

			// Create a new intermediate node
			newNode := NewRadixNode(matching, expirationTime, false)
			newNode.nodes[remainingPrefix[0]] = child

			// Update the parent's child reference
			node.nodes[matching[0]] = newNode

			if remainingWord == "" {
				newNode.isLeaf = true // Mark the new node as a leaf if no remaining word
			} else {
				newNode.Insert(remainingWord, expirationTime)
			}
		}
	}
}

// Find checks if a word exists in the Radix Tree.
// It traverses the tree based on the prefix matching until the word is fully matched or not found.
func (node *RadixNode) Find(word string) (time.Time, bool) {
	child, exists := node.nodes[word[0]]
	if !exists {
		return time.Time{}, false // No child node matches the first character
	}

	_, remainingPrefix, remainingWord := child.Match(word)

	if remainingPrefix != "" {
		return time.Time{}, false // The word can't match if there's leftover prefix
	} else if remainingWord == "" {
		return child.expirationTime, child.isLeaf // Word matches exactly if it's a leaf
	} else {
		return child.Find(remainingWord) // Continue searching in the child nodes
	}
}

// Delete removes a word from the Radix Tree if it exists.
// It adjusts the tree structure to ensure minimal nodes while maintaining correctness.
func (node *RadixNode) Delete(word string) bool {
	child, exists := node.nodes[word[0]]
	if !exists {
		return false // Word doesn't exist
	}

	_, remainingPrefix, remainingWord := child.Match(word)

	if remainingPrefix != "" {
		return false // Word can't exist if there's leftover prefix
	} else if remainingWord != "" {
		return child.Delete(remainingWord) // Continue deleting in child nodes
	} else if !child.isLeaf {
		return false // Node isn't a leaf, so the word doesn't exist
	} else {
		// Handle deletion of the node
		if len(child.nodes) == 0 {
			delete(node.nodes, word[0]) // Case 1: Node has no children, so delete it

			// Merge the current node with its only child if applicable
			if len(node.nodes) == 1 && !node.isLeaf {
				for _, mergingNode := range node.nodes {
					node.prefix += mergingNode.prefix
					node.isLeaf = mergingNode.isLeaf
					node.nodes = mergingNode.nodes
				}
			}
		} else if len(child.nodes) > 1 {
			child.isLeaf = false // Case 2: Node has multiple children, so mark it as non-leaf
		} else {
			// Case 3: Node has one child, merge it with the child
			for _, mergingNode := range child.nodes {
				child.prefix += mergingNode.prefix
				child.isLeaf = mergingNode.isLeaf
				child.nodes = mergingNode.nodes
			}
		}
		return true
	}
}

// PrintTree recursively prints the structure of the Radix Tree, showing prefixes and leaf nodes.
func (node *RadixNode) PrintTree(height int) {
	if node.prefix != "" {
		fmt.Printf("%s%s   %s\n", strings.Repeat("-", height), node.prefix, func() string {
			if node.isLeaf {
				return "(leaf)"
			}
			return ""
		}())
	}

	for _, child := range node.nodes {
		child.PrintTree(height + 1)
	}
}

// DeleteExpiredNodes traverses the Radix Tree and removes nodes with an expirationTime less than time.Now().
// It returns true if the current node itself should be deleted, otherwise false.
func (node *RadixNode) DeleteExpiredNodes() bool {
	now := time.Now()

	// Recursively check and delete expired nodes in the child nodes
	for char, child := range node.nodes {
		if child.DeleteExpiredNodes() {
			delete(node.nodes, char)
		}
	}

	// Check if the current node itself should be deleted
	if node.isLeaf && node.expirationTime.Before(now) {
		// Node is a leaf and expired
		return true
	}

	// If the node has no children and is not a leaf, it can be deleted
	if len(node.nodes) == 0 && !node.isLeaf {
		return true
	}

	// Optimize the tree if the current node has exactly one child
	if len(node.nodes) == 1 && !node.isLeaf {
		for _, child := range node.nodes {
			node.prefix += child.prefix
			node.isLeaf = child.isLeaf
			node.expirationTime = child.expirationTime
			node.nodes = child.nodes
		}
	}

	return false
}

// TestRadixTree runs various tests on the Radix Tree to verify its correctness.
// func TestRadixTree() bool {
// 	words := []string{"banana", "bananas", "bandana", "band", "apple", "all", "beast"}
// 	node := NewRadixNode("", false)
// 	node.InsertMany(words)

// 	// Verify all inserted words can be found
// 	for _, word := range words {
// 		if !node.Find(word) {
// 			return false
// 		}
// 	}

// 	// Verify non-existing words are not found
// 	if node.Find("bandanas") || node.Find("apps") {
// 		return false
// 	}

// 	// Verify deletions
// 	node.Delete("all")
// 	if node.Find("all") {
// 		return false
// 	}

// 	node.Delete("banana")
// 	if node.Find("banana") {
// 		return false
// 	}

// 	// Ensure other words are still present
// 	if !node.Find("bananas") {
// 		return false
// 	}

// 	return true
// }

// func main() {
// 	words := []string{"banana", "bananas", "bandanas", "bandana", "band", "apple", "all", "beast"}
// 	node := NewRadixNode("", false)
// 	node.InsertMany(words)

// 	fmt.Println("Words:", words)
// 	fmt.Println("Tree:")
// 	node.PrintTree(0)

// 	fmt.Println("\nRunning tests...")
// 	if TestRadixTree() {
// 		fmt.Println("All tests passed!")
// 	} else {
// 		fmt.Println("Some tests failed.")
// 	}
// }
