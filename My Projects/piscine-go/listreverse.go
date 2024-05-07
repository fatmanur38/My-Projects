package piscine

func ListReverse(l *List) {
	if l.Head == nil || l.Head.Next == nil {
		return
	}
	var prev *NodeL
	curr := l.Head
	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}
	l.Head = prev
}
