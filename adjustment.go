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
	"time"
)

const (
	baseVersion = 1
	bookName    = "Devops Tutorial"
	ga4tagId    = "G-ZJMMV9WFV8"
)

func main() {
	doAdjustment()
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

func doAdjustment() error {
	regex := regexp.MustCompile(`<title>(.*?)<\/title>`)

	basePath, _ := os.Getwd()
	bookPath := filepath.Join(basePath, "_book")

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
		metaReplacement = metaReplacement + `<meta http-equiv="content-language" content="en-US"/>`
		htmlString = strings.Replace(htmlString, metaToFind, metaReplacement, -1)

		// ==== inject github stars button
		buttonToFind := `</body>`
		buttonReplacement := `<div style="position: fixed; top: 10px; right: 30px; padding: 10px; background-color: rgba(255, 255, 255, 0.7);"><a class="github-button" href="https://github.com/novalagung" data-size="large" aria-label="Follow @novalagung on GitHub">Follow @novalagung</a><script async defer src="https://buttons.github.io/buttons.js"></script></div>` + buttonToFind
		htmlString = strings.Replace(htmlString, buttonToFind, buttonReplacement, -1)

		// ==== inject github stars js script
		buttonScriptToFind := `</head>`
		buttonScriptReplacement := `<script async defer src="https://buttons.github.io/buttons.js"></script>` + buttonScriptToFind
		htmlString = strings.Replace(htmlString, buttonScriptToFind, buttonScriptReplacement, -1)

		// ==== inject ga4
		ga4propertyToFind := `</head>`
		ga4propertyReplacement := `<script async src="https://www.googletagmanager.com/gtag/js?id=` + ga4tagId + `"></script>
		<script>
			window.dataLayer = window.dataLayer || [];
			function gtag(){dataLayer.push(arguments);}
			gtag('js', new Date());
			gtag('config', '` + ga4tagId + `');
		</script>` + ga4propertyToFind
		htmlString = strings.Replace(htmlString, ga4propertyToFind, ga4propertyReplacement, -1)

		// ===== inject fb pixel
		fbPixelToFind := `</head>`
		fbPixelReplacement := `<script>!function(f,b,e,v,n,t,s){if(f.fbq)return;n=f.fbq=function(){n.callMethod?n.callMethod.apply(n,arguments):n.queue.push(arguments)};if(!f._fbq)f._fbq=n;n.push=n;n.loaded=!0;n.version='2.0';n.queue=[];t=b.createElement(e);t.async=!0;t.src=v;s=b.getElementsByTagName(e)[0];s.parentNode.insertBefore(t,s)}(window,document,'script','https://connect.facebook.net/en_US/fbevents.js');fbq('init','1247398778924723');fbq('track','PageView');</script><noscript><imgheight="1"width="1"style="display:none"src="https://www.facebook.com/tr?id=1247398778924723&ev=PageView&noscript=1"/></noscript>` + fbPixelToFind
		htmlString = strings.Replace(htmlString, fbPixelToFind, fbPixelReplacement, -1)

		// ==== inject adjustment css
		adjustmentCSSToFind := `</head>`
		adjustmentCSSReplacement := `<link rel="stylesheet" href="/adjustment.css?v=` + getVersion() + `">` + adjustmentCSSToFind
		htmlString = strings.Replace(htmlString, adjustmentCSSToFind, adjustmentCSSReplacement, -1)

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
			Loc:        `https://devops.novalagung.com/` + each,
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

	// ==== root index file adjustment (for google search verification)
	indexHTMLPath := filepath.Join(basePath, "_book", "index.html")
	indexHTMLBuf, err := ioutil.ReadFile(indexHTMLPath)
	if err != nil {
		return err
	}
	indexHTMLString := string(indexHTMLBuf)

	err = ioutil.WriteFile(indexHTMLPath, []byte(indexHTMLString), os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("  ==>", siteMapPath)
	return nil
}

func getVersion() string {
	return fmt.Sprintf("%d.%s", baseVersion, time.Now().Format("2006.01.02.150405"))
}
