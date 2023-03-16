import {
  TextField,
  ColorField,
} from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/common";
import { Wrapper } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { generateColorGradient } from "@web/extensions/sidebar/utils/color";
import { isEqual } from "lodash";
import { useCallback, useEffect, useState } from "react";

import { BaseFieldProps } from "../types";

const PointColorGradient: React.FC<BaseFieldProps<"pointColorGradient">> = ({
  value,
  editMode,
  isActive,
  onUpdate,
}) => {
  const [colorGradient, setColorGradient] = useState(value);

  const handleFieldChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setColorGradient({ ...colorGradient, field: e.target.value });
    },
    [colorGradient],
  );

  const handleMinValueUpdate = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const min = !isNaN(parseFloat(e.currentTarget.value)) ? parseFloat(e.currentTarget.value) : 1;
      setColorGradient({ ...colorGradient, min });
    },
    [colorGradient],
  );

  const handleMaxValueUpdate = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const max = !isNaN(parseFloat(e.currentTarget.value)) ? parseFloat(e.currentTarget.value) : 1;
      setColorGradient({ ...colorGradient, max });
    },
    [colorGradient],
  );

  const handleStartColorUpdate = useCallback(
    (color: string) => {
      // colors should be hex for now
      if (color && /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/.test(color)) {
        setColorGradient({ ...colorGradient, startColor: color });
      }
    },
    [colorGradient],
  );

  const handleEndColorUpdate = useCallback(
    (color: string) => {
      // colors should be hex for now
      if (color && /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/.test(color)) {
        setColorGradient({ ...colorGradient, endColor: color });
      }
    },
    [colorGradient],
  );

  const handleStepUpdate = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const step = !isNaN(parseFloat(e.currentTarget.value))
        ? parseFloat(e.currentTarget.value)
        : 1;
      setColorGradient({ ...colorGradient, step });
    },
    [colorGradient],
  );

  const generateValues = (min: number, max: number, step: number) => {
    const result = [];
    for (let i = min; i <= max; i += step) {
      result.push(i);
    }
    return result;
  };

  const generateConditions = useCallback(
    (
      field: string,
      min: number,
      max: number,
      step: number,
      startColor: string,
      endColor: string,
    ) => {
      const values = generateValues(min, max, step);
      const count = values.length;
      const colors = generateColorGradient(startColor, endColor, count);
      const conditions: [string, string][] = [];
      const fieldName = "${" + field + "}";
      values.forEach((value, index) => {
        const cond: [string, string] = [`${fieldName} >= ${value}`, `'${colors[index]}'`];
        conditions.unshift(cond);
      });

      return conditions;
    },
    [],
  );

  useEffect(() => {
    if (!isActive || isEqual(colorGradient, value)) return;

    const { field, min, max, step, startColor, endColor } = colorGradient;
    if (
      field &&
      min !== undefined &&
      max !== undefined &&
      step !== undefined &&
      startColor &&
      endColor
    ) {
      const conditions = generateConditions(field, min, max, step, startColor, endColor);
      const timer = setTimeout(() => {
        onUpdate({
          ...colorGradient,
          override: {
            marker: {
              style: "point",
              pointColor: {
                expression: {
                  conditions,
                },
              },
            },
          },
        });
      }, 500);
      return () => {
        clearTimeout(timer);
      };
    }
  }, [isActive, value, onUpdate, colorGradient, generateConditions]);

  return editMode ? (
    <Wrapper>
      <TextField
        title="Field"
        titleWidth={82}
        defaultValue={colorGradient.field}
        onChange={handleFieldChange}
      />
      <TextField
        title="Min Value"
        titleWidth={82}
        defaultValue={colorGradient.min}
        onChange={handleMinValueUpdate}
      />
      <TextField
        title="Max Value"
        titleWidth={82}
        defaultValue={colorGradient.max}
        onChange={handleMaxValueUpdate}
      />
      <ColorField
        title="Start Color"
        titleWidth={82}
        color={colorGradient.startColor}
        onChange={handleStartColorUpdate}
      />
      <ColorField
        title="End Color"
        titleWidth={82}
        color={colorGradient.endColor}
        onChange={handleEndColorUpdate}
      />
      <TextField
        title="Step"
        titleWidth={82}
        defaultValue={colorGradient.step}
        onChange={handleStepUpdate}
      />
    </Wrapper>
  ) : null;
};

export default PointColorGradient;
