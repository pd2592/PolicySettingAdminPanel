package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var jsonString []byte

func AddRow(values []string, columns []string, table string) int64 {
	db = GetDB()
	// insert
	column := strings.Join(columns, ",")
	value := strings.Join(values, "','")
	stmt, err := db.Exec("INSERT into " + table + "(" + column + ") VALUES ('" + value + "')")
	checkErr(err)

	rowid, err := stmt.LastInsertId()
	checkErr(err)

	fmt.Println("Value Inserted into ", table)
	fmt.Println("RowId ", rowid)

	return rowid

}

func EditRow(formVal []string, columnVal []string, formCondVal []string, columnCondVal []string, table string) int64 {
	db = GetDB()
	//column := strings.Join(columns, ",")
	//value := strings.Join(values, "','")
	setStr := "SET"
	fmt.Println(len(columnVal))
	i := 0
	for i = 0; i < len(columnVal); i++ {
		setStr = setStr + " , " + columnVal[i] + " = '" + formVal[i] + "'"

	}
	setStr = strings.Replace(setStr, "SET ,", " SET", -2)

	conditionStr := createCondStr(formCondVal, columnCondVal)
	stmt, err := db.Exec("UPDATE " + table + setStr + conditionStr)

	checkErr(err)

	rowcnt, err := stmt.RowsAffected()
	checkErr(err)

	fmt.Println("Row(s) Updated ", rowcnt)

	return rowcnt

}

func QueryRow(values []string, columns []string, table string) int {

	db = GetDB()

	conditionStr := createCondStr(values, columns)

	// select
	var cnt int
	fmt.Println(">>>>>>>>>>>>", conditionStr)
	_ = db.QueryRow("select count(*) from " + table + conditionStr).Scan(&cnt)
	//fmt.Println(">>>>", cnt)
	fmt.Println(cnt)
	return cnt
}

func ListCity(citycatId string) string {
	db = GetDB()
	//var label string
	stmt, err := db.Query("SELECT City as label1, State as label2 from city_master where CityID IN (SELECT CityID from city_mapping WHERE CityCatID = '" + citycatId + "')")
	checkErr(err)

	return ParseRow(stmt)
}
func ListBundle(tablename string, formCondVal []string, columnCondVal []string) string {
	db = GetDB()
	conditionStr := createCondStr(formCondVal, columnCondVal)

	stmt, err := db.Query("select * from " + tablename + conditionStr)
	checkErr(err)
	return ParseRow(stmt)
}

func ListRow(tablename string, companyId string) string {
	db = GetDB()

	stmt, err := db.Query("select * from " + tablename + " where CompanyID = " + companyId)
	checkErr(err)
	return ParseRow(stmt)
}

func ParseRow(stmt *sql.Rows) string {

	columnNames, err := stmt.Columns()
	checkErr(err)
	rc := NewMapStringScan(columnNames)
	var slice = "["
	for stmt.Next() {
		err := rc.Update(stmt)
		checkErr(err)
		cv := rc.Get()
		//log.Printf("%#v\n\n >>>>>>", cv[columnNames[0]])
		jsonString, _ = json.Marshal(cv)

		slice += string(jsonString) + ","
		//fmt.Println(string(slice))

	}
	slice += "]"
	r := strings.NewReplacer(",]", "]")
	slice = r.Replace(slice)
	return slice
}

func CityMapCheck(cityId, companyId string) string {
	db = GetDB()
	// insert
	var citycat = ""
	_ = db.QueryRow("select CityCatID from city_mapping where cityID = '" + cityId + "' AND CityCatID IN (SELECT CityCatID from city_category where companyID = '" + companyId + "') ").Scan(&citycat)

	_ = db.QueryRow("select CityCatName from city_category where CityCatID = '" + citycat + "'").Scan(&citycat)

	fmt.Println("/////", citycat)
	//rowcnt, err := stmt.RowsAffected()
	//checkErr(err)
	//fmt.Println(rowcnt)
	return citycat
}

func getBenefitBundleID(benefitBundleCode string, companyId string) string {
	db = GetDB()
	var benefitBundleId string
	_ = db.QueryRow("select BenefitBundleID from policy_benefit_bundle where BenefitBundleCode = '" + benefitBundleCode + "' and CompanyID = '" + companyId + "'").Scan(&benefitBundleId)
	return benefitBundleId
}

func getId(table, requestId string, formCondVal []string, columnCondVal []string) string {
	db = GetDB()
	conditionStr := createCondStr(formCondVal, columnCondVal)

	var Id string
	_ = db.QueryRow("select " + requestId + " from " + table + conditionStr).Scan(&Id)
	return Id
}

func getBenefitBundleTypeMappingID(benefitBundleType string, benefitBundleID string) string {
	db = GetDB()
	var benefitBundleTypeMappingId string
	_ = db.QueryRow("select BenefitBundleTypeMappingID from benefit_bundle_type_mapping where BenefitTypeID = '" + benefitBundleType + "' and BenefitBundleID = '" + benefitBundleID + "'").Scan(&benefitBundleTypeMappingId)
	return benefitBundleTypeMappingId
}

// func GetDependency(tablename string) string {
// 	fmt.Println("I am inside get dependency")
// 	stmt, err := db.Query("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE REFERENCED_TABLE_SCHEMA = 'company_policy' AND REFERENCED_TABLE_NAME = '" + tablename + "'")
// 	checkErr(err)
// 	fmt.Println("ran query")

// 	var tables string
// 	fmt.Println("all set! go!")

// 	for stmt.Next() {
// 		var name string

// 		if err := stmt.Scan(&name); err != nil {
// 			log.Fatal(err)
// 		}
// 		tables = tables + "," + name
// 	}
// 	fmt.Printf("......", tables)
// 	return tables
// }

func DeleteById(tablename string, formCondVal []string, columnCondVal []string) {

	db = GetDB()
	_, err := db.Exec("delete from " + tablename + createCondStr(formCondVal, columnCondVal))
	checkErr(err)

	//res, err := stmt.Exec(formCondVal[0])
	checkErr(err)

	//affect, err := res.RowsAffected()
	checkErr(err)

	//fmt.Println(affect)

}

func createCondStr(formCondVal []string, columnCondVal []string) string {
	conditionStr := " WHERE"
	i := 0
	for i = 0; i < len(columnCondVal); i++ {
		conditionStr = conditionStr + " and " + columnCondVal[i] + " = '" + formCondVal[i] + "'"
	}
	conditionStr = strings.Replace(conditionStr, "WHERE and", " WHERE", -2)
	return conditionStr
}

func GetDB() *sql.DB {
	var err error
	if db == nil {
		db, err = sql.Open("mysql", "root:@/company_policy?charset=utf8")
		checkErr(err)
	}

	return db
}
