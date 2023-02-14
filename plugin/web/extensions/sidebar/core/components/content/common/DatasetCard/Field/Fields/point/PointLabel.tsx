import { Switch, Dropdown } from "@web/sharedComponents";
import { useState } from "react";

import { BaseFieldProps } from "../types";

import ColorField from "./common/ColorField";
import Field from "./common/Field";
import { TextInput, Wrapper } from "./common/styled";

const PointLabel: React.FC<BaseFieldProps<"pointLabel">> = ({ value, editMode, onUpdate }) => {
  const [pointLabel, setPointLabel] = useState(value);

  const handlePointLabelUpdate = () => {
    setPointLabel(value);
    onUpdate(value);
  };

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
        value={<TextInput defaultValue={pointLabel.fontSize} onChange={handlePointLabelUpdate} />}
      />
      <Field
        title="Font color"
        titleWidth={82}
        value={<ColorField color={pointLabel.fontColor} value={pointLabel.fontColor} />}
      />
      <Field
        title="Height"
        titleWidth={82}
        value={<TextInput defaultValue={pointLabel.height} onChange={handlePointLabelUpdate} />}
      />
      <Field
        title="Extrude"
        titleWidth={82}
        value={<Switch defaultChecked={pointLabel.extruded} onChange={handlePointLabelUpdate} />}
      />
      <Field
        title="Use Background"
        titleWidth={82}
        value={
          <Switch defaultChecked={pointLabel.useBackground} onChange={handlePointLabelUpdate} />
        }
      />
      <Field
        title="Background Color"
        titleWidth={82}
        value={<ColorField color={pointLabel.backgroundColor} value={pointLabel.backgroundColor} />}
      />
    </Wrapper>
  ) : null;
};
export default PointLabel;
