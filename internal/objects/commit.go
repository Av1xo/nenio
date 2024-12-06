package objects

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
	Mode string
	Hash string
	Name string
}

type Commit struct {
	Tree string
	Parent string
	Author string
	Timestamp string
	Message string
}
