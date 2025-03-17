package main

import (
	"fmt"
	"strings"
	"time"

	//"os"

	"flag"

	"github.com/tebeka/selenium"
)

const (
	chromeDriverPath = "./chromedriver"
	port             = 8080
)

func numberOfWeekInMonth() int {
	now := time.Now()
	_, w1 := time.Now().UTC().ISOWeek()
	_, w2 := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).UTC().ISOWeek()
	return w1 - w2 + 1
}

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

	userURL := fmt.Sprintf("https://%d:%s@intra.concords.com.tw/Site2/main/EIP_User.aspx", concordID, concordPW)
	wd.Get(userURL)
	name, err := wd.FindElement(selenium.ByCSSSelector, "#Label_name") //使用者姓名
	if err != nil {
		panic(err)
	}
	nameStr, _ := name.Text()
	nameStr = nameStr[1:10]
	fmt.Println(nameStr)

	homeURL := fmt.Sprintf("https://%d:%s@intra.concords.com.tw/Site2/MIS/Flow/DateTimeFlow.aspx", concordID, concordPW)
	wd.Get(homeURL)
	time.Sleep(1 * time.Second)
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

	dateStr, _ := date.Text()
	fmt.Println(dateStr)
	time1Str, _ := time1.Text()
	fmt.Println(time1Str)
	time2Str, _ := time2.Text()
	fmt.Println(time2Str)

	hrURL := fmt.Sprintf("https://%d:%s@intra.concords.com.tw/Site2/mis/hr/OverTimeListNew.aspx", concordID, concordPW)
	wd.Get(hrURL)
	time.Sleep(1 * time.Second)
	css := "#Calendar1 > tbody >"

	week := numberOfWeekInMonth()
	fmt.Println("第", week, "週")
	//看看第幾周, 去組合CSSSelector
	/*
		tr:nth-child(3) 第1周
		tr:nth-child(4) 第2周
		tr:nth-child(5) 第3周
		tr:nth-child(6) 第4周
		tr:nth-child(7) 第5周
	*/
	if week == 1 {
		css += " tr:nth-child(3)"
	} else if week == 2 {
		css += " tr:nth-child(4)"
	} else if week == 3 {
		css += " tr:nth-child(5)"
	} else if week == 4 {
		css += " tr:nth-child(6)"
	} else if week == 5 {
		css += " tr:nth-child(7)"
	}
	//看看星期幾, 去組合CSSSelector
	/*
		td:nth-child(1) //星期天
		td:nth-child(2) //星期一
		td:nth-child(3)
		td:nth-child(4)
		td:nth-child(5)
		td:nth-child(6)
		td:nth-child(7) //星期六
	*/
	if strings.Contains(dateStr, "一") {
		css += " td:nth-child(2)"
	} else if strings.Contains(dateStr, "二") {
		css += " td:nth-child(3)"
	} else if strings.Contains(dateStr, "三") {
		css += " td:nth-child(4)"
	} else if strings.Contains(dateStr, "四") {
		css += " td:nth-child(5)"
	} else if strings.Contains(dateStr, "五") {
		css += " td:nth-child(6)"
	} else if strings.Contains(dateStr, "六") {
		css += " td:nth-child(7)"
	} else {
		css += " td:nth-child(1)"
	}

	people, err := wd.FindElement(selenium.ByCSSSelector, css) //今天請假人數
	if err != nil {
		panic(err)
	}
	PeopleStr, _ := people.Text()
	fmt.Println("請假名單")
	fmt.Println("日期:", PeopleStr)

	workflag := true                          //有沒有上班
	if strings.Contains(PeopleStr, nameStr) { //有可能是請假或是代理人
		//判斷代理人
		b := "(" + nameStr + " 代)"
		if strings.Contains(PeopleStr, b) {
			workflag = true
		} else {
			workflag = false
		}
	}
	fmt.Println()

	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()

	//早上忘了打卡 幫主人打卡
	if workflag && time1Str == " " && hour == 8 && minute <= 10 {
		fmt.Println("小精靈出現了")
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

	//下班了 幫主人先打卡
	if workflag && hour == 17 {
		fmt.Println("小精靈出現了")
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
}
