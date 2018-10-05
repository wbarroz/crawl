package main

import(
    "github.com/gocolly/colly"
    "fmt"
    "regexp"
    "strings"
    "strconv"
)

type stock struct{
    papel string
    empresa string
    varDia float64
    valor float64
}

func main() {
    c := colly.NewCollector(
        colly.AllowedDomains("www.fundamentus.com.br",),
    )
    var (
        acronLink=regexp.MustCompile(`detalhes.*papel=[A-Z]{4}[0-9]{1,2}$`)
        mapOfMine=make(map[int]string)
        counter=0
        second=false
        isStockVal=true
        workStock stock
    )

    // Find and visit all links, "first sweep":
    c.OnHTML("a", func(e *colly.HTMLElement) {

        if(acronLink.MatchString(e.Attr("href"))){
            counter=0
            second=false
            isStockVal=true
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
            mapOfMine[counter]="varDia"
        /*case "?Valor de mercado":
            mapOfMine[counter]="valor"
            */
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
            //for key,value:=range(mapOfMine){
                //switch value{
                switch key{
                case "papel":
                    workStock.papel=e.Text
                case "empresa":
                    workStock.empresa=e.Text
                case "varDia":
                    workStock.varDia,_=strconv.ParseFloat(strings.Replace(strings.Replace(e.Text,",",".",1),"%","",1),64)
                }
            //}
        }
        //fmt.Printf("counter = %d\n",counter);
        counter+=1
    })

    c.OnHTML("tr td:nth-of-type(4)", func(e *colly.HTMLElement) {
        if isStockVal{
            isStockVal=false
            workStock.valor,_=strconv.ParseFloat(strings.Replace(e.Text,",",".",1),64)
            fmt.Printf("Valor de mercado(cotação) => %s\n",e.Text)
            fmt.Println(workStock);
        }
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

    //c.Visit("http://go-colly.org/")
    c.Visit("https://www.fundamentus.com.br/detalhes.php")
}
