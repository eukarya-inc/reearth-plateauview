import { Field } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/common";
import { Wrapper } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { Select } from "@web/sharedComponents";
import { useCallback, useState } from "react";

import { BaseFieldProps } from "../types";

const displayStyleOptions = [
  { value: "properties", label: "Properties" },
  { value: "description", label: "Description" },
];

const InfoboxStyle: React.FC<BaseFieldProps<"infoboxStyle">> = ({ value, editMode }) => {
  const [displayStyleValue, setDisplayStyleValue] = useState(value);

  const handleEventTypeChange = useCallback(
    (value: string) => {
      setDisplayStyleValue({ ...displayStyleValue, eventType: value });
    },
    [displayStyleValue],
  );

  return editMode ? (
    <Wrapper>
      <Field
        title="Display Style"
        titleWidth={88}
        noBorder
        value={
          <Select
            defaultValue={"select"}
            options={displayStyleOptions}
            style={{ width: "100%" }}
            value={displayStyleValue.eventType}
            onChange={handleEventTypeChange}
            getPopupContainer={trigger => trigger.parentElement ?? document.body}
          />
        }
      />
    </Wrapper>
  ) : null;
};

export default InfoboxStyle;
