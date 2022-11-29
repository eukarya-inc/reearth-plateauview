package indexer

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/qmuntal/gltf"
	b3dms "github.com/reearth/go3dtiles/b3dm"
	tiles "github.com/reearth/go3dtiles/tileset"
	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/util"
	"gonum.org/v1/gonum/mat"
)

type Indexer struct{}

func NewIndexer() *Indexer {
	return &Indexer{}
}

type Fs interface {
	Open(string) (fs.File, error)
}

type ResultData []map[string]string

func (indexer *Indexer) GenerateIndexes(config *Config, tilesetPath string, fsys Fs) (indexBuilders []IndexBuilder, resultData ResultData, errMsg error) {
	basePath := strings.Split(tilesetPath, "tileset.json")[0]
	ts, err := fsys.Open(tilesetPath)
	if err != nil {
		errMsg = fmt.Errorf("failed to open the tileset: %w", err)
		return
	}
	reader := tiles.NewTilsetReader(ts)
	tileset := new(tiles.Tileset)
	if err := reader.Decode(tileset); err != nil {
		errMsg = fmt.Errorf("failed to decode the tileset: %w", err)
		return
	}

	for property, config := range config.Indexes {
		indexBuilders = append(indexBuilders, createIndexBuilder(property, config))
	}

	features, err := ReadTilesetFeatures(tileset, config, basePath, fsys)
	if err != nil {
		errMsg = fmt.Errorf("failed to read features: %w", err)
		return
	}

	featureCount := len(features)
	log.Debugln("Number of features counted: ", featureCount)

	for idValue, tilsetFeature := range features {
		// taking all positionProperties map entries as string for better writing experience
		positionProperties := map[string]string{
			config.IdProperty: idValue,
			"Longitude":       strconv.FormatFloat(roundFloat(toDegrees(tilsetFeature.Position.Longitude), 5), 'g', -1, 64),
			"Latitude":        strconv.FormatFloat(roundFloat(toDegrees(tilsetFeature.Position.Latitude), 5), 'g', -1, 64),
			"Height":          strconv.FormatFloat(roundFloat(tilsetFeature.Position.Height, 3), 'g', -1, 64),
		}
		resultData = append(resultData, positionProperties)
		length := len(resultData)
		dataRowId := length - 1
		for _, b := range indexBuilders {
			switch t := b.(type) {
			case EnumIndexBuilder:
				if val, ok := tilsetFeature.Properties[t.Property]; ok && val != nil {
					t.AddIndexValue(dataRowId, val.(string))
				}
			default:
				continue
			}
		}
	}
	return
}

type TilesetFeature struct {
	Properties map[string]interface{}
	Position   Cartographic
}

func ReadTilesetFeatures(ts *tiles.Tileset, config *Config, basePath string, fsys Fs) (map[string]TilesetFeature, error) {
	uniqueFeatures := make(map[string]TilesetFeature)
	tilesetQueue := []*tiles.Tileset{ts}

	for _, tileset := range tilesetQueue {

		tilesetIterFn := func(tile *tiles.Tile, computedTransform *mat.Dense) error {
			tileUri, err := tile.Uri()
			if err != nil {
				return fmt.Errorf("failed to fetch uri of tile: %v", err)
			}
			contentPath := filepath.Join(basePath, tileUri)
			log.Debugln(tileUri)
			if strings.HasSuffix(tileUri, ".json") {
				childTileset, _ := tiles.Open(tileUri)
				tilesetQueue = append(tilesetQueue, childTileset)
				return nil
			}

			b3dmFile, err := fsys.Open(contentPath)
			if err != nil {
				return fmt.Errorf("failed to open b3dm file: %v", err)
			}
			reader := b3dms.NewB3dmReader(b3dmFile)
			b3dm := new(b3dms.B3dm)
			if err := reader.Decode(b3dm); err != nil {
				return err
			}
			featureTable := b3dm.GetFeatureTable()
			batchLength := featureTable.GetBatchLength()
			featureTableView := b3dm.GetFeatureTableView()
			batchTable := b3dm.GetBatchTable()
			batchTableProperties := batchTable.Data
			computedFeaturePositions := []Cartographic{}
			gltf := b3dm.GetModel()
			if gltf != nil {
				rtcTransform, err := getRtcTransform(featureTableView, gltf)
				if err != nil {
					return fmt.Errorf("failed to getRtcTransform: %v", err)
				}
				toZUpTransform := getZUpTransform()
				computedFeaturePositions, err = computeFeaturePositionsFromGltfVertices(
					gltf,
					computedTransform,
					rtcTransform,
					toZUpTransform,
					batchLength,
				)
				if err != nil {
					return fmt.Errorf("failed to open b3dm file: %v", err)
				}
			}

			for batchId := 0; batchId < batchLength; batchId++ {
				batchProperties := make(map[string]interface{})
				for name, values := range batchTableProperties {
					batchProperties[name] = nil
					if len(values) > 0 {
						batchProperties[name] = values[batchId]
					}
				}
				position := computedFeaturePositions[batchId]
				idValue := batchProperties[config.IdProperty].(string)
				uniqueFeatures[idValue] = TilesetFeature{
					Position:   position,
					Properties: batchProperties,
				}
			}

			return nil
		}
		err := ForEachTile(tileset, tilesetIterFn)
		if err != nil {
			return nil, fmt.Errorf("something went wrong at iterTile: %v", err)
		}
	}

	return uniqueFeatures, nil
}

