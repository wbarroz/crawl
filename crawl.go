package main

import(
    "github.com/gocolly/colly"
    "fmt"
    "regexp"
)

func main() {
    c := colly.NewCollector(
        colly.AllowedDomains("www.fundamentus.com.br",),
    )
    var (
        acronLink=regexp.MustCompile(`detalhes.*papel=[A-Z]{4}[0-9]{1,2}$`)
        mapOfMine=make(map[int]string)
        counter=0
        second=false
    )

    // Find and visit all links, "first sweep":
    c.OnHTML("a", func(e *colly.HTMLElement) {

        if(acronLink.MatchString(e.Attr("href"))){
            counter=0
            second=false
            e.Request.Visit(e.Attr("href"))
            //fmt.Printf("counter = %d\n",counter);
        }
    })


    c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
        //fmt.Println("First column of a table row:", e.Text)
        switch e.Text{
        case "?Papel":
            mapOfMine[counter]="papel"
        case "?Empresa":
            mapOfMine[counter]="empresa"
        case "Dia":
            mapOfMine[counter]="dia"
        case "?Valor de mercado":
            mapOfMine[counter]="valor"
        }
        //fmt.Printf("counter = %d\n",counter);
        counter+=1
    })

    c.OnHTML("tr td:nth-of-type(2)", func(e *colly.HTMLElement) {
        if!second{
            second=true
            counter=0
        }
        //fmt.Println("Second column of a table row:", e.Text)
        key,ok:=mapOfMine[counter]
        if ok{
            fmt.Printf("%s => %s\n",key,e.Text);
        }
        //fmt.Printf("counter = %d\n",counter);
        counter+=1
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

    //c.Visit("http://go-colly.org/")
    c.Visit("https://www.fundamentus.com.br/detalhes.php")
}
