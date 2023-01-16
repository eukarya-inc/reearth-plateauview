package searchindex

import (
	"context"
	"fmt"
	"io"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/eukarya-inc/reearth-plateauview/server/searchindex/indexer"
)

var builtinConfig = &indexer.Config{
	IdProperty: "gml_id",
	Indexes: map[string]indexer.Index{
		"名称":           {Kind: "enum"},
		"用途":           {Kind: "enum"},
		"住所":           {Kind: "enum"},
		"名建物利用現況_中分類称": {Kind: "enum"},
		"建物利用現況_小分類":   {Kind: "enum"},
	},
}

type Indexer struct {
	i      *indexer.Indexer
	config *indexer.Config
	cms    cms.Interface
	pid    string
}

func NewIndexer(cms cms.Interface, pid, base string) *Indexer {
	return &Indexer{
		i:      indexer.NewIndexer(builtinConfig, indexer.NewHTTPFS(nil, base), nil),
		config: builtinConfig,
		cms:    cms,
		pid:    pid,
	}
}

func (i *Indexer) BuildIndex(ctx context.Context, name string) (string, error) {
	res, err := i.i.Build()
	if err != nil {
		return "", fmt.Errorf("failed to build indexes: %v", err)
	}

	pr, pw := io.Pipe()

	aids := make(chan string)
	errs := make(chan error)
	go func() {
		aid, err := i.cms.UploadAssetDirectly(ctx, i.pid, fmt.Sprintf("%s_index.zip", name), pr)
		aids <- aid
		errs <- err
	}()

	zw := indexer.NewZipOutputFS("", pw)
	if err := indexer.NewWriter(i.config, zw).Write(res); err != nil {
		return "", fmt.Errorf("failed to save indexes: %v", err)
	}

	if err := pw.Close(); err != nil {
		return "", fmt.Errorf("failed to save indexes: %v", err)
	}

	aid := <-aids
	err = <-errs
	if err != nil {
		return "", fmt.Errorf("failed to save indexes: %w", err)
	}
	return aid, nil
}
