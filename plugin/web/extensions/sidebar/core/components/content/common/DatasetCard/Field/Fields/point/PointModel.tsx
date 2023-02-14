import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { BaseFieldProps } from "../types";

import Field from "./common/Field";

const PointModel: React.FC<BaseFieldProps<"pointModel">> = ({ value, editMode, onUpdate }) => {
  const [modelURL, setModelURL] = useState(value.modelURL ?? "");
  const [scale, setImageSize] = useState(value.scale);

  const handleURLUpdate = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setModelURL(e.currentTarget.value);
      onUpdate({
        ...value,
        modelURL: e.currentTarget.value,
      });
    },
    [value, onUpdate],
  );

  const handleScaleUpdate = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const size = !isNaN(parseFloat(e.currentTarget.value))
        ? parseFloat(e.currentTarget.value)
        : 1;
      setImageSize(size);
      onUpdate({
        ...value,
        scale,
      });
    },
    [onUpdate, value, scale],
  );

  return editMode ? (
    <Wrapper>
      <Field
        title="モデルURL"
        titleWidth={82}
        value={<TextInput defaultValue={modelURL} onChange={handleURLUpdate} />}
      />
      <Field
        title="Scale"
        titleWidth={82}
        value={<TextInput defaultValue={scale} onChange={handleScaleUpdate} />}
      />
    </Wrapper>
  ) : null;
};

export default PointModel;

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
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
