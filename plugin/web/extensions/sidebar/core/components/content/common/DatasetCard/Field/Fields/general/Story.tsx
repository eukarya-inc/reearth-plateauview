import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import { generateID, moveItemDown, moveItemUp } from "@web/extensions/sidebar/utils";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { BaseFieldProps } from "../types";

const Story: React.FC<BaseFieldProps<"story">> = ({ value, editMode, onUpdate }) => {
  const [stories, updateStories] = useState(value.stories);

  const handleStoryAdd = useCallback(() => {
    updateStories(s => {
      const newItem = {
        id: generateID(),
        title: "title",
      };
      const newStories = s ? [...s, newItem] : [newItem];
      onUpdate({ ...value, stories: newStories });
      return newStories;
    });
  }, [onUpdate, value]);

  const handleItemMoveUp = useCallback(
    (idx: number) => {
      updateStories(s => {
        const newStories = moveItemUp(idx, s) ?? s;
        onUpdate({ ...value, stories: newStories });
        return newStories;
      });
    },
    [onUpdate, value],
  );

  const handleItemMoveDown = useCallback(
    (idx: number) => {
      updateStories(s => {
        const newStories = moveItemDown(idx, s) ?? s;
        onUpdate({ ...value, stories: newStories });
        return newStories;
      });
    },
    [onUpdate, value],
  );

  const handleItemRemove = useCallback(
    (id: string) => {
      updateStories(s => {
        const newStories = s?.filter(st => st.id !== id);
        onUpdate({ ...value, stories: newStories });
        return newStories;
      });
    },
    [onUpdate, value],
  );

  const handleStoryTitleChange = useCallback(
    (index: number) => (e: React.ChangeEvent<HTMLInputElement>) => {
      updateStories(s => {
        if (!s) return s;
        const updatedStories = s;
        updatedStories[index].title = e.currentTarget.value;
        onUpdate({ ...value, stories: updatedStories });
        return updatedStories;
      });
    },
    [onUpdate, value],
  );

  const handleStoryEdit = useCallback(() => {}, []);
  const handleStoryShow = useCallback(() => {}, []);

  return editMode ? (
    <Wrapper>
      <AddButton text="New Story" onClick={handleStoryAdd} />
      {stories?.map((g, idx) => (
        <Item key={idx}>
          <ItemControls>
            <Icon icon="arrowUpThin" size={16} onClick={() => handleItemMoveUp(idx)} />
            <Icon icon="arrowDownThin" size={16} onClick={() => handleItemMoveDown(idx)} />
            <TrashIcon icon="trash" size={16} onClick={() => handleItemRemove(g.id)} />
          </ItemControls>
          <Field>
            <FieldTitle>Title </FieldTitle>
            <FieldValue>
              <TextInput defaultValue={g.title} onChange={handleStoryTitleChange(idx)} />
            </FieldValue>
          </Field>
          <EditButton onClick={handleStoryEdit}>
            <Icon icon="edit" size={14} />
            <Text>Edit Story</Text>
          </EditButton>
        </Item>
      ))}
    </Wrapper>
  ) : (
    <Wrapper>
      {stories?.map(story => (
        <StoryButton key={story.id} onClick={handleStoryShow}>
          <Icon icon="circledPlay" size={24} />
          <Text>{story.title}</Text>
        </StoryButton>
      ))}
    </Wrapper>
  );
};

export default Story;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Text = styled.p`
  margin: 0;
  line-height: 24px;
  font-weight: 400;
  font-size: 14px;
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
const TrashIcon = styled(Icon)<{ disabled?: boolean }>`
  ${({ disabled }) =>
    disabled &&
    `
      color: rgb(209, 209, 209);
      pointer-events: none;
    `}
`;

const EditButton = styled.div`
  margin-top: 12px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  background: #ffffff;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  padding: 5px;
  height: 32px;
  cursor: pointer;

  :hover {
    background: #f4f4f4;
  }
`;
const StoryButton = styled.div`
  display: flex;
  align-items: center;
  background: #ffffff;
  box-shadow: 0px 2px 8px rgba(0, 0, 0, 0.15);
  border-radius: 4px;
  padding: 12px;
  gap: 12px;
  height: 48px;
  cursor: pointer;
  :hover {
    background: #f4f4f4;
  }
`;
