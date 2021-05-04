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
	Rn      int
	Code    string
	Formula string
}
