import { Catalog, CatalogItem } from "@web/extensions/sidebar/core/components/hooks";

import File from "./File";
import Folder from "./Folder";

type Props = {
  catalog: Catalog | CatalogItem[];
  isMobile?: boolean;
  expandAll?: boolean;
  addedDatasetIds?: string[];
  selectedId?: string;
  nestLevel: number;
  onDatasetAdd: (dataset: CatalogItem) => void;
  onOpenDetails?: (item?: CatalogItem) => void;
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
  return (
    <>
      {!Array.isArray(catalog)
        ? Object.keys(catalog).map(loc => (
            <Folder
              key={loc}
              name={loc}
              nestLevel={nestLevel + 1}
              isMobile={isMobile}
              expandAll={expandAll}>
              <TreeBuilder
                catalog={catalog[loc]}
                addedDatasetIds={addedDatasetIds}
                selectedId={selectedId}
                nestLevel={nestLevel + 1}
                onDatasetAdd={onDatasetAdd}
                onOpenDetails={onOpenDetails}
                onSelect={onSelect}
              />
            </Folder>
          ))
        : catalog.map(item => (
            <File
              key={item.id}
              item={item}
              addedDatasetIds={addedDatasetIds}
              isMobile={isMobile}
              nestLevel={nestLevel + 1}
              selectedID={selectedId}
              onDatasetAdd={onDatasetAdd}
              onOpenDetails={onOpenDetails}
              onSelect={onSelect}
            />
          ))}
    </>
  );
};

export default TreeBuilder;
