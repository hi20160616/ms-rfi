package fetcher

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hi20160616/exhtml"
	"github.com/hi20160616/gears"
	"github.com/hi20160616/ms-rfi/configs"
	"github.com/pkg/errors"
)

func fetchLinks() ([]string, error) {
	rt := []string{}

	for _, rawurl := range configs.Data.MS.URL {
		links, err := getLinks(rawurl)
		if err != nil {
			return nil, err
		}
		rt = append(rt, links...)
	}
	rt = gears.StrSliceDeDupl(rt)
	rt = linksFilter(rt, `https://www.rfi.fr/cn/.+`)
	regs := []string{
		`/%E5%85%B3%E9%94%AE%E8%AF%8D/`,
		`/%E6%A1%A3%E6%A1%88%E8%B5%84%E6%96%99%E5%BA%93/`,
		`/%E6%B3%95%E5%B9%BF%E5%90%88%E4%BD%9C%E4%BC%99%E4%BC%B4/$`,
		`/%E4%B8%AD%E5%9B%BD/$`,
		`/%E6%B8%AF%E6%BE%B3%E5%8F%B0/$`,
		`/%E6%B3%95%E5%9B%BD/$`,
		`/%E4%BA%9A%E6%B4%B2/$`,
		`/%E9%9D%9E%E6%B4%B2/$`,
		`/%E7%BE%8E%E6%B4%B2/$`,
		`/%E6%AC%A7%E6%B4%B2/$`,
		`/%E4%B8%AD%E4%B8%9C/$`,
		`/%E4%BA%BA%E6%9D%83/$`,
		`/%E6%94%BF%E6%B2%BB/$`,
		`/%E7%BB%8F%E8%B4%B8/$`,
		`/%E7%A4%BE%E4%BC%9A/$`,
		`/%E7%A7%91%E6%8A%80%E4%B8%8E%E6%96%87%E5%8C%96/$`,
		`/%E6%9C%80%E6%96%B0%E6%B6%88%E6%81%AF/$`,
		`/%E5%9B%BD%E9%99%85/$`,
		`/%E4%BD%93%E8%82%B2/$`,
		`/%E6%B3%95%E6%96%B0%E7%A4%BE%E7%BB%8F%E6%B5%8E/$`,
		`/%E5%81%A5%E5%BA%B7/$`,
		`/%E5%90%8D%E6%B5%81/$`,
		`/%E7%89%B9%E5%88%AB%E4%B8%93%E9%A2%98/$`,
		`/%E8%A7%86%E9%A2%91$`,
		`/%E7%94%9F%E6%80%81/$`,
		`/%E5%86%8D%E6%AC%A1%E6%94%B6%E5%90%AC/$`,
		`/%E4%B8%93%E6%A0%8F%E6%A3%80%E7%B4%A2/$`,
		`/%E7%9B%B4%E6%92%AD$`,
		`/%E5%A6%82%E4%BD%95%E6%94%B6%E5%90%AC%E6%B3%95%E5%B9%BF$`,
		`/%E5%BA%94%E7%94%A8$`,
		`/%E8%A7%86%E9%A2%91$`,
		`/%E6%88%91%E4%BB%AC%E6%98%AF%E8%B0%81$`,
		`/%E6%BB%9A%E5%8A%A8%E6%96%B0%E9%97%BB/$`,
		`%E6%AC%A1%E6%92%AD%E9%9F%B3-%E5%8C%97%E4%BA%AC%E6%97%B6%E9%97%B4`,
	}
	for _, re := range regs {
		rt = kickOutLinksRegex(rt, re)
	}
	return rt, nil
}

// getLinksJson get links from a url that return json data.
func getLinksJson(rawurl string) ([]string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	raw, _, err := exhtml.GetRawAndDoc(u, 1*time.Minute)
	if err != nil {
		return nil, err
	}
	re := regexp.MustCompile(`"url":\s"(.*?)",`)
	rs := re.FindAllStringSubmatch(string(raw), -1)
	rt := []string{}
	for _, item := range rs {
		rt = append(rt, "https://"+u.Hostname()+item[1])
	}
	return gears.StrSliceDeDupl(rt), nil
}

func getLinks(rawurl string) ([]string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	links, err := exhtml.ExtractLinks(u.String())
	if err != nil {
		return nil, errors.WithMessagef(err, "[%s] cannot extract links from %s",
			configs.Data.MS.Title, rawurl)
	}
	return links, nil
	// return gears.StrSliceDeDupl(links), nil
}

// kickOutLinksMatchPath will kick out the links match the path,
func kickOutLinksRegex(links []string, reg string) []string {
	rt := []string{}
	re := regexp.MustCompile(reg)
	// path = "/" + url.QueryEscape(path) + "/"
	// path = url.QueryEscape(path)
	// reg = url.QueryEscape(reg)
	for _, link := range links {
		if !re.MatchString(link) {
			rt = append(rt, link)
		}
	}
	return rt
}

func linksFilter(links []string, regex string) []string {
	flinks := []string{}
	re := regexp.MustCompile(regex)
	s := strings.Join(links, "\n")
	flinks = re.FindAllString(s, -1)
	return flinks
}
