import { useEffect, useState } from "react";

import { DataCatalogGroup, DataCatalogItem } from "../../../../api/api";

import File from "./File";
import Folder from "./Folder";

type Props = {
  catalogItem: DataCatalogGroup | DataCatalogItem | (DataCatalogItem | DataCatalogGroup)[];
  isMobile?: boolean;
  expandAll?: boolean;
  addedDatasetDataIDs?: string[];
  selectedID?: string;
  nestLevel: number;
  showParentIndicator?: boolean;
  onShowParentIndicator?: React.Dispatch<React.SetStateAction<boolean>>;
  addDisabled: (dataID: string) => boolean;
  onDatasetAdd: (dataset: DataCatalogItem) => void;
  onOpenDetails?: (item?: DataCatalogItem) => void;
  onSelect?: (dataID: string) => void;
};

const TreeBuilder: React.FC<Props> = ({
  catalogItem,
  isMobile,
  expandAll,
  addedDatasetDataIDs,
  selectedID,
  nestLevel,
  showParentIndicator,
  onShowParentIndicator,
  addDisabled,
  onDatasetAdd,
  onOpenDetails,
  onSelect,
}) => {
  const [showIndicator, setShowIndicator] = useState(false);

  useEffect(() => {
    if (showIndicator && showParentIndicator !== showIndicator) {
      onShowParentIndicator?.(showIndicator);
    }
  }, [showIndicator, showParentIndicator, onShowParentIndicator]);

  return (
    <>
      {Array.isArray(catalogItem) ? (
        catalogItem.map(item =>
          "children" in item ? (
            <Folder
              key={item.name}
              name={item.name}
              nestLevel={nestLevel + 1}
              isMobile={isMobile}
              expandAll={expandAll}
              showIndicator={showIndicator}>
              <TreeBuilder
                catalogItem={item.children}
                addedDatasetDataIDs={addedDatasetDataIDs}
                selectedID={selectedID}
                nestLevel={nestLevel + 1}
                showParentIndicator={showIndicator}
                onShowParentIndicator={setShowIndicator}
                addDisabled={addDisabled}
                onDatasetAdd={onDatasetAdd}
                onOpenDetails={onOpenDetails}
                onSelect={onSelect}
              />
            </Folder>
          ) : (
            <TreeBuilder
              catalogItem={item}
              addedDatasetDataIDs={addedDatasetDataIDs}
              selectedID={selectedID}
              nestLevel={nestLevel + 1}
              showParentIndicator={showIndicator}
              onShowParentIndicator={setShowIndicator}
              addDisabled={addDisabled}
              onDatasetAdd={onDatasetAdd}
              onOpenDetails={onOpenDetails}
              onSelect={onSelect}
            />
          ),
        )
      ) : "children" in catalogItem ? (
        <Folder
          key={catalogItem.name}
          name={catalogItem.name}
          nestLevel={nestLevel + 1}
          isMobile={isMobile}
          expandAll={expandAll}
          showIndicator={showIndicator}>
          <TreeBuilder
            catalogItem={catalogItem.children}
            addedDatasetDataIDs={addedDatasetDataIDs}
            selectedID={selectedID}
            nestLevel={nestLevel + 1}
            showParentIndicator={showIndicator}
            onShowParentIndicator={setShowIndicator}
            addDisabled={addDisabled}
            onDatasetAdd={onDatasetAdd}
            onOpenDetails={onOpenDetails}
            onSelect={onSelect}
          />
        </Folder>
      ) : (
        <File
          item={catalogItem}
          isMobile={isMobile}
          nestLevel={nestLevel + 1}
          selectedID={selectedID}
          showParentIndicator={showParentIndicator}
          onShowParentIndicator={onShowParentIndicator}
          addDisabled={addDisabled}
          onDatasetAdd={onDatasetAdd}
          onOpenDetails={onOpenDetails}
          onSelect={onSelect}
        />
      )}
    </>
  );
};

export default TreeBuilder;
