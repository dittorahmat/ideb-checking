package main

import (
	"fmt"
	"time"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

// generateIdebPDF generates a PDF report from InputJSON data.
func generateIdebPDF(inputData *InputJSON) ([]byte, error) {
	// Create Maroto instance with better page configuration
	m := maroto.New()

	// Add header with logo space and title
	m.AddRow(20,
		col.New(2), // Left margin
		col.New(8).Add(text.New("LAPORAN INFORMASI DEBITUR (IDEB)", props.Text{
			Align:  align.Center,
			Size:   14,
			Style:  fontstyle.Bold,
		})),
		col.New(2), // Right margin
	)
	
	m.AddRow(15,
		col.New(2), // Left margin
		col.New(8).Add(text.New("OTORITAS JASA KEUANGAN", props.Text{
			Align:  align.Center,
			Size:   12,
		})),
		col.New(2), // Right margin
	)
	
	// Add a line separator
	m.AddRow(5, col.New(12).Add(line.New(props.Line{
		Color: &props.Color{
			Red:   0,
			Green: 0,
			Blue:  0,
		},
	})))
	m.AddRow(5)

	// Header Information in a more structured format with better spacing
	// Corporate Information Section
	m.AddRow(8, col.New(12).Add(text.New("INFORMASI PERUSAHAAN", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	})))
	m.AddRow(2, col.New(12).Add(line.New(props.Line{
		Color: &props.Color{
			Red:   0,
			Green: 0,
			Blue:  0,
		},
	})))
	
	// Format date for better readability
	requestDate := inputData.Data.Corporate.RequestDate
	if len(requestDate) >= 8 {
		year := requestDate[0:4]
		month := requestDate[4:6]
		day := requestDate[6:8]
		requestDate = fmt.Sprintf("%s-%s-%s", day, month, year)
	}
	
	// Parse the date to format it properly
	if t, err := time.Parse("20060102", inputData.Data.Corporate.RequestDate[0:8]); err == nil {
		requestDate = t.Format("02-01-2006")
	}

	m.AddRow(6,
		col.New(4).Add(text.New("Nama Perusahaan", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
		})),
		col.New(8).Add(text.New(": "+inputData.Data.Corporate.CorporateKeyWord.IdentityNumberName, props.Text{
			Size: 9,
		})),
	)
	
	m.AddRow(6,
		col.New(4).Add(text.New("NPWP", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
		})),
		col.New(8).Add(text.New(": "+inputData.Data.Corporate.CorporateDebtors[0].TaxId, props.Text{
			Size: 9,
		})),
	)
	
	m.AddRow(6,
		col.New(4).Add(text.New("Tempat Pendirian", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
		})),
		col.New(8).Add(text.New(": "+inputData.Data.Corporate.CorporateKeyWord.TestPlace, props.Text{
			Size: 9,
		})),
	)
	
	// Format establishment date
	estDate := inputData.Data.Corporate.CorporateDebtors[0].EstCertDate
	if len(estDate) >= 8 {
		if t, err := time.Parse("20060102", estDate[0:8]); err == nil {
			estDate = t.Format("02-01-2006")
		}
	}
	
	m.AddRow(6,
		col.New(4).Add(text.New("Tanggal Akte Pendirian", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
		})),
		col.New(8).Add(text.New(": "+estDate, props.Text{
			Size: 9,
		})),
	)
	
	// Add a separator
	m.AddRow(10)
	
	// Report Information Section
	m.AddRow(8, col.New(12).Add(text.New("INFORMASI LAPORAN", props.Text{
		Size:  11,
		Style: fontstyle.Bold,
	})))
	m.AddRow(2, col.New(12).Add(line.New(props.Line{
		Color: &props.Color{
			Red:   0,
			Green: 0,
			Blue:  0,
		},
	})))
	
	m.AddRow(6,
		col.New(4).Add(text.New("Kode Ref. Pengguna", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
		})),
		col.New(8).Add(text.New(": "+inputData.Data.Header.UserReferenceCode, props.Text{
			Size: 9,
		})),
	)
	
	m.AddRow(6,
		col.New(4).Add(text.New("Nomor Laporan", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
		})),
		col.New(8).Add(text.New(": "+inputData.Data.Corporate.ReportNumber, props.Text{
			Size: 9,
		})),
	)
	
	// Format latest data date
	latestDataDate := inputData.Data.Corporate.LatestDataYearMonth
	if len(latestDataDate) >= 6 {
		year := latestDataDate[0:4]
		month := latestDataDate[4:6]
		latestDataDate = fmt.Sprintf("%s-%s", month, year)
	}
	
	m.AddRow(6,
		col.New(4).Add(text.New("Posisi Data Terakhir", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
		})),
		col.New(8).Add(text.New(": "+latestDataDate, props.Text{
			Size: 9,
		})),
	)
	
	m.AddRow(6,
		col.New(4).Add(text.New("Tanggal Permintaan", props.Text{
			Size:  9,
			Style: fontstyle.Bold,
		})),
		col.New(8).Add(text.New(": "+requestDate, props.Text{
			Size: 9,
		})),
	)
	
	// Add a separator
	m.AddRow(15)
	
	// Data Pokok Debitur Section
	m.AddRow(10, col.New(12).Add(text.New("DATA POKOK DEBITUR", props.Text{
		Align: align.Center,
		Size:  12,
		Style: fontstyle.Bold,
	})))
	m.AddRow(10)

	// Corporate Debtor Information with improved formatting
	for i, debtor := range inputData.Data.Corporate.CorporateDebtors {
		// Add section header for each debtor with a line separator
		m.AddRow(2, col.New(12).Add(line.New(props.Line{
			Color: &props.Color{
				Red:   0,
				Green: 0,
				Blue:  0,
			},
		})))
		m.AddRow(8,
			col.New(12).Add(text.New(fmt.Sprintf("DEBITUR %d", i+1), props.Text{
				Size:  11,
				Style: fontstyle.Bold,
			})),
		)
		m.AddRow(2, col.New(12).Add(line.New(props.Line{
			Color: &props.Color{
				Red:   0,
				Green: 0,
				Blue:  0,
			},
		})))
		m.AddRow(7)

		// Debtor details in a structured format
		m.AddRow(6,
			col.New(4).Add(text.New("Nama Lengkap", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.FullName, props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("NPWP", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.TaxId, props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("Bentuk BU / Go Public", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(fmt.Sprintf(": %s / %s", debtor.CompanyTypeDesc, mapGoPublicFlag(debtor.GoPublicFlag)), props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("Tempat Pendirian", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.EstPlace, props.Text{
				Size: 9,
			})),
		)
		
		// Format establishment certificate date
		estCertDate := debtor.EstCertDate
		if len(estCertDate) >= 8 {
			if t, err := time.Parse("20060102", estCertDate[0:8]); err == nil {
				estCertDate = t.Format("02-01-2006")
			}
		}
		
		m.AddRow(6,
			col.New(4).Add(text.New("No / Tanggal Akte Pendirian", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(fmt.Sprintf(": %s / %s", debtor.EstCertNo, estCertDate), props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("Alamat", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.Address, props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("Kelurahan", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.SubDistrict, props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("Kecamatan", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.District, props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("Kabupaten / Kota", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.CityDesc, props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("Kode Pos", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.PostalCode, props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("Negara", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.CountryDesc, props.Text{
				Size: 9,
			})),
		)
		
		m.AddRow(6,
			col.New(4).Add(text.New("Bidang Usaha", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(": "+debtor.EconomicSectorDesc, props.Text{
				Size: 9,
			})),
		)
		
		// Format updated datetime
		updatedDatetime := debtor.UpdatedDatetime
		if len(updatedDatetime) >= 8 {
			if t, err := time.Parse("20060102", updatedDatetime[0:8]); err == nil {
				updatedDatetime = t.Format("02-01-2006")
			}
		}
		
		m.AddRow(6,
			col.New(4).Add(text.New("Pelapor / Tanggal Update", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(fmt.Sprintf(": %s / %s", debtor.MemberDesc, updatedDatetime), props.Text{
				Size: 9,
			})),
		)
		
		// Format rating date
		ratingDate := debtor.RatingDate
		if len(ratingDate) >= 8 {
			if t, err := time.Parse("20060102", ratingDate[0:8]); err == nil {
				ratingDate = t.Format("02-01-2006")
			}
		}
		
		m.AddRow(6,
			col.New(4).Add(text.New("Peringkat / Tgl Pemeringkatan", props.Text{
				Size:  9,
				Style: fontstyle.Bold,
			})),
			col.New(8).Add(text.New(fmt.Sprintf(": %s / %s", "-", ratingDate), props.Text{
				Size: 9,
			})),
		)

		// Pemilik / Pengurus section
		if len(debtor.OfficisSharehldrsGroups) > 0 {
			m.AddRow(12)
			m.AddRow(8, col.New(12).Add(text.New("PEMILIK / PENGURUS", props.Text{
				Align: align.Center,
				Size:  12,
				Style: fontstyle.Bold,
			})))
			m.AddRow(2, col.New(12).Add(line.New(props.Line{
				Color: &props.Color{
					Red:   0,
					Green: 0,
					Blue:  0,
				},
			})))
			m.AddRow(8)

			for _, group := range debtor.OfficisSharehldrsGroups {
				m.AddRow(7,
					col.New(12).Add(text.New("Pelapor : "+group.MemberDesc, props.Text{
						Size:  10,
						Style: fontstyle.Bold,
					})),
				)
				m.AddRow(6)

				// Display shareholders information
				for j, shareholder := range group.OfficisSharehldrs {
					// Add section header for each shareholder
					m.AddRow(7,
						col.New(12).Add(text.New(fmt.Sprintf("Pemilik/Pengurus %d", j+1), props.Text{
							Size:  10,
							Style: fontstyle.Bold,
						})),
					)
					m.AddRow(5)
					
					m.AddRow(6,
						col.New(4).Add(text.New("Nama", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+shareholder.IdentityNumberName, props.Text{
							Size: 9,
						})),
					)
					
					m.AddRow(6,
						col.New(4).Add(text.New("No. Identitas", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+shareholder.IdentityNumber, props.Text{
							Size: 9,
						})),
					)
					
					m.AddRow(6,
						col.New(4).Add(text.New("Jenis Kelamin", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+shareholder.GenderDesc, props.Text{
							Size: 9,
						})),
					)
					
					m.AddRow(6,
						col.New(4).Add(text.New("Jabatan", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+shareholder.JobPositionDesc, props.Text{
							Size: 9,
						})),
					)
					
					m.AddRow(6,
						col.New(4).Add(text.New("Pangsa Kepemilikan", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+shareholder.ShareOwnership+"%", props.Text{
							Size: 9,
						})),
					)
					
					m.AddRow(6,
						col.New(4).Add(text.New("Alamat", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+shareholder.Address, props.Text{
							Size: 9,
						})),
					)
					
					m.AddRow(6,
						col.New(4).Add(text.New("Kelurahan", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+shareholder.SubDistrict, props.Text{
							Size: 9,
						})),
					)
					
					m.AddRow(6,
						col.New(4).Add(text.New("Kecamatan", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+shareholder.District, props.Text{
							Size: 9,
						})),
					)
					
					m.AddRow(6,
						col.New(4).Add(text.New("Kabupaten / Kota", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+shareholder.CityDesc, props.Text{
							Size: 9,
						})),
					)
					
					m.AddRow(6,
						col.New(4).Add(text.New("Status Pengurus/Pemilik", props.Text{
							Size:  9,
							Style: fontstyle.Bold,
						})),
						col.New(8).Add(text.New(": "+mapShareholderStatus(shareholder.ShareholderStatusDesc), props.Text{
							Size: 9,
						})),
					)
					m.AddRow(8)
				}
			}
		}
		
		// Add separator between debtors except for the last one
		if i < len(inputData.Data.Corporate.CorporateDebtors)-1 {
			m.AddRow(15)
			m.AddRow(1, col.New(12).Add(line.New(props.Line{
				Color: &props.Color{
					Red:   0,
					Green: 0,
					Blue:  0,
				},
			})))
			m.AddRow(15)
		}
	}

	// Output the PDF
	pdf, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("Error generating PDF: %w", err)
	}

	return pdf.GetBytes(), nil
}

// mapGoPublicFlag maps the go public flag to a more readable format
func mapGoPublicFlag(flag string) string {
	switch flag {
	case "T":
		return "Tidak"
	case "Y":
		return "Ya"
	default:
		return flag
	}
}

// mapShareholderStatus maps the shareholder status to a more readable format
func mapShareholderStatus(status string) string {
	switch status {
	case "1":
		return "Aktif"
	case "0":
		return "Tidak Aktif"
	default:
		return status
	}
}
