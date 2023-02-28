import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useState, useEffect } from "react";

export type Props = {
  name: string;
  isMobile?: boolean;
  expandAll?: boolean;
  nestLevel: number;
  nodeKey: string;
  expandedKeys: string[];
  onExpand: (key: string) => void;
  children?: React.ReactNode;
};

const Folder: React.FC<Props> = ({
  name,
  isMobile,
  expandAll,
  nestLevel,
  nodeKey,
  expandedKeys,
  onExpand,
  children,
}) => {
  const [isOpen, open] = useState(false);

  useEffect(() => {
    open(() => expandAll || expandedKeys.includes(nodeKey));
  }, [expandAll, expandedKeys, nodeKey]);

  return (
    <Wrapper key={name} isOpen={isOpen}>
      <FolderItem nestLevel={nestLevel} onClick={() => onExpand(nodeKey)}>
        <NameWrapper isMobile={isMobile}>
          <Icon icon={isOpen ? "folderOpen" : "folder"} size={20} />
          <Name>{name}</Name>
        </NameWrapper>
      </FolderItem>
      {children}
    </Wrapper>
  );
};

export default Folder;

const Wrapper = styled.div<{ isOpen?: boolean }>`
  width: 100%;
  ${({ isOpen }) =>
    isOpen
      ? "height: 100%;"
      : `
  height: 29px; 
  overflow: hidden;
  `}
`;

const FolderItem = styled.div<{ nestLevel: number; selected?: boolean }>`
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
  margin: 0;
  user-select: none;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;
