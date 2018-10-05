package main

import(
    "github.com/gocolly/colly"
    "fmt"
    "regexp"
    "strings"
    "strconv"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

type stock struct{
    papel string
    empresa string
    varDia float64
    valor float64
}

func(this stock)is_bigger_than(that stock)bool{
    return this.valor>that.valor
}

func insort(input stock,array[]stock)[]stock{
    if(len(array)!=0){
        for i,item:=range(array){
            if input.is_bigger_than(item){
                tail:=array[i:]
                new_tail:=make([]stock,len(tail))
                for j,source:=range(tail){
                    new_tail[j]=source
                }
                fmt.Println("tail",new_tail)
                head:=append(array[0:i],input)
                new_head:=make([]stock,len(head))
                for j,source:=range(head){
                    new_head[j]=source
                }
                fmt.Println("head",new_head)
                fmt.Println("comb",append(new_head,new_tail...))
                return append(new_head,new_tail...)
            }
        }
    }
    return append(array,input)
}

func main() {
    c := colly.NewCollector(
        colly.AllowedDomains("www.fundamentus.com.br",),
        colly.MaxDepth(2),
    )
    var (
        acronLink=regexp.MustCompile(`detalhes.*papel=[A-Z]{4}[0-9]{1,2}$`)
        mapOfMine=make(map[int]string)
        counter=0
        second=false
        isStockVal=true
        workStock stock
        //sortedArray[10]stock
        sortedArray[]stock
        //limit=100
    )

    // Find and visit all links, "first sweep":
    c.OnHTML("a", func(e *colly.HTMLElement) {
        /*if limit==0{
            return
        }else{
            limit-=1
        }*/
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
            sortedArray=insort(workStock,sortedArray)
            if(len(sortedArray)>10){
                sortedArray=sortedArray[0:10]
            }
            fmt.Println(sortedArray);
        }
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

    //c.Visit("http://go-colly.org/")
    c.Visit("https://www.fundamentus.com.br/detalhes.php")

    /*c.OnScraped(func(r *colly.Response) {*/
        fmt.Println("Finished")//, r.Request.URL)
        db, err := sql.Open("mysql", "root:supersecret@tcp(127.0.0.1:3306)/redventures")
        if err != nil {
            panic(err.Error())
        }
        //defer db.Close()
        // Prepare statement for inserting data
        //stmtIns, err := db.Prepare("INSERT INTO stocks VALUES( ?,?, )") // ? = placeholder
        stmt, err := db.Prepare("INSERT stocks SET papel=?,empresa=?,varDia=?,valor=?")
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
        for _,source:=range sortedArray{
            stmt.Exec(source.papel,source.empresa,source.varDia,source.valor)
            if err != nil {
                panic(err.Error()) // proper error handling instead of panic in your app
            }
        }
        db.Close()
    //})
    fmt.Println(sortedArray);
}
