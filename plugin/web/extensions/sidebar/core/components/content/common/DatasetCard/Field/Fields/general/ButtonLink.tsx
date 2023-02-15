import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { BaseFieldProps } from "../types";

const ButtonLink: React.FC<BaseFieldProps<"buttonLink">> = ({ value, editMode, onUpdate }) => {
  const [CurrentButtonTitle, setCurrentButtonTitle] = useState(value.title);
  const [CurrentButtonLink, setCurrentButtonLink] = useState(value.link);

  const handleChangeButtonTitle = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setCurrentButtonTitle(e.currentTarget.value);
      onUpdate({
        ...value,
        title: CurrentButtonTitle,
      });
    },
    [CurrentButtonTitle, onUpdate, value],
  );

  const handleChangeButtonLink = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      let url = e.currentTarget.value;
      const prefix = "http://";
      if (!url.match(/^[a-zA-Z]+:\/\//)) {
        url = prefix + url;
      }
      setCurrentButtonLink(url);
      onUpdate({
        ...value,
        link: CurrentButtonLink,
      });
    },
    [CurrentButtonLink, onUpdate, value],
  );
  return editMode ? (
    <Wrapper>
      <Field>
        <FieldTitle>Title</FieldTitle>
        <FieldValue>
          <TextInput value={CurrentButtonTitle} onChange={handleChangeButtonTitle} />
        </FieldValue>
      </Field>

      <Field>
        <FieldTitle>Link</FieldTitle>
        <FieldValue>
          <TextInput value={CurrentButtonLink} onChange={handleChangeButtonLink} />
        </FieldValue>
      </Field>
    </Wrapper>
  ) : (
    <Wrapper>
      <Field>
        <FieldValue>
          <StyledButton onClick={() => window.open(CurrentButtonLink, "_blank")}>
            <Text>{CurrentButtonTitle}</Text>
          </StyledButton>
        </FieldValue>
      </Field>
    </Wrapper>
  );
};

export default ButtonLink;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

const Text = styled.p`
  margin: 0;
  font-weight: 400;
  font-size: 14px;
  line-height: 22px;
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

const StyledButton = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  padding: 1px 8px;
  background: #00bebe;
  color: #ffffff;
  box-shadow: 0px 2px 0px rgba(0, 0, 0, 0.043);
  border-radius: 2px;
  width: 100%;
  height: 100%;
  cursor: pointer;
`;
