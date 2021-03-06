package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type PB struct {
	BenefitBundleID string         `json:"benefitBundleId,omitempty"`
	BundleName      string         `json:"bundleName,omitempty"`
	BundleCode      string         `json:"bundleCode,omitempty"`
	CompanyID       string         `json:"companyId,omitempty"`
	MethType        string         `json:"methType,omitempty"`
	PolicyBundles   []PolicyBundle `json:"policybundles,omitempty"`
}

type PolicyBundle struct {
	BenefitTypeID        LabVal                `json:"benefitTypeId,omitempty"`
	Priority             string                `json:"priority,omitempty"`
	Benefits             []LabVal              `json:"benefits,omitempty"`
	CityCatAndAllowances []CityCatAndAllowance `json:"cityCatAndAllowances,omitempty"`
}

type CityCatAndAllowance struct {
	Label      string `json:"label,omitempty"`
	Value      string `json:"value,omitempty"`
	LimitSpent string `json:"limitSpent,omitempty"`
	Min        string `json:"min,omitempty"`
	Max        string `json:"max,omitempty"`
	Flex       string `json:"flex,omitempty"`
	FlexAmt    string `json:"flexAmt,omitempty"`
	StarCat    string `json:"starCat,omitempty"`
}

type CityCategoryMap struct {
	CompanyID string   `json:"companyId,omitempty"`
	CityCat   LabVal   `json:"cityCat,omitempty"`
	Cities    []LabVal `json:"cities,omitempty"`
}
type CityCategoryMapEdit struct {
	CompanyID string   `json:"companyId,omitempty"`
	CityCat   LabVal   `json:"cityCat,omitempty"`
	Cities    []LabVal `json:"cities,omitempty"`
	RemCities []LabVal `json:"remCities,omitempty"`
}

type LabVal struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

type Department struct {
	DepartmentID         string `json:"departmentId,omitempty"`
	DepartmentName       string `json:"departmentName,omitempty"`
	DepartmentCode       string `json:"departmentCode,omitempty"`
	TravelAgencyMasterID string `json:"travelAgencyMasterId,omitempty"`
}

type Designation struct {
	DesignationID        string `json:"designationId,omitempty"`
	DesignationName      string `json:"designationName,omitempty"`
	DesignationCode      string `json:"designationCode,omitempty"`
	HierarchyID          string `json:"hierarchyId,omitempty"`
	TravelAgencyMasterID string `json:"travelAgencyMasterId,omitempty"`
	BenefitBundleID      string `json:"benefitBundleId,omitempty"`
	Department           string `json:"department,omitempty"`
	// CreateDate           time.Time
	// UpdateDate           time.Time
}

type AssignDesig struct {
	DesignationName    LabVal   `json:"designationName,omitempty"`
	TravelAgencyUserId []string `json:"travelAgencyUserId,omitempty"`
}

type AssignDesignations struct {
	DesignationName LabVal
	EmpDetails      []EmpDetail
}
type EmpDetail struct {
	TravelAgencyUserId string `json:"travelAgencyUserId,omitempty"`
	VirtualName        string `json:"virtualName,omitempty"`
	Email              string `json:"email,omitempty"`
	PersonalEmail      string `json:"personalEmail,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Designation        string `json:"designation,omitempty"`
}

//editted
type AllEmp struct {
	TravelAgencyUserId string `json:"travelAgencyUserId,omitempty"`
	VirtualName        string `json:"virtualName,omitempty"`
	Email              string `json:"email,omitempty"`
	PersonalEmail      string `json:"personalEmail,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Department         LabVal `json:"department,omitempty"`
	Designation        LabVal `json:"designation,omitempty"`
	BenefitBundle      LabVal `json:"benefitBundle,omitempty"`
}

type Bundle struct {
	BenefitBundle LabVal `json:"benefitBundle"`
	HierarchyId   string `json:"hierarchyId"`
}

type DeleteEmps struct {
	TravelAgencyUserId []string `json:"travelAgencyUserId"`
}

