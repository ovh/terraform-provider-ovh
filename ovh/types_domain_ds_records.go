package ovh

type DsRecordAlgorithm int

var DsRecordAlgorithmValuesMap = map[string]DsRecordAlgorithm{
	"RSASHA1":            5,
	"RSASHA1_NSEC3_SHA1": 7,
	"RSASHA256":          8,
	"RSASHA512":          10,
	"ECDSAP256SHA256":    13,
	"ECDSAP384SHA384":    14,
	"ED25519":            15,
}

var DsRecordAlgorithmLabelsMap = map[DsRecordAlgorithm]string{
	5:  "RSASHA1",
	7:  "RSASHA1_NSEC3_SHA1",
	8:  "RSASHA256",
	10: "RSASHA512",
	13: "ECDSAP256SHA256",
	14: "ECDSAP384SHA384",
	15: "ED25519",
}

type DsRecordFlag int

var DsRecordFlagValuesMap = map[string]DsRecordFlag{
	"ZONE_SIGNING_KEY": 256,
	"KEY_SIGNING_KEY":  257,
}

var DsRecordFlagLabelsMap = map[DsRecordFlag]string{
	256: "ZONE_SIGNING_KEY",
	257: "KEY_SIGNING_KEY",
}

type DomainDsRecords struct {
	Domain    string           `json:"domain"`
	DsRecords []DomainDsRecord `json:"ds_records"`
}

func (v DomainDsRecords) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["domain"] = v.Domain
	obj["ds_records"] = []map[string]interface{}{}

	for _, dsRecord := range v.DsRecords {
		obj["ds_records"] = append(obj["ds_records"].([]map[string]interface{}), dsRecord.ToMap())
	}

	return obj
}

type DomainDsRecord struct {
	Id        int               `json:"id,omitempty"`
	Algorithm DsRecordAlgorithm `json:"algorithm"`
	Flags     DsRecordFlag      `json:"flags"`
	PublicKey string            `json:"publickey"`
	Tag       int               `json:"tag"`
}

func (v DomainDsRecord) ToMap() map[string]interface{} {
	obj := make(map[string]interface{})

	obj["algorithm"] = DsRecordAlgorithmLabelsMap[v.Algorithm]
	obj["flags"] = DsRecordFlagLabelsMap[v.Flags]
	obj["public_key"] = v.PublicKey
	obj["tag"] = v.Tag

	return obj
}

type DomainDsRecordsUpdateOpts struct {
	DsRecords []DomainDsRecord `json:"keys"`
}
