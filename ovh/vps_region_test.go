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

// vpsTestCatalogDatacenters is the vps_datacenter configuration of plan
// vps-2025-model1 in the public order catalog (FR, 2026-11-25).
var vpsTestCatalogDatacenters = []string{"BHS", "DE", "GRA", "SBG", "SGP", "SYD", "UK", "WAW", "YNM", "EU-SOUTH-MIL", "EU-WEST-RBX"}

func TestVPSDatacenterCodeForRegion(t *testing.T) {
	cases := []struct {
		region  string
		allowed []string
		code    string
		ok      bool
	}{
		{"eu-west-gra", []string{"GRA", "SBG"}, "GRA", true}, // city code
		{"eu-west-rbx", []string{"EU-WEST-RBX"}, "EU-WEST-RBX", true},
		{"eu-west-rbx", vpsTestCatalogDatacenters, "EU-WEST-RBX", true}, // region name before city code
		{"eu-west-lim", vpsTestCatalogDatacenters, "DE", true},          // exception
		{"eu-west-eri", vpsTestCatalogDatacenters, "UK", true},          // exception
		{"eu-west-gra", []string{"SBG"}, "", false},                     // not orderable in this plan
		{"eu-west-xyz", vpsTestCatalogDatacenters, "", false},           // unknown region
		// A city code shared with another region is refused: GRA reads back
		// as eu-west-gra, MIL as eu-south-mil.
		{"eu-west-gra-snc", []string{"GRA"}, "", false},
		{"eu-south-lz-mil", []string{"MIL"}, "", false},
		{"eu-south-mil", []string{"MIL"}, "MIL", true},
	}

	for _, c := range cases {
		code, ok := vpsDatacenterCodeForRegion(c.region, vpsTestLocations, c.allowed)
		if code != c.code || ok != c.ok {
			t.Errorf("region %q, allowed %v: got (%q, %t), want (%q, %t)", c.region, c.allowed, code, ok, c.code, c.ok)
		}
		// Both ways agree: the code reads back as the region.
		if ok {
			if back, _ := vpsRegionFromDatacenterCode(code, vpsTestLocations); back != c.region {
				t.Errorf("region %q: code %q reads back as %q", c.region, code, back)
			}
		}
	}
}
