import { DataCatalogGroup, DataCatalogItem } from "../../../../api/api";

import File from "./File";
import Folder from "./Folder";

type Props = {
  catalog: DataCatalogGroup | DataCatalogItem | (DataCatalogItem | DataCatalogGroup)[];
  isMobile?: boolean;
  expandAll?: boolean;
  addedDatasetIds?: string[];
  selectedId?: string;
  nestLevel: number;
  onDatasetAdd: (dataset: DataCatalogItem) => void;
  onOpenDetails?: (item?: DataCatalogItem) => void;
  onSelect?: (id: string) => void;
};

const TreeBuilder: React.FC<Props> = ({
  catalog,
  isMobile,
  expandAll,
  addedDatasetIds,
  selectedId,
  nestLevel,
  onDatasetAdd,
  onOpenDetails,
  onSelect,
}) => {
  console.log("CATALOG IN TREE BUILDER: ", catalog);
  return (
    <>
      {Array.isArray(catalog) ? (
        catalog.map(item => {
          "children" in item ? (
            <Folder
              key={item.name}
              name={item.name}
              nestLevel={nestLevel + 1}
              isMobile={isMobile}
              expandAll={expandAll}>
              <TreeBuilder
                catalog={item}
                addedDatasetIds={addedDatasetIds}
                selectedId={selectedId}
                nestLevel={nestLevel + 1}
                onDatasetAdd={onDatasetAdd}
                onOpenDetails={onOpenDetails}
                onSelect={onSelect}
              />
            </Folder>
          ) : (
            <TreeBuilder
              catalog={item}
              addedDatasetIds={addedDatasetIds}
              selectedId={selectedId}
              nestLevel={nestLevel + 1}
              onDatasetAdd={onDatasetAdd}
              onOpenDetails={onOpenDetails}
              onSelect={onSelect}
            />
          );
        })
      ) : "children" in catalog ? (
        <Folder
          key={catalog.name}
          name={catalog.name}
          nestLevel={nestLevel + 1}
          isMobile={isMobile}
          expandAll={expandAll}>
          <TreeBuilder
            catalog={catalog}
            addedDatasetIds={addedDatasetIds}
            selectedId={selectedId}
            nestLevel={nestLevel + 1}
            onDatasetAdd={onDatasetAdd}
            onOpenDetails={onOpenDetails}
            onSelect={onSelect}
          />
        </Folder>
      ) : (
        <File
          item={catalog}
          addedDatasetIds={addedDatasetIds}
          isMobile={isMobile}
          nestLevel={nestLevel + 1}
          selectedID={selectedId}
          onDatasetAdd={onDatasetAdd}
          onOpenDetails={onOpenDetails}
          onSelect={onSelect}
        />
      )}
    </>
  );
};

export default TreeBuilder;
