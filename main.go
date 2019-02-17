package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
)

var symbol string
var header string

var monthsOfYear []Month

type Month struct {
	name  string
	space int
}

var sunday, monday, tuesday, wednesday, thursday, friday, saturday []string

func main() {
	app := cli.NewApp()
	app.Name = "gengrass"
	app.Usage = "print github contributions"
	app.Version = "1.0.0"

	app.Action = func(c *cli.Context) {
		command(c)
	}

	app.Run(os.Args)
}

func command(c *cli.Context) {
	if c.NArg() == 0 {
		return
	} else if c.NArg() == 2 {
		setSymbol(c.Args()[1])
	} else {
		symbol = " ■"
	}
	username := c.Args()[0]
	user := fmt.Sprintf("https://github.com/users/%s/contributions", username)

	res, err := http.Get(user)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		fmt.Printf("%s is not found", username)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}
	header = doc.Find("h2[class='f4 text-normal mb-2']").Text()

	lastMonthP := 0
	doc.Find(".month").Each(func(i int, s *goquery.Selection) {
		x, _ := s.Attr("x")
		p, _ := strconv.Atoi(x)

		var space int
		if i == 0 {
			if p/20 >= 1 {
				space = 2
			} else {
				space = 0
			}
		} else {
			space = (p-lastMonthP)/10 + 1
			if space == 3 {
				space = 1
			} else if space == 5 {
				space = 5
			} else {
				space = 7
			}
		}
		monthsOfYear = append(monthsOfYear, Month{name: s.Text(), space: space})
		lastMonthP = p
	})

	doc.Find(".js-calendar-graph-svg g g").Each(func(_ int, s *goquery.Selection) {
		sun, _ := s.Find("rect").Eq(0).Attr("fill")
		sunday = append(sunday, sun)
		mon, _ := s.Find("rect").Eq(1).Attr("fill")
		monday = append(monday, mon)
		tue, _ := s.Find("rect").Eq(2).Attr("fill")
		tuesday = append(tuesday, tue)
		wed, _ := s.Find("rect").Eq(3).Attr("fill")
		wednesday = append(wednesday, wed)
		thu, _ := s.Find("rect").Eq(4).Attr("fill")
		thursday = append(thursday, thu)
		fri, _ := s.Find("rect").Eq(5).Attr("fill")
		friday = append(friday, fri)
		sat, _ := s.Find("rect").Eq(6).Attr("fill")
		saturday = append(saturday, sat)
	})
	execute()
}

func setSymbol(input string) {
	s := strings.TrimSpace(input)
	if len(s) != 0 {
		symbol = " " + s[:1]
	} else {
		symbol = " ■"
	}
}

func execute() {
	fmt.Print("\n" + strings.Repeat("=", 120) + "\n")
	fmt.Println(header)
	printMonth()
	printContributions("Sun", sunday)
	printContributions("Mon", monday)
	printContributions("Tue", tuesday)
	printContributions("Wed", wednesday)
	printContributions("Thu", thursday)
	printContributions("Fri", friday)
	printContributions("Sat", saturday)
	fmt.Print("\n" + strings.Repeat("=", 120) + "\n\n")
}

func printMonth() {
	fmt.Print(strings.Repeat(" ", 10))
	for _, m := range monthsOfYear {
		var pos string
		for i := 0; i < m.space; i++ {
			pos += " "
		}
		fmt.Print(pos + m.name)
	}
	fmt.Println()
}

func printContributions(dayOfWeek string, array []string) {
	fmt.Print(strings.Repeat(" ", 6) + dayOfWeek)
	for _, val := range array {
		var str string
		switch val {
		case "#196127":
			str += "\x1b[31m" + symbol + "\x1b[0m"
		case "#239a3b":
			str += "\x1b[35m" + symbol + "\x1b[0m"
		case "#7bc96f":
			str += "\x1b[36m" + symbol + "\x1b[0m"
		case "#c6e48b":
			str += "\x1b[33m" + symbol + "\x1b[0m"
		default:
			str += symbol
		}
		fmt.Print(str)
	}

	fmt.Print("\n")
}
