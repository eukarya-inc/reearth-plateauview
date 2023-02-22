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
  compare,
} from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

import { BaseFieldProps, Cond } from "../../types";

import PointColorItem from "./PointColorItem";

const PointColor: React.FC<BaseFieldProps<"pointColor">> = ({
  dataID,
  value,
  editMode,
  isActive,
  onUpdate,
}) => {
  const [pointColors, updatePointColors] = useState(value.pointColors);

  const operandOptions = [{ value: "pointSize", label: "size" }];

  const handleMoveUp = useCallback(
    (idx: number) => {
      updatePointColors(c => {
        const newPointColors = moveItemUp(idx, c) ?? c;
        onUpdate({
          ...value,
          pointColors: newPointColors,
        });
        return newPointColors;
      });
    },
    [value, onUpdate],
  );

  const handleMoveDown = useCallback(
    (idx: number) => {
      updatePointColors(c => {
        const newPointColors = moveItemDown(idx, c) ?? c;
        onUpdate({
          ...value,
          pointColors: newPointColors,
        });
        return newPointColors;
      });
    },
    [onUpdate, value],
  );

  const handleAdd = useCallback(() => {
    updatePointColors(c => {
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
        pointColors: value.pointColors ? [...value.pointColors, newPointColor] : [newPointColor],
      });
      return c ? [...c, newPointColor] : [newPointColor];
    });
  }, [value, onUpdate]);

  const handleRemove = useCallback(
    (idx: number) => {
      updatePointColors(c => {
        const newPointColors = removeItem(idx, c) ?? c;
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
      const newPointColors = [...(c ?? [])];
      newPointColors.splice(index, 1, item);
      onUpdate({
        ...value,
        pointColors: newPointColors,
      });
      return newPointColors;
    });
  };

  useEffect(() => {
    if (!isActive || !dataID) return;
    const timer = setTimeout(() => {
      pointColors?.forEach(item => {
        const operand = 1; // should be something like item[pointSize]
        const pointColor = compare(operand, item.condition.operator, item.condition.value)
          ? item.color
          : "";
        postMsg({
          action: "updateDatasetInScene",
          payload: {
            dataID,
            update: { marker: { style: "point", pointColor } },
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
          update: { marker: undefined },
        },
      });
    };
  }, [dataID, isActive, pointColors]);

  return editMode ? (
    <Wrapper>
      {pointColors?.map((c, idx) => (
        <PointColorItem
          key={idx}
          index={idx}
          item={c}
          operandOptions={operandOptions}
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
