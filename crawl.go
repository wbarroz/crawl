package main

//All the stuff we need:
import(
    "github.com/gocolly/colly" //This lib does most of the heavy-lifting
    "fmt"
    "regexp"                //For pattern filtering of the fetches
    "strings"
    "strconv"
    _ "github.com/go-sql-driver/mysql" //These two for access mysql
    "database/sql"
)

//Local struct resembling the stocks table in DBMS
type stock struct{
    papel string
    empresa string
    varDia float64
    valor float64
}

//More methods like this one can be very useful
//in the future; now mostly cosmetic
func(this stock)is_bigger_than(that stock)bool{
    return this.valor>that.valor
}

//An important piece of work is done here:
//the descending sorting of the stocks.
//Each fetched stock is compared against the
//present set,returning in an ordered set at the end
func insort(input stock,array[]stock)[]stock{
    if(len(array)!=0){
        for i,item:=range(array){
            if input.is_bigger_than(item){
                tail:=array[i:]
                new_tail:=make([]stock,len(tail))
                //it SHOULD BE a simple cut and paste(using slicing
                //or copy)but THIS is needed(set apart AND directly
                //address memory):
                for j,source:=range(tail){
                    new_tail[j]=source
                }
                fmt.Println("tail",new_tail)
                head:=append(array[0:i],input)
                new_head:=make([]stock,len(head))
                //here we go again:
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

    //This colly library rox! First, we initialize the collector object:
    c := colly.NewCollector(
        colly.AllowedDomains("www.fundamentus.com.br",), //Only here!
        colly.MaxDepth(2), //No need to fetch the whole site
    )
    var (
        acronLink=regexp.MustCompile(`detalhes.*papel=[A-Z]{4}[0-9]{1,2}$`) //only the pages we need
        mapOfMine=make(map[int]string)
        counter=0
        second=false
        isStockVal=true
        workStock stock
        sortedArray[]stock
    )

    //And register each of the callbacks, the most important being
    //the "OnHTML" ones, that brings the page content:

    // Here, we set up for find and visit the intended links:
    c.OnHTML("a", func(e *colly.HTMLElement) {
        //search all links from this page, matching
        //regex defined criterium before
        if(acronLink.MatchString(e.Attr("href"))){
            counter=0
            second=false
            isStockVal=true
            e.Request.Visit(e.Attr("href"))
            //fmt.Printf("counter = %d\n",counter);
        }
    })


    //Most of the interesting field tags where at first column:
    c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
        //When they match, we save them in a dictionary:
        switch e.Text{
        case "?Papel":
            mapOfMine[counter]="papel"
        case "?Empresa":
            mapOfMine[counter]="empresa"
        case "Dia":
            mapOfMine[counter]="varDia"
        }
        //fmt.Printf("counter = %d\n",counter);
        counter+=1
    })

    //Conversely, the respective fields are at the second column...
    c.OnHTML("tr td:nth-of-type(2)", func(e *colly.HTMLElement) {
        //...and they're ordered according to dictionary:
        if!second{
            second=true
            counter=0
        }
        //fmt.Println("Second column of a table row:", e.Text)
        key,ok:=mapOfMine[counter]
        if ok{
            fmt.Printf("%s => %s\n",key,e.Text);
            //Filling the blanks:
            switch key{
            case "papel":
                workStock.papel=e.Text
            case "empresa":
                workStock.empresa=e.Text
            case "varDia":
                workStock.varDia,_=strconv.ParseFloat(strings.Replace(strings.Replace(e.Text,",",".",1),"%","",1),64)
            }
        }
        counter+=1 //next row
    })

    //The ultimate field: stock value(always the fourth column of first row):
    c.OnHTML("tr td:nth-of-type(4)", func(e *colly.HTMLElement) {
        if isStockVal{
            isStockVal=false //only the first
            workStock.valor,_=strconv.ParseFloat(strings.Replace(e.Text,",",".",1),64)
            fmt.Printf("Valor de mercado(cotação) => %s\n",e.Text)
            fmt.Println(workStock);
            //Checking in stock info in the sorting process:
            sortedArray=insort(workStock,sortedArray)
            //The sorted array may be bigger than 10...
            if(len(sortedArray)>10){
                //...if it's the case, we cut the tail:
                sortedArray=sortedArray[0:10]
            }
            fmt.Println(sortedArray);
        }
    })

    //Each request start fire up this one:
    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL) //check visiting links
    })

    //And this method start up everything:
    c.Visit("https://www.fundamentus.com.br/detalhes.php")

    //As we're operating in synchronous mode(but it may be assync too),
    //the code after "Visit" call occurs after EVERYTHING finishes.

    //Now we already have the top ten, so we can commit them to the DBMS:
        fmt.Println("Finished")//, r.Request.URL)
        db, err := sql.Open("mysql", "root:supersecret@tcp(127.0.0.1:3306)/redventures")
        if err != nil {
            panic(err.Error())
        }
        stmt, err := db.Prepare("INSERT stocks SET papel=?,empresa=?,varDia=?,valor=?")
        if err != nil {
            panic(err.Error())
        }
        for _,source:=range sortedArray{
            stmt.Exec(source.papel,source.empresa,source.varDia,source.valor)
            if err != nil {
                panic(err.Error())
            }
        }
        db.Close()

    fmt.Println(sortedArray);
}
