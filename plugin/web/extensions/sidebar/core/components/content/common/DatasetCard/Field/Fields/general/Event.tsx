import { Field } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/common";
import { TextInput } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { useCallback, useEffect, useState } from "react";

import { BaseFieldProps } from "../types";

const Event: React.FC<BaseFieldProps<"event">> = ({ value, editMode, isActive, onUpdate }) => {
  const [type, setType] = useState(value.eventType);

  const handleSizeUpdate = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setType(e.currentTarget.value);
  }, []);

  useEffect(() => {
    if (!isActive || type === value.type) return;
    const timer = setTimeout(() => {
      onUpdate({
        ...value,
        eventType: type,
        override: { marker: { style: "point", type: type } },
      });
    }, 500);
    return () => {
      clearTimeout(timer);
    };
  }, [isActive, type, value, onUpdate]);

  return editMode ? (
    <Field
      title="Event type"
      titleWidth={82}
      value={<TextInput defaultValue={type} onChange={handleSizeUpdate} />}
    />
  ) : null;
};

export default Event;
