package main

import "gorm.io/gorm"

const (
	SearchTypeInternal = "internal"
	SearchTypeLive     = "live"
)

type BaseRequest struct {
	gorm.Model
	NomorReferensiPengguna         string `json:"nomor_referensi_pengguna"`
	TujuanPenggunaan               string `json:"tujuan_penggunaan"`
	NomorIdentitas                 string `json:"nomor_identitas"`
	PermintaanFasilitasOutstanding bool   `json:"permintaan_fasilitas_outstanding"`
	SearchType                     string `json:"search_type"`
	StatusAksi                     string `json:"status_aksi" gorm:"default:'Dalam Proses'"`
}

type Request struct {
	BaseRequest
	JenisIdentitas string `json:"jenis_identitas"`
}

func (Request) TableName() string {
	return "getDebtorExactIndividual"
}

type CorporateRequest struct {
	BaseRequest
}

func (CorporateRequest) TableName() string {
	return "getDebtorExactCorporate"
}

// CombinedRequestPayload represents a unified payload for both individual and corporate requests.
// It's used to unmarshal the incoming JSON in a single pass.

type CombinedRequestPayload struct {
	NomorReferensiPengguna         string `json:"nomor_referensi_pengguna"`
	TujuanPenggunaan               string `json:"tujuan_penggunaan"`
	NomorIdentitas                 string `json:"nomor_identitas"`
	PermintaanFasilitasOutstanding bool   `json:"permintaan_fasilitas_outstanding"`
	SearchType                     string `json:"search_type"`
	RequestType                    string `json:"request_type"`
	JenisIdentitas                 string `json:"jenis_identitas,omitempty"` // Only for individual
}

type GetIdeb struct {
	gorm.Model
	NomorReferensiPengguna string `json:"nomor_referensi_pengguna"`
	NomorIdentitas         string `json:"nomor_identitas"`
	Data                   string `json:"data"` // Store the full JSON as a string
}

func (GetIdeb) TableName() string {
	return "get_idebs"
}

// Structs for parsing input.json
type InputJSON struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Remark  string `json:"remark"`
	RawData []byte `json:"-"` // To store the raw JSON data
	Data    struct {
		Header    Header    `json:"header"`
		Corporate Corporate `json:"corporate"`
	} `json:"data"`
}

type Header struct {
	UserReferenceCode        string `json:"userReferenceCode"`
	ResultDate               string `json:"resultDate"`
	InquiryId                string `json:"inquiryId"`
	InquiryUserId            string `json:"inquiryUserId"`
	InquiryCreatedBy         string `json:"inquiryCreatedBy"`
	InquiryMemberCode        string `json:"inquiryMemberCode"`
	InquiryOfficeCode        string `json:"inquiryOfficeCode"`
	ReportRequestPurposeCode string `json:"reportRequestPurposeCode"`
	InquiryDate              string `json:"inquiryDate"`
	DataSetTotal             string `json:"dataSetTotal"`
	DataSetNumber            string `json:"dataSetNumber"`
}

type Corporate struct {
	ReportNumber        string `json:"reportNumber"`
	LatestDataYearMonth string `json:"latestDataYearMonth"`
	RequestDate         string `json:"requestDate"`
	CorporateKeyWord    struct {
		IdentityNumberName string `json:"identityNumberName"`
		TestPlace          string `json:"testPlace"`
		RecordStatusFlag   string `json:"recordStatusFlag"`
	} `json:"corporateKeyWord"`
	CorporateDebtors []CorporateDebtor `json:"corporateDebtors"`
}

type CorporateDebtor struct {
	IdentityNumberName      string                   `json:"identityNumberName"`
	FullName                string                   `json:"fullName"`
	TaxId                   string                   `json:"taxId"`
	CompanyType             string                   `json:"companyType"`
	CompanyTypeDesc         string                   `json:"companyTypeDesc"`
	EstPlace                string                   `json:"estPlace"`
	EstCertNo               string                   `json:"estCertNo"`
	EstCertDate             string                   `json:"estCertDate"`
	Member                  string                   `json:"member"`
	MemberDesc              string                   `json:"memberDesc"`
	UpdatedDatetime         string                   `json:"updatedDatetime"`
	Address                 string                   `json:"address"`
	SubDistrict             string                   `json:"subDistrict"`
	District                string                   `json:"district"`
	City                    string                   `json:"city"`
	CityDesc                string                   `json:"cityDesc"`
	PostalCode              string                   `json:"postalCode"`
	Country                 string                   `json:"country"`
	CountryDesc             string                   `json:"countryDesc"`
	LatestAddCertNo         string                   `json:"latestAddCertNo"`
	LatestAddCertDate       string                   `json:"latestAddCertDate"`
	EconomicSector          string                   `json:"economicSector"`
	EconomicSectorDesc      string                   `json:"economicSectorDesc"`
	RatingDate              string                   `json:"ratingDate"`
	CreatedDatetime         string                   `json:"createdDatetime"`
	GoPublicFlag            string                   `json:"goPublicFlag"`
	OfficisSharehldrsGroups []OfficisSharehldrsGroup `json:"officisSharehldrsGroups"`
}

type OfficisSharehldrsGroup struct {
	Member            string              `json:"member"`
	MemberDesc        string              `json:"memberDesc"`
	OfficisSharehldrs []OfficisSharehldrs `json:"officisSharehldrs"`
}

type OfficisSharehldrs struct {
	IdentityNumberName    string `json:"identityNumberName"`
	IdentityNumber        string `json:"identityNumber"`
	Gender                string `json:"gender"`
	GenderDesc            string `json:"genderDesc"`
	JobPosition           string `json:"jobPosition"`
	JobPositionDesc       string `json:"jobPositionDesc"`
	ShareOwnership        string `json:"shareOwnership"`
	Address               string `json:"address"`
	District              string `json:"district"`
	City                  string `json:"city"`
	CityDesc              string `json:"cityDesc"`
	ShareholderStatus     string `json:"shareholderStatus"`
	ShareholderStatusDesc string `json:"shareholderStatusDesc"`
	SubDistrict           string `json:"subDistrict"`
}
