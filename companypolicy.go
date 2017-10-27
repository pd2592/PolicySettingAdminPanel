package main

import (
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	gomail "gopkg.in/gomail.v2"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/citycategory", CityCategory)
	//router.HandleFunc("/citymapping", CityMapping)
	router.HandleFunc("/listCity", ListCities)
	router.HandleFunc("/cityCatMap", CityCatMap)
	router.HandleFunc("/citymapping/list", ListCityCatCities)
	//router.HandleFunc("/citymapping/add", AddCities)
	//router.HandleFunc("/citymapping/delete", RemoveCities)
	router.HandleFunc("/cityCatMap/edit", EditCityCatMap)

	router.HandleFunc("/benefitbundle", BenefitBundle)
	router.HandleFunc("/listBenefitBundle", ListBenefitBundle)
	router.HandleFunc("/getBenefitBundle", GetBenefitBundle)
	router.HandleFunc("/listBundleRequirement", ListBundleRequirement)
	router.HandleFunc("/benefitbundle/delete", DeleteBundle)

	router.HandleFunc("/department/create", CreateDepartment)
	router.HandleFunc("/department/edit", EditDepartment)
	router.HandleFunc("/department/get", GetDepartment)
	router.HandleFunc("/department/list", ListDepartment)
	router.HandleFunc("/department/delete", DeleteDepartment)

	router.HandleFunc("/designation/create", CreateDesignation)
	router.HandleFunc("/designation/edit", EditDesignation)
	router.HandleFunc("/designation/get", GetDesignation)
	router.HandleFunc("/designation/list", ListDesignation)
	router.HandleFunc("/designation/dep/list", ListDesignationByDep)
	router.HandleFunc("/designation/delete", DeleteDesignation)
	router.HandleFunc("/bundle/des", BenefitBundleByDes)

	router.HandleFunc("/designation/assign", AssignDesignation)
	router.HandleFunc("/uploadEmp", UploadEmployees)
	router.HandleFunc("/empassign/list", ListUnassignedEmp)
	router.HandleFunc("/emp/add", AddEmployee)
	router.HandleFunc("/emp/get", GetEmployee)
	router.HandleFunc("/emp/all", ListAllEmp)
	router.HandleFunc("/emp/edit", EditEmployee)

	router.HandleFunc("/emp/deact", DeactivateEmployee)
	router.HandleFunc("/emp/search", SearchByName)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func CreateDesignation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	Designation := UnmarshalDesignation(string(body))
	if checkIfExist([]string{Designation.TravelAgencyMasterID, Designation.Department}, []string{"travelAgencyMasterId", "departmentMasterId"}, "department") && (Designation.BenefitBundleID == "" || checkIfExist([]string{Designation.TravelAgencyMasterID, Designation.BenefitBundleID}, []string{"CompanyID", "BenefitBundleID"}, "policy_benefit_bundle")) {
		createTime := time.Now()
		updateTime := time.Now()
		var bundleId string
		if Designation.BenefitBundleID == "" {
			bundleId = "0"
		} else {
			bundleId = Designation.BenefitBundleID
		}
		Form := []string{Designation.DesignationName, Designation.DesignationCode, Designation.HierarchyID, Designation.Department, Designation.TravelAgencyMasterID, bundleId, createTime.String(), updateTime.String()}
		Columns := []string{"designationName", "designationCode", "hierarchyId", "department", "travelAgencyMasterId", "benefitBundleId", "createdDate", "updatedDate"}

		Result := serv("designationmaster", "create", "", Form, Columns, nil, nil, nil, nil)
		if Result == "Record already exists" {
			fmt.Fprintln(w, "Designation already exists")
		} else if len(Result) <= 5 {
			fmt.Fprintln(w, "Changes Saved successfully")

		} else {
			fmt.Fprintln(w, Result)
		}
	} else {
		fmt.Println("Invalid Company(travelAgencyMaster)Id or Department or BundleID")

	}
}

func EditDesignation(w http.ResponseWriter, r *http.Request) {
	var Result string
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	Designation := UnmarshalDesignation(string(body))

	if checkIfExist([]string{Designation.DesignationID}, []string{"designationMasterId"}, "designationmaster") {
		updateTime := time.Now()
		FormVal := []string{Designation.DesignationName, Designation.DesignationCode, Designation.HierarchyID, Designation.Department, Designation.TravelAgencyMasterID, Designation.BenefitBundleID, updateTime.String(), Designation.DesignationID}
		ColumnsVal := []string{"designationName", "designationCode", "hierarchyId", "department", "travelAgencyMasterId", "benefitBundleId", "updatedDate", "designationMasterId"}
		FormCondVal := []string{Designation.DesignationID}
		ColumnCondVal := []string{"designationMasterId"}
		Result = serv("designationmaster", "edit", Designation.TravelAgencyMasterID, nil, nil, FormVal, ColumnsVal, FormCondVal, ColumnCondVal)
		fmt.Fprintln(w, Result)
	} else {
		Result = "Invalid designationMasterId"
		fmt.Println(Result)
	}
}

func GetDesignation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	designationID := r.FormValue("designationId")
	//methType := r.FormValue("methodType")
	FormCondVal := []string{designationID}
	ColumnCondVal := []string{"designationMasterId"}
	Result := serv("designationmaster", "list", "", nil, nil, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)
}

func ListDesignation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")
	Result := serv("designationmaster", "list", travelAgencyMasterId, nil, nil, nil, nil, nil, nil)
	fmt.Fprintln(w, Result)
}

func ListDesignationByDep(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")
	departmentID := r.FormValue("departmentId")
	//methType := r.FormValue("methodType")
	FormCondVal := []string{departmentID}
	ColumnCondVal := []string{"department"}
	Result := serv("designationmaster", "list", travelAgencyMasterId, nil, nil, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)
}

