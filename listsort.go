package piscine

// NodeI struct'ı, bağlı listelerde kullanılacak düğüm yapısını tanımlar.
type NodeI struct {
	Data int
	Next *NodeI
}

// ListSort fonksiyonu, bağlı listeyi küçükten büyüğe sıralar.
func ListSort(l *NodeI) *NodeI {
	if l == nil || l.Next == nil {
		return l
	}

	// Bağlı listeyi iki parçaya böler
	left, right := splitList(l)

	// Her iki parçayı sırala
	left = ListSort(left)
	right = ListSort(right)

	// İki sıralı parçayı birleştir
	return mergeLists(left, right)
}

// splitList fonksiyonu, bağlı listeyi iki parçaya böler.
func splitList(l *NodeI) (*NodeI, *NodeI) {
	var slow, fast, prev *NodeI
	slow, fast = l, l

	for fast != nil && fast.Next != nil {
		prev = slow
		slow = slow.Next
		fast = fast.Next.Next
	}

	if prev != nil {
		prev.Next = nil
	}

	return l, slow
}

// mergeLists fonksiyonu, iki sıralı bağlı listeyi birleştirir.
func mergeLists(left, right *NodeI) *NodeI {
	var result, current, temp *NodeI
	result = nil

	for left != nil && right != nil {
		if left.Data <= right.Data {
			temp = left
			left = left.Next
		} else {
			temp = right
			right = right.Next
		}

		if result == nil {
			result = temp
			current = result
		} else {
			current.Next = temp
			current = temp
		}
	}

	if left != nil {
		current.Next = left
	} else if right != nil {
		current.Next = right
	}

	return result
}
