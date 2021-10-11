package fetcher

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/hi20160616/exhtml"
	"github.com/hi20160616/ms-rfi/configs"
	"github.com/pkg/errors"
)

func TestFetchTitle(t *testing.T) {
	tests := &[]struct {
		url   string
		title string
	}{
		{"https://www.rfi.fr/cn/%E6%B8%AF%E6%BE%B3%E5%8F%B0/20211011-%E9%A9%AC%E8%8B%B1%E4%B9%9D%E8%B4%A8%E7%96%91%E8%94%A1%E8%8B%B1%E6%96%87%E6%83%B3%E6%94%B9%E5%8F%98%E5%9B%BD%E5%8F%B7", "马英九质疑蔡英文想改变国号"},
		{"https://www.rfi.fr/cn/%E4%B8%AD%E5%9B%BD/20211011-%E4%B8%AD%E5%9B%BD%E4%BC%81%E7%B4%AB%E9%87%91%E5%B8%83%E5%B1%80%E9%94%82%E8%B5%84%E6%BA%90-9-6%E4%BA%BF%E5%8A%A0%E5%85%83%E6%94%B6%E8%B4%AD%E5%8A%A0%E5%9B%BD%E6%96%B0%E9%94%82%E5%85%AC%E5%8F%B8-%E4%B8%A4%E6%9C%88%E5%86%85%E7%AC%AC%E4%BA%8C%E5%AE%97", "马英九质疑蔡英文想改变国号"},
	}
	if err := configs.Reset("../../"); err != nil {
		t.Error(err)
	}

	for _, tc := range *tests {
		a := NewArticle()
		u, err := url.Parse(tc.url)
		if err != nil {
			t.Error(err)
			return
		}
		// Dail
		a.U = u
		a.raw, a.doc, err = exhtml.GetRawAndDoc(a.U, timeout)
		if err != nil {
			t.Error(err)
		}

		if err != nil {
			t.Error(err)
			return
		}
		title, err := a.fetchTitle()
		if err != nil {
			t.Error(err)
			return
		}
		fmt.Println(title)
	}
}

func TestFetchUpdateTime(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{"https://www.rfi.fr/cn/%e6%b8%af%e6%be%b3%e5%8f%b0/20211011-%e9%a9%ac%e8%8b%b1%e4%b9%9d%e8%b4%a8%e7%96%91%e8%94%a1%e8%8b%b1%e6%96%87%e6%83%b3%e6%94%b9%e5%8f%98%e5%9b%bd%e5%8f%b7", "2021-10-11 15:57:05 +0800 UTC"},
		{"https://www.rfi.fr/cn/%E4%B8%AD%E5%9B%BD/20211011-%E4%B8%AD%E5%9B%BD%E4%BC%81%E7%B4%AB%E9%87%91%E5%B8%83%E5%B1%80%E9%94%82%E8%B5%84%E6%BA%90-9-6%E4%BA%BF%E5%8A%A0%E5%85%83%E6%94%B6%E8%B4%AD%E5%8A%A0%E5%9B%BD%E6%96%B0%E9%94%82%E5%85%AC%E5%8F%B8-%E4%B8%A4%E6%9C%88%E5%86%85%E7%AC%AC%E4%BA%8C%E5%AE%97", "2021-10-11 15:44:29 +0800 UTC"},
	}
	var err error
	for _, tc := range tests {
		a := NewArticle()
		a.U, err = url.Parse(tc.url)
		if err != nil {
			t.Error(err)
		}
		// Dail
		a.raw, a.doc, err = exhtml.GetRawAndDoc(a.U, timeout)
		if err != nil {
			t.Error(err)
		}
		tt, err := a.fetchUpdateTime()
		if err != nil {
			if !errors.Is(err, ErrTimeOverDays) {
				t.Error(err)
			}
		}
		ttt := tt.AsTime()
		got := shanghai(ttt)
		if got.String() != tc.want {
			t.Errorf("\nwant: %s\n got: %s", tc.want, got.String())
		}
	}
}

func TestFetchContent(t *testing.T) {
	tests := []struct {
		url  string
		want string
	}{
		{
			"https://www.rfi.fr/cn/%e6%b8%af%e6%be%b3%e5%8f%b0/20211011-%e9%a9%ac%e8%8b%b1%e4%b9%9d%e8%b4%a8%e7%96%91%e8%94%a1%e8%8b%b1%e6%96%87%e6%83%b3%e6%94%b9%e5%8f%98%e5%9b%bd%e5%8f%b7",
			"2021-06-02 15:44:33 +0800 UTC",
		},
	}
	var err error

	for _, tc := range tests {
		a := NewArticle()
		a.U, err = url.Parse(tc.url)
		if err != nil {
			t.Error(err)
		}
		// Dail
		a.raw, a.doc, err = exhtml.GetRawAndDoc(a.U, timeout)
		if err != nil {
			t.Error(err)
		}
		c, err := a.fetchContent()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(c)
	}
}

func TestFetchArticle(t *testing.T) {
	tests := []struct {
		url string
		err error
	}{
		{"https://www.rfi.fr/cn/%E4%B8%AD%E5%9B%BD/20201114-%E6%8D%B7%E5%85%8B%E5%86%9B%E6%83%85%E7%A0%94%E7%A9%B6%E8%AD%A6%E5%91%8A%E4%B8%AD%E7%BE%8E%E4%BF%84%E4%BA%89%E9%9C%B8%E6%88%96%E5%BC%95%E7%88%86%E7%AC%AC3%E6%AC%A1%E4%B8%96%E7%95%8C%E5%A4%A7%E6%88%98-%E5%8C%97%E4%BA%AC%E6%8C%87%E4%B8%8D%E4%BA%89%E9%9C%B8", ErrTimeOverDays},
		{"https://www.rfi.fr/cn/%e6%b8%af%e6%be%b3%e5%8f%b0/20211011-%e9%a9%ac%e8%8b%b1%e4%b9%9d%e8%b4%a8%e7%96%91%e8%94%a1%e8%8b%b1%e6%96%87%e6%83%b3%e6%94%b9%e5%8f%98%e5%9b%bd%e5%8f%b7", nil},
	}
	for _, tc := range tests {
		a := NewArticle()
		a, err := a.fetchArticle(tc.url)
		if err != nil {
			if !errors.Is(err, ErrTimeOverDays) {
				t.Error(err)
				return
			} else {
				fmt.Println("ignored test: ", tc.url)
			}
		} else {
			fmt.Println("pass test:\n", a.Content)
		}
	}
}