func DeleteDesignation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	designationID := r.FormValue("designationId")
	//methType := r.FormValue("methodType")
	FormCondVal := []string{designationID}
	ColumnCondVal := []string{"designationMasterId"}
	Result := serv("designationmaster", "delete", "", nil, nil, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)
}

func CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var Result string
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	DepartmentVar := UnmarshalDepartment(string(body))

	if checkIfExist([]string{DepartmentVar.TravelAgencyMasterID}, []string{"travelAgencyMasterId"}, "travelagencymaster") {
		if Validate(DepartmentVar.DepartmentName) || Validate(DepartmentVar.TravelAgencyMasterID) {
			fmt.Fprintln(w, "Department name must be filled")
		} else {
			Form := []string{DepartmentVar.DepartmentName, DepartmentVar.DepartmentCode, DepartmentVar.TravelAgencyMasterID}
			Columns := []string{"departmentName", "departmentCode", "travelAgencyMasterId"}

			Result = serv("department", "create", "", Form, Columns, nil, nil, nil, nil)
			if Result == "Record already exists" {
				fmt.Fprintln(w, "Department already exists")
			} else if len(Result) <= 5 {
				fmt.Fprintln(w, "Department created successfully")
			} else {
				fmt.Fprintln(w, Result)
			}
		}
	} else {
		fmt.Println("Invalid Company(travelAgencyMaster)Id")
	}
}

func EditDepartment(w http.ResponseWriter, r *http.Request) {
	var Result string
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)

	DepartmentVar := UnmarshalDepartment(string(body))

	if checkIfExist([]string{DepartmentVar.DepartmentID}, []string{"departmentMasterId"}, "department") {
		FormVal := []string{DepartmentVar.DepartmentName, DepartmentVar.DepartmentCode, DepartmentVar.TravelAgencyMasterID, DepartmentVar.DepartmentID}
		ColumnVal := []string{"departmentName", "departmentCode", "travelAgencyMasterId", "departmentMasterId"}
		FormCondVal := []string{DepartmentVar.DepartmentID}
		ColumnCondVal := []string{"departmentMasterId"}
		Result = serv("department", "edit", DepartmentVar.TravelAgencyMasterID, nil, nil, FormVal, ColumnVal, FormCondVal, ColumnCondVal)
		fmt.Fprintln(w, Result)
	} else {
		Result = "Invalid departmentId"
		fmt.Println(Result)
	}
}

func GetDepartment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	departmentId := r.FormValue("departmentId")
	//methType := r.FormValue("methodType")
	FormCondVal := []string{departmentId}
	ColumnCondVal := []string{"departmentMasterId"}
	Result := serv("department", "list", "", nil, nil, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)
}

func ListDepartment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")
	Result := serv("department", "list", travelAgencyMasterId, nil, nil, nil, nil, nil, nil)
	fmt.Fprintln(w, Result)
}

func DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	departmentID := r.FormValue("departmentId")
	//methType := r.FormValue("methodType")
	FormCondVal := []string{departmentID}
	ColumnCondVal := []string{"departmentMasterId"}
	Result := serv("department", "delete", "", nil, nil, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)
}

func ListCities(w http.ResponseWriter, r *http.Request) {
	Result := serv("citymaster", "list", "", nil, nil, nil, nil, nil, nil)
	fmt.Fprintln(w, Result)

}
func CityCatMap(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	var Result1 string
	countcity := 0
	CityCatPar := UnmarshalJsonCityCat(string(body))

	if Validate(CityCatPar.CityCat.Label) {
		fmt.Fprintln(w, "City Category name must not be Empty!! ;)")
	} else {
		Form := []string{CityCatPar.CityCat.Label, CityCatPar.CompanyID}
		Columns := []string{"CityCatName", "CompanyId"}
		Result := serv("city_category", "create", CityCatPar.CompanyID, Form, Columns, nil, nil, nil, nil)
		if Result != "Record already exists" {
			for _, cities := range CityCatPar.Cities {

				Form := []string{Result, cities.Value}
				Columns := []string{"CityCatID", "CityID"}
				Result1 = serv("city_mapping", "create", CityCatPar.CompanyID, Form, Columns, nil, nil, nil, nil)
				log.Print(Result1)
				countcity++
			}
			fmt.Fprintln(w, "City Category successfully created with "+strconv.Itoa(countcity)+" choosen cities")

		} else {
			fmt.Fprintln(w, "City Category Already Exists")
		}
	}
}
func EditCityCatMap(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	res := ""
	CityCatPar := UnmarshalJsonCityCatEdit(string(body))
	fmt.Println(CityCatPar)
	if Validate(CityCatPar.CityCat.Label) {
		fmt.Fprintln(w, "City Category name must not be Empty!! ;)")
	} else {
		FormVal := []string{CityCatPar.CityCat.Label}
		ColumnVal := []string{"CityCatName"}
		FormCondVal := []string{CityCatPar.CityCat.Value, CityCatPar.CompanyID}
		ColumnCondVal := []string{"CityCatID", "CompanyID"}
		res = serv("city_category", "edit", CityCatPar.CompanyID, nil, nil, FormVal, ColumnVal, FormCondVal, ColumnCondVal)
		//fmt.Fprintln(w, Result1)

		for _, cities := range CityCatPar.Cities {

			Form := []string{CityCatPar.CityCat.Value, cities.Value}
			Columns := []string{"CityCatID", "CityID"}
			FormCondVal := []string{CityCatPar.CityCat.Value}
			ColumnCondVal := []string{"CityCatID"}
			res = serv("city_mapping", "create", CityCatPar.CompanyID, Form, Columns, nil, nil, FormCondVal, ColumnCondVal)
			//fmt.Fprintln(w, Result1)
		}
		for _, remcities := range CityCatPar.RemCities {

			//Form := []string{CityCatPar.CityCat.Value, cities.Value}
			//Columns := []string{"CityCatID", "CityID"}
			FormCondVal := []string{remcities.Value, CityCatPar.CityCat.Value}
			ColumnCondVal := []string{"CityID", "CityCatID"}
			res = serv("city_mapping", "delete", "", nil, nil, nil, nil, FormCondVal, ColumnCondVal)
			//fmt.Fprintln(w, Result1)
		}
		if res != "" {
			fmt.Fprintln(w, "All Changes Saved")
		}
	}
}

