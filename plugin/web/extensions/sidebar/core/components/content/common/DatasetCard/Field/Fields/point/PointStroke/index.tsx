import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import {
  ButtonWrapper,
  Wrapper,
} from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import {
  generateID,
  moveItemDown,
  moveItemUp,
  removeItem,
  postMsg,
} from "@web/extensions/sidebar/utils";
import { useCallback, useState, useEffect } from "react";

import { BaseFieldProps, Cond } from "../../types";

import PointStrokeItem from "./PointStrokeItem";

const PointStroke: React.FC<BaseFieldProps<"pointStroke">> = ({
  dataID,
  value,
  editMode,
  isActive,
  onUpdate,
}) => {
  const [items, updateItems] = useState(value.items);

  const handleMoveUp = useCallback(
    (idx: number) => {
      if (idx === 0) return;
      updateItems(c => {
        const newItems = moveItemUp(idx, c) ?? c;
        onUpdate({
          ...value,
          items: newItems,
        });
        return newItems;
      });
    },
    [value, onUpdate],
  );

  const handleMoveDown = useCallback(
    (idx: number) => {
      if (items && idx >= items.length - 1) return;
      updateItems(c => {
        const newItems = moveItemDown(idx, c) ?? c;
        onUpdate({
          ...value,
          items: newItems,
        });
        return newItems;
      });
    },
    [items, onUpdate, value],
  );

  const handleAdd = useCallback(() => {
    updateItems(c => {
      const newItem: {
        strokeColor: string;
        strokeWidth: number;
        condition: Cond<string | number>;
      } = {
        strokeColor: "",
        strokeWidth: 0,
        condition: {
          key: generateID(),
          operator: "=",
          operand: "width",
          value: 1,
        },
      };
      onUpdate({
        ...value,
        items: value.items ? [...value.items, newItem] : [newItem],
      });
      return c ? [...c, newItem] : [newItem];
    });
  }, [value, onUpdate]);

  const handleRemove = useCallback(
    (idx: number) => {
      updateItems(c => {
        const newItems = removeItem(idx, c) ?? c;
        onUpdate({
          ...value,
          items: newItems,
        });
        return newItems;
      });
    },
    [value, onUpdate],
  );

  const handleItemUpdate = (
    item: { condition: Cond<string | number>; strokeColor: string; strokeWidth: number },
    index: number,
  ) => {
    updateItems(c => {
      const newItems = [...(c ?? [])];
      newItems.splice(index, 1, item);
      onUpdate({
        ...value,
        items: newItems,
      });
      return newItems;
    });
  };

  useEffect(() => {
    if (!isActive || !dataID) return;
    const timer = setTimeout(() => {
      postMsg({
        action: "updateDatasetInScene",
        payload: {
          dataID,
          update: { marker: { style: "image", items } },
        },
      });
    }, 500);
    return () => {
      clearTimeout(timer);
    };
  }, [dataID, isActive, items]);

  return editMode ? (
    <Wrapper>
      {items?.map((c, idx) => (
        <PointStrokeItem
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

export default PointStroke;
