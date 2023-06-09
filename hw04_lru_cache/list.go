package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	first  *ListItem
	last   *ListItem
	length int
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.first
}

func (l *list) Back() *ListItem {
	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	ret := &ListItem{v, nil, nil}
	if l.first == nil {
		l.first = ret
		l.last = ret
	} else {
		ret.Next = l.first
		l.first.Prev = ret
		l.first = ret
	}
	l.length++
	return ret
}

func (l *list) PushBack(v interface{}) *ListItem {
	ret := &ListItem{v, nil, nil}
	if l.last == nil {
		l.last = ret
		l.first = ret
	} else {
		ret.Prev = l.last
		l.last.Next = ret
		l.last = ret
	}
	l.length++
	return ret
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.first = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	if i.Next == nil {
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.first == i {
		return
	}
	if i.Next == nil {
		l.last = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	i.Prev.Next = i.Next
	i.Next = l.first
	l.first.Prev = i
	l.first = i
}

func NewList() List {
	return new(list)
}
