package main

import (
	"fmt"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// generateIdebPDF generates a PDF report from InputJSON data.
func generateIdebPDF(inputData *InputJSON) ([]byte, error) {
	m := maroto.New()

	// Add content to the PDF based on inputData and SARIPARI-PERTIWI-ABADI-b.pdf layout
	m.AddRow(10, col.New(12).Add(text.New("Laporan Informasi Debitur (iDeb)", props.Text{Align: align.Center, Size: 16, Style: fontstyle.Bold})))
	m.AddRow(7, col.New(12).Add(text.New("Otoritas Jasa Keuangan", props.Text{Align: align.Center, Size: 12})))
	m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))

	// Header Information
	m.AddRow(5, col.New(3).Add(text.New("Nama", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.CorporateKeyWord.IdentityNumberName, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("NPWP", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.CorporateDebtors[0].TaxId, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Tempat Pendirian", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.CorporateKeyWord.TestPlace, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Tanggal Akte Pendirian", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.CorporateDebtors[0].EstCertDate, props.Text{Size: 8})))

	m.AddRow(5, col.New(3).Add(text.New("Kode Ref. Pengguna", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Header.UserReferenceCode, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Nomor Laporan", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.ReportNumber, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Posisi Data Terakhir", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.LatestDataYearMonth, props.Text{Size: 8})))
	m.AddRow(5, col.New(3).Add(text.New("Tanggal Permintaan", props.Text{Size: 8})), col.New(9).Add(text.New(inputData.Data.Corporate.RequestDate, props.Text{Size: 8})))

	m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))
	m.AddRow(7, col.New(12).Add(text.New("DATA POKOK DEBITUR", props.Text{Align: align.Center, Size: 12, Style: fontstyle.Bold})))
	m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))

	// Corporate Debtor Information
	for _, debtor := range inputData.Data.Corporate.CorporateDebtors {
		m.AddRow(5, col.New(3).Add(text.New("Nama Lengkap", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.FullName, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("NPWP", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.TaxId, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Bentuk BU / Go Public", props.Text{Size: 8})), col.New(9).Add(text.New(fmt.Sprintf("%s / %s", debtor.CompanyTypeDesc, debtor.GoPublicFlag), props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Tempat Pendirian", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.EstPlace, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("No / Tanggal Akta Pendirian", props.Text{Size: 8})), col.New(9).Add(text.New(fmt.Sprintf("%s / %s", debtor.EstCertNo, debtor.EstCertDate), props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Alamat", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.Address, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Kelurahan", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.SubDistrict, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Kecamatan", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.District, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Kabupaten / Kota", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.CityDesc, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Kode Pos", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.PostalCode, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Negara", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.CountryDesc, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Bidang Usaha", props.Text{Size: 8})), col.New(9).Add(text.New(debtor.EconomicSectorDesc, props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Pelapor / Tanggal Update", props.Text{Size: 8})), col.New(9).Add(text.New(fmt.Sprintf("%s / %s", debtor.MemberDesc, debtor.UpdatedDatetime), props.Text{Size: 8})))
		m.AddRow(5, col.New(3).Add(text.New("Peringkat / Tgl Pemeringkatan", props.Text{Size: 8})), col.New(9).Add(text.New(fmt.Sprintf("%s / %s", debtor.RatingDate, debtor.RatingDate), props.Text{Size: 8})))
		m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))

		// Pemilik / Pengurus
		if len(debtor.OfficisSharehldrsGroups) > 0 {
			m.AddRow(7, col.New(12).Add(text.New("Pemilik / Pengurus", props.Text{Align: align.Center, Size: 12, Style: fontstyle.Bold})))
			m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))

			for _, group := range debtor.OfficisSharehldrsGroups {
				m.AddRow(5, col.New(12).Add(text.New("Pelapor : "+group.MemberDesc, props.Text{Size: 9, Style: fontstyle.Bold})))
				for _, shareholder := range group.OfficisSharehldrs {
					m.AddRow(5, col.New(3).Add(text.New("Nama", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.IdentityNumberName, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("No. Identitas", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.IdentityNumber, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("Jenis Kelamin", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.GenderDesc, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("Jabatan", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.JobPositionDesc, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("Pangsa Kepemilikan", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.ShareOwnership+"%", props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("Alamat", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.Address, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("Kelurahan", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.SubDistrict, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("Kecamatan", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.District, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("Kabupaten / Kota", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.CityDesc, props.Text{Size: 8})))
					m.AddRow(5, col.New(3).Add(text.New("Status Pengurus/Pemilik", props.Text{Size: 8})), col.New(9).Add(text.New(shareholder.ShareholderStatusDesc, props.Text{Size: 8})))
					m.AddRow(10, col.New(12).Add(text.New("", props.Text{})))
				}
			}
		}
		m.AddRow(15, col.New(12)) // Add some space between debtor sections
	}

	// Output the PDF
	pdf, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("Error generating PDF: %w", err)
	}

	return pdf.GetBytes(), nil
}
