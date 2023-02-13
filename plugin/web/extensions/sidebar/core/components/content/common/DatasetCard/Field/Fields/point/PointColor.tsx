import AddButton from "@web/extensions/sidebar/core/components/content/common/AddButton";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { BaseFieldProps, Cond } from "../types";

import ConditionField from "./ConditionField";
import Field from "./Field";
import ItemControls from "./ItemControls";

function array_move(arr: any[], old_index: number, new_index: number) {
  if (new_index >= arr.length) {
    let k = new_index - arr.length + 1;
    while (k--) {
      arr.push(undefined);
    }
  }
  arr.splice(new_index, 0, arr.splice(old_index, 1)[0]);
}

const PointColor: React.FC<BaseFieldProps<"pointColor">> = ({ value, editMode, onUpdate }) => {
  const [pointColors, updatePointColors] = useState(value.pointColors);

  const handleMoveUp = useCallback(
    (idx: number) => {
      if (idx === 0) return;
      updatePointColors(c => {
        let newPointColors: { condition: Cond<number>; color: string }[] | undefined = undefined;
        if (c) {
          newPointColors = c;
          array_move(newPointColors, idx, idx - 1);
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

  return editMode ? (
    <Wrapper>
      {value.pointColors?.map((c, idx) => (
        <Item key={idx}>
          <ItemControls
            index={idx}
            handleMoveDown={handleMoveDown}
            handleMoveUp={handleMoveUp}
            handleRemove={handleRemove}
          />
          {/* {legend.style === "icon" && (
            <Field>
              <FieldTitle>URL</FieldTitle>
              <FieldValue>
                <TextInput value={item.url} />
              </FieldValue>
            </Field>
          )} */}
          <ConditionField title="if" fieldGap={8} condition={c.condition} />
          <Field
            title="色"
            titleWidth={82}
            value={
              <>
                <ColorBlock color={c.color} />
                <TextInput value={c.color} />
              </>
            }
          />
        </Item>
      ))}
      <ButtonWrapper>
        <AddButton text="Add Condition" height={24} onClick={handleAdd} />
      </ButtonWrapper>
    </Wrapper>
  ) : null;
};

export default PointColor;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

// const StyledDropdownButton = styled.div`
//   display: flex;
//   justify-content: space-between;
//   align-items: center;
//   width: 100%;
//   align-content: center;
//   padding: 0 16px;
//   cursor: pointer;
// `;

const Item = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  padding: 8px;
`;

const TextInput = styled.input.attrs({ type: "text" })`
  height: 100%;
  width: 100%;
  flex: 1;
  padding: 0 12px;
  border: none;
  outline: none;

  :focus {
    border: none;
  }
`;

const ColorBlock = styled.div<{ color: string; legendStyle?: "circle" | "square" | "line" }>`
  width: 30px;
  height: ${({ legendStyle }) => (legendStyle === "line" ? "3px" : "30px")};
  background: ${({ color }) => color ?? "#d9d9d9"};
  border-radius: ${({ legendStyle }) =>
    legendStyle
      ? legendStyle === "circle"
        ? "50%"
        : legendStyle === "line"
        ? "5px"
        : "2px"
      : "1px 0 0 1px"};
`;

// const StyledImg = styled.img`
//   width: 30px;
//   height: 30px;
// `;

const ButtonWrapper = styled.div`
  width: 125px;
  align-self: flex-end;
`;
