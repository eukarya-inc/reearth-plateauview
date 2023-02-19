import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import { swap } from "@web/extensions/sidebar/utils";
import { useCallback, useState } from "react";

import { BaseFieldProps, Cond } from "../types";

import { ButtonWrapper, Wrapper } from "./commonComponents";
import PointColorItem from "./PointColorItem";

const PointColor: React.FC<BaseFieldProps<"pointColor">> = ({ value, editMode, onUpdate }) => {
  const [pointColors, updatePointColors] = useState(value.pointColors);

  const handleMoveUp = useCallback(
    (idx: number) => {
      if (idx === 0) return;
      updatePointColors(prevPointColors => {
        const copy = [...(prevPointColors ?? [])];
        swap(copy, idx, idx - 1);
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
      updatePointColors(prevPointColors => {
        const copy = [...(prevPointColors ?? [])];
        swap(copy, idx, idx + 1);
        onUpdate({
          ...value,
          pointColors: copy,
        });
        return copy;
      });
    },
    [value, pointColors, onUpdate],
  );

  const handleAdd = useCallback(() => {
    updatePointColors(c => {
      const newPointColor: { condition: Cond<number>; color: string } = {
        condition: {
          key: Math.random().toString(16).slice(2),
          operator: "=",
          operand: "width",
          value: 1,
        },
        color: "",
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
      updatePointColors(prevState => {
        const copy = [...(prevState ?? [])].filter((_, idx2) => idx2 != idx);
        onUpdate({
          ...value,
          pointColors: copy,
        });
        return copy;
      });
    },
    [value, onUpdate],
  );

  const handleItemUpdate = (item: { condition: Cond<number>; color: string }, index: number) => {
    updatePointColors(prevState => {
      const copy = [...(prevState ?? [])];
      copy.splice(index, 1, item);
      onUpdate({
        ...value,
        pointColors: copy,
      });
      return copy;
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
