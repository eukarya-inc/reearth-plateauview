import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { BaseFieldProps } from "../types";

const StyleCode: React.FC<BaseFieldProps<"styleCode">> = ({ editMode }) => {
  const [code, editCode] = useState<string>();

  const handleEditCode = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    editCode(e.currentTarget.value);
  }, []);

  return editMode ? (
    <Wrapper>
      <Field>
        <FieldTitle>Title</FieldTitle>
        <FieldValue>
          <TextInput value={code} onChange={handleEditCode} />
        </FieldValue>
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
