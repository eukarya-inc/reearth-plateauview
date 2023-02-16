import { Select, InputNumber } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useState } from "react";

import { Cond } from "../../types";
import { FieldTitle, FieldValue, FieldWrapper } from "../commonComponents";

const operatorOptions = [
  { value: ">", label: ">" },
  { value: "<", label: "<" },
  { value: ">=", label: ">=" },
  { value: "<=", label: "<=" },
  { value: "=", label: "=" },
];

const operandOptions = [
  { value: "height", label: "height" },
  { value: "width", label: "width" },
];

type Props = {
  title: string;
  fieldGap?: number;
  condition: Cond<any>;
  onChange?: (condition: Cond<any>) => void;
};

const ConditionField: React.FC<Props> = ({ title, fieldGap, condition, onChange }) => {
  const [cond, setCond] = useState<Cond<any>>(condition);

  const handleOperandChange = (operand: any) => {
    setCond(prevCond => {
      const copy = { ...prevCond, operand };
      onChange?.(copy);
      return copy;
    });
  };

  const handleOperatorChange = (operator: any) => {
    setCond(prevCond => {
      const copy = { ...prevCond, operator };
      onChange?.(copy);
      return copy;
    });
  };

  const handleValueChange = (value: any) => {
    setCond(prevCond => {
      const copy = { ...prevCond, value };
      onChange?.(copy);
      return copy;
    });
  };

  return (
    <FieldWrapper gap={fieldGap}>
      <FieldTitle>{title}</FieldTitle>
      <FieldValue noBorder>
        <Select
          options={operandOptions}
          style={{ width: "100%" }}
          value={cond.operand}
          onChange={handleOperandChange}
        />
      </FieldValue>
      <FieldValue noBorder>
        <Select
          options={operatorOptions}
          style={{ width: "100%" }}
          value={cond.operator}
          onChange={handleOperatorChange}
        />
      </FieldValue>
      <FieldValue>
        <NumberInput value={cond.value} onChange={handleValueChange} />
      </FieldValue>
    </FieldWrapper>
  );
};

export default ConditionField;

const NumberInput = styled(InputNumber)`
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
