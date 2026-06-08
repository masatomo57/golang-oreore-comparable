
// 構造体自体をジェネリックにする（最初から型が決まっている）
query := NewQuery[User](db)
query.Where("age > 20").Find() // 常にUserを返す

query := NewQuery(db)
// Findだけが型を知っていれば良い
users, err := query.Where("age > 20").Find[User]()
posts, err := query.Where("status = 1").Find[Post]()

// 型T（int）のストリーム
stream := iterator.From([]int{1, 2, 3})

// intからstringへの変換は、関数で行うしかなかった
strStream := iterator.Map[int, string](stream, func(i int) string {
    return strconv.Itoa(i)
})

// メソッドが新しい型 U を受け取れるようになる！
// func (s Stream[T]) Map[U any](f func(T) U) Stream[U]
result := iterator.From([]int{1, 2, 3}).
    Map[string](func(i int) string { return strconv.Itoa(i) })