// func AddCities(w http.ResponseWriter, r *http.Request) {
// 	body, err := ioutil.ReadAll(r.Body)
// 	checkErr(err)

// 	CityCatPar := UnmarshalJsonCityCat(string(body))
// 	//Form := []string{CityCatPar.CityCat.Value, CityCatPar.CompanyID}
// 	//Columns := []string{"CityCatName", "CompanyId"}
// 	//Result := serv("city_category", "create", CityCatPar.CompanyID, Form, Columns, nil, nil, nil, nil)

// 	for _, cities := range CityCatPar.Cities {

// 		Form := []string{CityCatPar.CityCat.Value, cities.Value}
// 		Columns := []string{"CityCatID", "CityID"}
// 		FormCondVal := []string{CityCatPar.CityCat.Value}
// 		ColumnCondVal := []string{"CityCatID"}
// 		Result1 := serv("city_mapping", "create", CityCatPar.CompanyID, Form, Columns, nil, nil, FormCondVal, ColumnCondVal)
// 		fmt.Fprintln(w, Result1)
// 	}

// }
// func RemoveCities(w http.ResponseWriter, r *http.Request) {
// 	body, err := ioutil.ReadAll(r.Body)
// 	checkErr(err)

// 	CityCatPar := UnmarshalJsonCityCat(string(body))
// 	//Form := []string{CityCatPar.CityCat.Value, CityCatPar.CompanyID}
// 	//Columns := []string{"CityCatName", "CompanyId"}
// 	//Result := serv("city_category", "create", CityCatPar.CompanyID, Form, Columns, nil, nil, nil, nil)

// 	for _, cities := range CityCatPar.Cities {

// 		//Form := []string{CityCatPar.CityCat.Value, cities.Value}
// 		//Columns := []string{"CityCatID", "CityID"}
// 		FormCondVal := []string{cities.Value}
// 		ColumnCondVal := []string{"CityMappingID"}
// 		Result1 := serv("city_mapping", "delete", "", nil, nil, nil, nil, FormCondVal, ColumnCondVal)
// 		fmt.Fprintln(w, Result1)
// 	}
// }

func ListCityCatCities(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	Form := []string{r.FormValue("cityCatId")}
	Columns := []string{"CityCatID"}
	companyID := r.FormValue("companyId")
	Result := serv("city_mapping", "list", companyID, Form, Columns, nil, nil, nil, nil)
	fmt.Fprintln(w, Result)

}

func CityCategory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Manipulate CityCategory stuff inside me!!")
	table := "city_category"
	m := make(map[string]string)

	m["cityCatName"] = "CityCatName" //mapping formNames to Database Column names
	m["companyId"] = "CompanyID"
	m["cityCatId"] = "CityCatID"
	err := r.ParseForm()
	checkErr(err)

	Form := []string{r.FormValue("cityCatName"), r.FormValue("companyId")}
	Columns := []string{m["cityCatName"], m["companyId"]}

	FormVal := []string{r.FormValue("cityCatName")}
	ColumnVal := []string{m["cityCatName"]}
	FormCondVal := []string{r.FormValue("cityCatId")}
	ColumnCondVal := []string{m["cityCatId"]}
	companyID := r.FormValue("companyId")

	methType := r.FormValue("methodType")

	fmt.Println(r.FormValue("cityCatId"))
	fmt.Println(r.FormValue("methodType"))

	Result := serv(table, methType, companyID, Form, Columns, FormVal, ColumnVal, FormCondVal, ColumnCondVal)

	fmt.Fprintln(w, Result)
}

// func CityMapping(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("City Mapping stuff inside me!!")
// 	table := "city_mapping"
// 	m := make(map[string]string)

// 	m["cityCatId"] = "CityCatID" //mapping formNames to Database Column names
// 	m["cityId"] = "CityID"
// 	m["companyId"] = "CompanyID"
// 	m["cityMappingId"] = "CityMappingID"

// 	err := r.ParseForm()
// 	checkErr(err)

// 	companyID := r.FormValue("companyId")
// 	methType := r.FormValue("methodType")

// 	Form := []string{r.FormValue("cityCatId"), r.FormValue("cityId")}
// 	Columns := []string{m["cityCatId"], m["cityId"]}
// 	//FormVal := []string{r.FormValue("cityCatName")}
// 	//ColumnVal := []string{m["cityCatName"]}
// 	FormCondVal := []string{r.FormValue("cityMappingId")}
// 	ColumnCondVal := []string{m["cityMappingId"]}

// 	Result := serv(table, methType, companyID, Form, Columns, nil, nil, FormCondVal, ColumnCondVal)
// 	fmt.Fprintln(w, Result)

// 	for key, values := range r.Form { // range over map
// 		if strings.Contains(key, "cityId[") {
// 			for key, value := range values { // range over []string
// 				fmt.Println(key, value)
// 				Form := []string{r.FormValue("cityCatId"), value}
// 				Columns := []string{m["cityCatId"], m["cityId"]}
// 				fmt.Println(Form, "  ", Columns)

