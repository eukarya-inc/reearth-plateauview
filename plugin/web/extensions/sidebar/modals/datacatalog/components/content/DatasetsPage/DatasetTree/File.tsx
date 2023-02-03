import { CatalogItem } from "@web/extensions/sidebar/core/processCatalog";
import { Button, Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useMemo } from "react";

export type Props = {
  item: CatalogItem;
  addedDatasetIds?: string[];
  selectedDataset?: CatalogItem;
  nestLevel: number;
  selected: boolean;
  onDatasetAdd: (dataset: CatalogItem) => void;
  onOpenDetails?: (item?: CatalogItem) => void;
  onSelect?: (id: string) => void;
};

const File: React.FC<Props> = ({
  item,
  addedDatasetIds,
  selectedDataset,
  nestLevel,
  selected,
  onDatasetAdd,
  onOpenDetails,
  onSelect,
}) => {
  const handleClick = useCallback(() => {
    onDatasetAdd(item);
  }, [item, onDatasetAdd]);

  const handleOpenDetails = useCallback(() => {
    if (item.type === "group") return;
    onOpenDetails?.(item);
    onSelect?.(item.id);
  }, [item, onOpenDetails, onSelect]);

  const addDisabled = useMemo(() => {
    return !!addedDatasetIds?.find(
      id => selectedDataset?.type === "item" && id === selectedDataset.id,
    );
  }, [addedDatasetIds, selectedDataset]);

  return (
    <>
      {item.type === "item" && (
        <Wrapper key={item.id} nestLevel={nestLevel} selected={selected}>
          <NameWrapper onClick={handleOpenDetails}>
            <Icon icon="file" size={20} />
            <Name>{item.cityName ?? item.name}</Name>
          </NameWrapper>
          <Button
            type="link"
            icon={<StyledIcon icon="plusCircle" selected={selected ?? false} />}
            onClick={handleClick}
            disabled={addDisabled}
          />
        </Wrapper>
      )}
    </>
  );
};

export default File;

const Wrapper = styled.div<{ nestLevel: number; selected?: boolean }>`
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
  gap: 8px;
  min-height: 29px;
  ${({ selected }) =>
    selected &&
    `
  background: #00BEBE;
  color: white;
  `}

  padding-left: ${({ nestLevel }) => (nestLevel ? `${nestLevel * 8}px` : "8px")};
  padding-right: 8px;
  cursor: pointer;

  :hover {
    background: #00bebe;
    color: white;
  }
`;

const NameWrapper = styled.div`
  display: flex;
`;

const Name = styled.p`
  margin: 0 0 0 8px;
  user-select: none;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
  width: 200px;
`;

const StyledIcon = styled(Icon)<{ selected: boolean }>`
  color: ${({ selected }) => (selected ? "#ffffff" : "#00bebe")};
  ${Wrapper}:hover & {
    color: #ffffff;
  }
`;
