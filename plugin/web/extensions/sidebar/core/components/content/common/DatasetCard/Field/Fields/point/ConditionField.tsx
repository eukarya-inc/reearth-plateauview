import { Icon, Dropdown, Menu } from "@web/sharedComponents";
import { styled } from "@web/theme";

import { Cond } from "../types";

const operators: { [key: string]: string } = {
  greater: ">",
  less: "<",
  greaterEqual: ">=",
  lessEqual: "<=",
  equal: "=",
};

type Props = {
  title: string;
  titleWidth?: number;
  fieldGap?: number;
  condition: Cond<any>;
};

const ConditionField: React.FC<Props> = ({ title, fieldGap, condition }) => {
  const menu = (
    <Menu
      items={Object.keys(operators).map(op => {
        return {
          key: op,
          label: (
            <p style={{ margin: 0 }} onClick={() => console.log(op)}>
              {operators[op]}
            </p>
          ),
        };
      })}
    />
  );

  return (
    <Field gap={fieldGap}>
      <FieldTitle>{title}</FieldTitle>
      <FieldValue>
        <Dropdown overlay={menu} placement="bottom" trigger={["click"]}>
          <StyledDropdownButton>
            <p style={{ margin: 0 }}>{condition.operand}</p>
            <Icon icon="arrowDownSimple" size={12} />
          </StyledDropdownButton>
        </Dropdown>
      </FieldValue>
      <FieldValue>
        <Dropdown overlay={menu} placement="bottom" trigger={["click"]}>
          <StyledDropdownButton>
            <p style={{ margin: 0 }}>{condition.operator}</p>
            <Icon icon="arrowDownSimple" size={12} />
          </StyledDropdownButton>
        </Dropdown>
      </FieldValue>
      <FieldValue>
        <NumberInput value={condition.value} />
      </FieldValue>
    </Field>
  );
};

export default ConditionField;

const Text = styled.p`
  margin: 0;
`;

const Field = styled.div<{ gap?: number }>`
  display: flex;
  align-items: center;
  ${({ gap }) => gap && `gap: ${gap}px;`}
  height: 32px;
`;

const FieldTitle = styled(Text) <{ width?: number }>`
  ${({ width }) => width && `width: ${width}px;`}
`;

const FieldValue = styled.div`
  display: flex;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  flex: 1;
  height: 100%;
  width: 100%;
`;

const NumberInput = styled.input.attrs({ type: "number" })`
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

const StyledDropdownButton = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  align-content: center;
  padding: 0 16px;
  cursor: pointer;
`;
