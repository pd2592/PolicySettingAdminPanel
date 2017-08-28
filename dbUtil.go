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
func ListAllCity() string {
	db = GetDB()
	var cities []LabVal
	var city LabVal
	//var label string
	stmt, err := db.Query("SELECT CityId as value, City as label from city_master")
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&city.Value, &city.Label)
		checkErr(err)

		city = LabVal{
			Label: city.Label,
			Value: city.Value,
		}
		cities = append(cities, city)
	}
	b, err := json.Marshal(cities)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
}

func ListCityCat(companyId string) string {
	db = GetDB()

	stmt, err := db.Query("SELECT CityCatName as label, CityCatID as value from city_category where CompanyID = '" + companyId + "'")
	checkErr(err)
	return ParseRow(stmt)
}

func ListCity(citycatId string) string {
	db = GetDB()
	var cities []LabVal
	var labval LabVal
	var cityCatName string
	//var label string
	stmt, err := db.Query("SELECT cmas.City, cmap.CityMappingID from city_mapping as cmap JOIN city_master as cmas ON cmas.CityID = cmap.CityID where cmap.CityCatID = '" + citycatId + "'")
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&labval.Label, &labval.Value)
		checkErr(err)

		labval = LabVal{
			Label: labval.Label,
			Value: labval.Value,
		}
		cities = append(cities, labval)
	}
	err = db.QueryRow("SELECT CityCatName from city_category where CityCatID = '" + citycatId + "'").Scan(&cityCatName)
	checkErr(err)

	labval = LabVal{
		Label: cityCatName,
		Value: citycatId,
	}
	cityCategoryMap := CityCategoryMap{
		CompanyID: "",
		CityCat:   labval,
		Cities:    cities,
	}
	b, err := json.Marshal(cityCategoryMap)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
}
func ListBundleRequirements(tablename, companyId string) string {
	db = GetDB()
	//var labval LabVal
	var benefittype LabVal
	var policybundle PolicyBundle
	var mybundles []PolicyBundle
	var cityCatAndAllowance CityCatAndAllowance

	stmt, err := db.Query("SELECT BenefitTypeID, BenefitTypeName from benefit_type_master")
	checkErr(err)

	for stmt.Next() {
		err := stmt.Scan(&benefittype.Value, &benefittype.Label)
		checkErr(err)
		benefittype = LabVal{
			Label: benefittype.Label,
			Value: benefittype.Value,
		}
		stmt1, err := db.Query("SELECT BenefitID, BenefitName FROM benefit_master WHERE BenefitTypeID='" + benefittype.Value + "'")

		var benefit LabVal
		var cityCatAndAllowances []CityCatAndAllowance
		//var benefit LabVal
		var benefits []LabVal
		for stmt1.Next() {
			err := stmt1.Scan(&benefit.Value, &benefit.Label)
			checkErr(err)
			benefit = LabVal{
				Label: benefit.Label,
				Value: benefit.Value,
			}
			benefits = append(benefits, benefit)
		}
		stmt2, err := db.Query("SELECT CityCatName, CityCatID from city_category where CompanyID = '" + companyId + "'")
		for stmt2.Next() {
			err := stmt2.Scan(&cityCatAndAllowance.Label, &cityCatAndAllowance.Value)
			checkErr(err)
			cityCatAndAllowance = CityCatAndAllowance{
				Label:      cityCatAndAllowance.Label,
				Value:      cityCatAndAllowance.Value,
				LimitSpent: "",
				Min:        "",
				Max:        "",
				Flex:       "",
				FlexAmt:    "",
				StarCat:    "",
			}
			cityCatAndAllowances = append(cityCatAndAllowances, cityCatAndAllowance)
			//benefits = append(benefits, labval)
		}
		policybundle = PolicyBundle{
			BenefitTypeID:        benefittype,
			Priority:             policybundle.Priority,
			Benefits:             benefits,
			CityCatAndAllowances: cityCatAndAllowances,
		}
		mybundles = append(mybundles, policybundle)
	}
	//var pb PB

	var pb = PB{
		BenefitBundleID: "",
		BundleName:      "",
		BundleCode:      "",
		CompanyID:       companyId,
		MethType:        "",
		PolicyBundles:   mybundles,
	}
	b, err := json.Marshal(pb)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
}
func ListBundleDetail(tablename string, formCondVal []string, columnCondVal []string) string {
	db = GetDB()
	//conditionStr := createCondStr(formCondVal, columnCondVal)
	var pb PB
	var policybundle PolicyBundle
	var labval LabVal
	var benefittype LabVal
	var cityCatAndAllowance CityCatAndAllowance
	var mybundles []PolicyBundle
	var benefitBundleTypeMappingId string

	conditionStr := createCondStr(formCondVal, columnCondVal)

	stmt, err := db.Query("select BenefitBundleID, BenefitBundleName, BenefitBundleCode, CompanyID from policy_benefit_bundle " + conditionStr)
	checkErr(err)
	for stmt.Next() {
		err := stmt.Scan(&pb.BenefitBundleID, &pb.BundleName, &pb.BundleCode, &pb.CompanyID)
		checkErr(err)
	}
	stmt1, err := db.Query("SELECT bbtm.BenefitBundleTypeMappingID, bbtm.BenefitTypeID, btm.BenefitTypeName, bbtm.Priority from benefit_bundle_type_mapping as bbtm JOIN benefit_type_master as btm ON bbtm.BenefitTypeID = btm.BenefitTypeID WHERE bbtm.BenefitBundleID = '" + pb.BenefitBundleID + "'")
	checkErr(err)

	for stmt1.Next() {
		err := stmt1.Scan(&benefitBundleTypeMappingId, &labval.Value, &labval.Label, &policybundle.Priority)
		checkErr(err)

		benefittype = LabVal{
			Label: labval.Label,
			Value: labval.Value,
		}
		stmt2, err := db.Query("SELECT btbm.BenefitID, bm.BenefitName FROM bundle_type_benefit_mapping as btbm JOIN benefit_master as bm ON btbm.BenefitID = bm.BenefitID  WHERE BenefitBundleTypeMappingID='" + benefitBundleTypeMappingId + "'")

		var cityCatAndAllowances []CityCatAndAllowance
		var benefits []LabVal
		for stmt2.Next() {
			err := stmt2.Scan(&labval.Value, &labval.Label)
			checkErr(err)
			labval = LabVal{
				Label: labval.Label,
				Value: labval.Value,
			}
			benefits = append(benefits, labval)
		}
		stmt3, err := db.Query("SELECT ct.CityCatName, btam.CityCatID, btam.LimitSpend, btam.MaxAmount, btam.MinAmount, btam.Flexibility, btam.FlexAmount, btam.StarCat FROM benefit_type_allowance_mapping as btam JOIN city_category as ct ON btam.CityCatID = ct.CityCatID WHERE BenefitBundleTypeMappingID = '" + benefitBundleTypeMappingId + "'")
		for stmt3.Next() {
			err := stmt3.Scan(&cityCatAndAllowance.Label, &cityCatAndAllowance.Value, &cityCatAndAllowance.LimitSpent, &cityCatAndAllowance.Max, &cityCatAndAllowance.Min, &cityCatAndAllowance.Flex, &cityCatAndAllowance.FlexAmt, &cityCatAndAllowance.StarCat)
			checkErr(err)
			cityCatAndAllowance = CityCatAndAllowance{
				Label:      cityCatAndAllowance.Label,
				Value:      cityCatAndAllowance.Value,
				LimitSpent: cityCatAndAllowance.LimitSpent,
				Min:        cityCatAndAllowance.Min,
				Max:        cityCatAndAllowance.Max,
				Flex:       cityCatAndAllowance.Flex,
				FlexAmt:    cityCatAndAllowance.FlexAmt,
				StarCat:    cityCatAndAllowance.StarCat,
			}
			cityCatAndAllowances = append(cityCatAndAllowances, cityCatAndAllowance)
			//benefits = append(benefits, labval)
		}

		policybundle = PolicyBundle{
			BenefitTypeID:        benefittype,
			Priority:             policybundle.Priority,
			Benefits:             benefits,
			CityCatAndAllowances: cityCatAndAllowances,
		}
		mybundles = append(mybundles, policybundle)

	}

	pb = PB{
		BenefitBundleID: pb.BenefitBundleID,
		BundleName:      pb.BundleName,
		BundleCode:      pb.BundleCode,
		CompanyID:       pb.CompanyID,
		MethType:        "",
		PolicyBundles:   mybundles,
	}
	b, err := json.Marshal(pb)
	checkErr(err)
	//fmt.Println(string(b))
	return string(b)
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

// func getBenefitBundleID(benefitBundleCode string, companyId string) string {
// 	db = GetDB()
// 	var benefitBundleId string
// 	_ = db.QueryRow("select BenefitBundleID from policy_benefit_bundle where BenefitBundleCode = '" + benefitBundleCode + "' and CompanyID = '" + companyId + "'").Scan(&benefitBundleId)
// 	return benefitBundleId
// }
func ListBundles(table, companyId string) string {
	db = GetDB()
	var bundle LabVal
	var bundlelist []LabVal
	stmt, err := db.Query("select BenefitBundleName as label, BenefitBundleID as value from " + table + " where CompanyID = '" + companyId + "'")
	checkErr(err)

	for stmt.Next() {
		err := stmt.Scan(&bundle.Label, &bundle.Value)
		checkErr(err)
		bundle = LabVal{
			Label: bundle.Label,
			Value: bundle.Value,
		}
		bundlelist = append(bundlelist, bundle)
	}
	b, err := json.Marshal(bundlelist)
	return string(b)
}

func getId(table, requestId string, formCondVal []string, columnCondVal []string) string {
	db = GetDB()
	conditionStr := createCondStr(formCondVal, columnCondVal)

	var Id string
	_ = db.QueryRow("select " + requestId + " from " + table + conditionStr).Scan(&Id)
	return Id
}

// func getBenefitBundleTypeMappingID(benefitBundleType string, benefitBundleID string) string {
// 	db = GetDB()
// 	var benefitBundleTypeMappingId string
// 	_ = db.QueryRow("select BenefitBundleTypeMappingID from benefit_bundle_type_mapping where BenefitTypeID = '" + benefitBundleType + "' and BenefitBundleID = '" + benefitBundleID + "'").Scan(&benefitBundleTypeMappingId)
// 	return benefitBundleTypeMappingId
// }

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
