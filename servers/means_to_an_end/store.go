package means_to_an_end

type Store interface {
	Query(minTime int32, maxTime int32) int32
	Insert(timestamp int32, price int32)
}

type InMemoryStore struct {
	store map[int32]int32 // timestamp -> price
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[int32]int32),
	}
}

func (s *InMemoryStore) Query(minTime int32, maxTime int32) int32 {
	var total int64 = 0
	var count int32 = 0
	for timestamp, price := range s.store {
		if timestamp >= minTime && timestamp <= maxTime {
			total += int64(price)
			count++
		}
	}
	return CalculateAverage(total, count)
}

func CalculateAverage(total int64, count int32) int32 {
	if count == 0 {
		return 0
	}
	return int32(total / int64(count))
}

func (s *InMemoryStore) Insert(timestamp int32, price int32) {
	// only insert if the timestamp does not already exist
	if _, exists := s.store[timestamp]; !exists {
		s.store[timestamp] = price
	}
}