// 				Result := serv(table, methType, companyID, Form, Columns, nil, nil, nil, nil)
// 				fmt.Fprintln(w, Result)

// 			}
// 		}
// 	}
// }

func BenefitBundle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Benefit stuff inside me!!")
	//var wg sync.WaitGroup
	//table := "city_mapping"
	m := make(map[string]string)
	m["benefitBundleName"] = "BenefitBundleName" //mapping formNames to Database Column names
	m["benefitBundleCode"] = "BenefitBundleCode"
	m["benefitBundleId"] = "BenefitBundleID"
	m["benefitBundleTypeMappingId"] = "BenefitBundleTypeMappingID"
	m["bundleTypeBenefitMappingId"] = "BundleTypeBenefitMappingID"
	m["companyId"] = "CompanyID"
	m["benefitTypeId"] = "BenefitTypeID"
	m["priority"] = "Priority"
	m["benefitId"] = "BenefitID"
	m["cityCatId"] = "CityCatID"
	m["limitSpend"] = "LimitSpend"
	m["maxAmount"] = "MaxAmount"
	m["minAmount"] = "MinAmount"
	m["flexibility"] = "Flexibility"
	m["flexAmount"] = "FlexAmount"
	m["starCat"] = "StarCat"

	//err := r.ParseForm()
	//checkErr(err)

	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	//log.Println(string(body))

	PolicyBundle := UnmarshalJsonPolicyBundle(string(body))
	//fmt.Println(PolicyBundle)
	fmt.Println(PolicyBundle.BundleName)

	//defer wg.Done()
	table := "policy_benefit_bundle"

	fmt.Println("table   ", table)
	var bundleId string
	var Result string

	if Validate(PolicyBundle.BundleName) || Validate(PolicyBundle.BundleCode) {
		fmt.Fprintln(w, "Policy Bundle Name/Code must not be Empty!! ;)")
	} else {
		Form := []string{PolicyBundle.BundleName, PolicyBundle.BundleCode, PolicyBundle.CompanyID}
		Columns := []string{m["benefitBundleName"], m["benefitBundleCode"], m["companyId"]}
		FormVal := []string{PolicyBundle.BundleName, PolicyBundle.BundleCode}
		ColumnVal := []string{m["benefitBundleName"], m["benefitBundleCode"]}
		FormCondVal := []string{PolicyBundle.BenefitBundleID}
		ColumnCondVal := []string{m["benefitBundleId"]}
		if PolicyBundle.MethType == "edit" {
			Result = serv(table, "delete", PolicyBundle.CompanyID, Form, Columns, FormVal, ColumnVal, FormCondVal, ColumnCondVal)
			Result = serv(table, "create", PolicyBundle.CompanyID, Form, Columns, FormVal, ColumnVal, FormCondVal, ColumnCondVal)

		} else {
			Result = serv(table, PolicyBundle.MethType, PolicyBundle.CompanyID, Form, Columns, FormVal, ColumnVal, FormCondVal, ColumnCondVal)
			if Result == "Record already exists" {
				fmt.Fprintln(w, "BundleCode/BundleName already exists. Enter Unique Bundle code for each bundle!!")
			}
		}

		//bundleId = Result
		//	time.Sleep(time.Second * 3)
		if PolicyBundle.MethType == "list" {
			fmt.Println(",,,,,///////////,,,,,,")
			bundleId = PolicyBundle.BenefitBundleID
		} else if Result != "Record already exists" {
			bundleId = Result
		} else {
			bundleId = ""
			//bundleId = getId(table, "BenefitBundleID", Form, Columns)
			//Result = getBenefitBundleID(r.FormValue("benefitBundleCode"), companyID)
		}

		table = "benefit_bundle_type_mapping"
		var benefittypemappingId string
		var Result1 string
		listbenefitbundlestr := ""
		for i := range PolicyBundle.PolicyBundles {
			fmt.Println(PolicyBundle.PolicyBundles[i])
			//defer wg.Done()
			if bundleId != "" {
				fmt.Println("BenefitBundleID   ", bundleId)
				Form1 := []string{bundleId, PolicyBundle.PolicyBundles[i].BenefitTypeID.Value, PolicyBundle.PolicyBundles[i].Priority}
				Columns1 := []string{m["benefitBundleId"], m["benefitTypeId"], m["priority"]}
				FormVal := []string{PolicyBundle.PolicyBundles[i].BenefitTypeID.Value, PolicyBundle.PolicyBundles[i].Priority}
				ColumnVal := []string{m["benefitTypeId"], m["priority"]}
				FormCondVal := []string{bundleId}
				ColumnCondVal := []string{m["benefitBundleId"]}
				Result1 = serv(table, "create", PolicyBundle.CompanyID, Form1, Columns1, FormVal, ColumnVal, FormCondVal, ColumnCondVal)

				fmt.Println("benefittypemappingId ", Result1)
				benefittypemappingId = Result1

				//time.Sleep(time.Second * 2)
				// if PolicyBundle.MethType != "list" {
				// 	if Result1 != "Record already exists" {
				// 		benefittypemappingId = Result1
				// 	} else {
				// 		benefittypemappingId = getId(table, "BenefitBundleTypeMappingID", Form1, Columns1)

				// 		//benefittypemappingId = Result1
				// 	}
				// }
			} else {
				benefittypemappingId = ""
			}

			for _, benefits := range PolicyBundle.PolicyBundles[i].Benefits {
				table := "bundle_type_benefit_mapping"
				//	fmt.Println(">>>>>>>>>>>>>>>>", benefittypemappingId)
				if benefittypemappingId != "" {
					var Result2 string

					Form := []string{benefittypemappingId, benefits.Value}
					Columns := []string{"BenefitBundleTypeMappingID", m["benefitId"]}
					FormVal := []string{benefits.Value}
					ColumnVal := []string{m["benefitId"]}
					FormCondVal := []string{benefittypemappingId}
					ColumnCondVal := []string{m["benefitBundleTypeMappingId"]}
					fmt.Println(Form, "  ", Columns)

					Result2 = serv(table, "create", PolicyBundle.CompanyID, Form, Columns, FormVal, ColumnVal, FormCondVal, ColumnCondVal)

					//	allowance_mapping(a)
					//time.Sleep(time.Second * 1)
					fmt.Println(Result2)
					//	time.Sleep(time.Second * 3)

				} else {
					fmt.Println("")
				}

			}
			for j := range PolicyBundle.PolicyBundles[i].CityCatAndAllowances {
				//defer wg.Done()

				table := "benefit_type_allowance_mapping"
				fmt.Println(benefittypemappingId)
				if benefittypemappingId != "" && checkIfExist([]string{PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].Value, PolicyBundle.CompanyID}, []string{m["cityCatId"], m["companyId"]}, "city_category") {
					var Result3 string

					limit := PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].LimitSpent
					Form1 := []string{benefittypemappingId, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].Value, limit, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].Max, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].Min, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].Flex, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].FlexAmt, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].StarCat}
					Columns1 := []string{"BenefitBundleTypeMappingID", m["cityCatId"], m["limitSpend"], m["maxAmount"], m["minAmount"], m["flexibility"], m["flexAmount"], m["starCat"]}
					FormCondVal := []string{benefittypemappingId}
					ColumnCondVal := []string{m["benefitBundleTypeMappingId"]}
					if PolicyBundle.MethType == "edit" {
						Result3 = serv(table, "create", PolicyBundle.CompanyID, Form1, Columns1, nil, nil, nil, nil)
						listbenefitbundlestr = Result3
					} else {
						Result3 = serv(table, PolicyBundle.MethType, PolicyBundle.CompanyID, Form1, Columns1, nil, nil, FormCondVal, ColumnCondVal)
						listbenefitbundlestr = Result3
					}
					//fmt.Println(values)
					//}

					//time.Sleep(time.Second * 1)
					fmt.Println(Result3)
					//	time.Sleep(time.Second * 3)

				} else {
					fmt.Fprintln(w, "")
				}

			}
		}
		if listbenefitbundlestr != "" {
			fmt.Fprintln(w, "Policy Bundle Saved Successfully")
		}
	}
}

