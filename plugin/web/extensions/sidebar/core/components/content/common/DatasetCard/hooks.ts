import { DataCatalogItem, Group, Template } from "@web/extensions/sidebar/core/types";
import { generateID } from "@web/extensions/sidebar/utils";
import { useCallback } from "react";

import generateFieldComponentsList from "./Field/fieldHooks";

export default ({
  dataset,
  templates,
  inEditor,
  onDatasetUpdate,
}: {
  dataset: DataCatalogItem;
  templates?: Template[];
  inEditor?: boolean;
  onDatasetUpdate: (dataset: DataCatalogItem) => void;
}) => {
  const handleFieldAdd =
    (property: any) =>
    ({ key }: { key: string }) => {
      if (!inEditor) return;
      onDatasetUpdate?.({
        ...dataset,
        components: [
          ...(dataset.components ?? []),
          {
            id: generateID(),
            type: key.includes("template") ? "template" : key,
            ...property,
          },
        ],
      });
    };

  const handleFieldUpdate = useCallback(
    (id: string) => (property: any) => {
      if (!inEditor) return;
      const newDatasetComponents = dataset.components ? [...dataset.components] : [];
      const componentIndex = newDatasetComponents?.findIndex(c => c.id === id);

      if (!newDatasetComponents || componentIndex === undefined) return;

      newDatasetComponents[componentIndex] = property;

      onDatasetUpdate?.({
        ...dataset,
        components: newDatasetComponents,
      });
    },
    [dataset, inEditor, onDatasetUpdate],
  );

  const handleFieldRemove = useCallback(
    (id: string) => {
      if (!inEditor) return;
      const newDatasetComponents = dataset.components ? [...dataset.components] : [];
      const componentIndex = newDatasetComponents?.findIndex(c => c.id === id);

      if (!newDatasetComponents || componentIndex === undefined) return;

      newDatasetComponents.splice(componentIndex, 1);

      onDatasetUpdate?.({
        ...dataset,
        components: newDatasetComponents,
      });
    },
    [dataset, inEditor, onDatasetUpdate],
  );

  const handleGroupsUpdate = useCallback(
    (fieldID: string) => (groups: Group[], selectedGroupID?: string) => {
      if (!inEditor) return;

      const newDatasetComponents = dataset.components ? [...dataset.components] : [];
      const componentIndex = newDatasetComponents.findIndex(c => c.id === fieldID);

      if (newDatasetComponents.length > 0 && componentIndex !== undefined) {
        newDatasetComponents[componentIndex].group = selectedGroupID;
      }

      onDatasetUpdate?.({
        ...dataset,
        components: newDatasetComponents,
        fieldGroups: groups,
      });
    },
    [dataset, inEditor, onDatasetUpdate],
  );

  const generalFields: FieldDropdownItem = {
    idealZoom: {
      name: fieldName["idealZoom"],
      onClick: handleFieldAdd({
        position: {
          lng: 0,
          lat: 0,
          height: 0,
          pitch: 0,
          heading: 0,
          roll: 0,
        },
      }),
    },
    description: {
      name: fieldName["description"],
      onClick: handleFieldAdd({}),
    },
    legend: {
      name: fieldName["legend"],
      onClick: handleFieldAdd({
        style: "square",
        items: [{ title: "hey", color: "red" }],
      }),
    },
    switchGroup: {
      name: fieldName["switchGroup"],
      onClick: handleFieldAdd({
        title: "Switch Group",
        groups: dataset.fieldGroups[0]
          ? [{ id: generateID(), title: "新グループ1", fieldGroupID: dataset.fieldGroups[0].id }]
          : [],
      }),
    },
    buttonLink: {
      name: fieldName["buttonLink"],
      onClick: handleFieldAdd({}),
    },
    story: {
      name: fieldName["story"],
      onClick: handleFieldAdd({}),
    },
  };
  const fieldComponentsList = generateFieldComponentsList({
    fieldGroups: dataset.fieldGroups,
    templates,
    onFieldAdd: handleFieldAdd,
  });

  return {
    fieldComponentsList,
    handleFieldUpdate,
    handleFieldRemove,
    handleGroupsUpdate,
  };
};
