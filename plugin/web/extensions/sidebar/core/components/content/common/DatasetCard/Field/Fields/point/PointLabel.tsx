import { ChangeEvent, useCallback, useState } from "react";

import { BaseFieldProps, Fields } from "../types";

import ColorField from "./common/ColorField";
import Field from "./common/Field";
import SelectField from "./common/SelectField";
import { TextInput, Wrapper } from "./common/styled";
import SwitchField from "./common/SwitchField";

const options = [
  { value: "Option1", label: "Option1" },
  { value: "Option2", label: "Option2" },
];

const PointLabel: React.FC<BaseFieldProps<"pointLabel">> = ({ value, editMode, onUpdate }) => {
  const [pointLabel, setPointLabel] = useState(value);

  const handleFontSizeUpdate = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      const fontSize = !isNaN(parseFloat(e.currentTarget.value))
        ? parseFloat(e.currentTarget.value)
        : 1;
      setPointLabel(pointLabel => {
        const newPointLabel: Fields["pointLabel"] = {
          ...pointLabel,
          fontSize,
        };
        onUpdate({
          ...pointLabel,
          fontSize,
        });
        return newPointLabel;
      });
    },
    [onUpdate],
  );

  const handleFieldChange = (field: any) => {
    setPointLabel(pointLabel => {
      const newPointLabel: Fields["pointLabel"] = {
        ...pointLabel,
        field,
      };
      onUpdate({
        ...pointLabel,
        field,
      });
      return newPointLabel;
    });
  };

  return editMode ? (
    <Wrapper>
      <SelectField
        title="Choose field"
        titleWidth={82}
        onChange={handleFieldChange}
        options={options}
      />
      <Field
        title="Font size"
        titleWidth={82}
        value={<TextInput defaultValue={pointLabel.fontSize} onChange={handleFontSizeUpdate} />}
      />
      <ColorField title="Font color" titleWidth={82} color={pointLabel.fontColor} />
      <Field
        title="Height"
        titleWidth={82}
        value={<TextInput defaultValue={pointLabel.height} />}
      />
      <SwitchField title="Extrude" titleWidth={82} checked={pointLabel.extruded} />
      <SwitchField title="Use Background" titleWidth={82} checked={pointLabel.useBackground} />
      <ColorField title="Background color" titleWidth={82} color={pointLabel.backgroundColor} />
    </Wrapper>
  ) : null;
};
export default PointLabel;