func ListBenefitBundle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside me list policy bundle")
	err := r.ParseForm()
	checkErr(err)

	companyId := r.FormValue("companyId")
	//methType := r.FormValue("methodType")
	//	FormCondVal := []string{companyId}
	//	ColumnCondVal := []string{"CompanyID"}
	Result3 := serv("policy_benefit_bundle", "list", companyId, nil, nil, nil, nil, nil, nil)
	fmt.Fprintln(w, Result3)
}

func ListBundleRequirement(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside me list initial requirement for policy bundle")
	err := r.ParseForm()
	checkErr(err)
	companyId := r.FormValue("companyId")
	//FormCondVal := []string{companyId}
	//ColumnCondVal := []string{"CompanyID"}
	Result := serv("benefit_type_master", "list", companyId, nil, nil, nil, nil, nil, nil)
	fmt.Fprintln(w, Result)
}
func GetBenefitBundle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside me request get policy bundle")
	err := r.ParseForm()
	checkErr(err)
	benefitBundleId := r.FormValue("benefitBundleId")
	FormCondVal := []string{benefitBundleId}
	ColumnCondVal := []string{"BenefitBundleID"}
	Result := serv("policy_benefit_bundle", "list", "", nil, nil, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)
}

func DeleteBundle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside me delete policy bundle")
	err := r.ParseForm()
	checkErr(err)
	benefitBundleId := r.FormValue("benefitBundleId")

	FormCondVal := []string{benefitBundleId}
	ColumnCondVal := []string{"BenefitBundleID"}
	Result := serv("policy_benefit_bundle", "delete", "", nil, nil, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Employee Editing .....")
	err := r.ParseForm()
	checkErr(err)
	travelAgencyUsersId := r.FormValue("travelAgencyUserId")

	Result := GetEmployeedetails("travelAgencyUsers", travelAgencyUsersId)
	fmt.Fprintln(w, Result)
}
func EditEmployee(w http.ResponseWriter, r *http.Request) {
	var MemberProfileId string
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	//log.Println(string(body))

	EmployeeDet := UnmarshalAddEmployee(string(body))

	TravelagencyusersId := EmployeeDet.TravelAgencyUserId
	designationText := getId("designationamaster", "designationName", []string{EmployeeDet.Designation}, []string{"designationMasterId"})
	FormVal := []string{EmployeeDet.CompanyID, MemberProfileId, EmployeeDet.CompanyName, EmployeeDet.CompanyMail, EmployeeDet.PersonalMail, designationText, EmployeeDet.Designation, EmployeeDet.EmployeeName, EmployeeDet.MobileNo}
	ColumnVal := []string{"travelAgencyMasterId", "memberProfileId", "travelAgencyNameTemp", "email", "personalEmail", "designation", "designationId", "virtualName", "mobile"}
	FormCondVal := []string{TravelagencyusersId}
	ColumnCondVal := []string{"travelagencyusersId"}
	Result := serv("travelagencyusers", "edit", EmployeeDet.CompanyID, nil, nil, FormVal, ColumnVal, FormCondVal, ColumnCondVal)
	fmt.Println(Result)
	maxdate := getId("user_employment_track", "desgStartDate", []string{"3", TravelagencyusersId}, []string{"travelAgencyMasterId", "travelAgencyUsersId"})

	t, err := time.Parse("2006-01-02T00:00:00Z", maxdate)

	trackId := getId("user_employment_track", "trackMasterId", []string{"3", TravelagencyusersId, t.String()}, []string{"travelAgencyMasterId", "travelAgencyUsersId", "desgStartDate"})
	fmt.Println(trackId)

	//Result1 := getId("user_employment_track", "trackMasterId", []string{"3", Result}, []string{"travelAgencyMasterId", "travelAgencyUserId"})

	t, err = time.Parse("02 Jan 2006", EmployeeDet.StartDate)
	checkErr(err)
	fmt.Println(TravelagencyusersId, "3", t.String(), EmployeeDet.Designation)
	Form1 := []string{TravelagencyusersId, "3", t.String(), EmployeeDet.Designation}
	Columns1 := []string{"travelAgencyUsersId", "travelAgencyMasterId", "desgStartDate", "designationId"}
	FormCondVal = []string{trackId}
	ColumnCondVal = []string{"trackMasterId"}
	Result1 := serv("user_employment_track", "edit", "3", nil, nil, Form1, Columns1, FormCondVal, ColumnCondVal)
	fmt.Println(Result1)

	fmt.Fprintln(w, "Rows Updated")

}

