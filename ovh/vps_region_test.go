package ovh

import "testing"

// vpsTestLocations is an extract of GET /v2/location: the regions that share a
// city code, the ones API v1 names with a legacy code, and a recent region.
var vpsTestLocations = []vpsLocation{
	{Name: "eu-west-gra", CityCode: "GRA", SpecificType: "STANDARD"},
	{Name: "eu-west-gra-snc", CityCode: "GRA", SpecificType: "SNC"},
	{Name: "eu-west-lim", CityCode: "LIM", SpecificType: "STANDARD"},
	{Name: "eu-west-eri", CityCode: "ERI", SpecificType: "STANDARD"},
	{Name: "eu-south-mil", CityCode: "MIL", SpecificType: "STANDARD"},
	{Name: "eu-south-lz-mil", CityCode: "MIL", SpecificType: "LZ"},
	{Name: "eu-west-rbx", CityCode: "RBX", SpecificType: "STANDARD"},
}

func TestVPSRegionFromDatacenterCode(t *testing.T) {
	cases := []struct {
		code   string
		region string
		ok     bool
	}{
		{"GRA", "eu-west-gra", true},         // the SNC region is left aside
		{"DE", "eu-west-lim", true},          // exception
		{"UK", "eu-west-eri", true},          // exception
		{"EU-WEST-RBX", "eu-west-rbx", true}, // region name in upper case
		{"EU-SOUTH-MIL", "eu-south-mil", true},
		{"MIL", "eu-south-mil", true}, // STANDARD preferred to LZ
		{"XYZ", "", false},
	}

	for _, c := range cases {
		region, ok := vpsRegionFromDatacenterCode(c.code, vpsTestLocations)
		if region != c.region || ok != c.ok {
			t.Errorf("code %q: got (%q, %t), want (%q, %t)", c.code, region, ok, c.region, c.ok)
		}
	}
}

func TestVPSRegionFromDatacenterCodeAmbiguous(t *testing.T) {
	// Two non-SNC regions share the city code and none is STANDARD: no guess.
	locations := []vpsLocation{
		{Name: "eu-west-lz-abc", CityCode: "ABC", SpecificType: "LZ"},
		{Name: "eu-west-abc-bkp", CityCode: "ABC", SpecificType: "BACKUP"},
	}
	if region, ok := vpsRegionFromDatacenterCode("ABC", locations); ok {
		t.Errorf("code %q: got (%q, true), want no region", "ABC", region)
	}
}
