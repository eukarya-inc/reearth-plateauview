import { Icon, Dropdown, Menu } from "@web/sharedComponents";
import { styled } from "@web/theme";

import { BaseField as BaseFieldProps } from "..";

import useHooks, { LegendGradientType, LegendStyleType } from "./hooks";

type Props = BaseFieldProps<"legendGradient"> & {
  value: LegendGradientType;
  editMode?: boolean;
};

const LegendGradient: React.FC<Props> = ({ value, editMode }) => {
  const {
    legendStyles,
    Fields,
    currentLegendGradient,
    handleStyleChange,
    handleChooseField,
    handleStepChange,
    handleStartColorChange,
    handleEndColorChange,
  } = useHooks(value);

  const stylesMenu = (
    <Menu
      items={Object.keys(legendStyles).map(ls => {
        return {
          key: ls,
          label: (
            <p style={{ margin: 0 }} onClick={() => handleStyleChange(ls as LegendStyleType)}>
              {legendStyles[ls]}
            </p>
          ),
        };
      })}
    />
  );

  const FieldsMenu = (
    <Menu
      items={Object.keys(Fields).map(ls => {
        return {
          key: ls,
          label: (
            <p style={{ margin: 0 }} onClick={() => handleChooseField(ls)}>
              {Fields[ls]}
            </p>
          ),
        };
      })}
    />
  );
  return editMode ? (
    <Wrapper>
      <Field>
        <FieldTitle>Legend style</FieldTitle>
        <FieldValue>
          <Dropdown overlay={stylesMenu} placement="bottom" trigger={["click"]}>
            <StyledDropdownButton>
              <p style={{ margin: 0 }}>{legendStyles[currentLegendGradient.style]}</p>
              <Icon icon="arrowDownSimple" size={12} />
            </StyledDropdownButton>
          </Dropdown>
        </FieldValue>
      </Field>
      <Field>
        <FieldTitle>Choose Field</FieldTitle>
        <FieldValue>
          <Dropdown overlay={FieldsMenu} placement="bottom" trigger={["click"]}>
            <StyledDropdownButton>
              <p style={{ margin: 0 }}>{"-"}</p>
              <Icon icon="arrowDownSimple" size={12} />
            </StyledDropdownButton>
          </Dropdown>
        </FieldValue>
      </Field>

      {currentLegendGradient.items && (
        <>
          <Field>
            <FieldTitle>Start color</FieldTitle>
            <FieldValue>
              <ColorBlock color={currentLegendGradient.items[0]?.color} />
              <TextInput
                defaultValue={currentLegendGradient.items[0]?.color}
                onChange={handleStartColorChange}
              />
            </FieldValue>
          </Field>
          <Field>
            <FieldTitle>End color</FieldTitle>
            <FieldValue>
              <ColorBlock
                color={currentLegendGradient.items[currentLegendGradient.items.length - 1]?.color}
              />
              <TextInput
                defaultValue={
                  currentLegendGradient.items[currentLegendGradient.items.length - 1].color
                }
                onChange={handleEndColorChange}
              />
            </FieldValue>
          </Field>
        </>
      )}
      <Field>
        <FieldTitle> Step width</FieldTitle>
        <FieldValue>
          <TextInput
            defaultValue={currentLegendGradient.items.length}
            type={"number"}
            onChange={handleStepChange}
          />
        </FieldValue>
      </Field>
    </Wrapper>
  ) : (
    <Wrapper>
      {currentLegendGradient.items?.map((item, idx) => (
        <Field key={idx} gap={12}>
          <ColorBlock color={item?.color} legendStyle={currentLegendGradient.style} />
          <Text>{item?.title}</Text>
        </Field>
      ))}
    </Wrapper>
  );
};
export default LegendGradient;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
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

const Text = styled.p`
  margin: 0;
`;

const Field = styled.div<{ gap?: number }>`
  display: flex;
  align-items: center;
  ${({ gap }) => gap && `gap: ${gap}px;`}
  height: 32px;
`;

const FieldTitle = styled(Text)`
  width: 82px;
`;

const FieldValue = styled.div`
  display: flex;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  flex: 1;
  height: 100%;
  width: 100%;
`;

const TextInput = styled.input`
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
const ColorBlock = styled.div<{
  color?: string;
  legendStyle?: "circle" | "square" | "line";
}>`
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