func DeactivateEmployee(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Employee Deactivating .....")

	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	//log.Println(string(body))
	Emps := UnmarshalDeleteEmps(string(body))

	for i := range Emps.TravelAgencyUserId {
		Form1 := []string{"0"}
		Columns1 := []string{"status"}
		FormCondVal := []string{Emps.TravelAgencyUserId[i]}
		ColumnCondVal := []string{"travelAgencyUsersId"}
		Result1 := serv("travelagencyusers", "edit", "", nil, nil, Form1, Columns1, FormCondVal, ColumnCondVal)
		fmt.Fprintln(w, Result1)

	}
}

func AddEmployee(w http.ResponseWriter, r *http.Request) {
	var MemberProfileId string
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	//log.Println(string(body))
	pass := RandomPass()
	EmployeeDet := UnmarshalAddEmployee(string(body))
	//fmt.Println(PolicyBundle)
	fmt.Println(EmployeeDet)

	if checkIfExist([]string{EmployeeDet.PersonalMail}, []string{"email"}, "memberprofile") {
		fmt.Println("Record exist inside member profile")
		MemberProfileId = getId("memberprofile", "memberProfileId", []string{EmployeeDet.PersonalMail}, []string{"email"})

	} else {
		Form := []string{EmployeeDet.PersonalMail, "", EmployeeDet.MobileNo, pass}
		Columns := []string{"email", "phone", "mobile", "password"}
		MemberProfileId = serv("memberprofile", "create", "3", Form, Columns, nil, nil, nil, nil)
		//MemberProfileId = "0"
	}
	Form := []string{EmployeeDet.CompanyMail}
	Column := []string{"email"}
	var Result string
	if !checkIfExist(Form, Column, "travelagencyusers") {
		Form := []string{EmployeeDet.CompanyID, MemberProfileId, EmployeeDet.CompanyName, EmployeeDet.CompanyMail, EmployeeDet.PersonalMail, EmployeeDet.Designation, EmployeeDet.EmployeeName, EmployeeDet.MobileNo, pass, EmployeeDet.BenefitBundleId}
		Columns := []string{"travelAgencyMasterId", "memberProfileId", "travelAgencyNameTemp", "email", "personalEmail", "designationId", "virtualName", "mobile", "password", "benefitBundleId"}

		Result = serv("travelagencyusers", "create", EmployeeDet.CompanyID, Form, Columns, nil, nil, nil, nil)
		sendMail(EmployeeDet.CompanyMail, pass, EmployeeDet.CompanyMail)
	} else {
		TravelagencyusersId := getId("travelagencyusers", "travelagencyusersId", []string{EmployeeDet.CompanyMail}, []string{"email"})
		FormVal := []string{EmployeeDet.CompanyID, MemberProfileId, EmployeeDet.CompanyName, EmployeeDet.CompanyMail, EmployeeDet.PersonalMail, EmployeeDet.Designation, EmployeeDet.EmployeeName, EmployeeDet.MobileNo, EmployeeDet.BenefitBundleId}
		ColumnVal := []string{"travelAgencyMasterId", "memberProfileId", "travelAgencyNameTemp", "email", "personalEmail", "designationId", "virtualName", "mobile", "benefitBundleId"}
		FormCondVal := []string{TravelagencyusersId}
		ColumnCondVal := []string{"travelagencyusersId"}
		serv("travelagencyusers", "edit", EmployeeDet.CompanyID, nil, nil, FormVal, ColumnVal, FormCondVal, ColumnCondVal)
		Result = TravelagencyusersId

	}
	var Result1 string
	fmt.Println(EmployeeDet.CompanyID, Result, EmployeeDet.Designation)
	if !checkIfExist([]string{EmployeeDet.CompanyID, Result, EmployeeDet.Designation}, []string{"travelAgencyMasterId", "travelAgencyUsersId", "designationId"}, "user_employment_track") {
		t, err := time.Parse("02 Jan 2006", EmployeeDet.StartDate)
		checkErr(err)
		Form1 := []string{Result, "3", MemberProfileId, "s1", t.String(), "1", EmployeeDet.Designation}
		Columns1 := []string{"travelAgencyUsersId", "travelAgencyMasterId", "memberProfileId", "travelAgencyName", "desgStartDate", "active", "designationId"}
		Result1 = serv("user_employment_track", "create", "3", Form1, Columns1, nil, nil, nil, nil)
		fmt.Fprintln(w, "Rows Updated")

	} else {

		fmt.Println("Record Exists in user_employment_track")
		Result1 = getId("user_employment_track", "trackMasterId", []string{"3", Result}, []string{"travelAgencyMasterId", "travelAgencyUserId"})
		fmt.Println("trackMasterId", Result)
		fmt.Fprintln(w, "Record Exists")

	}
	fmt.Println(Result1)
}

