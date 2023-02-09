import { CatalogItem } from "@web/extensions/sidebar/core/processCatalog";
import { useMemo } from "react";

import File from "./File";
import Folder from "./Folder";

type Props = {
  item: CatalogItem;
  addedDatasetIds?: string[];
  selectedId?: string;
  nestLevel: number;
  onDatasetAdd: (dataset: CatalogItem) => void;
  onOpenDetails?: (item?: CatalogItem) => void;
  onSelect?: (id: string) => void;
};

const TreeBuilder: React.FC<Props> = ({
  item,
  addedDatasetIds,
  selectedId,
  nestLevel,
  onDatasetAdd,
  onOpenDetails,
  onSelect,
}) => {
  const selected = useMemo(
    () => (item.type !== "group" ? selectedId === item.id : false),
    [selectedId, item],
  );

  return item.type === "group" ? (
    <Folder
      item={item}
      addedDatasetIds={addedDatasetIds}
      selectedId={selectedId}
      nestLevel={nestLevel + 1}
      onDatasetAdd={onDatasetAdd}
      onOpenDetails={onOpenDetails}
      onSelect={onSelect}
    />
  ) : (
    <File
      key={item.name}
      item={item}
      addedDatasetIds={addedDatasetIds}
      nestLevel={nestLevel}
      selected={selected}
      onDatasetAdd={onDatasetAdd}
      onOpenDetails={onOpenDetails}
      onSelect={onSelect}
    />
  );
};

export default TreeBuilder;
