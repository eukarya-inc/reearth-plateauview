import { Select, InputNumber } from "@web/sharedComponents";
import { styled } from "@web/theme";

import { Cond } from "../../types";
import { FieldTitle, FieldValue, FieldWrapper } from "../commonComponents";

const operatorOptions = [
  { value: "greater", label: ">" },
  { value: "less", label: "<" },
  { value: "greaterEqual", label: ">=" },
  { value: "lessEqual", label: "<=" },
  { value: "equal", label: "=" },
];

const operandOptions = [
  { value: "height", label: "Height" },
  { value: "Width", label: "Width" },
];

type Props = {
  title: string;
  fieldGap?: number;
  condition: Cond<any>;
  onChange?: (condition: Cond<any>) => void;
};

const ConditionField: React.FC<Props> = ({ title, fieldGap, condition, onChange }) => {
  const handleOperandChange = (value: any) => {
    const cond = { ...condition, value };
    onChange?.(cond);
  };

  const handleOperatorChange = (value: any) => {
    const cond = { ...condition, value };
    onChange?.(cond);
  };

  const handleValueChange = (value: any) => {
    const cond = { ...condition, value };
    onChange?.(cond);
  };

  return (
    <FieldWrapper gap={fieldGap}>
      <FieldTitle>{title}</FieldTitle>
      <FieldValue noBorder>
        <Select options={operandOptions} style={{ width: "100%" }} onChange={handleOperandChange} />
      </FieldValue>
      <FieldValue noBorder>
        <Select
          options={operatorOptions}
          style={{ width: "100%" }}
          onChange={handleOperatorChange}
        />
      </FieldValue>
      <FieldValue>
        <NumberInput value={condition.value} onChange={handleValueChange} />
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
