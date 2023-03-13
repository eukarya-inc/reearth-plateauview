import { Field } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/common";
import { Wrapper } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { Select } from "@web/sharedComponents";
import { useCallback, useState } from "react";

import { BaseFieldProps } from "../types";

const displayStyleOptions = [
  { value: "table", label: "Properties" },
  { value: "html", label: "Description" },
];

const InfoboxStyle: React.FC<BaseFieldProps<"infoboxStyle">> = ({ value, editMode, onUpdate }) => {
  const [displayStyleValue, setDisplayStyleValue] = useState<"table" | "html">(value.displayStyle);

  const handleEventTypeChange = useCallback(
    (selectedProperty: "table" | "html") => {
      setDisplayStyleValue(selectedProperty);
      onUpdate({
        ...value,
        displayStyle: selectedProperty,
        override: {
          infobox: {
            content: {
              type: selectedProperty,
              value: "override",
            },
          },
        },
      });
    },
    [onUpdate, value],
  );

  return editMode ? (
    <Wrapper>
      <Field
        title="Display Style"
        titleWidth={88}
        noBorder
        value={
          <Select
            defaultValue={"table"}
            options={displayStyleOptions}
            style={{ width: "100%" }}
            value={displayStyleValue}
            onChange={handleEventTypeChange}
            getPopupContainer={trigger => trigger.parentElement ?? document.body}
          />
        }
      />
    </Wrapper>
  ) : null;
};

export default InfoboxStyle;
