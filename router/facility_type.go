package main

import (
	"strings"
	"math/rand"
)

type FacilityType = int

const (
	TR FacilityType = iota // Represents facility types as enum of ints starting from 0
	Lab
	Theatre
)

// Extract facility type from confirmationId or facility name, if invalid facility
// route to a random server to evenly distribute invalid requests
func extractFacilityType(str string) FacilityType {
	if strings.Contains(str, "TR") {
		return TR
	} else if strings.Contains(str, "LAB") {
		return Lab
	} else if strings.Contains(str, "THEATRE") {
		return Theatre
	} 
	facilities := []FacilityType{TR, Lab, Theatre}
	return facilities[rand.Intn(len(facilities))]
}