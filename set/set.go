package set

type Set[T comparable] struct{
    setMap map[T]bool
}

func New[T comparable]() *Set[T]{
    return &Set[T]{
        setMap: make(map[T]bool),
    }
}

func (set Set[T]) Add(value T){
    set.setMap[value] = true
}
func (set Set[T]) Remove(value T){
    delete(set.setMap, value)
}
func (set Set[T]) Has(value T) bool{
    return set.setMap[value]
}