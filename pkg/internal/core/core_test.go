package core

import (
	"fmt"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	path := "../../../cmd/rack/plugins/example/test-1.so"
	p, err := loader(path)
	if err != nil {
		t.Fatal(err.Error())
	}
	testSym, err := p.Lookup("Test")
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = testSym.(func(int, int) error)(10, 11); err != nil {
		t.Fatal(err.Error())
	}

	myIntSym, err := p.Lookup("MyInt")
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Printf("myIntSym=%d\n", *myIntSym.(*int))
	p = nil
	//======================================

	// go func() {
	// 	path2 := "example/test2.so"
	// 	ctx = context.Background()

	// 	p2, err := d.Load(ctx, path2)
	// 	if err != nil {
	// 		t.Fatal(err.Error())
	// 	}
	// 	testSym, err = p2.Lookup("Test")
	// 	if err != nil {
	// 		t.Fatal(err.Error())
	// 	}
	// 	if err = testSym.(func(int, int) error)(12, 13); err != nil {
	// 		t.Fatal(err.Error())
	// 	}
	// }()

	path2 := "../../../cmd/rack/plugins/example/test-2.so"

	p2, err := loader(path2)
	if err != nil {
		t.Fatal(err.Error())
	}
	testSym2, err := p2.Lookup("Test")
	if err != nil {
		t.Fatal(err.Error())
	}
	if err = testSym2.(func(int, int) error)(12, 13); err != nil {
		t.Fatal(err.Error())
	}

	time.Sleep(2 * time.Second)

	if err = testSym.(func(int, int) error)(10, 11); err != nil {
		t.Fatal(err.Error())
	}

}
