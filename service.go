package main

import (
	"fmt"
	"strconv"
)

func serv(table string, methType string, companyID string, Form []string, Columns []string, FormVal []string, ColumnVal []string, FormCondVal []string, ColumnCondVal []string) string {
	var result string

	switch methType {
	case "create": //input should be column values with name in form. It will add record if it is not already exists
		{
			//FormVal := []string{r.FormValue("cityCatName"), r.FormValue("companyId")}
			//Columns := []string{m["cityCatName"], m["companyId"]}
			fmt.Println("Inside create")
			fmt.Println(Form, "   ", Columns)
			//Form1 := []string{Form[1], Form[2]}
			//Columns1 := []string{Columns[1], Columns[2]}
			//Form2 := []string{Form[0], Form[2]}
			//Columns2 := []string{Columns[0], Columns[2]}
			if (table == "city_mapping") && CityMapCheck(Form[1], companyID) != "" {

				result = "Record already exists in " + CityMapCheck(Form[1], companyID) + " category"
			} else if (table == "policy_benefit_bundle") && (QueryRow([]string{Form[1], Form[2]}, []string{Columns[1], Columns[2]}, table) != 0 || QueryRow([]string{Form[0], Form[2]}, []string{Columns[0], Columns[2]}, table) != 0) {
				fmt.Println("inside else if ", Form, "  ", Columns)
				//	fmt.Println("inside else if ", Form1, "  ", Columns1)

				// fmt.Println("here table is :", table)
				// if QueryRow(Form, Columns, table) != 0 {
				// 	fmt.Println("next next check!!! ")
				// 	Form1 := RemoveIndex(Form, 0)
				// 	Columns1 := RemoveIndex(Columns, 0)
				// 	if QueryRow(Form1, Columns1, table) != 0 {
				result = "Record already exists"
				// } else {
				// 	RowCount := AddRow(Form, Columns, table)
				// 	result = strconv.FormatInt(RowCount, 10)
				// 	break
				// }
				//} // else {
				// 	RowCount := AddRow(Form, Columns, table)
				// 	result = strconv.FormatInt(RowCount, 10)
				// }
			} else {
				fmt.Println("one check!!! ")
				fmt.Println(Form, "  ", Columns)
				if QueryRow(Form, Columns, table) != 0 {
					fmt.Println("here table issss :", table)
					result = "Record already exists"

				} else {
					RowCount := AddRow(Form, Columns, table)
					result = strconv.FormatInt(RowCount, 10)
				}
				//result = "Record already exists"
			}

		}
	case "edit":
		{
			// FormVal := []string{r.FormValue("cityCatName")}
			// ColumnVal := []string{m["cityCatName"]}
			// FormCondVal := []string{r.FormValue("cityCatId")}
			// ColumnCondVal := []string{m["cityCatId"]}
			if QueryRow(FormCondVal, ColumnCondVal, table) == 0 {
				result = " Record(s) Not Exists"

			} else {
				RowCount := EditRow(FormVal, ColumnVal, FormCondVal, ColumnCondVal, table)
				result = strconv.FormatInt(RowCount, 10) + " Record(s) Updated"

			}
		}
	case "list": //input should be companyID and tablename
		{
			if table == "city_mapping" {
				result = ListCity(Form[0])
			} else if table == "policy_benefit_bundle" || table == "benefit_bundle_type_mapping" || table == "bundle_type_benefit_mapping" || table == "benefit_type_allowance_mapping" {
				result = ListBundle(table, FormCondVal, ColumnCondVal)
			} else {
				result = ListRow(table, companyID)
			}
			//w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		}
	case "delete":
		{
			fmt.Println("Delete Operation")
			DeleteById(table, FormCondVal, ColumnCondVal)
		}
	}
	return result
}

// func RemoveIndex(s []string, index int) []string {
// 	return append(s[:index], s[index+1:]...)
// }