func UploadEmployees(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside Upload emp")

	err := r.ParseForm()
	checkErr(err)
	EmpCsvUrl := r.FormValue("importCsv")
	fmt.Println("importCsv : ", EmpCsvUrl)

	file, err := os.Open(EmpCsvUrl)

	//C:\HostingSpaces\admin\amazingnature.in\wwwroot\hobse\demo\uploads\emp.csv

	//file, err := os.Open("C:\\HostingSpaces\\admin\\amazingnature.in\\wwwroot\\hobse\\demo\\uploads\\emp.csv")
	checkErr(err)
	// automatically call Close() at the end of current method
	defer file.Close()
	//
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	checkErr(err)

	reader.Comma = ','
	lineCount := 0
	totalrecords := 0
	for i, record := range lines {
		var Result, Result1 string

		totalrecords++
		// read just one record, but we could ReadAll() as well
		if i == 0 {
			// skip header line
			continue
		}
		//	record, err := line.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		pass := RandomPass()

		Form := []string{"3", record[0], record[1], record[2], record[3], record[4], record[5], record[6]}
		Column := []string{"travelAgencyMasterId", "travelAgencyNameTemp", "virtualName", "designation", "personalEmail", "email", "phone", "mobile"}

		if !checkIfExist(Form, Column, "travelagencyusers") {
			if record[4] == "" && record[3] == "" {
				fmt.Println("profile cant be created")
				failText := "profile cant be created"
				Form := []string{"3", record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], failText}
				Columns := []string{"travelAgencyMasterId", "travelAgencyNameTemp", "virtualName", "designation", "personalEmail", "companyEmail", "phone", "mobile", "startDate", "failReason"}
				serv("travelagencyuserfail_log", "create", "3", Form, Columns, nil, nil, nil, nil)

			} else if record[3] != "" && !checkIfExist([]string{record[3]}, []string{"email"}, "memberprofile") {
				fmt.Println("If not exist in memberprofile")
				Form := []string{record[3], record[5], record[6], pass}
				Columns := []string{"email", "phone", "mobile", "password"}
				Result = serv("memberprofile", "create", "3", Form, Columns, nil, nil, nil, nil)
			} else if checkIfExist([]string{record[3]}, []string{"email"}, "memberprofile") {
				fmt.Println("Record exist inside member profile")
				Result = getId("memberprofile", "memberProfileId", []string{record[3]}, []string{"email"})
			} else {
				fmt.Println("Record not to be updated")
			}
			failText := "fail"
			if !(record[4] == "" && record[3] == "") {
				if record[4] != "" && checkIfExist([]string{"3", record[4]}, []string{"travelAgencyMasterId", "email"}, "travelagencyusers") {
					fmt.Println(record[4], "  Record Exists in travel agency usr")
					Result = getId("travelagencyusers", "travelAgencyUsersId", []string{"3", record[4]}, []string{"travelAgencyMasterId", "email"})
					fmt.Println("travelAgencyUserId", Result)
				} else {
					// record is an array of string so is directly printable
					fmt.Println("Record", lineCount, "is", record, "and has", len(record), "fields")
					// and we can iterate on top of that
					//for i := 0; i < len(record); i++ {
					fmt.Println(" >>>", record[0]+record[1]+record[2]+record[3])
					if record[4] == "" || ExtensionCheck(record[4], "3") {
						fmt.Println("extension exists")
						Form := []string{"3", record[0], record[1], record[2], record[3], record[4], record[5], record[6], pass}
						Columns := []string{"travelAgencyMasterId", "travelAgencyNameTemp", "virtualName", "designation", "personalEmail", "email", "phone", "mobile", "password"}
						Result = serv("travelagencyusers", "create", "3", Form, Columns, nil, nil, nil, nil)
						lineCount += 1

					} else {
						failText = "Email extension is not correct"
						Form := []string{"3", record[0], record[1], record[2], record[3], record[4], record[5], record[6], record[7], failText}
						Columns := []string{"travelAgencyMasterId", "travelAgencyNameTemp", "virtualName", "designation", "personalEmail", "companyEmail", "phone", "mobile", "startDate", "failReason"}
						Result = serv("travelagencyuserfail_log", "create", "3", Form, Columns, nil, nil, nil, nil)
					}

				}

				if failText == "fail" {

					if checkIfExist([]string{"3", Result}, []string{"travelAgencyMasterId", "travelAgencyUserId"}, "user_employment_track") {
						fmt.Println("Record Exists in user_employment_track")
						Result1 = getId("user_employment_track", "trackMasterId", []string{"3", Result}, []string{"travelAgencyMasterId", "travelAgencyUserId"})
						fmt.Println("trackMasterId", Result)
					} else {
						t, err := time.Parse("1/2/2006", record[7])
						checkErr(err)
						Form1 := []string{Result, "3", "s1", t.String(), "1"}
						Columns1 := []string{"travelAgencyUsersId", "travelAgencyMasterId", "travelAgencyName", "desgStartDate", "active"}
						Result1 = serv("user_employment_track", "create", "3", Form1, Columns1, nil, nil, nil, nil)
					}

					var memberProfileID string

					if checkIfExist([]string{record[3]}, []string{"email"}, "memberprofile") {
						memberProfileID = getId("memberprofile", "memberProfileId", []string{record[3]}, []string{"email"})
						FormVal := []string{memberProfileID}
						ColumnsVal := []string{"memberProfileId"}
						FormCondVal := []string{Result}
						ColumnCondVal := []string{"travelAgencyUsersId"}
						FormCondVal1 := []string{Result1}
						ColumnCondVal1 := []string{"trackMasterId"}
						go serv("travelagencyusers", "edit", "3", nil, nil, FormVal, ColumnsVal, FormCondVal, ColumnCondVal)
						go serv("user_employment_track", "edit", "3", nil, nil, FormVal, ColumnsVal, FormCondVal1, ColumnCondVal1)

						//go UpdateMemberPro()
					} else {
						sendMail(record[4], pass, record[4])
						fmt.Println("sending mail to " + record[4] + " id to sign up with hobse using personal details")
					}
				}
			}
		}
		//}
		//fmt.Fprintln(w, Result)
	}
	fmt.Fprintln(w, lineCount, " record inserted out of ", totalrecords-1)

}