type EditEmploye struct {
	TravelAgencyUserId string `json:"travelAgencyUserId,omitempty"`
	CompanyName        string `json:"companyName,omitempty"`
	VirtualName        string `json:"virtualName,omitempty"`
	Email              string `json:"email,omitempty"`
	PersonalEmail      string `json:"personalEmail,omitempty"`
	HierarchyId        string `json:"hierarchyId,omitempty"`
	Mobile             string `json:"mobile,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Department         LabVal `json:"department,omitempty"`
	Designation        LabVal `json:"designation,omitempty"`
	BenefitBundle      LabVal `json:"benefitBundle,omitempty"`
}

type mapStringScan struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final result
	row      map[string]string
	colCount int
	colNames []string
}

//editted
type AddIndivEmployee struct {
	CompanyID          string `json:"companyId,omitempty"`
	TravelAgencyUserId string `json:"travelAgencyUserId,omitempty"`
	CompanyMail        string `json:"companyMail,omitempty"`
	CompanyName        string `json:"companyName,omitempty"`
	Department         string `json:"department,omitempty"`
	Designation        string `json:"designation,omitempty"`
	EmployeeName       string `json:"employeeName,omitempty"`
	MobileNo           string `json:"mobileNo,omitempty"`
	PersonalMail       string `json:"personalMail,omitempty"`
	StartDate          string `json:"startDate,omitempty"`
	BenefitBundleId    string `json:"benefitBundleId,omitempty"`
}

func MarshalBenefitByDes(bundleId, bundleName, hierarchyId string) string {
	var BundleVar Bundle
	var labval LabVal
	labval = LabVal{
		Label: bundleName,
		Value: bundleId,
	}
	BundleVar = Bundle{
		BenefitBundle: labval,
		HierarchyId:   hierarchyId,
	}
	b, err := json.Marshal(BundleVar)
	checkErr(err)
	return string(b)
}

func UnmarshalAddEmployee(jsonStr string) *AddIndivEmployee {
	res := &AddIndivEmployee{}
	err := json.Unmarshal([]byte(jsonStr), res)
	checkErr(err)
	//fmt.Println(res)
	return res
}

// func UnmarshalAllEmp(jsonStr string) *AllEmp {
// 	res := &AllEmp{}
// 	err := json.Unmarshal([]byte(jsonStr), res)
// 	checkErr(err)
// 	//fmt.Println(res)
// 	return res
// }

//edited
func UnmarshalDeleteEmps(jsonStr string) *DeleteEmps {
	res := &DeleteEmps{}
	err := json.Unmarshal([]byte(jsonStr), res)
	checkErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalAssignDesig(jsonStr string) *AssignDesig {
	res := &AssignDesig{}
	err := json.Unmarshal([]byte(jsonStr), res)
	checkErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalAssignDesignation(jsonStr string) *AssignDesignations {
	res := &AssignDesignations{}
	err := json.Unmarshal([]byte(jsonStr), res)
	checkErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalDesignation(jsonStr string) *Designation {
	res := &Designation{}
	err := json.Unmarshal([]byte(jsonStr), res)
	checkErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalDepartment(jsonStr string) *Department {
	res := &Department{}
	err := json.Unmarshal([]byte(jsonStr), res)
	checkErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalJsonCityCatEdit(jsonStr string) *CityCategoryMapEdit {
	res := &CityCategoryMapEdit{}
	err := json.Unmarshal([]byte(jsonStr), res)
	checkErr(err)
	fmt.Println(res)
	return res
}

func UnmarshalJsonCityCat(jsonStr string) *CityCategoryMap {
	res := &CityCategoryMap{}
	err := json.Unmarshal([]byte(jsonStr), res)
	checkErr(err)
	//fmt.Println(res)
	return res
}

func UnmarshalJsonPolicyBundle(jsonStr string) *PB {
	res := &PB{}
	err := json.Unmarshal([]byte(jsonStr), res)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func NewMapStringScan(columnNames []string) *mapStringScan {
	lenCN := len(columnNames)
	s := &mapStringScan{
		cp:       make([]interface{}, lenCN),
		row:      make(map[string]string, lenCN),
		colCount: lenCN,
		colNames: columnNames,
	}
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
	}
	return s
}

func (s *mapStringScan) Update(rows *sql.Rows) error {
	if err := rows.Scan(s.cp...); err != nil {
		return err
	}

	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {
			s.row[s.colNames[i]] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return fmt.Errorf("Cannot convert index %d column %s to type *sql.RawBytes", i, s.colNames[i])
		}
	}
	return nil
}

func (s *mapStringScan) Get() map[string]string {
	return s.row
}
