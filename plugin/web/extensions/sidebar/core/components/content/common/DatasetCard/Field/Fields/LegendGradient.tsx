import { Icon, Dropdown, Menu } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { BaseField as BaseFieldProps } from ".";

type LegendStyleType = "square" | "circle" | "line" | "icon";

const legendStyles: { [key: string]: string } = {
  square: "四角",
  circle: "丸",
  line: "線",
  icon: "アイコン",
};

const Fields: { [key: string]: string } = {
  unDefiend: "-",
};

type LegendItem = {
  title: string;
  color: string;
  url?: string;
};

type LegendGradient = {
  id?: string;
  style: LegendStyleType;
  items?: LegendItem[];
};

type Props = BaseFieldProps<"legendGradient"> & {
  value: LegendGradient;
  editMode?: boolean;
};
const LegendGradient: React.FC<Props> = ({ value, editMode }) => {
  const [LegendGradient, updateLegend] = useState<LegendGradient>(value);

  const handleStyleChange = useCallback((style: LegendStyleType) => {
    updateLegend(l => {
      return {
        ...l,
        style,
      };
    });
  }, []);

  const handleChooseField = useCallback((field: string) => {
    console.log(field);
  }, []);

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
              <p style={{ margin: 0 }}>{legendStyles[LegendGradient.style]}</p>
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

      {LegendGradient.style === "icon" && (
        <Field>
          <FieldTitle> Image URL</FieldTitle>
          <FieldValue>
            <TextInput value={"item.url"} />
          </FieldValue>
        </Field>
      )}
      {LegendGradient.items && (
        <>
          <Field>
            <FieldTitle>Start color</FieldTitle>
            <FieldValue>
              <ColorBlock color={LegendGradient.items[LegendGradient.items.length - 1].color} />
              <TextInput value={LegendGradient.items[LegendGradient.items.length - 1].color} />
            </FieldValue>
          </Field>
          <Field>
            <FieldTitle>End color</FieldTitle>
            <FieldValue>
              <ColorBlock color={LegendGradient?.items[0].color} />
              <TextInput value={LegendGradient?.items[0].color} />
            </FieldValue>
          </Field>
        </>
      )}
      <Field>
        <FieldTitle> Step width</FieldTitle>
        <FieldValue>
          <TextInput value={100} />
        </FieldValue>
      </Field>
    </Wrapper>
  ) : (
    <Wrapper>
      {LegendGradient.items?.map((item, idx) => (
        <Field key={idx} gap={12}>
          {LegendGradient.style === "icon" ? (
            <StyledImg src={item.url} />
          ) : (
            <ColorBlock color={item.color} legendStyle={LegendGradient.style} />
          )}
          <Text>{item.title}</Text>
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

const TextInput = styled.input.attrs({ type: "text" })`
  height: 100%;
  width: 100%;
  flex: 1;
  padding: 0 12px;
  border: none;
  outline: none;
  color: rgba(0, 0, 0, 0.25);
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

const StyledImg = styled.img`
  width: 30px;
  height: 30px;
`;
