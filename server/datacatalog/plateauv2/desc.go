package plateauv2

import (
	"path"
	"strings"

	"github.com/reearth/reearthx/util"
)

type Description struct {
	NameOverride   string
	TypeOverride   string
	TypeEnOverride string
	AreaOverride   string
	LayerOverride  []string
	Desc           string
}

func descFromAsset(an AssetName, descs []string) Description {
	if len(descs) == 0 {
		return Description{}
	}

	assetName := an.String()
	fn := strings.TrimSuffix(assetName, path.Ext(assetName))
	for _, desc := range descs {
		b, a, ok := strings.Cut(desc, "\n")
		if ok && strings.Contains(b, fn) {
			tags, rest := extractTags(strings.TrimSpace(a))
			return Description{
				NameOverride:   tags["name"],
				TypeOverride:   tags["type"],
				TypeEnOverride: tags["type_en"],
				AreaOverride:   tags["area"],
				LayerOverride:  multipleValues(tags["layer"]),
				Desc:           rest,
			}
		}
	}

	return Description{}
}

func extractTags(s string) (map[string]string, string) {
	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")
	tags := map[string]string{}

	last := -1
	for i, l := range lines {
		if l != "" && !strings.HasPrefix(l, "@") {
			break
		}

		if l == "" {
			last = i
			continue
		}

		l = strings.TrimSpace(strings.TrimPrefix(l, "@"))
		k, v, found := strings.Cut(l, ":")
		if !found {
			break
		}

		tags[k] = strings.TrimSpace(v)
		last = i
	}

	if last == -1 {
		return tags, s
	}

	rest := strings.TrimSpace(strings.Join(lines[last+1:], "\n"))
	return tags, rest
}

func multipleValues(v string) []string {
	vv := strings.Split(v, ",")
	if len(vv) == 0 || len(vv) == 1 && vv[0] == "" {
		return nil
	}
	return util.Map(vv, strings.TrimSpace)
}
