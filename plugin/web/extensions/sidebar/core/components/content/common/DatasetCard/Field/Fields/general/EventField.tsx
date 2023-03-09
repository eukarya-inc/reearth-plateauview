import {
  Field,
  TextField,
} from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/common";
import { Wrapper } from "@web/extensions/sidebar/core/components/content/common/DatasetCard/Field/commonComponents";
import { Select } from "@web/sharedComponents";
import { isEqual } from "lodash";
import { useCallback, useEffect, useMemo, useState } from "react";

import { BaseFieldProps } from "../types";

const eventTypeOptions = [{ value: "select", label: "Select Feature" }];

const triggerEventOptions = [{ value: "jump-to-url", label: "Jump To URL" }];

const urlTypeOptions = [
  { value: "manual", label: "Manual" },
  { value: "from-data", label: "From Data" },
];

const EventField: React.FC<BaseFieldProps<"eventField">> = ({
  value,
  editMode,
  isActive,
  onUpdate,
}) => {
  const [eventValue, setEventValue] = useState(value);

  const showURL = useMemo(() => eventValue.urlType === "manual", [eventValue.urlType]);
  const showField = useMemo(() => eventValue.urlType === "from-data", [eventValue.urlType]);

  const handleEventTypeChange = useCallback(
    (value: string) => {
      console.log(value);
      setEventValue({ ...eventValue, eventType: value });
    },
    [eventValue],
  );

  const handleTriggerEventChange = useCallback(
    (value: string) => {
      console.log(value);
      setEventValue({ ...eventValue, triggerEvent: value });
    },
    [eventValue],
  );

  const handleURLTypeChange = useCallback(
    (value: typeof eventValue.urlType) => {
      console.log(value);
      setEventValue({ ...eventValue, urlType: value });
    },
    [eventValue],
  );

  const handleURLChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      console.log(e.target.value);
      setEventValue({ ...eventValue, url: e.target.value });
    },
    [eventValue],
  );

  const handleFieldChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      console.log(e.target.value);
      setEventValue({ ...eventValue, field: e.target.value });
    },
    [eventValue],
  );

  useEffect(() => {
    if (!isActive || isEqual(eventValue, value)) return;
    const timer = setTimeout(() => {
      onUpdate({
        ...eventValue,
        override: {
          marker: {
            imageSize: 2,
          },
          events: {
            select: {
              openUrl: {
                urlKey: eventValue.url,
              },
            },
          },
        },
      });
    }, 500);
    return () => {
      clearTimeout(timer);
    };
  }, [isActive, value, onUpdate, eventValue.eventType, eventValue]);

  return editMode ? (
    <Wrapper>
      <Field
        title="Event Type"
        titleWidth={88}
        noBorder
        value={
          <Select
            defaultValue={"select-feature"}
            options={eventTypeOptions}
            style={{ width: "100%" }}
            value={eventValue.eventType}
            onChange={handleEventTypeChange}
            getPopupContainer={trigger => trigger.parentElement ?? document.body}
          />
        }
      />
      <Field
        title="Trigger Event"
        titleWidth={88}
        noBorder
        value={
          <Select
            defaultValue={"jump-to-url"}
            options={triggerEventOptions}
            style={{ width: "100%" }}
            value={eventValue.triggerEvent}
            onChange={handleTriggerEventChange}
            getPopupContainer={trigger => trigger.parentElement ?? document.body}
          />
        }
      />
      <Field
        title="URL Type"
        titleWidth={88}
        noBorder
        value={
          <Select
            defaultValue={"manual"}
            options={urlTypeOptions}
            style={{ width: "100%" }}
            value={eventValue.urlType}
            onChange={handleURLTypeChange}
            getPopupContainer={trigger => trigger.parentElement ?? document.body}
          />
        }
      />
      {showURL && (
        <TextField
          title="URL"
          titleWidth={88}
          defaultValue={eventValue.url}
          onChange={handleURLChange}
        />
      )}
      {showField && (
        <TextField
          title="Choose Field"
          titleWidth={88}
          defaultValue={eventValue.field}
          onChange={handleFieldChange}
        />
      )}
    </Wrapper>
  ) : null;
};

export default EventField;
