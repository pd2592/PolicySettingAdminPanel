package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/citycategory", CityCategory)
	router.HandleFunc("/citymapping", CityMapping)
	router.HandleFunc("/benefitbundle", BenefitBundle)

	log.Fatal(http.ListenAndServe(":8080", router))
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

	Result := serv(table, methType, companyID, Form, Columns, FormVal, ColumnVal, FormCondVal, ColumnCondVal)

	fmt.Fprintln(w, Result)
}

func CityMapping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("City Mapping stuff inside me!!")
	table := "city_mapping"
	m := make(map[string]string)

	m["cityCatId"] = "CityCatID" //mapping formNames to Database Column names
	m["cityId"] = "CityID"
	m["companyId"] = "CompanyID"
	m["cityMappingId"] = "CityMappingID"

	err := r.ParseForm()
	checkErr(err)

	companyID := r.FormValue("companyId")
	methType := r.FormValue("methodType")

	Form := []string{r.FormValue("cityCatId"), r.FormValue("cityId")}
	Columns := []string{m["cityCatId"], m["cityId"]}
	//FormVal := []string{r.FormValue("cityCatName")}
	//ColumnVal := []string{m["cityCatName"]}
	FormCondVal := []string{r.FormValue("cityMappingId")}
	ColumnCondVal := []string{m["cityMappingId"]}

	Result := serv(table, methType, companyID, Form, Columns, nil, nil, FormCondVal, ColumnCondVal)
	fmt.Fprintln(w, Result)

	for key, values := range r.Form { // range over map
		if strings.Contains(key, "cityId[") {
			for key, value := range values { // range over []string
				fmt.Println(key, value)
				Form := []string{r.FormValue("cityCatId"), value}
				Columns := []string{m["cityCatId"], m["cityId"]}
				fmt.Println(Form, "  ", Columns)

				Result := serv(table, methType, companyID, Form, Columns, nil, nil, nil, nil)
				fmt.Fprintln(w, Result)

			}
		}
	}
}

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

	err := r.ParseForm()
	checkErr(err)
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
	}
	fmt.Fprintln(w, Result)
	//bundleId = Result
	//	time.Sleep(time.Second * 3)
	if PolicyBundle.MethType == "list" {
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
			if PolicyBundle.MethType == "edit" {
				Result1 = serv(table, "create", PolicyBundle.CompanyID, Form1, Columns1, FormVal, ColumnVal, FormCondVal, ColumnCondVal)
			} else {
				Result1 = serv(table, PolicyBundle.MethType, PolicyBundle.CompanyID, Form1, Columns1, FormVal, ColumnVal, FormCondVal, ColumnCondVal)

			}
			fmt.Fprintln(w, Result1)
			//time.Sleep(time.Second * 2)
			if PolicyBundle.MethType == "list" {
				Result1 = getId(table, "BenefitBundleTypeMappingID", FormCondVal, ColumnCondVal)

				benefittypemappingId = Result1
			} else if Result1 != "Record already exists" {
				benefittypemappingId = Result1
			} else {
				Result1 = getId(table, "BenefitBundleTypeMappingID", Form1, Columns1)

				benefittypemappingId = Result1
			}
		} else {
			benefittypemappingId = ""
		}

		for _, benefits := range PolicyBundle.PolicyBundles[i].Benefits {
			table := "bundle_type_benefit_mapping"
			fmt.Println(">>>>>>>>>>>>>>>>", benefittypemappingId)
			if benefittypemappingId != "" {
				var Result2 string

				Form := []string{benefittypemappingId, benefits.Value}
				Columns := []string{"BenefitBundleTypeMappingID", m["benefitId"]}
				FormVal := []string{benefits.Value}
				ColumnVal := []string{m["benefitId"]}
				FormCondVal := []string{benefittypemappingId}
				ColumnCondVal := []string{m["benefitBundleTypeMappingId"]}
				fmt.Println(Form, "  ", Columns)
				if PolicyBundle.MethType == "edit" {
					Result2 = serv(table, "create", PolicyBundle.CompanyID, Form, Columns, FormVal, ColumnVal, FormCondVal, ColumnCondVal)
				} else {
					Result2 = serv(table, PolicyBundle.MethType, PolicyBundle.CompanyID, Form, Columns, FormVal, ColumnVal, FormCondVal, ColumnCondVal)

				}
				//	allowance_mapping(a)
				//time.Sleep(time.Second * 1)
				fmt.Fprintln(w, Result2)
				//	time.Sleep(time.Second * 3)

			} else {
				fmt.Fprintln(w, "")
			}

		}
		for j := range PolicyBundle.PolicyBundles[i].CityCatAndAllowances {
			//defer wg.Done()

			table := "benefit_type_allowance_mapping"
			fmt.Println(benefittypemappingId)
			if benefittypemappingId != "" {
				var Result3 string

				Form1 := []string{benefittypemappingId, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].Value, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].LimitSpent, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].Min, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].Max, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].Flex, PolicyBundle.PolicyBundles[i].CityCatAndAllowances[j].FlexAmt}
				Columns1 := []string{"BenefitBundleTypeMappingID", m["cityCatId"], m["limitSpend"], m["maxAmount"], m["minAmount"], m["flexibility"], m["flexAmount"]}
				FormCondVal := []string{benefittypemappingId}
				ColumnCondVal := []string{m["benefitBundleTypeMappingId"]}
				if PolicyBundle.MethType == "edit" {
					Result3 = serv(table, "create", PolicyBundle.CompanyID, Form1, Columns1, nil, nil, nil, nil)
				} else {
					Result3 = serv(table, PolicyBundle.MethType, PolicyBundle.CompanyID, Form1, Columns1, nil, nil, FormCondVal, ColumnCondVal)

				}
				//fmt.Println(values)
				//}

				//time.Sleep(time.Second * 1)
				fmt.Fprintln(w, Result3)
				//	time.Sleep(time.Second * 3)

			} else {
				fmt.Fprintln(w, "")
			}

		}
	}
}

func checkErr(err error) {

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		os.Exit(1)
	}
}