func AssignDesignation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside me assign designation")
	body, err := ioutil.ReadAll(r.Body)
	checkErr(err)
	//log.Println(string(body))

	AssignDesignation := UnmarshalAssignDesig(string(body))
	//fmt.Println(PolicyBundle)
	fmt.Println(AssignDesignation)

	benefitBundleId := getId("designationmaster", "benefitBundleId", []string{AssignDesignation.DesignationName.Value}, []string{"designationMasterId"})
	hierarchyId := getId("designationmaster", "hierarchyId", []string{AssignDesignation.DesignationName.Value}, []string{"designationMasterId"})

	for i := range AssignDesignation.TravelAgencyUserId {
		fmt.Println(AssignDesignation.DesignationName.Label, AssignDesignation.DesignationName.Value)
		FormVal := []string{AssignDesignation.DesignationName.Label, AssignDesignation.DesignationName.Value, benefitBundleId, hierarchyId}
		ColumnVal := []string{"designation", "designationId", "benefitBundleId", "hierarchyId"}
		FormCondVal := []string{AssignDesignation.TravelAgencyUserId[i]}
		ColumnCondVal := []string{"travelAgencyUsersId"}
		Result := serv("travelagencyusers", "edit", "", nil, nil, FormVal, ColumnVal, FormCondVal, ColumnCondVal)
		fmt.Println(Result)

	}
	//	fmt.Println(Result)
}

func BenefitBundleByDes(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)

	designationId := r.FormValue("designationId")
	//methType := r.FormValue("methodType")
	FormCondVal := []string{designationId}
	ColumnCondVal := []string{"designationMasterId"}
	bundleId := getId("designationmaster", "benefitBundleId", FormCondVal, ColumnCondVal)
	bundleName := getId("policy_benefit_bundle", "BenefitBundleName", []string{bundleId}, []string{"BenefitBundleID"})
	hierarchyId := getId("designationmaster", "hierarchyId", FormCondVal, ColumnCondVal)
	Res := MarshalBenefitByDes(bundleId, bundleName, hierarchyId)

	fmt.Fprintln(w, Res)
}

func SearchByName(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get list of all assigned emps")
	err := r.ParseForm()
	checkErr(err)
	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")
	travelAgencyUserName := r.FormValue("searchText")

	res := SeachAllEmp("travelagencyusers", travelAgencyMasterId, travelAgencyUserName)
	fmt.Fprintln(w, res)
}

func ListAllEmp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get list of all unassigned emps")
	err := r.ParseForm()
	checkErr(err)
	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")
	FormCondVal := []string{travelAgencyMasterId}
	ColumnCondVal := []string{"travelAgencyMasterId"}
	Result := serv("travelagencyusers", "list", "", nil, nil, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)
}

func ListUnassignedEmp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside me request get policy bundle")
	err := r.ParseForm()
	checkErr(err)
	travelAgencyMasterId := r.FormValue("travelAgencyMasterId")
	FormCondVal := []string{travelAgencyMasterId}
	ColumnCondVal := []string{"travelAgencyMasterId"}
	Result := serv("travelagencyusers", "list", travelAgencyMasterId, nil, nil, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)
}

func ExtensionCheck(str, travelAgencyMasterId string) bool {
	emailExtension := getId("travelagencymaster", "officialemailExtension", []string{travelAgencyMasterId}, []string{"travelAgencyMasterId"})
	if strings.Contains(str, emailExtension) {
		return true
	} else {
		return false
	}
}
func Validate(str string) bool {
	if str == "" {
		return true
	} else {
		return false
	}
}

func checkIfExist(value []string, column []string, table string) bool {
	if QueryRow(value, column, table) != 0 {
		return true
	} else {
		return false
	}
}
func checkErr(err error) {

	if err != nil {
		fmt.Println(err)
		//log.Fatal(err)
		//os.Exit(1)
	}
}

func RandomPass() string {
	n := 5
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		fmt.Println("Cant generate pass")
	}
	s := fmt.Sprintf("%X", b)
	return s
}

func sendMail(to, username, password string) {
	from := "priyanka@infonixweblab.com"
	subject := "SignIn with Hobse"
	To := to
	text := "Please sign in with hobse using this credentials. Username : " + username + ", Password : " + password
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", To)
	m.SetHeader("subject", subject)
	m.SetBody("text/html", text)
	d := gomail.NewDialer("smtp.zoho.com", 587, from, "priyanka123")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	fmt.Println("mail sent")
}
