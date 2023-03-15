import {
  ConfigData,
  FieldComponent,
  SwitchGroup,
} from "../core/components/content/common/DatasetCard/Field/Fields/types";
import { Template } from "../core/types";

export const getActiveFieldIDs = (
  components?: FieldComponent[],
  selectedGroup?: string,
  config?: ConfigData[],
  templates?: Template[],
) =>
  flattenComponents(components, templates)
    ?.filter(
      c =>
        !selectedGroup ||
        !c.group ||
        c.type === "switchGroup" ||
        (c.group && c.group === selectedGroup),
    )
    ?.filter(c => !(!config && c.type === "switchDataset"))
    ?.map(c => c.id);

export const flattenComponents = (components?: FieldComponent[], templates?: Template[]) =>
  components?.reduce((a: FieldComponent[], c?: FieldComponent) => {
    if (!c) return a;
    if (c.type === "template") {
      return [...a, c, ...(templates?.find(t => t.id === c.templateID)?.components ?? [])];
    } else {
      return [...a, c];
    }
  }, []);

export const getDefaultGroup = (components?: FieldComponent[], templates?: Template[]) => {
  if (!components) return;

  const switchGroupComponents = flattenComponents(components, templates)?.filter(
    c => c.type === "switchGroup",
  ) as SwitchGroup[] | undefined;

  if (switchGroupComponents && switchGroupComponents.length > 0) {
    return switchGroupComponents[switchGroupComponents.length - 1].groups[0].fieldGroupID;
  }
};
