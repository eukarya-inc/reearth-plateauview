import {
  ConfigData,
  FieldComponent,
} from "../core/components/content/common/DatasetCard/Field/Fields/types";

export const getActiveFieldIDs = (
  components?: FieldComponent[],
  selectedGroup?: string,
  config?: ConfigData[],
) =>
  flattenComponents(components)
    ?.filter(
      c =>
        !selectedGroup ||
        !c.group ||
        c.type === "switchGroup" ||
        (c.group && c.group === selectedGroup),
    )
    ?.filter(c => !(!config && c.type === "switchDataset"))
    ?.map(c => c.id);

export const flattenComponents = (components?: FieldComponent[]) =>
  components?.reduce((a: FieldComponent[], c?: FieldComponent) => {
    if (!c) return a;
    if (c.type === "template") {
      return [...a, c, ...(c.components ?? [])];
    } else {
      return [...a, c];
    }
  }, []);
