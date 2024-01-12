package piscine

func ListRemoveIf(l *List, data_ref interface{}) {
	current := l.Head
	var previous *NodeL
	for current != nil {
		if current.Data == data_ref {
			if previous == nil {
				l.Head = current.Next
			} else {
				previous.Next = current.Next
			}
		} else {
			previous = current
		}
		current = current.Next
	}
}
