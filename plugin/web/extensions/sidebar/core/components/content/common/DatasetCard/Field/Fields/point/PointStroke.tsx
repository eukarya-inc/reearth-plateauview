import AddButton from "@web/extensions/sidebar/core/components/content/common/AddButton";
import { array_move } from "@web/extensions/sidebar/utils";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { BaseFieldProps, Cond, Expression, Fields } from "../types";

import ConditionField from "./common/ConditionField";
import Field from "./common/Field";
import ItemControls from "./common/ItemControls";

const PointStroke: React.FC<BaseFieldProps<"pointStroke">> = ({ value, editMode, onUpdate }) => {
  const [conditions, updateConditions] = useState(value.conditions);

  const handleMoveUp = useCallback(
    (idx: number) => {
      if (idx === 0) return;
      updateConditions(c => {
        let newConditions: Fields["pointStroke"]["conditions"] = undefined;

        if (c) {
          newConditions = c;
          array_move(newConditions, idx, idx - 1);
        }
        onUpdate({
          ...value,
          conditions: newConditions,
        });
        return newConditions;
      });
    },
    [value, onUpdate],
  );

  const handleMoveDown = useCallback(
    (idx: number) => {
      if (conditions && idx >= conditions.length - 1) return;
      updateConditions(c => {
        let newConditions: Fields["pointStroke"]["conditions"] = undefined;
        if (c) {
          newConditions = c;
          array_move(newConditions, idx, idx + 1);
        }
        onUpdate({
          ...value,
          conditions: newConditions,
        });
        return newConditions;
      });
    },
    [conditions, onUpdate, value],
  );

  const handleAdd = useCallback(() => {
    updateConditions(c => {
      const newCondition: {
        expression: Expression;
        strokeColor: string;
        strokeWidth: number;
      } = {
        expression: {
          conditions: [
            {
              key: "ARGH",
              operator: "=",
              operand: "AField",
              value: 1,
            },
          ],
        },
        strokeColor: "brown",
        strokeWidth: 10,
      };
      onUpdate({
        ...value,
        conditions: value.conditions ? [...value.conditions, newCondition] : [newCondition],
      });
      return c ? [...c, newCondition] : [newCondition];
    });
  }, [value, onUpdate]);

  const handleRemove = useCallback(
    (idx: number) => {
      updateConditions(c => {
        let newConditions: Fields["pointStroke"]["conditions"] = undefined;
        if (c) {
          newConditions = c.filter((_, idx2) => idx2 != idx);
        }
        onUpdate({
          ...value,
          conditions: newConditions,
        });
        return newConditions;
      });
    },
    [value, onUpdate],
  );

  // this is hard-codded condition
  const condition: Cond<number> = {
    key: "ARGH",
    operator: "=",
    operand: "AField",
    value: 1,
  };

  return editMode ? (
    <Wrapper>
      {value.conditions?.map((c, idx) => (
        <Item key={idx}>
          <ItemControls
            index={idx}
            handleMoveDown={handleMoveDown}
            handleMoveUp={handleMoveUp}
            handleRemove={handleRemove}
          />
          <ConditionField title="if" fieldGap={8} condition={condition} />
          <Field
            title="strokeColor"
            titleWidth={82}
            value={
              <>
                <ColorBlock color={c.strokeColor} />
                <TextInput value={c.strokeColor} />
              </>
            }
          />
          <Field title="strokeWidth" titleWidth={82} value={<TextInput value={c.strokeWidth} />} />
        </Item>
      ))}
      <ButtonWrapper>
        <AddButton text="Add Condition" height={24} onClick={handleAdd} />
      </ButtonWrapper>
    </Wrapper>
  ) : null;
};

export default PointStroke;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

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

const ButtonWrapper = styled.div`
  width: 125px;
  align-self: flex-end;
`;
