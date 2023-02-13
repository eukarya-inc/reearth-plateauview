import AddButton from "@web/extensions/sidebar/core/components/content/common/AddButton";
import { Icon, Dropdown, Menu } from "@web/sharedComponents";
import { styled } from "@web/theme";

import { BaseFieldProps } from "../../types";

import useHooks from "./hooks";

const SwitchGroup: React.FC<BaseFieldProps<"switchGroup">> = ({
  value,
  editMode,
  fieldGroups,
  onUpdate,
}) => {
  const {
    title,
    groupItems,
    selectedGroup,
    handleTitleChange,
    handleGroupChoose,
    handleItemGroupChange,
    handleItemTitleChange,
    handleItemAdd,
    handleItemRemove,
    handleItemMoveUp,
    handleItemMoveDown,
  } = useHooks({
    value,
    fieldGroups,
    onUpdate,
  });

  const menu = (
    <Menu
      items={groupItems?.map((gi, idx) => {
        return {
          key: idx,
          label: (
            <p style={{ margin: 0 }} onClick={() => handleGroupChoose(gi.groupID)}>
              {gi.title}
            </p>
          ),
        };
      })}
    />
  );

  return editMode ? (
    <Wrapper>
      <Field>
        <FieldTitle>タイトル</FieldTitle>
        <FieldValue>
          <TextInput defaultValue={title} onChange={handleTitleChange} />
        </FieldValue>
      </Field>
      <AddButton text="Add Item" onClick={handleItemAdd} />
      {value.groups?.map((g, idx) => (
        <Item key={idx}>
          <ItemControls>
            <Icon icon="arrowUpThin" size={16} onClick={() => handleItemMoveUp(idx)} />
            <Icon icon="arrowDownThin" size={16} onClick={() => handleItemMoveDown(idx)} />
            <Icon icon="trash" size={16} onClick={() => handleItemRemove(idx)} />
          </ItemControls>
          <Field>
            <FieldTitle>グループ</FieldTitle>
            <FieldValue>
              <SelectWrapper
                defaultValue={g.title}
                onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
                  handleItemGroupChange(
                    idx,
                    fieldGroups?.find(fg => fg.name === e.target.value)?.id,
                  );
                }}>
                {fieldGroups?.map((fg, idx) => (
                  <option
                    key={fg.id}
                    defaultChecked={idx === 0}
                    defaultValue={fg.name}
                    disabled={fg.id === g.groupID ? true : false}>
                    {fg.name}
                  </option>
                ))}
              </SelectWrapper>
            </FieldValue>
          </Field>
          <Field>
            <FieldTitle>名前</FieldTitle>
            <FieldValue>
              <TextInput
                defaultValue={g.title}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  handleItemTitleChange(e.target.value, idx);
                }}
              />
            </FieldValue>
          </Field>
        </Item>
      ))}
    </Wrapper>
  ) : (
    <Wrapper>
      <Field>
        <FieldTitle>{title}</FieldTitle>
        <FieldValue>
          <Dropdown overlay={menu} placement="bottom" trigger={["click"]}>
            <StyledDropdownButton>
              <p style={{ margin: 0 }}>{selectedGroup ? selectedGroup.name : "-"}</p>
              <Icon icon="arrowDownSimple" size={12} />
            </StyledDropdownButton>
          </Dropdown>
        </FieldValue>
      </Field>
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

const SelectWrapper = styled.select`
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  align-content: center;
  padding: 0 16px;
  cursor: pointer;
`;
