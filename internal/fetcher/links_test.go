package fetcher

import (
	"fmt"
	"testing"
)

func TestFetchLinks(t *testing.T) {
	ls, err := fetchLinks()
	if err != nil {
		t.Error(err)
		return
	}
	for _, e := range ls {
		fmt.Println(e)
	}
	fmt.Println(len(ls))
}
