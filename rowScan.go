package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type PB struct {
	BenefitBundleID string         `json:"benefitBundleId"`
	BundleName      string         `json:"bundleName"`
	BundleCode      string         `json:"bundleCode"`
	CompanyID       string         `json:"companyId"`
	MethType        string         `json:"methType"`
	PolicyBundles   []PolicyBundle `json:"policybundles"`
}

type PolicyBundle struct {
	BenefitTypeID        LabVal                `json:"benefitTypeId"`
	Priority             string                `json:"priority"`
	Benefits             []LabVal              `json:"benefits"`
	CityCatAndAllowances []CityCatAndAllowance `json:"cityCatAndAllowances"`
}

type CityCatAndAllowance struct {
	Label      string `json:"label"`
	Value      string `json:"value"`
	LimitSpent string `json:"limitSpent"`
	Min        string `json:"min"`
	Max        string `json:"max"`
	Flex       string `json:"flex"`
	FlexAmt    string `json:"flexAmt"`
	StarCat    string `json:"starCat"`
}

type CityCategoryMap struct {
	CompanyID string   `json:"companyId"`
	CityCat   LabVal   `json:"cityCat"`
	Cities    []LabVal `json:"cities"`
}

type LabVal struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type mapStringScan struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final result
	row      map[string]string
	colCount int
	colNames []string
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
