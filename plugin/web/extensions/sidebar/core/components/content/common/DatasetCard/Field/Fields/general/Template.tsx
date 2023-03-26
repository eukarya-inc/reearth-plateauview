import { mergeOverrides } from "@web/extensions/sidebar/core/components/utils";
import { Template } from "@web/extensions/sidebar/core/types";
import { Select } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { isEqual } from "lodash";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import FieldComponent from "../..";
import { BaseFieldProps, FieldComponent as FieldComponentType } from "../types";

const Template: React.FC<BaseFieldProps<"template">> = ({
  value,
  dataID,
  editMode,
  activeIDs,
  templates,
  selectedGroup,
  configData,
  onUpdate,
  onCurrentGroupUpdate,
  onSceneUpdate,
}) => {
  const [fieldComponents, setFieldComponents] = useState<FieldComponentType[] | undefined>();

  const hasTemplates = useMemo(() => templates && templates.length > 0, [templates]);
  const currentTemplates = useRef<Template[] | undefined>();

  useEffect(() => {
    if (currentTemplates.current !== templates) {
      currentTemplates.current = templates;

      const newFieldComponents = hasTemplates
        ? templates?.find(t => t.id === value.templateID)?.components ?? templates?.[0].components
        : undefined;
      setFieldComponents(newFieldComponents);

      const cleanseOverride = mergeOverrides("cleanse", fieldComponents);
      onUpdate({
        ...value,
        userSettings: {
          components: newFieldComponents,
          override: cleanseOverride,
        },
      });
    } else {
      const newFieldComponents = value.userSettings?.components?.length
        ? value.userSettings?.components
        : hasTemplates
        ? templates?.find(t => t.id === value.templateID)?.components ?? templates?.[0].components
        : undefined;

      if (newFieldComponents && !isEqual(newFieldComponents, fieldComponents)) {
        setFieldComponents(newFieldComponents);

        onUpdate({
          ...value,
          userSettings: {
            components: newFieldComponents,
          },
        });
      }
    }
  }, [activeIDs, fieldComponents, templates, hasTemplates, value.templateID, value, onUpdate]);

  const handleTemplateChange = useCallback(
    (id: string) => {
      const cleanseOverride = mergeOverrides("cleanse", fieldComponents);

      onUpdate({
        ...value,
        templateID: id,
        userSettings: {
          components: templates?.find(t => t.id === id)?.components ?? [],
          override: cleanseOverride,
        },
      });
    },
    [value, fieldComponents, templates, onUpdate],
  );

  const handleFieldUpdate = useCallback(
    (id: string) => (property: any) => {
      const newComponents = [...(fieldComponents ?? [])];

      const componentIndex = newComponents?.findIndex(c => c.id === id);

      if (!newComponents || componentIndex === undefined) return;

      newComponents[componentIndex] = { ...property };

      onUpdate({
        ...value,
        userSettings: {
          components: newComponents,
          override: fieldComponents?.[componentIndex].cleanseOverride,
        },
      });
    },
    [value, fieldComponents, onUpdate],
  );

  const templateOptions = useMemo(
    () =>
      hasTemplates
        ? templates?.map(t => {
            return {
              value: t.id,
              label: t.name,
            };
          })
        : [{ value: "-", label: "-" }],
    [templates, hasTemplates],
  );

  return (
    <>
      {editMode ? (
        <div>
          <Title>テンプレート</Title>
          {hasTemplates ? (
            <Select
              options={templateOptions}
              style={{ width: "100%", alignItems: "center", height: "32px" }}
              value={value.templateID ?? templates?.[0].id}
              onChange={handleTemplateChange}
              getPopupContainer={trigger => trigger.parentElement ?? document.body}
            />
          ) : (
            <Text>保存されているテンプレートがないです。</Text>
          )}
        </div>
      ) : (
        fieldComponents
          ?.filter(t => activeIDs?.includes(t.id))
          ?.map(tc => (
            <FieldComponent
              key={tc.id}
              field={tc}
              editMode={editMode}
              dataID={dataID}
              activeIDs={activeIDs}
              isActive={!!activeIDs?.find(id => id === tc.id)}
              templates={templates}
              selectedGroup={selectedGroup}
              configData={configData}
              onUpdate={handleFieldUpdate}
              onCurrentGroupUpdate={onCurrentGroupUpdate}
              onSceneUpdate={onSceneUpdate}
            />
          ))
      )}
    </>
  );
};

export default Template;

const Title = styled.p`
  font-size: 12px;
  color: rgba(0, 0, 0, 0.85);
  margin: 0 0 4px 0;
`;

const Text = styled.p`
  margin: 0;
`;
