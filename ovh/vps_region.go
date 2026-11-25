package ovh

import "strings"

// vpsLocation is the part of a GET /v2/location item used to find the region
// of a VPS.
type vpsLocation struct {
	Name         string `json:"name"`
	CityCode     string `json:"cityCode"`
	SpecificType string `json:"specificType"`
}

// vpsDatacenterExceptions maps the regions whose VPS datacenter code (API v1,
// order configuration vps_datacenter) is neither the city code nor the region
// name in upper case. Frozen: new regions use the region name in upper case.
var vpsDatacenterExceptions = map[string]string{
	"eu-west-lim": "DE",
	"eu-west-eri": "UK",
}

// vpsRegionFromDatacenterCode returns the region (common.RegionEnum, such as
// eu-west-gra) of an API v1 VPS datacenter code (such as GRA or EU-WEST-RBX).
//
// In order: the exceptions, then the region whose name in upper case is the
// code, then the only non-SNC region with this city code, or the only STANDARD
// one among several. ok is false when no region, or more than one, matches.
func vpsRegionFromDatacenterCode(code string, locations []vpsLocation) (region string, ok bool) {
	for region, exceptionCode := range vpsDatacenterExceptions {
		if exceptionCode == code {
			return region, true
		}
	}

	for _, location := range locations {
		if strings.ToUpper(location.Name) == code {
			return location.Name, true
		}
	}

	var sameCity []vpsLocation
	for _, location := range locations {
		if location.CityCode == code && location.SpecificType != "SNC" {
			sameCity = append(sameCity, location)
		}
	}
	if len(sameCity) == 1 {
		return sameCity[0].Name, true
	}

	region = ""
	for _, location := range sameCity {
		if location.SpecificType == "STANDARD" {
			if region != "" {
				return "", false
			}
			region = location.Name
		}
	}
	return region, region != ""
}
