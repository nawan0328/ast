package main

import "database/sql"

var (
	db    *sql.DB
	stmt  *sql.Stmt
	stock Stock
	rule  Rule
)

const (
	select_code      = "SELECT code FROM code_info where status=0"
	insert_price     = "INSERT INTO price_history VALUES(?,?,?,?,?,?,?)"
	select_code_rule = "SELECT code, name FROM code_info WHERE market_cap > 60000 AND status=0"
	select_rule1     = "SELECT start, close, high, low, quantity FROM price_history WHERE date <= ? and code = ? order by date desc limit 10"
)

type Stock struct {
	Date  string `json:"date"`
	Rule1 []Rule `json:"rule1"`
	Rule2 []Rule `json:"rule2"`
	Rule3 []Rule `json:"rule3"`
	Rule4 []Rule `json:"rule4"`
}

type Rule struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Start_limit int    `json:"start_limit"`
	High_limit  int    `json:"high_limit"`
}

type JsonObject struct {
	//	ResultCode string
	Result ResultType
}
type ResultType struct {
	//	Time            int
	Areas []AreasType
	//	PollingInterval int
}
type AreasType struct {
	Datas []DatasType
	//	Name  string
}

type DatasType struct {
	//Ms string  `json:"ms"` //Market status
	//Nm string `json:"nm"` //Stock Name
	//Cd string  `json:"cd"` //Stock Code
	Sv int `json:"ov"` //Opening Value
	Cv int `json:"nv"` //Close Value
	Hv int `json:"hv"` //High Value
	Lv int `json:"lv"` //Low Value
	//Rf int `json:"rf"` //rise fall check 1=high_limit, 2=rise, 3=nothinh, 4=low_limit, 5=fall
	//Yv int     `json:"pcv"` //Yesterday Value
	//Dv int     `json:"cv"` //Diff Value
	//Dr float32 `json:"cr"` //Diff Rate
	Aq int `json:"aq"` //Quantity trade volume
}
