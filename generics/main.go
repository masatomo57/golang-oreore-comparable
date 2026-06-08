package main

import "fmt"

// ============================================================
// 1. 関数 (Function)
// ============================================================

func Map[T, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// ============================================================
// 2. 構造体 (Struct)
// ============================================================

type Pair[T, U any] struct {
	First  T
	Second U
}

// ============================================================
// 3. インターフェース (Interface)
// ============================================================

type Container[T any] interface {
	Add(T)
	Get(int) T
	Len() int
}

// Container[T] を実装する具体型
type SliceContainer[T any] struct{ items []T }

func (c *SliceContainer[T]) Add(v T)    { c.items = append(c.items, v) }
func (c *SliceContainer[T]) Get(i int) T { return c.items[i] }
func (c *SliceContainer[T]) Len() int   { return len(c.items) }

// ============================================================
// 4. スライス型 (Slice type)
// ============================================================

type Stack[T any] []T

func (s *Stack[T]) Push(v T) { *s = append(*s, v) }
func (s *Stack[T]) Pop() T {
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

// ============================================================
// 5. マップ型 (Map type)
// ============================================================

type Registry[K comparable, V any] map[K]V

func (r Registry[K, V]) Set(k K, v V) { r[k] = v }
func (r Registry[K, V]) Get(k K) (V, bool) {
	v, ok := r[k]
	return v, ok
}

// ============================================================
// 6. チャネル型 (Channel type)
// ============================================================

type Queue[T any] chan T

func NewQueue[T any](size int) Queue[T] { return make(Queue[T], size) }
func (q Queue[T]) Send(v T)             { q <- v }
func (q Queue[T]) Recv() T              { return <-q }

// ============================================================
// 7. ポインタ型 (Pointer type)
// ============================================================

type Optional[T any] struct {
	ptr *T
}

func Some[T any](v T) Optional[T] { return Optional[T]{ptr: &v} }
func (o Optional[T]) IsNone() bool { return o.ptr == nil }
func (o Optional[T]) Value() T     { return *o.ptr }

// ============================================================
// 8. 配列型 (Array type)
// ============================================================

type Vec3[T any] [3]T

func (v Vec3[T]) At(i int) T { return v[i] }

// ============================================================
// 9. 関数型 (Function type)
// ============================================================

type Transformer[T, U any] func(T) U

func (t Transformer[T, U]) Apply(v T) U { return t(v) }

// ============================================================
// 10. 型エイリアス的な名前付き型 (Named type)
// ============================================================

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

type SortedSlice[T Ordered] []T

func (s SortedSlice[T]) Min() T { return s[0] }
func (s SortedSlice[T]) Max() T { return s[len(s)-1] }

// ============================================================
// 11. ネストした型パラメータ (Nested type parameters)
// ============================================================

type Result[T any] struct {
	value T
	err   error
}

func Ok[T any](v T) Result[T]      { return Result[T]{value: v} }
func Err[T any](e error) Result[T] { return Result[T]{err: e} }
func (r Result[T]) Unwrap() T      { return r.value }
func (r Result[T]) IsErr() bool    { return r.err != nil }

// Result[T] を値に持つ構造体 (型パラメータの合成)
type Cache[K comparable, V any] struct {
	store Registry[K, Result[V]]
}

func NewCache[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{store: make(Registry[K, Result[V]])}
}
func (c Cache[K, V]) Put(k K, v V) { c.store.Set(k, Ok(v)) }
func (c Cache[K, V]) Fetch(k K) (V, bool) {
	r, ok := c.store.Get(k)
	if !ok || r.IsErr() {
		var zero V
		return zero, false
	}
	return r.Unwrap(), true
}

// ============================================================
// 12. メソッドに型パラメータを持たせる (コンパイルエラー)
// ============================================================

type Converter struct{}

// これをしても怒られない
func (c Converter) Convert[T any](v T) T {
	return v
}

// ============================================================
// main
// ============================================================

func main() {
	// 1. 関数
	doubled := Map([]int{1, 2, 3}, func(n int) int { return n * 2 })
	fmt.Println("Map:", doubled)

	// 2. 構造体
	p := Pair[string, int]{First: "age", Second: 30}
	fmt.Println("Pair:", p)

	// 3. インターフェース (Container[T] を実装した SliceContainer を渡す)
	var c Container[string] = &SliceContainer[string]{}
	c.Add("hello")
	c.Add("world")
	fmt.Println("Container:", c.Get(0), c.Len())

	// 4. スライス型
	var s Stack[int]
	s.Push(10)
	s.Push(20)
	fmt.Println("Stack Pop:", s.Pop())

	// 5. マップ型
	r := make(Registry[string, int])
	r.Set("a", 1)
	v, _ := r.Get("a")
	fmt.Println("Registry:", v)

	// 6. チャネル型
	q := NewQueue[string](1)
	q.Send("hello")
	fmt.Println("Queue:", q.Recv())

	// 7. ポインタ型
	opt := Some(42)
	fmt.Println("Optional:", opt.Value())

	// 8. 配列型
	vec := Vec3[float64]{1.0, 2.0, 3.0}
	fmt.Println("Vec3:", vec.At(1))

	// 9. 関数型
	tr := Transformer[int, string](func(n int) string {
		return fmt.Sprintf("num=%d", n)
	})
	fmt.Println("Transformer:", tr.Apply(7))

	// 10. 名前付き型
	ss := SortedSlice[int]{1, 3, 5, 7}
	fmt.Println("SortedSlice Min/Max:", ss.Min(), ss.Max())

	// 11. ネスト
	cache := NewCache[string, int]()
	cache.Put("x", 99)
	val, ok := cache.Fetch("x")
	fmt.Println("Cache:", val, ok)
}
