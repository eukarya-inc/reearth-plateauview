import { styled } from "@web/theme";

import { BaseFieldProps } from "../../types";

import useHooks from "./hooks";

const StyleCode: React.FC<BaseFieldProps<"styleCode">> = ({
  value,
  dataID,
  editMode,
  onUpdate,
}) => {
  const { code, handleEditCode } = useHooks({ value, dataID, onUpdate });

  return editMode ? (
    <Wrapper>
      <Field>
        <CodeEditor value={code} onChange={handleEditCode} />
      </Field>
    </Wrapper>
  ) : null;
};

export default StyleCode;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Field = styled.div<{ gap?: number }>`
  display: flex;
  align-items: center;
  ${({ gap }) => gap && `gap: ${gap}px;`}
  height: 160px;
  padding: 8px;
  gap: 8px;
`;

const CodeEditor = styled.textarea`
  height: 144px;
  width: 280px;
  flex: 1;
  padding: 0 12px;
  border: none;
  overflow: auto;
  background: #f3f3f3;
  outline: none;
  resize: none;
  :focus {
    border: none;
  }
`;
