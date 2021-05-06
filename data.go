package main

type Config struct {
	Devices map[string]Device
	Port    string
}

type Device struct {
	Rn      int
	Uid     string
	Crn     int
	Address int
	Schema  []Param
}

type Param struct {
	Rn      int    `json:"rn"`
	Code    string `json:"code"`
	Formula string `json:"formula"`
}

type Project struct {
	Rn     int    `json:"rn"`
	Code   string `json:"code"`
	Params []struct {
		Prn int     `json:"prn"`
		Val float64 `json:"val"`
	} `json:"params"`
}

type Address struct {
	Rn      int    `json:"rn"`
	Project int    `json:"prn"`
	Code    string `json:"code"`
	Params  []struct {
		Prn int     `json:"prn"`
		Val float64 `json:"val"`
	} `json:"params"`
}

func Difference(slice1, slice2 []interface{}) []interface{} {
	var diff []interface{}
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}
	return diff
}

func Equal(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
