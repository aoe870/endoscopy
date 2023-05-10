package files

type Node struct {
	Name         string
	Path         string
	FileType     FileType
	LocationPath string
}

func (node *Node) readNode() {
	if node.FileType == Ctalogue {
		return
	}

}
