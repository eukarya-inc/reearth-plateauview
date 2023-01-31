import AddButton from "@web/extensions/sidebar/core/components/content/common/AddButton";
import { Icon, Dropdown, Menu } from "@web/sharedComponents";
import { styled } from "@web/theme";

import { BaseField as BaseFieldProps } from "..";

import useHooks, { SwitchGroupObj } from "./hooks";

type Props = BaseFieldProps<"legend"> & {
  value: SwitchGroupObj;
  editMode?: boolean;
};
const SwitchGroup: React.FC<Props> = ({ value, editMode }) => {
  const {
    groups,
    groupsTitle,
    handleMoveDown,
    handleMoveUp,
    handleAdd,
    handleRemove,
    handleTitleChange,
  } = useHooks(value);

  const menu = (
    <Menu
      items={Object.keys(groups.groups).map(ls => {
        return {
          key: ls,
          label: <p style={{ margin: 0 }}>{ls}</p>,
        };
      })}
    />
  );

  return editMode ? (
    <Wrapper>
      <Field>
        <FieldTitle>題名</FieldTitle>
        <FieldValue>
          <TextInput defaultValue={groupsTitle} onChange={handleTitleChange} />
        </FieldValue>
      </Field>
      <AddButton text="項目" onClick={handleAdd} />
      {groups.groups?.map((item, idx) => (
        <Item key={idx}>
          <ItemControls>
            <Icon icon="arrowUpThin" size={16} onClick={() => handleMoveUp(idx)} />
            <Icon icon="arrowDownThin" size={16} onClick={() => handleMoveDown(idx)} />
            <Icon icon="trash" size={16} onClick={() => handleRemove(idx)} />
          </ItemControls>
          <Field>
            <FieldTitle>グループ</FieldTitle>
            <FieldValue>
              <Dropdown overlay={menu} placement="bottom" trigger={["click"]}>
                <StyledDropdownButton>
                  <p style={{ margin: 0 }}>{groups.groups[0].title}</p>
                  <Icon icon="arrowDownSimple" size={12} />
                </StyledDropdownButton>
              </Dropdown>
            </FieldValue>
          </Field>
          <Field>
            <FieldTitle>題名</FieldTitle>
            <FieldValue>
              <TextInput value={item.title} />
            </FieldValue>
          </Field>
        </Item>
      ))}
    </Wrapper>
  ) : (
    <Wrapper>
      {/* {legend.items?.map((item, idx) => (
        <Field key={idx} gap={12}>
          {legend.style === "icon" ? (
            <StyledImg src={item.url} />
          ) : (
            <ColorBlock color={item.color} legendStyle={legend.style} />
          )}
          <Text>{item.title}</Text>
        </Field>
      ))} */}
      <div>test</div>
    </Wrapper>
  );
};

export default SwitchGroup;

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

const Item = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  padding: 8px;
`;

const ItemControls = styled.div`
  display: flex;
  justify-content: right;
  gap: 4px;
  cursor: pointer;
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

  :focus {
    border: none;
  }
`;

// const ColorBlock = styled.div<{ color: string; legendStyle?: "circle" | "square" | "line" }>`
//   width: 30px;
//   height: ${({ legendStyle }) => (legendStyle === "line" ? "3px" : "30px")};
//   background: ${({ color }) => color ?? "#d9d9d9"};
//   border-radius: ${({ legendStyle }) =>
//     legendStyle
//       ? legendStyle === "circle"
//         ? "50%"
//         : legendStyle === "line"
//         ? "5px"
//         : "2px"
//       : "1px 0 0 1px"};
// `;

// const StyledImg = styled.img`
//   width: 30px;
//   height: 30px;
// `;
