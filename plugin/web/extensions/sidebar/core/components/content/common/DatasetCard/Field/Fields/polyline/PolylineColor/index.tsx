import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import {
  ButtonWrapper,
  Wrapper,
} from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { generateID, moveItemDown, moveItemUp, removeItem } from "@web/extensions/sidebar/utils";
import { useCallback, useState } from "react";

import { BaseFieldProps, Cond } from "../../types";

import PolylineColorItem from "./PolylineColorItem";

const PolylineColor: React.FC<BaseFieldProps<"polylineColor">> = ({
  value,
  editMode,
  onUpdate,
}) => {
  const [items, updateItems] = useState(value.items);

  const handleMoveUp = useCallback(
    (idx: number) => {
      updateItems(c => {
        const newPointColors = moveItemUp(idx, c) ?? c;
        onUpdate({
          ...value,
          items: newPointColors,
        });
        return newPointColors;
      });
    },
    [value, onUpdate],
  );

  const handleMoveDown = useCallback(
    (idx: number) => {
      updateItems(c => {
        const newPointColors = moveItemDown(idx, c) ?? c;
        onUpdate({
          ...value,
          items: newPointColors,
        });
        return newPointColors;
      });
    },
    [onUpdate, value],
  );

  const handleAdd = useCallback(() => {
    updateItems(c => {
      const newPointColor: { condition: Cond<number>; color: string } = {
        condition: {
          key: generateID(),
          operator: "=",
          operand: "width",
          value: 1,
        },
        color: "",
      };
      onUpdate({
        ...value,
        items: value.items ? [...value.items, newPointColor] : [newPointColor],
      });
      return c ? [...c, newPointColor] : [newPointColor];
    });
  }, [value, onUpdate]);

  const handleRemove = useCallback(
    (idx: number) => {
      updateItems(c => {
        const newPointColors = removeItem(idx, c) ?? c;
        onUpdate({
          ...value,
          items: newPointColors,
        });
        return newPointColors;
      });
    },
    [value, onUpdate],
  );

  const handleItemUpdate = (item: { condition: Cond<number>; color: string }, index: number) => {
    updateItems(c => {
      const newPointColors = [...(c ?? [])];
      newPointColors.splice(index, 1, item);
      onUpdate({
        ...value,
        items: newPointColors,
      });
      return newPointColors;
    });
  };

  return editMode ? (
    <Wrapper>
      {items?.map((c, idx) => (
        <PolylineColorItem
          key={idx}
          index={idx}
          item={c}
          handleMoveDown={handleMoveDown}
          handleMoveUp={handleMoveUp}
          handleRemove={handleRemove}
          onItemUpdate={handleItemUpdate}
        />
      ))}
      <ButtonWrapper>
        <AddButton text="Add Condition" height={24} onClick={handleAdd} />
      </ButtonWrapper>
    </Wrapper>
  ) : null;
};

export default PolylineColor;
