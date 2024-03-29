package cmsintegrationv2

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestFMEID(t *testing.T) {
	i := fmeID{ItemID: "item", AssetID: "asset", ProjectID: "project"}
	assert.Equal(t, i, lo.Must(parseFMEID(i.String("aaa"), "aaa")))
	_, err := parseFMEID(i.String("aaa"), "aaa2")
	assert.Same(t, ErrInvalidFMEID, err)
}

func TestFMEResult_GetResult(t *testing.T) {
	r, u := fmeResult{
		Results: map[string]any{
			"*":                                         "https://example.com",
			"_dic":                                      "https://example.com/04100_sendai-shi_2022_citygml_1/dic.json",
			"bldg_lod1":                                 "https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod1.zip",
			"bldg_lod2":                                 "https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2.zip",
			"bldg_lod2_low_resolution":                  "https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2_low_resolution.zip",
			"bldg_lod2_no_texture":                      "https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2_no_texture.zip",
			"fld/natl/natorigawa_hirosegawa-2_l1":       "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-2_l1.zip",
			"fld/natl/natorigawa_hirosegawa-2_l2":       "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-2_l2.zip",
			"fld/natl/natorigawa_natorigawa_l1":         "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_natorigawa_l1.zip",
			"fld/natl/natorigawa_natorigawa_l2":         "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_natorigawa_l2.zip",
			"fld/natl/natorigawa_zarugawa_l2":           "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_zarugawa_l2.zip",
			"fld/pref/nanakitagawa_nanakitagawa_l1":     "https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_nanakitagawa_l1.zip",
			"fld/pref/nanakitagawa_nanakitagawa_l2":     "https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_nanakitagawa_l2.zip",
			"fld/pref/nanakitagawa_umedagawa_l1":        "https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_umedagawa_l1.zip",
			"fld/pref/nanakitagawa_umedagawa_l2":        "https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_umedagawa_l2.zip",
			"fld/pref/natorigawa_hirosegawa-1_l1":       "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-1_l1.zip",
			"fld/pref/natorigawa_hirosegawa-1_l2":       "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-1_l2.zip",
			"fld/pref/natorigawa_kyuuzarugawa_l1":       "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_kyuuzarugawa_l1.zip",
			"fld/pref/natorigawa_kyuuzarugawa_l2":       "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_kyuuzarugawa_l2.zip",
			"fld/pref/natorigawa_masudagawa_l1":         "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_masudagawa_l1.zip",
			"fld/pref/natorigawa_masudagawa_l2":         "https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_masudagawa_l2.zip",
			"fld/pref/sunaoshikawa_sunaoshikawa_l2":     "https://example.com/04100_sendai-shi_2022_citygml_1/sunaoshikawa_sunaoshikawa_l2.zip",
			"lsld":                                      "https://example.com/04100_sendai-shi_2022_citygml_1/lsld.zip",
			"tnm/04_1":                                  "https://example.com/04100_sendai-shi_2022_citygml_1/04_1.zip",
			"urf/AreaClassification":                    "https://example.com/04100_sendai-shi_2022_citygml_1/AreaClassification.zip",
			"urf/FirePreventionDistrict":                "https://example.com/04100_sendai-shi_2022_citygml_1/FirePreventionDistrict.zip",
			"urf/HeightControlDistrict":                 "https://example.com/04100_sendai-shi_2022_citygml_1/HeightControlDistrict.zip",
			"urf/HighLevelUseDistrict":                  "https://example.com/04100_sendai-shi_2022_citygml_1/HighLevelUseDistrict.zip",
			"urf/LandscapeZone":                         "https://example.com/04100_sendai-shi_2022_citygml_1/LandscapeZone.zip",
			"urf/ParkingPlaceDevelopmentZone":           "https://example.com/04100_sendai-shi_2022_citygml_1/ParkingPlaceDevelopmentZone.zip",
			"urf/PortZone":                              "https://example.com/04100_sendai-shi_2022_citygml_1/PortZone.zip",
			"urf/ScenicDistrict":                        "https://example.com/04100_sendai-shi_2022_citygml_1/ScenicDistrict.zip",
			"urf/SpecialGreenSpaceConservationDistrict": "https://example.com/04100_sendai-shi_2022_citygml_1/SpecialGreenSpaceConservationDistrict.zip",
			"urf/SpecialUrbanRenaissanceDistrict":       "https://example.com/04100_sendai-shi_2022_citygml_1/SpecialUrbanRenaissanceDistrict.zip",
			"urf/SpecialUseDistrict":                    "https://example.com/04100_sendai-shi_2022_citygml_1/SpecialUseDistrict.zip",
			"urf/UrbanPlanningArea":                     "https://example.com/04100_sendai-shi_2022_citygml_1/UrbanPlanningArea.zip",
			"urf/UseDistrict":                           "https://example.com/04100_sendai-shi_2022_citygml_1/UseDistrict.zip",
			"tran_lod1":                                 "tran_lod1",
			"tran_lod2":                                 "tran_lod2",
			"tran_lod3":                                 "tran_lod3",
			"frn":                                       "frn",
			"veg":                                       "veg",
			"luse":                                      "luse",
			"htd":                                       "htd",
			"ifld":                                      "ifld",
			"unknown":                                   "hoge",
		},
	}.GetResult()
	assert.Equal(t, fmeResultAssets{
		All: "https://example.com",
		Dic: "https://example.com/04100_sendai-shi_2022_citygml_1/dic.json",
		Bldg: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2_low_resolution.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2_no_texture.zip",
		},
		Fld: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_nanakitagawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_nanakitagawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_umedagawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_umedagawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-1_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-1_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-2_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-2_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_kyuuzarugawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_kyuuzarugawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_masudagawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_masudagawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_natorigawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_natorigawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_zarugawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/sunaoshikawa_sunaoshikawa_l2.zip",
		},
		Lsld: "https://example.com/04100_sendai-shi_2022_citygml_1/lsld.zip",
		Tnm: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/04_1.zip",
		},
		Urf: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/AreaClassification.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/FirePreventionDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/HeightControlDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/HighLevelUseDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/LandscapeZone.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/ParkingPlaceDevelopmentZone.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/PortZone.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/ScenicDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/SpecialGreenSpaceConservationDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/SpecialUrbanRenaissanceDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/SpecialUseDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/UrbanPlanningArea.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/UseDistrict.zip",
		},
		Tran: []string{"tran_lod1", "tran_lod2", "tran_lod3"},
		Frn:  "frn",
		Veg:  "veg",
		Luse: "luse",
		Htd:  []string{"htd"},
		Ifld: []string{"ifld"},
	}, r)
	assert.Equal(t, []string{"unknown"}, u)
}

