package datacatalogv3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAssetName(t *testing.T) {
	tests := []struct {
		name string
		args string
		want *AssetName
	}{
		{
			name: "citygml",
			args: "26100_kyoto-shi_city_2022_citygml_3",
			want: &AssetName{
				CityCode:    "26100",
				CityName:    "kyoto-shi",
				Provider:    "city",
				Year:        2022,
				Format:      "citygml",
				UpdateCount: 3,
			},
		},
		{
			name: "citygml with op",
			args: "26100_kyoto-shi_city_2022_citygml_3_op",
			want: &AssetName{
				CityCode:    "26100",
				CityName:    "kyoto-shi",
				Provider:    "city",
				Year:        2022,
				Format:      "citygml",
				UpdateCount: 3,
			},
		},
		{
			name: "bldg",
			args: "26100_kyoto-shi_city_2023_citygml_1_op_bldg_3dtiles_lod2",
			want: &AssetName{
				CityCode:    "26100",
				CityName:    "kyoto-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Normal: &AssetNameExNormal{
						Type:      "bldg",
						Format:    "3dtiles",
						LOD:       2,
						NoTexture: false,
					},
					Ex: "bldg_3dtiles_lod2",
				},
			},
		},
		{
			name: "bldg with ward",
			args: "26100_kyoto-shi_city_2023_citygml_1_op_bldg_3dtiles_26103_sakyo-ku_lod2_no_texture",
			want: &AssetName{
				CityCode:    "26100",
				CityName:    "kyoto-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Normal: &AssetNameExNormal{
						Type:      "bldg",
						Format:    "3dtiles",
						WardCode:  "26103",
						WardName:  "sakyo-ku",
						LOD:       2,
						NoTexture: true,
					},
					Ex: "bldg_3dtiles_26103_sakyo-ku_lod2_no_texture",
				},
			},
		},
		{
			name: "tran",
			args: "33211_bizen-shi_city_2023_citygml_1_op_tran_mvt_lod1",
			want: &AssetName{
				CityCode:    "33211",
				CityName:    "bizen-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Normal: &AssetNameExNormal{
						Type:   "tran",
						Format: "mvt",
						LOD:    1,
					},
					Ex: "tran_mvt_lod1",
				},
			},
		},
		{
			name: "fld",
			args: "40202_omuta-shi_city_2023_citygml_1_op_fld_natl_yabegawa_haegawa_3dtiles_l1",
			want: &AssetName{
				CityCode:    "40202",
				CityName:    "omuta-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Fld: &AssetNameExFld{
						Type:      "fld",
						Admin:     "natl",
						River:     "yabegawa_haegawa",
						Format:    "3dtiles",
						L:         1,
						NoTexture: false,
					},
					Ex: "fld_natl_yabegawa_haegawa_3dtiles_l1",
				},
			},
		},
		{
			name: "fld notexture",
			args: "40202_omuta-shi_city_2023_citygml_1_op_fld_natl_yabegawa_haegawa_3dtiles_l2_no_texture",
			want: &AssetName{
				CityCode:    "40202",
				CityName:    "omuta-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Fld: &AssetNameExFld{
						Type:      "fld",
						Admin:     "natl",
						River:     "yabegawa_haegawa",
						Format:    "3dtiles",
						L:         2,
						NoTexture: true,
					},
					Ex: "fld_natl_yabegawa_haegawa_3dtiles_l2_no_texture",
				},
			},
		},
		{
			name: "tnm",
			args: "40202_omuta-shi_city_2023_citygml_1_op_tnm_40_1_3dtiles",
			want: &AssetName{
				CityCode:    "40202",
				CityName:    "omuta-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Tnm: &AssetNameExTnm{
						Type:   "tnm",
						Name:   "40_1",
						Format: "3dtiles",
					},
					Ex: "tnm_40_1_3dtiles",
				},
			},
		},
		{
			name: "tnm notexture",
			args: "40202_omuta-shi_city_2023_citygml_1_op_tnm_40_1_3dtiles_no_texture",
			want: &AssetName{
				CityCode:    "40202",
				CityName:    "omuta-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Tnm: &AssetNameExTnm{
						Type:      "tnm",
						Name:      "40_1",
						Format:    "3dtiles",
						NoTexture: true,
					},
					Ex: "tnm_40_1_3dtiles_no_texture",
				},
			},
		},
		{
			name: "urf",
			args: "40202_omuta-shi_city_2023_citygml_1_op_urf_AreaClassification_mvt_lod1",
			want: &AssetName{
				CityCode:    "40202",
				CityName:    "omuta-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Urf: &AssetNameExUrf{
						Type:   "urf",
						Name:   "AreaClassification",
						Format: "mvt",
						LOD:    1,
					},
					Ex: "urf_AreaClassification_mvt_lod1",
				},
			},
		},
		{
			name: "urf without lod",
			args: "40202_omuta-shi_city_2023_citygml_1_op_urf_AreaClassification_mvt",
			want: &AssetName{
				CityCode:    "40202",
				CityName:    "omuta-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Urf: &AssetNameExUrf{
						Type:   "urf",
						Name:   "AreaClassification",
						Format: "mvt",
					},
					Ex: "urf_AreaClassification_mvt",
				},
			},
		},
		{
			name: "invalid ex",
			args: "26100_kyoto-shi_city_2023_citygml_1_op_exex",
			want: &AssetName{
				CityCode:    "26100",
				CityName:    "kyoto-shi",
				Provider:    "city",
				Year:        2023,
				Format:      "citygml",
				UpdateCount: 1,
				Ex: AssetNameEx{
					Ex: "exex",
				},
			},
		},
		{
			name: "invalid",
			args: "aaaaa",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ParseAssetName(tt.args))
		})
	}
}

func TestParseAssetURLs(t *testing.T) {
	name := "26100_kyoto-shi_city_2022_citygml_3"
	assert.Equal(t, []*AssetName{
		ParseAssetName(name),
	}, ParseAssetUrls([]string{
		fmt.Sprintf("https://example.com/%s.zip", name),
	}))
}
