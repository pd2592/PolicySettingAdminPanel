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
			//fmt.Println(Form, "   ", Columns)

			if checkDuplicate(table, methType, companyID, Form, Columns) != "0" {
				result = checkDuplicate(table, methType, companyID, Form, Columns)
			} else {
				fmt.Println("one check!!! ")
				//fmt.Println(Form, "  ", Columns)
				if QueryRow(Form, Columns, table) != 0 {
					result = "Record already exists"

				} else {
					RowId := AddRow(Form, Columns, table)
					result = strconv.FormatInt(RowId, 10)
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
			result1 := checkDuplicate(table, methType, companyID, FormVal, ColumnVal)
			if result1 != "0" {
				fmt.Println("Edit")
				result = result1
			} else {
				if QueryRow(FormCondVal, ColumnCondVal, table) == 0 {
					result = " Record(s) Not Exists"

				} else {
					RowCount := EditRow(FormVal, ColumnVal, FormCondVal, ColumnCondVal, table)
					result = strconv.FormatInt(RowCount, 10) + " Record(s) Updated"

				}
			}
		}
	case "list": //input should be companyID and tablename
		{
			c := true
			switch c {
			case table == "city_mapping":
				{
					result = ListCity(Form[0])
				}
			case table == "citymaster":
				{
					result = ListAllCity()
				}
			case table == "city_category":
				{
					result = ListCityCat(companyID)
				}
			case table == "benefit_type_master" && ColumnCondVal == nil:
				{
					result = ListBundleRequirements(table, companyID)
				}
			case table == "policy_benefit_bundle" && ColumnCondVal == nil:
				{
					result = ListBundles(table, companyID)
				}
			case table == "policy_benefit_bundle" && companyID == "":
				{
					result = ListBundleDetail(table, FormCondVal, ColumnCondVal)
				}
			case (table == "department" && companyID == ""):
				{
					result = ListDepartmentDetail(table, FormCondVal, ColumnCondVal)
				}
			case table == "department" && ColumnCondVal == nil:
				{
					result = ListDepartments(table, companyID)
				}
			case table == "designationmaster" && companyID == "":
				{
					result = ListDesignaionDetail(table, FormCondVal, ColumnCondVal)
				}
			case table == "designationmaster" && ColumnCondVal == nil:
				{
					result = ListDesignations(table, companyID)
				}
			case table == "designationmaster" && ColumnCondVal != nil && companyID != "":
				{
					result = ListDesignationsByDep(table, companyID, FormCondVal, ColumnCondVal)
				}
			default:
				{
					fmt.Println("dedault")
					result = ListRow(table, companyID)
				}
			}
			//w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		}
	case "delete":
		{
			fmt.Println("Delete Operation")
			result = DeleteById(table, FormCondVal, ColumnCondVal)

		}
	}
	return result
}

func checkDuplicate(table string, methType string, companyID string, Form []string, Columns []string) string {
	var message string
	fmt.Println(Form)
	c := true

	fmt.Println(table == "policy_benefit_bundle" && (QueryRow([]string{Form[1], Form[2]}, []string{Columns[1], Columns[2]}, table) != 0 || QueryRow([]string{Form[0], Form[2]}, []string{Columns[0], Columns[2]}, table) != 0))
	//fmt.Println
	switch c {
	case (table == "policy_benefit_bundle") && (QueryRow([]string{Form[1], Form[2]}, []string{Columns[1], Columns[2]}, table) != 0 || QueryRow([]string{Form[0], Form[2]}, []string{Columns[0], Columns[2]}, table) != 0):
		{
			//fmt.Println("inside else if 1", Form, "  ", Columns)
			message = "Record already exists"

		}
	case (table == "department") && (QueryRow([]string{Form[0], Form[1], Form[2]}, []string{Columns[0], Columns[1], Columns[2]}, table) != 0):
		{
			//fmt.Println("inside else if 2", Form, "  ", Columns)
			message = "Department already saved!!"
		}
	case (table == "department") && (QueryRow([]string{Form[0], Form[2]}, []string{Columns[0], Columns[2]}, table) != 0 && QueryRow([]string{Form[1], Form[2]}, []string{Columns[1], Columns[2]}, table) != 0):
		{
			//fmt.Println("inside else if 3", Form, "  ", Columns)
			message = "Department Exists!!"
		}
	case (table == "designationmaster") && methType == "edit" && (QueryRowEdit([]string{Form[0], Form[1], Form[4]}, []string{Columns[0], Columns[1], Columns[4]}, Form[7], table) != 0):
		{
			//fmt.Println("inside else if 3", Form, "  ", Columns)
			message = "Designation already saved!!"
		}
	case (table == "designationmaster") && methType == "edit" && QueryRowEdit([]string{Form[0], Form[4]}, []string{Columns[0], Columns[4]}, Form[7], table) != 0:
		{
			//fmt.Println("inside else if 4", Form, "  ", Columns)
			message = "Designation exists!!"
		}
	case (table == "designationmaster") && methType == "edit" && QueryRowEdit([]string{Form[1], Form[4]}, []string{Columns[1], Columns[4]}, Form[7], table) != 0:
		{
			//fmt.Println("inside else if 5", Form, "  ", Columns)
			message = "Designation Code exists!!"
		}
	case (table == "designationmaster") && methType == "create" && (QueryRow([]string{Form[0], Form[1], Form[4]}, []string{Columns[0], Columns[1], Columns[4]}, table) != 0):
		{
			//fmt.Println("inside else if 6", Form, "  ", Columns)
			message = "Designation already saved!!"
		}
	case (table == "designationmaster") && methType == "create" && QueryRow([]string{Form[0], Form[4]}, []string{Columns[0], Columns[4]}, table) != 0:
		{
			//fmt.Println("inside else if 7", Form, "  ", Columns)
			message = "Designation exists!!"
		}
	case (table == "designationmaster") && methType == "create" && QueryRow([]string{Form[1], Form[4]}, []string{Columns[1], Columns[4]}, table) != 0:
		{
			//fmt.Println("inside else if 8", Form, "  ", Columns)
			message = "Designation Code exists!!"
		}
	default:
		{
			message = "0"
		}
	}
	return message
}

// func RemoveIndex(s []string, index int) []string {
// 	return append(s[:index], s[index+1:]...)
// }
