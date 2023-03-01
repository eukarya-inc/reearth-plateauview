import { DataCatalogItem, Group, Template } from "@web/extensions/sidebar/core/types";
import { generateID } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

import generateFieldComponentsList from "./Field/fieldHooks";

export default ({
  dataset,
  templates,
  inEditor,
  onDatasetUpdate,
  onOverride,
}: {
  dataset: DataCatalogItem;
  templates?: Template[];
  inEditor?: boolean;
  onDatasetUpdate: (dataset: DataCatalogItem, cleanseOverride?: any) => void;
  onOverride?: (dataID: string, activeIDs?: string[]) => void;
}) => {
  const [selectedGroup, setGroup] = useState<string>();
  const [activeComponentIDs, setActiveIDs] = useState<string[] | undefined>();

  useEffect(() => {
    const newActiveIDs = selectedGroup
      ? (!dataset.components?.find(c => c.type === "switchGroup") || !dataset.fieldGroups
          ? dataset.components
          : dataset.components.filter(
              c => (c.group && c.group === selectedGroup) || c.type === "switchGroup",
            )
        )
          ?.filter(c => !(!dataset.config?.data && c.type === "switchDataset"))
          ?.map(c => c.id)
      : dataset.components?.map(c => c.id);

    console.log("USE EFFECT IN DATASET CARDDD CALLEDDD");
    if (newActiveIDs !== activeComponentIDs) {
      console.log(
        "CHANGING ACTIVE IDS AND ON GROUP OVERRIDDEEEE",
        selectedGroup,
        activeComponentIDs,
      );
      setActiveIDs(newActiveIDs);
    }
  }, [selectedGroup, dataset.components]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    if (activeComponentIDs) {
      console.log("OVERRIDING", activeComponentIDs);
      onOverride?.(dataset.dataID, activeComponentIDs);
    }
  }, [activeComponentIDs]); // eslint-disable-line react-hooks/exhaustive-deps

  const handleCurrentGroupUpdate = useCallback(
    (fieldGroupID?: string) => {
      if (fieldGroupID === selectedGroup) return;
      console.log("fieldGroupID current grou pupdate", fieldGroupID);
      setGroup(fieldGroupID);
    },
    [selectedGroup],
  );

  const handleFieldAdd =
    (property: any) =>
    ({ key }: { key: string }) => {
      if (!inEditor) return;
      const newField = {
        id: generateID(),
        type: key.includes("template") ? "template" : key,
        ...property,
      };

      onDatasetUpdate?.({
        ...dataset,
        components: [...(dataset.components ?? []), newField],
      });
    };

  const handleFieldUpdate = useCallback(
    (id: string) => (property: any) => {
      const newDatasetComponents = dataset.components ? [...dataset.components] : [];
      const componentIndex = newDatasetComponents?.findIndex(c => c.id === id);

      if (!newDatasetComponents || componentIndex === undefined) return;

      newDatasetComponents[componentIndex] = property;

      onDatasetUpdate?.({
        ...dataset,
        components: newDatasetComponents,
      });
    },
    [dataset, onDatasetUpdate],
  );

  const handleFieldRemove = useCallback(
    (id: string) => {
      if (!inEditor) return;
      const newDatasetComponents = dataset.components ? [...dataset.components] : [];
      const componentIndex = newDatasetComponents?.findIndex(c => c.id === id);

      if (!newDatasetComponents || componentIndex === undefined) return;

      const removedComponent = newDatasetComponents.splice(componentIndex, 1)[0];

      console.log("REMOVED COMPONENT: ", removedComponent);
      if (removedComponent.type === "switchGroup") {
        handleCurrentGroupUpdate(undefined);
      }

      onDatasetUpdate?.(
        {
          ...dataset,
          components: newDatasetComponents,
        },
        removedComponent.cleanseOverride,
      );
    },
    [dataset, inEditor, onDatasetUpdate, handleCurrentGroupUpdate],
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
    templates,
    onFieldAdd: handleFieldAdd,
  });

  return {
    activeComponentIDs,
    fieldComponentsList,
    handleFieldUpdate,
    handleFieldRemove,
    handleCurrentGroupUpdate,
    handleGroupsUpdate,
  };
};
