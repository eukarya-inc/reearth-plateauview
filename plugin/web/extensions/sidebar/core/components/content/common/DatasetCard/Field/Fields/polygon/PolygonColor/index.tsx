import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import {
  ButtonWrapper,
  Wrapper,
} from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { generateID, moveItemDown, moveItemUp, removeItem } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

import { stringifyCondition } from "../../../utils";
import { BaseFieldProps, Cond } from "../../types";

import PolygonColorItem from "./PolygonColorItem";

const PolygonColor: React.FC<BaseFieldProps<"polygonColor">> = ({
  dataID,
  value,
  editMode,
  onUpdate,
}) => {
  const [items, updateItems] = useState(value.items);

  const handleMoveUp = useCallback((idx: number) => {
    updateItems(c => {
      const newitems = moveItemUp(idx, c) ?? c;
      return newitems;
    });
  }, []);

  const handleMoveDown = useCallback((idx: number) => {
    updateItems(c => {
      const newitems = moveItemDown(idx, c) ?? c;
      return newitems;
    });
  }, []);

  const handleAdd = useCallback(() => {
    updateItems(c => {
      const newPointColor: { condition: Cond<any>; color: string } = {
        condition: {
          key: generateID(),
          operator: "",
          operand: "",
          value: "",
        },
        color: "",
      };
      return c ? [...c, newPointColor] : [newPointColor];
    });
  }, []);

  const handleRemove = useCallback((idx: number) => {
    updateItems(c => {
      const newitems = removeItem(idx, c) ?? c;
      return newitems;
    });
  }, []);

  const handleItemUpdate = (item: { condition: Cond<number>; color: string }, index: number) => {
    updateItems(c => {
      const newitems = [...(c ?? [])];
      newitems.splice(index, 1, item);
      return newitems;
    });
  };

  useEffect(() => {
    if (!dataID || value.items === items) return;

    const timer = setTimeout(() => {
      const fillColorConditions: [string, string][] = [["true", 'color("white")']];
      const fillConditions: [string, string][] = [["true", "true"]];
      items?.forEach(item => {
        const resFillColor = "color" + `("${item.color}")`;
        const cond = stringifyCondition(item.condition);
        fillColorConditions.unshift([cond, resFillColor]);
        fillConditions.unshift([cond, cond]);
        onUpdate({
          ...value,
          items,
          override: {
            polygon: {
              fill: { expression: { conditions: fillConditions } },
              fillColor: { expression: { conditions: fillColorConditions } },
            },
          },
        });
      });
    }, 500);
    return () => {
      clearTimeout(timer);
    };
  }, [dataID, items, value, onUpdate]);

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
