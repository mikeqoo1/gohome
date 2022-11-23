package main

import (
	"fmt"
	"time"

	//"os"

	"flag"

	"github.com/tebeka/selenium"
)

const (
	chromeDriverPath = "./chromedriver"
	port             = 8080
)

func main() {
	var concordID int
	var concordPW string
	flag.IntVar(&concordID, "u", 0, "帳號 默認為0")
	flag.StringVar(&concordPW, "p", "", "密碼 默認為空")
	flag.Parse()
	opts := []selenium.ServiceOption{
		// Enable fake XWindow session.
		// selenium.StartFrameBuffer(),
		//selenium.Output(os.Stderr), // Output debug information to STDERR
	}

	// Enable debug info.
	// selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	homeURL := fmt.Sprintf("https://%d:%s@intra.concords.com.tw/Site2/MIS/Flow/DateTimeFlow.aspx", concordID, concordPW)
	wd.Get(homeURL)
	time.Sleep(3 * time.Second)
	// page, _ := wd.PageSource()
	// fmt.Println(page)
	date, err := wd.FindElement(selenium.ByCSSSelector, "#DateTimeTable > tbody > tr:nth-child(1) > td.tCell_0 > span") //日期
	if err != nil {
		panic(err)
	}
	time1, err := wd.FindElement(selenium.ByCSSSelector, "#DateTimeTable > tbody > tr:nth-child(1) > td.tCell_1 > span") //上班
	if err != nil {
		panic(err)
	}
	time2, err := wd.FindElement(selenium.ByCSSSelector, "#DateTimeTable > tbody > tr:nth-child(1) > td.tCell_2 > span") //下班
	if err != nil {
		panic(err)
	}
	fmt.Println(date.Text())
	fmt.Println(time1.Text())
	fmt.Println(time2.Text())

	//打卡
	URL := fmt.Sprintf("https://%d:%s@intra.concords.com.tw/site2/main/RunCard.aspx", concordID, concordPW)
	wd.Get(URL)
	btn, err := wd.FindElement(selenium.ByID, "btnSelf")
	if err != nil {
		panic(err)
	}
	if err := btn.Click(); err != nil {
		panic(err)
	}
}
