package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

type PB struct {
	BenefitBundleID string
	BundleName      string
	BundleCode      string
	CompanyID       string
	MethType        string
	PolicyBundles   []PolicyBundle
}

type PolicyBundle struct {
	BenefitTypeID        LabVal
	Priority             string
	Benefits             []LabVal
	CityCatAndAllowances []CityCatAndAllowance
}

type CityCatAndAllowance struct {
	Label      string
	Value      string
	LimitSpent string
	Min        string
	Max        string
	Flex       string
	FlexAmt    string
}
type LabVal struct {
	Label string
	Value string
}

type mapStringScan struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final result
	row      map[string]string
	colCount int
	colNames []string
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
