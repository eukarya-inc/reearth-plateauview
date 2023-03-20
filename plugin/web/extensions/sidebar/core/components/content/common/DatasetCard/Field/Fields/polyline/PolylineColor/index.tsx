import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import {
  ButtonWrapper,
  Wrapper,
} from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { generateID, moveItemDown, moveItemUp, removeItem } from "@web/extensions/sidebar/utils";
import { styled, commonStyles } from "@web/theme";
import { useCallback, useEffect, useState, useRef } from "react";

import { stringifyCondition } from "../../../utils";
import { BaseFieldProps, Cond } from "../../types";

import PolylineColorItem from "./PolylineColorItem";

const PolylineColor: React.FC<BaseFieldProps<"polylineColor">> = ({
  value,
  editMode,
  onUpdate,
}) => {
  const [items, updateItems] = useState(value.items);

  const handleMoveUp = useCallback((idx: number) => {
    updateItems(c => {
      const newItems = moveItemUp(idx, c) ?? c;
      return newItems;
    });
  }, []);

  const handleMoveDown = useCallback((idx: number) => {
    updateItems(c => {
      const newItems = moveItemDown(idx, c) ?? c;
      return newItems;
    });
  }, []);

  const handleAdd = useCallback(() => {
    updateItems(c => {
      const newItem: { condition: Cond<any>; color: string } = {
        condition: {
          key: generateID(),
          operator: "===",
          operand: true,
          value: true,
        },
        color: "",
      };
      return c ? [...c, newItem] : [newItem];
    });
  }, []);

  const handleRemove = useCallback((idx: number) => {
    updateItems(c => {
      const newItems = removeItem(idx, c) ?? c;
      return newItems;
    });
  }, []);

  const handleItemUpdate = (item: { condition: Cond<number>; color: string }, index: number) => {
    updateItems(c => {
      const newItems = [...(c ?? [])];
      newItems.splice(index, 1, item);
      return newItems;
    });
  };

  const generateOverride = useCallback((items: typeof value.items) => {
    const strokeColorConditions: [string, string][] = [["true", 'color("white")']];
    items?.forEach(item => {
      const resStrokeColor = "color" + `("${item.color}")`;
      const cond = stringifyCondition(item.condition);
      strokeColorConditions.unshift([cond, resStrokeColor]);
    });
    return {
      polyline: {
        strokeColor: {
          expression: {
            conditions: strokeColorConditions,
          },
        },
      },
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const [override, updateOverride] = useState<{ polyline: any }>(generateOverride(value.items));
  const valueRef = useRef(value);
  const itemsRef = useRef(items);
  const onUpdateRef = useRef(onUpdate);
  valueRef.current = value;
  itemsRef.current = items;
  onUpdateRef.current = onUpdate;

  const handleApply = useCallback(() => {
    updateOverride(generateOverride(items));
  }, [generateOverride, items]);

  useEffect(() => {
    onUpdateRef.current({
      ...valueRef.current,
      items: itemsRef.current,
      override,
    });
  }, [override]);

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
      <Button onClick={handleApply}>Apply</Button>
    </Wrapper>
  ) : null;
};

const Button = styled.div`
  ${commonStyles.simpleButton}
`;

export default PolylineColor;
