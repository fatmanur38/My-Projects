package piscine

// SortedListMerge fonksiyonu, iki sıralı bağlı listeyi birleştirir.
func SortedListMerge(n1 *NodeI, n2 *NodeI) *NodeI {
	// İki liste de boşsa, boş bir liste döndür
	if n1 == nil {
		return n2
	} else if n2 == nil {
		return n1
	}

	var mergedList *NodeI
	var current *NodeI

	// İki liste üzerinde dolaşarak küçükten büyüğe doğru birleştir
	for n1 != nil && n2 != nil {
		if n1.Data < n2.Data {
			if mergedList == nil {
				mergedList = n1
				current = mergedList
			} else {
				current.Next = n1
				current = n1
			}
			n1 = n1.Next
		} else {
			if mergedList == nil {
				mergedList = n2
				current = mergedList
			} else {
				current.Next = n2
				current = n2
			}
			n2 = n2.Next
		}
	}

	// Eğer bir liste tamamen bitmediyse, kalan kısmı ekle
	if n1 != nil {
		current.Next = n1
	} else if n2 != nil {
		current.Next = n2
	}

	return mergedList
}