func computeFeaturePositionsFromGltfVertices(gltf *gltf.Document, tileTransform, rtcTransform, toZUpTransform *mat.Dense, batchLength int) ([]Cartographic, error) {
	nodes := gltf.Nodes
	if nodes == nil {
		return nil, errors.New("nodes are nil")
	}
	meshes := gltf.Meshes
	if meshes == nil {
		return nil, errors.New("meshes are nil")
	}
	accessors := gltf.Accessors
	if accessors == nil {
		return nil, errors.New("accesors are nil")
	}
	bufferViews := gltf.BufferViews
	if bufferViews == nil {
		return nil, errors.New("bufferViews are nil")
	}

	batchIdPositions := make([][]Cartographic, batchLength)

	for _, node := range nodes {
		mesh := meshes[*node.Mesh]
		primitives := mesh.Primitives
		nodeMatrix := eyeMat(4)
		if len(node.Matrix) > 0 {
			nodeMatrix = mat4FromGltfNodeMatrix(node.Matrix)
		}

		modelMatrix := eyeMat(4)
		modelMatrix = mat4MultiplyTransformation(modelMatrix, tileTransform)
		modelMatrix = mat4MultiplyTransformation(modelMatrix, rtcTransform)
		modelMatrix = mat4MultiplyTransformation(modelMatrix, toZUpTransform)
		modelMatrix = mat4MultiplyTransformation(modelMatrix, nodeMatrix)

		for _, primitive := range primitives {
			attributes := primitive.Attributes
			_BATCHID := attributes["_BATCHID"]
			POSITION := attributes["POSITION"]

			count := accessors[POSITION].Count
			for i := uint32(0); i < count; i++ {
				// If the gltf vertices are tagged with BATCHID, store the positions at
				// the respective BATCHID. Otherwise store everything under a single
				// BATCHID=0
				var batchIdValue interface{}
				if _BATCHID == 0 {
					batchIdValue = 0
				} else {
					batchIdValue = b3dms.ReadGltfValueAt(gltf, _BATCHID, i)[0]
				}

				batchId, err := getInt(batchIdValue)
				if err != nil {
					return nil, fmt.Errorf("getInt failed: %w", err)
				}
				result := b3dms.ReadGltfValueAt(gltf, POSITION, i)
				points, err := util.TryMap(result, getFloat)
				if err != nil {
					return nil, fmt.Errorf("map failed: %w", err)
				}
				localPosition := cartesianFromSlice(points)
				worldPosition := multiplyMat4ByPoint(modelMatrix, localPosition)
				cartographic, err := cartographicFromCartesian3(worldPosition)
				if err != nil {
					return nil, fmt.Errorf("failed to convert cartesian to cartographic: %w", err)
				}
				if batchIdPositions[batchId] == nil {
					batchIdPositions[batchId] = []Cartographic{}
				}

				if cartographic != nil {
					batchIdPositions[batchId] = append(batchIdPositions[batchId], *cartographic)
				}
			}
		}
	}

	featurePositions := []Cartographic{}

	for _, positions := range batchIdPositions {
		height := []float64{}
		for _, carto := range positions {
			height = append(height, carto.Height)
		}
		minHeight, maxHeight := minMaxOfSlice(height)
		featureHeight := maxHeight - minHeight
		rectangle := rectangleFromCartographicArray(positions)
		position := rectangle.center()
		position.Height = featureHeight

		featurePositions = append(featurePositions, *position)
	}

	return featurePositions, nil
}

type TileIterFn func(*tiles.Tile, *mat.Dense) error

func ForEachTile(ts *tiles.Tileset, iterFn func(tile *tiles.Tile, computedTransform *mat.Dense) error) error {
	root := &ts.Root

	var iterTile TileIterFn
	iterTile = func(tile *tiles.Tile, parentTransform *mat.Dense) error {
		computedTransform := parentTransform
		if tile.Transform != nil {
			test := tile.Transform[:]
			computedTransform.Mul(parentTransform, mat.NewDense(4, 4, test))
		}
		err := iterFn(tile, computedTransform)
		if err != nil {
			return fmt.Errorf("something wrong at iterFn: %v", err)
		}
		if (tile.Children != nil) && len(*tile.Children) != 0 {
			for _, child := range *tile.Children {
				err = iterTile(&child, computedTransform)
				if err != nil {
					return fmt.Errorf("something went wrong at iterTile: %v", err)
				}
			}
		}
		return nil
	}

	err := iterTile(root, eyeMat(4))
	if err != nil {
		return fmt.Errorf("something went wrong at iterTile: %v", err)
	}

	return nil
}
