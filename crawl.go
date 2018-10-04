package main

import(
    "github.com/gocolly/colly"
    "fmt"
    "regexp"
)

func main() {
    c := colly.NewCollector()
    var acronLink = regexp.MustCompile(`detalhes.*papel=[A-Z]{4}[0-9]{1,2}$`)

    // Find and visit all links
    c.OnHTML("a", func(e *colly.HTMLElement) {

        if(acronLink.MatchString(e.Attr("href"))){
            e.Request.Visit(e.Attr("href"))
        }
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

    //c.Visit("http://go-colly.org/")
    c.Visit("https://www.fundamentus.com.br/detalhes.php")
}
