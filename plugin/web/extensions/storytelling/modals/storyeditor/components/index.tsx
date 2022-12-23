import { styled } from "@web/theme";

import useHooks from "./hooks";

type Props = {};

const StoryEditor: React.FC<Props> = () => {
  const { titleRef, descriptionRef, onSave, onCancel } = useHooks();

  return (
    <Wrapper>
      <TitleInput placeholder="Title" ref={titleRef} />
      <ContentInput placeholder="Content" ref={descriptionRef} />
      <Actions>
        <Button primary onClick={onSave}>
          Save
        </Button>
        <Button onClick={onCancel}>Cancel</Button>
      </Actions>
    </Wrapper>
  );
};

const Wrapper = styled.div`
  position: relative;
  width: 100%;
  height: 100%;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 12px;
`;

const TitleInput = styled.input`
  display: block;
  width: 100%;
  height: 32px;
  padding: 4px 12px;
  background-color: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  font-size: 14px;
  line-height: 24px;
  outline: none;
  flex-shrink: 0;
`;

const ContentInput = styled.textarea`
  display: block;
  width: 100%;
  height: 100%;
  padding: 4px 12px;
  background-color: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  font-size: 14px;
  line-height: 22px;
  outline: none;
  resize: none;
`;

const Actions = styled.div`
  display: flex;
  flex-direction: row-reverse;
  gap: 12px;
`;

const Button = styled.div<{ primary?: boolean }>`
  display: flex;
  align-items: center;
  justify-content: center;
  height: 29px;
  padding: 4px 12px;
  border-radius: 4px;
  background-color: ${({ primary }) => (primary ? "#00bebe" : "#d1d1d1")};
  color: #fff;
  font-size: 14px;
  line-height: 21px;
  cursor: pointer;
`;

export default StoryEditor;
