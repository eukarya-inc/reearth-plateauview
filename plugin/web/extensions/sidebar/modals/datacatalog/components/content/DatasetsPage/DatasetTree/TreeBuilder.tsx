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
  nodeKey: string;
  expandedKeys: string[];
  addDisabled: (dataID: string) => boolean;
  onDatasetAdd: (dataset: DataCatalogItem) => void;
  onOpenDetails?: (item?: DataCatalogItem) => void;
  onSelect?: (dataID: string) => void;
  setExpandedKeys: React.Dispatch<React.SetStateAction<string[]>>;
  setSelectedKey: (key: string) => void;
};

const TreeBuilder: React.FC<Props> = ({
  catalogItem,
  isMobile,
  expandAll,
  addedDatasetDataIDs,
  selectedID,
  nestLevel,
  nodeKey,
  expandedKeys,
  addDisabled,
  onDatasetAdd,
  onOpenDetails,
  onSelect,
  setExpandedKeys,
  setSelectedKey,
}) => {
  return (
    <>
      {Array.isArray(catalogItem) ? (
        catalogItem.map((item, index) =>
          "children" in item ? (
            <Folder
              key={item.name}
              name={item.name}
              nestLevel={nestLevel + 1}
              nodeKey={nodeKey + "-" + index}
              expandedKeys={expandedKeys}
              isMobile={isMobile}
              expandAll={expandAll}
              setExpandedKeys={setExpandedKeys}>
              <TreeBuilder
                catalogItem={item.children}
                addedDatasetDataIDs={addedDatasetDataIDs}
                selectedID={selectedID}
                nestLevel={nestLevel + 1}
                nodeKey={nodeKey + "-" + index}
                expandedKeys={expandedKeys}
                addDisabled={addDisabled}
                onDatasetAdd={onDatasetAdd}
                onOpenDetails={onOpenDetails}
                onSelect={onSelect}
                setExpandedKeys={setExpandedKeys}
                setSelectedKey={setSelectedKey}
              />
            </Folder>
          ) : (
            <TreeBuilder
              catalogItem={item}
              addedDatasetDataIDs={addedDatasetDataIDs}
              selectedID={selectedID}
              nestLevel={nestLevel + 1}
              nodeKey={nodeKey + "-" + index}
              expandedKeys={expandedKeys}
              addDisabled={addDisabled}
              onDatasetAdd={onDatasetAdd}
              onOpenDetails={onOpenDetails}
              onSelect={onSelect}
              setExpandedKeys={setExpandedKeys}
              setSelectedKey={setSelectedKey}
            />
          ),
        )
      ) : "children" in catalogItem ? (
        <Folder
          key={catalogItem.name}
          name={catalogItem.name}
          nestLevel={nestLevel + 1}
          nodeKey={nodeKey}
          expandedKeys={expandedKeys}
          isMobile={isMobile}
          expandAll={expandAll}
          setExpandedKeys={setExpandedKeys}>
          <TreeBuilder
            catalogItem={catalogItem.children}
            addedDatasetDataIDs={addedDatasetDataIDs}
            selectedID={selectedID}
            nestLevel={nestLevel + 1}
            nodeKey={nodeKey}
            expandedKeys={expandedKeys}
            addDisabled={addDisabled}
            onDatasetAdd={onDatasetAdd}
            onOpenDetails={onOpenDetails}
            onSelect={onSelect}
            setExpandedKeys={setExpandedKeys}
            setSelectedKey={setSelectedKey}
          />
        </Folder>
      ) : (
        <File
          item={catalogItem}
          isMobile={isMobile}
          nestLevel={nestLevel}
          selectedID={selectedID}
          nodeKey={nodeKey}
          addDisabled={addDisabled}
          onDatasetAdd={onDatasetAdd}
          onOpenDetails={onOpenDetails}
          onSelect={onSelect}
          setSelectedKey={setSelectedKey}
        />
      )}
    </>
  );
};

export default TreeBuilder;
