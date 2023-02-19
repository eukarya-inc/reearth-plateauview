import { useCallback } from "react";

import { Cond } from "../types";

import { ColorField, ConditionField, ItemControls } from "./common";
import { Item } from "./commonComponents";

const PointStrokeItem: React.FC<{
  index: number;
  item: { condition: Cond<string | number>; strokeColor: string; strokeWidth: number };
  handleMoveDown: (index: number) => void;
  handleMoveUp: (index: number) => void;
  handleRemove: (index: number) => void;
  onItemUpdate: (
    item: { condition: Cond<string | number>; strokeColor: string; strokeWidth: number },
    index: number,
  ) => void;
}> = ({ index, item, handleMoveDown, handleMoveUp, handleRemove, onItemUpdate }) => {
  const handleStrokeColorUpdate = useCallback(
    (color: string) => {
      if (color) {
        const copy = { ...item, strokeColor: color };
        onItemUpdate(copy, index);
      }
    },
    [index, item, onItemUpdate],
  );

  const handleConditionUpdate = useCallback(
    (condition: Cond<number>) => {
      if (condition) {
        const copy = { ...item, condition };
        onItemUpdate(copy, index);
      }
    },
    [index, item, onItemUpdate],
  );

  return (
    <Item>
      <ItemControls
        index={index}
        handleMoveDown={handleMoveDown}
        handleMoveUp={handleMoveUp}
        handleRemove={handleRemove}
      />
      <ConditionField
        title="if"
        fieldGap={8}
        condition={item.condition}
        onChange={handleConditionUpdate}
      />
      <ColorField
        title="strokeColor"
        titleWidth={82}
        color={item.strokeColor}
        onChange={handleStrokeColorUpdate}
      />
    </Item>
  );
};

export default PointStrokeItem;
