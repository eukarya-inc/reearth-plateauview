import { DataCatalogItem, Group } from "@web/extensions/sidebar/core/types";
import { generateID } from "@web/extensions/sidebar/utils";
import { useCallback } from "react";

import generateFieldComponentsList from "./Field/fieldHooks";

export default ({
  dataset,
  inEditor,
  onDatasetUpdate,
}: {
  dataset: DataCatalogItem;
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
            type: key,
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

  const fieldComponentsList = generateFieldComponentsList({
    fieldGroups: dataset.fieldGroups,
    onFieldAdd: handleFieldAdd,
  });

  return {
    fieldComponentsList,
    handleFieldUpdate,
    handleFieldRemove,
    handleGroupsUpdate,
  };
};
