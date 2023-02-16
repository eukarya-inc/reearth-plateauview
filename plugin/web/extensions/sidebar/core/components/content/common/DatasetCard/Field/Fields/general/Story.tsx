import AddButton from "@web/extensions/sidebar/core/components/content/common/DatasetCard/AddButton";
import { array_move, generateID } from "@web/extensions/sidebar/utils";
import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { BaseFieldProps, StoryItem } from "../types";

const Story: React.FC<BaseFieldProps<"story">> = ({ value, editMode, onUpdate }) => {
  const [, updateState] = useState<any>();
  const forceUpdate = useCallback(() => updateState({}), []);

  const handleStoryAdd = useCallback(() => {
    const newStory: StoryItem = {
      id: generateID(),
      title: "title",
    };
    onUpdate({ ...value, stories: value.stories ? [...value.stories, newStory] : [newStory] });
    forceUpdate();
  }, [forceUpdate, onUpdate, value]);

  const handleItemMoveUp = useCallback(
    (idx: number) => {
      if (idx === 0 || !value.stories) return;
      const newStories = [...value.stories];
      array_move(newStories, idx, idx - 1);
      onUpdate({ ...value, stories: newStories });
      forceUpdate();
    },
    [onUpdate, forceUpdate, value],
  );

  const handleItemMoveDown = useCallback(
    (idx: number) => {
      if (!value.stories || idx >= value.stories.length - 1) return;
      const newStories = [...value.stories];
      array_move(newStories, idx, idx + 1);
      onUpdate({ ...value, stories: newStories });
      forceUpdate();
    },
    [onUpdate, forceUpdate, value],
  );

  const handleItemRemove = useCallback(
    (id: string) => {
      const newStories = value.stories?.filter(st => st.id !== id);
      if (!newStories) return;
      onUpdate({ ...value, stories: newStories });
      forceUpdate();
    },
    [onUpdate, forceUpdate, value],
  );

  const handleStoryTitleChange = useCallback(
    (title: string, index: number) => {
      if (!value.stories) return;
      const updatedStories = value.stories;
      updatedStories[index].title = title;
      onUpdate({ ...value, stories: updatedStories });
      forceUpdate();
    },
    [onUpdate, forceUpdate, value],
  );

  const handleStoryEdit = useCallback(() => {}, []);
  const handleStoryShow = useCallback(() => {}, []);

  return editMode ? (
    <Wrapper>
      <AddButton text="New Story" onClick={handleStoryAdd} />
      {value.stories?.map((g, idx) => (
        <Item key={idx}>
          <ItemControls>
            <Icon icon="arrowUpThin" size={16} onClick={() => handleItemMoveUp(idx)} />
            <Icon icon="arrowDownThin" size={16} onClick={() => handleItemMoveDown(idx)} />
            <TrashIcon icon="trash" size={16} onClick={() => handleItemRemove(g.id)} />
          </ItemControls>
          <Field>
            <FieldTitle>Title </FieldTitle>
            <FieldValue>
              <TextInput
                defaultValue={g.title}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  handleStoryTitleChange(e.target.value, idx);
                }}
              />
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
      {value.stories?.map(story => (
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
