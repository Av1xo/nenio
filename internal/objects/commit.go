package objects

import (
	"fmt"

	"lukechampine.com/blake3"
)

/*
    branch → M ↔ M1 ↔ M2 ↔ M3 ↔ M4 ↔ M5 ↔ MD ↔ MD1
				↘						↗
				branch ↔ D ↔ D1 ↔ D2 ↔ merge


	Directed Acyclic Graph


	COMMIT HASH
	PARENT HASH

	METADATA

	TREE HASH

*/

type TreeEntry struct {
	Hash string
	Name string
	Data IndexEntry
}

type Commit struct {
	Hash      string
	Tree      string
	Parent    string
	Author    string
	Timestamp string
	Message   string
}

func GenerateCommitHash(commit *Commit) string {
	data := commit.Tree + commit.Parent + commit.Author + commit.Timestamp + commit.Message
	hash := blake3.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}