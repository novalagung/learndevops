package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	doAdjustment("en-US")
	doAdjustment("id")
}

type sitemapURL struct {
	Text       string `xml:",chardata"`
	Loc        string `xml:"loc"`
	Changefreq string `xml:"changefreq"`
	Priority   string `xml:"priority"`
}

type sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	Text    string       `xml:",chardata"`
	Xmlns   string       `xml:"xmlns,attr"`
	News    string       `xml:"news,attr"`
	Xhtml   string       `xml:"xhtml,attr"`
	Mobile  string       `xml:"mobile,attr"`
	Image   string       `xml:"image,attr"`
	URL     []sitemapURL `xml:"url"`
}

func doAdjustment(isoLang string) error {
	lang := strings.Split(isoLang, "-")[0]

	bookName := "Devops Tutorial"
	adClient := "ca-pub-1417781814120840"

	regex := regexp.MustCompile(`<title>(.*?)<\/title>`)

	basePath, _ := os.Getwd()
	bookPath := filepath.Join(basePath, "_book", lang)

	files := make([]string, 0)
	err := filepath.Walk(bookPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(info.Name()) != ".html" {
			return nil
		}

		files = append(files, filepath.Base(info.Name()))

		buf, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		htmlString := string(buf)

		// ==== remove invalid lang tag for EPUB validation
		htmlString = strings.Replace(htmlString, ` lang="" xml:lang=""`, "", -1)

		// ==== adjust title for SEO purpose
		oldTitle := regex.FindString(htmlString)
		oldTitle = strings.Replace(oldTitle, "<title>", "", -1)
		oldTitle = strings.Replace(oldTitle, "</title>", "", -1)
		isLandingPage := oldTitle == "Introduction · GitBook"
		newTitle := oldTitle
		if isLandingPage {
			newTitle = bookName
		} else {
			newTitle = strings.Replace(newTitle, "· GitBook", fmt.Sprintf("- %s", bookName), -1)
		}
		htmlString = strings.Replace(htmlString, oldTitle, newTitle, -1)

		// ==== adjust meta for SEO purpose
		metaToFind := `<meta content=""name="description">`
		metaReplacement := ""
		if isLandingPage {
			metaReplacement = `<meta content="` + bookName + `" name="description">`
		}
		metaReplacement = metaReplacement + `<meta http-equiv="content-language" content="` + isoLang + `"/><script data-ad-client="` + adClient + `" async src="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js"></script><script>(adsbygoogle = window.adsbygoogle || []).push({ google_ad_client: "` + adClient + `", enable_page_level_ads: true }); </script>`
		htmlString = strings.Replace(htmlString, metaToFind, metaReplacement, -1)

		// ==== inject github stars button
		buttonToFind := `</body>`
		buttonReplacement := `<div style="position: fixed; top: 10px; right: 30px; padding: 10px; background-color: rgba(255, 255, 255, 0.7);"><a class="github-button" href="https://github.com/novalagung" data-size="large" aria-label="Follow @novalagung on GitHub">Follow @novalagung</a><script async defer src="https://buttons.github.io/buttons.js"></script></div></body>`
		htmlString = strings.Replace(htmlString, buttonToFind, buttonReplacement, -1)

		// ==== inject github stars js script
		buttonScriptToFind := `</head>`
		buttonScriptReplacement := `<script async defer src="https://buttons.github.io/buttons.js"></script></head>`
		htmlString = strings.Replace(htmlString, buttonScriptToFind, buttonScriptReplacement, -1)

		// ==== update file
		err = ioutil.WriteFile(path, []byte(htmlString), info.Mode())
		if err != nil {
			return err
		}

		fmt.Println("  ==>", path)

		return nil
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	// ==== sitemap adjustment
	x := sitemap{
		Xmlns:  "http://www.sitemaps.org/schemas/sitemap/0.9",
		News:   "http://www.google.com/schemas/sitemap-news/0.9",
		Xhtml:  "http://www.w3.org/1999/xhtml",
		Mobile: "http://www.google.com/schemas/sitemap-mobile/1.0",
		Image:  "http://www.google.com/schemas/sitemap-image/1.1",
		URL:    make([]sitemapURL, 0),
	}
	for _, each := range files {
		x.URL = append(x.URL, sitemapURL{
			Loc:        `https://devops.novalagung.com/` + lang + `/` + each,
			Changefreq: "daily",
			Priority:   "0.5",
		})
	}

	buf, err := xml.Marshal(x)
	if err != nil {
		return err
	}

	siteMapPath := fmt.Sprintf("%s/sitemap.xml", bookPath)
	err = ioutil.WriteFile(siteMapPath, buf, 0644)
	if err != nil {
		return err
	}

	fmt.Println("  ==>", siteMapPath)
	return nil
}
