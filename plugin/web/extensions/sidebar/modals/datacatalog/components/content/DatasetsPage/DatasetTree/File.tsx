import { DataCatalogItem } from "@web/extensions/sidebar/core/types";
import { postMsg } from "@web/extensions/sidebar/utils";
import { getNameFromPath } from "@web/extensions/sidebar/utils/file";
import { Button, Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useEffect, useMemo } from "react";

export type Props = {
  item: DataCatalogItem;
  isMobile?: boolean;
  nestLevel: number;
  selectedID?: string;
  addDisabled: (dataID: string) => boolean;
  onDatasetAdd: (dataset: DataCatalogItem) => void;
  onOpenDetails?: (item?: DataCatalogItem) => void;
  onSelect?: (dataID: string) => void;
  setExpandedKeys: React.Dispatch<React.SetStateAction<string[]>>;
};

const File: React.FC<Props> = ({
  item,
  isMobile,
  nestLevel,
  selectedID,
  addDisabled,
  onDatasetAdd,
  onOpenDetails,
  onSelect,
  setExpandedKeys,
}) => {
  const handleClick = useCallback(() => {
    onDatasetAdd(item);
  }, [item, onDatasetAdd]);

  const handleOpenDetails = useCallback(() => {
    onOpenDetails?.(item);
    onSelect?.(item.dataID);
  }, [item, onOpenDetails, onSelect]);

  const selected = useMemo(
    () => (item.type !== "group" ? selectedID === item.id : false),
    [selectedID, item],
  );

  useEffect(() => {
    const { selectedDataset } = window as any;
    if (selectedDataset) {
      onOpenDetails?.(selectedDataset);
      onSelect?.(selectedDataset.dataID);
      if (selected && item.path) {
        const expandedKeys = [...item.path];
        setExpandedKeys(expandedKeys);
        postMsg({
          action: "storageSave",
          payload: { key: "expandedKeys", value: expandedKeys },
        });
      }
      setTimeout(() => {
        (window as any).selectedDataset = undefined;
      }, 500);
    }
  }, [item.path, onOpenDetails, onSelect, selected, setExpandedKeys]);

  const name = useMemo(() => getNameFromPath(item.name), [item.name]);

  return (
    <Wrapper nestLevel={nestLevel} selected={selected}>
      <NameWrapper isMobile={isMobile} onClick={handleOpenDetails}>
        <Icon icon="file" size={20} />
        {!item.public && <UnpublishedIndicator />}
        <Name>{name}</Name>
      </NameWrapper>
      <StyledButton
        type="link"
        icon={<StyledIcon icon="plusCircle" selected={selected ?? false} />}
        onClick={handleClick}
        disabled={addDisabled(item.dataID)}
      />
    </Wrapper>
  );
};

export default File;

const Wrapper = styled.div<{ nestLevel: number; selected?: boolean }>`
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
  gap: 8px;
  height: 29px;
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

const NameWrapper = styled.div<{ isMobile?: boolean }>`
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
`;

const Name = styled.p`
  width: 150px;
  margin: 0;
  user-select: none;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

const StyledButton = styled(Button)<{ disabled: boolean }>`
  display: ${({ disabled }) => (disabled ? "none" : "initial")};
`;

const StyledIcon = styled(Icon)<{ selected: boolean }>`
  color: ${({ selected }) => (selected ? "#ffffff" : "#00bebe")};
  ${Wrapper}:hover & {
    color: #ffffff;
  }
`;

const UnpublishedIndicator = styled.div`
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #6d6d6d;
`;
