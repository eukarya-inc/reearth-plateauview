import { styled } from "@web/theme";
import { ChangeEvent, useCallback, useState } from "react";

import { BaseFieldProps, Fields } from "../types";

import ColorField from "./common/ColorField";
import SelectField from "./common/SelectField";
import { Wrapper } from "./common/styled";
import SwitchField from "./common/SwitchField";
import TextField from "./common/TextField";

const options = [
  { value: "Option1", label: "Option1" },
  { value: "Option2", label: "Option2" },
];

const PointLabel: React.FC<BaseFieldProps<"pointLabel">> = ({ value, editMode, onUpdate }) => {
  const [pointLabel, setPointLabel] = useState(value);

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

  const handleFontColorUpdate = useCallback(
    (color: string) => {
      if (color) {
        setPointLabel(pointLabel => {
          const newPointLabel: Fields["pointLabel"] = {
            ...pointLabel,
            fontColor: color,
          };
          onUpdate({
            ...pointLabel,
            fontColor: color,
          });
          return newPointLabel;
        });
      }
    },
    [onUpdate],
  );

  const handleHeightUpdate = useCallback(
    (e: ChangeEvent<HTMLInputElement>) => {
      const height = !isNaN(parseFloat(e.currentTarget.value))
        ? parseFloat(e.currentTarget.value)
        : 1;
      setPointLabel(pointLabel => {
        const newPointLabel: Fields["pointLabel"] = {
          ...pointLabel,
          height,
        };
        onUpdate({
          ...pointLabel,
          height,
        });
        return newPointLabel;
      });
    },
    [onUpdate],
  );

  const handleExtrudedChange = (extruded: boolean) => {
    setPointLabel(pointLabel => {
      const newPointLabel: Fields["pointLabel"] = {
        ...pointLabel,
        extruded,
      };
      onUpdate({
        ...pointLabel,
        extruded,
      });
      return newPointLabel;
    });
  };

  const handleUseBackgroundChange = (useBackground: boolean) => {
    setPointLabel(pointLabel => {
      const newPointLabel: Fields["pointLabel"] = {
        ...pointLabel,
        useBackground,
      };
      onUpdate({
        ...pointLabel,
        useBackground,
      });
      return newPointLabel;
    });
  };

  const handleBackgroundColorUpdate = useCallback(
    (color: string) => {
      if (color) {
        setPointLabel(pointLabel => {
          const newPointLabel: Fields["pointLabel"] = {
            ...pointLabel,
            backgroundColor: color,
          };
          onUpdate({
            ...pointLabel,
            backgroundColor: color,
          });
          return newPointLabel;
        });
      }
    },
    [onUpdate],
  );

  return editMode ? (
    <Wrapper>
      <SelectField
        title="Choose field"
        titleWidth={82}
        options={options}
        onChange={handleFieldChange}
      />
      <TextField
        title="Font size"
        titleWidth={82}
        defaultValue={pointLabel.fontSize}
        suffix={<Suffix>px</Suffix>}
        onChange={handleFontSizeUpdate}
      />
      <ColorField
        title="Font color"
        titleWidth={82}
        color={pointLabel.fontColor}
        onChange={handleFontColorUpdate}
      />
      <TextField
        title="Height"
        titleWidth={82}
        defaultValue={pointLabel.height}
        suffix={<Suffix>m</Suffix>}
        onChange={handleHeightUpdate}
      />
      <SwitchField
        title="Extruded"
        titleWidth={82}
        checked={pointLabel.extruded}
        onChange={handleExtrudedChange}
      />
      <SwitchField
        title="Use Background"
        titleWidth={82}
        checked={pointLabel.useBackground}
        onChange={handleUseBackgroundChange}
      />
      {pointLabel.useBackground && (
        <ColorField
          title="Background color"
          titleWidth={82}
          color={pointLabel.backgroundColor}
          onChange={handleBackgroundColorUpdate}
        />
      )}
    </Wrapper>
  ) : null;
};
export default PointLabel;

export const Suffix = styled.span`
  color: rgba(0, 0, 0, 0.45);
`;
