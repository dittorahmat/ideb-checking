package main

import "gorm.io/gorm"

type Request struct {
	gorm.Model
	NomorReferensiPengguna       string `json:"nomor_referensi_pengguna"`
	TujuanPenggunaan             string `json:"tujuan_penggunaan"`
	JenisIdentitas               string `json:"jenis_identitas"`
	NomorIdentitas               string `json:"nomor_identitas"`
	PermintaanFasilitasOutstanding bool   `json:"permintaan_fasilitas_outstanding"`
	SearchType                   string `json:"search_type"`
	StatusAksi                   string `json:"status_aksi" gorm:"default:'Dalam Proses'"`
}

func (Request) TableName() string {
	return "getDebtorExactIndividual"
}

type CorporateRequest struct {
	gorm.Model
	NomorReferensiPengguna       string `json:"nomor_referensi_pengguna"`
	TujuanPenggunaan             string `json:"tujuan_penggunaan"`
	NomorIdentitas               string `json:"nomor_identitas"`
	PermintaanFasilitasOutstanding bool   `json:"permintaan_fasilitas_outstanding"`
	SearchType                   string `json:"search_type"`
	StatusAksi                   string `json:"status_aksi" gorm:"default:'Dalam Proses'"`
}

func (CorporateRequest) TableName() string {
	return "getDebtorExactCorporate"
}
