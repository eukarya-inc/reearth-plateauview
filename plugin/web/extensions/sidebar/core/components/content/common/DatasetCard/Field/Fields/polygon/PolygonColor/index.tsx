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
import { useCallback, useEffect, useState } from "react";

import { stringifyCondition } from "../../../utils";
import { BaseFieldProps, Cond } from "../../types";

import PolygonColorItem from "./PolygonColorItem";

const PolygonColor: React.FC<BaseFieldProps<"polygonColor">> = ({
  dataID,
  value,
  editMode,
  isActive,
  onUpdate,
}) => {
  const [items, updateItems] = useState(value.items);

  const handleMoveUp = useCallback(
    (idx: number) => {
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
      updateItems(c => {
        const newItems = moveItemDown(idx, c) ?? c;
        onUpdate({
          ...value,
          items: newItems,
        });
        return newItems;
      });
    },
    [onUpdate, value],
  );

  const handleAdd = useCallback(() => {
    updateItems(c => {
      const newItem: { condition: Cond<any>; color: string } = {
        condition: {
          key: generateID(),
          operator: "",
          operand: "",
          value: "",
        },
        color: "",
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

  const handleItemUpdate = (item: { condition: Cond<number>; color: string }, index: number) => {
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
      const fillColorConditions: [string, string][] = [["true", 'color("white")']];
      const fillConditions: [string, string][] = [["true", "true"]];
      items?.forEach(item => {
        const resFillColor = "color" + `("${item.color}")`;
        const cond = stringifyCondition(item.condition);
        fillColorConditions.unshift([cond, resFillColor]);
        fillConditions.unshift([cond, cond]);
        postMsg({
          action: "updateDatasetInScene",
          payload: {
            dataID,
            update: {
              polygon: {
                fill: { expression: { conditions: fillConditions } },
                fillColor: { expression: { conditions: fillColorConditions } },
              },
            },
          },
        });
      });
    }, 500);
    return () => {
      clearTimeout(timer);
      postMsg({
        action: "updateDatasetInScene",
        payload: {
          dataID,
          update: { polygon: undefined },
        },
      });
    };
  }, [dataID, isActive, items]);

  return editMode ? (
    <Wrapper>
      {items?.map((c, idx) => (
        <PolygonColorItem
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

export default PolygonColor;
