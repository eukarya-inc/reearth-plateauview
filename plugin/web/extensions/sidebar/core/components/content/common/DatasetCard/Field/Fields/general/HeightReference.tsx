import { Field } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/common";
import { Wrapper } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { Select } from "@web/sharedComponents";
import { useCallback, useState } from "react";

import { BaseFieldProps } from "../types";

const heightReferenceOptions = [
  { value: "clamp", label: "Clamp to ground" },
  { value: "relative", label: "Relative to ground" },
  { value: "none", label: "None" },
];

const HeightReference: React.FC<BaseFieldProps<"heightReference">> = ({
  value,
  editMode,
  onUpdate,
}) => {
  const [heightReferenceType, setHeightReferenceType] = useState<"clamp" | "relative" | "none">(
    value.heightReferenceType ?? "clamp",
  );

  const handleEventTypeChange = useCallback(
    (selectedProperty: "clamp" | "relative" | "none") => {
      setHeightReferenceType(selectedProperty);
      onUpdate({
        ...value,
        heightReferenceType: selectedProperty,
        override: {
          resource: {
            clampToGround: selectedProperty === "clamp" ? true : false,
          },
          marker: {
            heightReference: selectedProperty,
          },
          polygon: {
            heightReference: selectedProperty,
          },
          polyline: {
            clampToGround: selectedProperty === "clamp" ? true : false,
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
            defaultValue={null}
            options={heightReferenceOptions}
            style={{ width: "100%" }}
            value={heightReferenceType}
            onChange={handleEventTypeChange}
            getPopupContainer={trigger => trigger.parentElement ?? document.body}
          />
        }
      />
    </Wrapper>
  ) : null;
};

export default HeightReference;
