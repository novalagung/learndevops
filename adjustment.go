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
	bookName := "Devops Tutorial"
	adClient := "ca-pub-1417781814120840"

	regex := regexp.MustCompile(`<title>(.*?)<\/title>`)

	basePath, _ := os.Getwd()
	bookPath := filepath.Join(basePath, "_book")

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
		isLandingPage := false
		newTitle := oldTitle
		if newTitle == "Introduction · GitBook" {
			isLandingPage = true
			newTitle = bookName
		} else {
			if titleParts := strings.Split(newTitle, "."); len(titleParts) > 2 {
				actualTitle := strings.TrimSpace(titleParts[2])

				if strings.Contains(actualTitle, "Go") || strings.Contains(actualTitle, "Golang") {
					// do nothing
				} else {
					titleParts[2] = fmt.Sprintf(" Golang %s", actualTitle)
				}

				newTitle = strings.Join(titleParts, ".")
			}

			newTitle = strings.Replace(newTitle, "· GitBook", fmt.Sprintf("- %s", bookName), -1)
		}
		htmlString = strings.Replace(htmlString, oldTitle, newTitle, -1)

		// ==== adjust meta for SEO purpose
		metaToFind := `<meta content=""name="description">`
		metaReplacement := metaToFind
		if isLandingPage {
			metaReplacement = `<meta content="Devops Tutorial" name="description">`
		}
		htmlString = strings.Replace(htmlString, metaToFind, fmt.Sprintf(`%s<script data-ad-client="%s" async src="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js"></script><script>(adsbygoogle = window.adsbygoogle || []).push({ google_ad_client: "%s", enable_page_level_ads: true }); </script>`, metaReplacement, adClient, adClient), -1)

		// ==== inject github stars button
		buttons := `<div style="position: fixed; top: 10px; right: 30px; padding: 10px; background-color: rgba(255, 255, 255, 0.7);">
			<a class="github-button" href="https://github.com/novalagung" data-size="large" aria-label="Follow @novalagung on GitHub">Follow @novalagung</a>
			<script async defer src="https://buttons.github.io/buttons.js"></script>
		</div>`
		htmlString = strings.Replace(htmlString, `</body>`, fmt.Sprintf("%s</body>", buttons), -1)

		buttonScript := `<script async defer src="https://buttons.github.io/buttons.js"></script>`
		htmlString = strings.Replace(htmlString, `</head>`, fmt.Sprintf("%s</head>", buttonScript), -1)

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

	// sitemap adjustment
	adjustSitemap("en")
	adjustSitemap("id")
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

func adjustSitemap(lang string) error {
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

		return nil
	})
	if err != nil {
		return err
	}

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
			Loc:        fmt.Sprintf("https://devops.novalagung.com/en/%s", each),
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
