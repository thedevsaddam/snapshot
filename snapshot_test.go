package snapshot

import (
	"sync"
	"testing"
	"time"
)

type user struct {
	Name, Email string
	Phone       int
	CreatedAt   time.Time
}

var userCollection *Collection

func TestMain(m *testing.M) {
	userCollection, _ = New("users")
	m.Run()
}

func TestCollection_Put(t *testing.T) {
	john := user{"John Doe", "john.doe@mail.com", 9898787, time.Now()}
	err := userCollection.Put("john", &john)
	if err != nil {
		t.Error("Failed to Put in collection!")
	}
}

func TestCollection_Get(t *testing.T) {
	john := user{}
	err := userCollection.Get("john", &john)
	if err != nil {
		t.Error("Failed to Get from collection!")
	}
	if john.Name != "John Doe" {
		t.Error("Failed to Get correct data")
	}
}

func TestCollection_Has(t *testing.T) {
	if !userCollection.Has("john") {
		t.Error("Failed to check using Has method!")
	}
}

func TestCollection_List(t *testing.T) {
	list, err := userCollection.List()
	if err != nil {
		t.Error(err)
	}
	if len(list) != 1 {
		t.Error("Failed to get collection keys list!")
	}
}

func TestCollection_TotalItem(t *testing.T) {
	if userCollection.TotalItem() != 1 {
		t.Error("Failed to count total item number!")
	}
}

func TestCollection_Remove(t *testing.T) {
	err := userCollection.Remove("john")
	if err != nil {
		t.Error("Failed to remove from collection!")
	}
}

func TestCollection_Flush(t *testing.T) {
	err := userCollection.Flush()
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkCollection_Put(b *testing.B) {
	john := user{"John Doe", "john.doe@mail.com", 9898787, time.Now()}
	for n := 0; n < b.N; n++ {
		go userCollection.Put("john", &john)
	}
}

func BenchmarkCollection_Get(b *testing.B) {
	john := user{}
	for n := 0; n < b.N; n++ {
		go userCollection.Get("john", &john)
	}
}

func BenchmarkCollection_Concurency(b *testing.B) {
	group := &sync.WaitGroup{}
	group.Add(2)
	go func() {
		for i := 0; i < b.N; i++ {
			john := user{"John Doe", "john.doe@mail.com", 9898787, time.Now()}
			userCollection.Put("john", john)
		}
		group.Done()
	}()
	go func() {
		for i := 0; i < b.N; i++ {
			john := user{}
			userCollection.Get("john", john)
		}
		group.Done()
	}()
	group.Wait()
}

func BenchmarkCollection_Has(b *testing.B) {
	for n := 0; n < b.N; n++ {
		userCollection.Has("john")
	}
}

func BenchmarkCollection_List(b *testing.B) {
	for n := 0; n < b.N; n++ {
		userCollection.List()
	}
}

func BenchmarkCollection_TotalItem(b *testing.B) {
	for n := 0; n < b.N; n++ {
		userCollection.TotalItem()
	}
}

func BenchmarkCollection_Remove(b *testing.B) {
	for n := 0; n < b.N; n++ {
		userCollection.Remove("john")
	}
}

func BenchmarkCollection_Flush(b *testing.B) {
	for n := 0; n < b.N; n++ {
		userCollection.Flush()
	}
}
