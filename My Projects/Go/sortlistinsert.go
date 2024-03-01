package piscine

// SortListInsert fonksiyonu, sıralı bağlı listeye yeni bir eleman ekler.
func SortListInsert(l *NodeI, data_ref int) *NodeI {
	newNode := &NodeI{Data: data_ref}

	// Başlangıç durumu: Eğer bağlı liste boşsa veya eklenen eleman listenin başına ekleniyorsa
	if l == nil || data_ref < l.Data {
		newNode.Next = l
		return newNode
	}

	// Diğer durumlar: Liste üzerinde gezinerek doğru konuma ekleme yap
	current := l
	for current.Next != nil && data_ref > current.Next.Data {
		current = current.Next
	}

	newNode.Next = current.Next
	current.Next = newNode

	return l
}
