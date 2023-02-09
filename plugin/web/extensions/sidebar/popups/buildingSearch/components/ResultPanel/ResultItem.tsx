import { styled } from "@web/theme";
import { useCallback } from "react";

import { Result } from "../../types";

type Props = {
  item: Result;
  selected: Result[];
  onSelect: (selected: Result[]) => void;
};

const ResultItem: React.FC<Props> = ({ item, selected, onSelect }) => {
  const onClick = useCallback(() => {
    if (selected.find(s => s.gml_id === item.gml_id)) {
      onSelect(selected.filter(s => s.gml_id !== item.gml_id));
    } else {
      onSelect([item]);
    }
  }, [onSelect, item, selected]);
  return (
    <StyledResultItem onClick={onClick} active={!!selected.find(s => s.gml_id === item.gml_id)}>
      {item.gml_id}
    </StyledResultItem>
  );
};

const StyledResultItem = styled.div<{ active: boolean }>`
  display: flex;
  align-items: center;
  width: 100%;
  height: 38px;
  padding: 0 12px;
  border-bottom: 1px solid #d9d9d9;
  font-size: 12px;
  background: ${({ active }) => (active ? "var(--theme-color)" : "#fff")};
  color: ${({ active }) => (active ? "#fff" : "#000")};
  cursor: pointer;

  &:last-child {
    border-bottom: none;
  }
`;

export default ResultItem;
