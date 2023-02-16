import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import { array_move } from "@web/extensions/sidebar/utils";
import { useCallback, useState } from "react";

import { BaseFieldProps, Cond } from "../types";

import { ButtonWrapper, Wrapper } from "./commonComponents";
import PointColorItem from "./PointColorItem";

const PointColor: React.FC<BaseFieldProps<"pointColor">> = ({ value, editMode, onUpdate }) => {
  const [pointColors, updatePointColors] = useState(value.pointColors);

  const swapElements = (array: any[], index1: number, index2: number) => {
    const temp = array[index1];
    array[index1] = array[index2];
    array[index2] = temp;
  };

  const handleMoveUp = useCallback(
    (idx: number) => {
      if (idx === 0) return;
      updatePointColors(prevPointColors => {
        const copy = [...(prevPointColors ?? [])];
        swapElements(copy, idx, idx - 1);
        onUpdate({
          ...value,
          pointColors: copy,
        });
        return copy;
      });
    },
    [value, onUpdate],
  );

  const handleMoveDown = useCallback(
    (idx: number) => {
      if (pointColors && idx >= pointColors.length - 1) return;
      updatePointColors(c => {
        let newPointColors: { condition: Cond<number>; color: string }[] | undefined = undefined;
        if (c) {
          newPointColors = c;
          array_move(newPointColors, idx, idx + 1);
        }
        onUpdate({
          ...value,
          pointColors: newPointColors,
        });
        return newPointColors;
      });
    },
    [value, pointColors, onUpdate],
  );

  const handleAdd = useCallback(() => {
    updatePointColors(c => {
      const newPointColor: { condition: Cond<number>; color: string } = {
        condition: {
          key: "ARGH",
          operator: "=",
          operand: "AField",
          value: 1,
        },
        color: "brown",
      };
      onUpdate({
        ...value,
        pointColors: value.pointColors ? [...value.pointColors, newPointColor] : [newPointColor],
      });
      return c ? [...c, newPointColor] : [newPointColor];
    });
  }, [value, onUpdate]);

  const handleRemove = useCallback(
    (idx: number) => {
      updatePointColors(c => {
        let newPointColors: { condition: Cond<number>; color: string }[] | undefined = undefined;
        if (c) {
          newPointColors = c.filter((_, idx2) => idx2 != idx);
        }
        onUpdate({
          ...value,
          pointColors: newPointColors,
        });
        return newPointColors;
      });
    },
    [value, onUpdate],
  );

  const handleItemUpdate = (item: { condition: Cond<number>; color: string }, index: number) => {
    updatePointColors(c => {
      let newPointColors: { condition: Cond<number>; color: string }[] | undefined = undefined;
      if (c) {
        newPointColors = c;
        newPointColors.splice(index, 1, item);
      }
      onUpdate({
        ...value,
        pointColors: newPointColors,
      });
      return newPointColors;
    });
  };

  return editMode ? (
    <Wrapper>
      {pointColors?.map((c, idx) => (
        <PointColorItem
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

export default PointColor;
