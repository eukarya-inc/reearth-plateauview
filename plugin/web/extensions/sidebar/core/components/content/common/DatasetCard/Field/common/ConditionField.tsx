import { Select } from "@web/sharedComponents";
import { useEffect, useState } from "react";

import { FieldTitle, FieldValue, FieldWrapper, NumberInput } from "../commonComponents";
import { Cond } from "../Fields/types";

const operatorOptions = [
  { value: ">", label: ">" },
  { value: "<", label: "<" },
  { value: ">=", label: ">=" },
  { value: "<=", label: "<=" },
  { value: "=", label: "=" },
];

type Props = {
  title: string;
  fieldGap?: number;
  condition: Cond<any>;
  operandOptions: { value: string; label: string }[];
  onChange?: (condition: Cond<any>) => void;
};

const ConditionField: React.FC<Props> = ({
  title,
  fieldGap,
  condition,
  operandOptions,
  onChange,
}) => {
  const [cond, setCond] = useState<Cond<any>>(condition);

  useEffect(() => {
    setCond(condition);
  }, [condition]);

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