func TestFMEResultAsset_Entries(t *testing.T) {
	a := fmeResultAssets{
		All: "https://example.com",
		Dic: "https://example.com/04100_sendai-shi_2022_citygml_1/dic.json",
		Bldg: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2_low_resolution.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2_no_texture.zip",
		},
		Fld: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_nanakitagawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_nanakitagawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_umedagawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_umedagawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-1_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-1_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-2_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-2_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_kyuuzarugawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_kyuuzarugawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_masudagawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_masudagawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_natorigawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_natorigawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_zarugawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/sunaoshikawa_sunaoshikawa_l2.zip",
		},
		Lsld: "https://example.com/04100_sendai-shi_2022_citygml_1/lsld.zip",
		Tnm: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/04_1.zip",
		},
		Urf: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/AreaClassification.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/FirePreventionDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/HeightControlDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/HighLevelUseDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/LandscapeZone.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/ParkingPlaceDevelopmentZone.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/PortZone.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/ScenicDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/SpecialGreenSpaceConservationDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/SpecialUrbanRenaissanceDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/SpecialUseDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/UrbanPlanningArea.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/UseDistrict.zip",
		},
		Tran: []string{"tran_lod1", "tran_lod2", "tran_lod3"},
		Frn:  "frn",
		Veg:  "veg",
		Luse: "luse",
		Htd:  []string{"htd"},
		Ifld: []string{"ifld"},
	}

	e := []lo.Entry[string, []string]{
		{Key: "all", Value: []string{"https://example.com"}},
		{Key: "dictionary", Value: []string{"https://example.com/04100_sendai-shi_2022_citygml_1/dic.json"}},
		{Key: "bldg", Value: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2_low_resolution.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/bldg_lod2_no_texture.zip",
		}},
		{Key: "tran", Value: []string{
			"tran_lod1",
			"tran_lod2",
			"tran_lod3",
		}},
		{Key: "frn", Value: []string{
			"frn",
		}},
		{Key: "luse", Value: []string{
			"luse",
		}},
		{Key: "veg", Value: []string{
			"veg",
		}},
		{Key: "lsld", Value: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/lsld.zip",
		}},
		{Key: "fld", Value: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_nanakitagawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_nanakitagawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_umedagawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/nanakitagawa_umedagawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-1_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-1_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-2_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_hirosegawa-2_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_kyuuzarugawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_kyuuzarugawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_masudagawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_masudagawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_natorigawa_l1.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_natorigawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/natorigawa_zarugawa_l2.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/sunaoshikawa_sunaoshikawa_l2.zip",
		}},
		{Key: "tnm", Value: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/04_1.zip",
		}},
		{Key: "htd", Value: []string{
			"htd",
		}},
		{Key: "ifld", Value: []string{
			"ifld",
		}},
		{Key: "urf", Value: []string{
			"https://example.com/04100_sendai-shi_2022_citygml_1/AreaClassification.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/FirePreventionDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/HeightControlDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/HighLevelUseDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/LandscapeZone.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/ParkingPlaceDevelopmentZone.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/PortZone.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/ScenicDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/SpecialGreenSpaceConservationDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/SpecialUrbanRenaissanceDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/SpecialUseDistrict.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/UrbanPlanningArea.zip",
			"https://example.com/04100_sendai-shi_2022_citygml_1/UseDistrict.zip",
		}},
	}

	assert.Equal(t, e, a.Entries())
}
