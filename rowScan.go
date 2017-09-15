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

type mapStringScan struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final result
	row      map[string]string
	colCount int
	colNames []string
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
