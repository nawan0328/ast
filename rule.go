package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func getRule() {

	var code, name string
	//stock = Stock{} //Initialization

	rows, err := db.Query(select_code_rule)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&code, &name)
		if err != nil {
			log.Fatal(err)
		}
		s, ned := getData(code)
		if ned != nil { //not enough data
			continue
		}
		rule1Check(code, name, s)
		rule2Check(code, name, s)
		rule3Check(code, name, s)
		rule4Check(code, name, s)
	}
}

func getData(code string) (s []DatasType, ned error) {

	s = make([]DatasType, 10)

	rows, err := db.Query(select_rule1, stock.Date, code)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for i := 0; rows.Next(); i++ {
		err := rows.Scan(&s[i].Sv, &s[i].Cv, &s[i].Hv, &s[i].Lv, &s[i].Aq)
		if err != nil {
			log.Fatal(err)
		}
	}

	if s[4].Sv == 0 && s[4].Hv == 0 && s[4].Lv == 0 && s[4].Aq == 0 {
		ned = fmt.Errorf("not enough data")
		return
	}

	ned = nil //not enough data

	return
}

func rule1Check(code string, name string, s []DatasType) {

	if s[0].Aq < 500000 || s[1].Aq < 500000 {
		return //check quantity
	}
	per := (s[1].Cv - s[1].Sv) * 100 / s[1].Cv
	if per < 5 {
		return // Candle is less than 5 percent
	}
	if s[1].Cv > s[0].Cv || s[1].Cv > s[0].Sv {
		return // check gap increase
	}
	//if s[1].Cv > s[0].Lv {
	//	return //gap touch
	//}

	var gap, high_limit int

	if s[0].Cv < s[0].Sv { // DOWN
		gap = s[0].Cv - s[1].Cv
		high_limit = s[0].Cv
	} else { // UP or SAME
		gap = s[0].Sv - s[1].Cv
		high_limit = s[0].Sv
	}

	if (gap * 1000 / s[0].Cv) < 18 { // Gap is less than 1.8 percent
		return
	}

	avg := (s[4].Cv + s[3].Cv + s[2].Cv + s[1].Cv + s[0].Cv) / 5
	if avg > s[1].Cv {
		fmt.Println("5 DAY AVG DOWN EXECPTED : ", code, name, avg, s[0].Cv)
		return
	}

	rule = Rule{code, name, s[1].Cv, high_limit}
	stock.Rule1 = append(stock.Rule1, rule)
}

func rule2Check(code string, name string, s []DatasType) {

	//if s[0].Aq < 500000 || s[1].Aq < 500000 {
	//	return //check quantity
	//}
	if s[2].Aq < 500000 {
		return
	}
	if s[0].Cv < 1000 {
		//fmt.Println(" LESS THAN 1000 EXECPTED : ", code, name)
		return
	}
	if s[2].Cv >= s[2].Sv || s[1].Cv >= s[1].Sv || s[0].Cv >= s[0].Sv {
		return //triple down check
	}
	if s[2].Aq < s[1].Aq || s[1].Aq < s[0].Aq {
		return //count check
	}
	if s[2].Cv <= s[1].Cv || s[1].Cv <= s[0].Cv {
		return //close value chack
	}
	if s[2].Lv <= s[1].Lv || s[1].Lv <= s[0].Lv {
		return //low value check
	}
	if (s[2].Sv-s[2].Cv) < (s[1].Sv-s[1].Cv) || (s[2].Sv-s[2].Cv) < (s[0].Sv-s[0].Cv) {
		fmt.Println("s[2] NO lAGEST EXECPTED : ", code, name)
		return //s[2] lagest check
	}
	per := (s[2].Sv - s[2].Cv) * 1000 / s[2].Sv
	if per < 30 || per >= 130 {
		fmt.Println("LOWER THAN 3% or HIGHER THAN 13% EXECPTED : ", code, name, per)
		return // Candle is less than 5 percent
	}
	if s[3].Sv >= s[3].Cv && s[4].Sv >= s[4].Cv {
		fmt.Println("3, 4 DAY DOWN EXECPTED : ", code, name)
		return //3, 4  down
	}
	//fmt.Println(code, name, s[2].Lv, s[1].Lv, s[0].Lv)

	rule = Rule{code, name, 0, s[0].Cv}
	stock.Rule2 = append(stock.Rule2, rule)
}

func rule3Check(code string, name string, s []DatasType) {
	//if s[2].Aq < 500000 && s[1].Aq < 500000 && s[0].Aq < 500000 {
	//	return
	//}
	if s[5].Cv <= s[5].Sv || s[4].Cv <= s[4].Sv || s[3].Cv <= s[3].Sv {
		return
	}
	if s[2].Cv >= s[2].Sv || s[1].Cv >= s[1].Sv || s[0].Cv <= s[0].Sv {
		return
	}
	/*
		if s[0].Sv-s[0].Cv > s[0].Cv-s[0].Lv {
			return //down tail longer than body check
		}
		if s[0].Cv-s[0].Lv < s[1].Cv-s[1].Lv || s[1].Cv-s[1].Lv < s[2].Cv-s[2].Lv {
			return //down tail check
		}
		if s[0].Hv-s[0].Sv > s[0].Sv-s[0].Cv || s[1].Hv-s[1].Sv > s[1].Sv-s[1].Cv || s[2].Hv-s[2].Sv > s[2].Sv-s[2].Cv {
			return //up tail check
		}
	*/

	/*
		per := (s[4].Cv - s[4].Sv) * 1000 / s[4].Sv
		if per < 50 || per >= 200 {
			return // Candle is less than 5 percent
		}
	*/

	per := (s[4].Cv - s[4].Sv) * 1000 / s[4].Sv
	if per < 50 {
		fmt.Println(code, name, "s[4] candle body is less than 5 percent", per)
		return // Candle is more than 5 percent
	}

	per = (s[5].Cv - s[5].Sv) * 1000 / s[5].Cv
	if per > 10 {
		fmt.Println(code, name, "s[5] candle body is more than 3 percent", per)
		return // Candle is less than 1 percent
	}

	per = (s[0].Cv - s[0].Sv) * 1000 / s[0].Cv
	if per > 36 {
		fmt.Println(code, name, "s[0] candle body is more than 3.6 percent", per)
		return // Candle is less than 3 percent
	}

	if s[1].Cv < s[4].Cv || s[2].Cv < s[4].Cv {
		fmt.Println(code, name, "s[1],s[2] are less than down limit")
		return //down check
	}

	if s[4].Cv >= s[3].Sv {
		fmt.Println(code, name, "no gap s[4],s[3]")
		return //gap check
	}
	if (s[4].Cv+s[4].Sv)/2 > s[2].Lv || (s[4].Cv+s[4].Sv)/2 > s[1].Lv {
		fmt.Println(code, name, "downtail touched s[4]'s half")
		return //downtail check
	}

	per = (s[0].Sv - s[0].Lv) * 1000 / s[0].Cv
	if per > 10 {
		fmt.Println(code, name, "downtail is more than 1 percent", per)
		return // downtail is less than 1 percent
	}

	rule = Rule{code, name, s[0].Cv, 0}
	stock.Rule3 = append(stock.Rule3, rule)
}
func rule4Check(code string, name string, s []DatasType) {
}
