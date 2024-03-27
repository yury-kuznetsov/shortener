package uricoder

import (
	"context"
	"fmt"

	"github.com/yury-kuznetsov/shortener/internal/storage/memory"
)

func Example() {
	s := memory.NewStorage()
	coder := NewCoder(s)

	code, _ := coder.ToCode(context.Background(), "https://ya.ru", 0)
	uri, _ := coder.ToURI(context.Background(), code, 0)
	fmt.Println(uri)

	// Output:
	// https://ya.ru
}
