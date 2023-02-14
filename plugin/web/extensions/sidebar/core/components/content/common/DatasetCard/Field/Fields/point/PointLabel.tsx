import { Switch, Dropdown } from "@web/sharedComponents";
import { ChangeEvent, useCallback, useState } from "react";

import { BaseFieldProps, Fields } from "../types";

import ColorField from "./common/ColorField";
import Field from "./common/Field";
import { TextInput, Wrapper } from "./common/styled";

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

  return editMode ? (
    <Wrapper>
      <Field
        title="Choose field"
        titleWidth={82}
        value={<Dropdown placement="bottom" trigger={["click"]} />}
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
      <Field
        title="Extrude"
        titleWidth={82}
        value={<Switch defaultChecked={pointLabel.extruded} />}
      />
      <Field
        title="Use Background"
        titleWidth={82}
        value={<Switch defaultChecked={pointLabel.useBackground} />}
      />
      <ColorField title="Background color" titleWidth={82} color={pointLabel.backgroundColor} />
    </Wrapper>
  ) : null;
};
export default PointLabel;
