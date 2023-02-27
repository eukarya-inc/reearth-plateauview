// import { useCallback, useMemo, useState } from "react";

// import { DataCatalogItem, Template } from "../types";

// export default ({
//   selectedDatasets,
//   templates,
// }: {
//   selectedDatasets?: DataCatalogItem[];
//   templates?: Template[];
// }) => {
//   //
//   // Overrides mgr

//   const [selectedGroup, setGroup] = useState<string>();

//   const handleCurrentGroupChange = useCallback((fieldGroupID: string) => {
//     setGroup(fieldGroupID);
//   }, []);

//   const activeComponentIDs = useMemo(
//     () =>
//       selectedDatasets
//         ?.map(sd =>
//           (!sd.components?.find(c => c.type === "switchGroup") || !sd.fieldGroups
//             ? sd.components
//             : sd.components.filter(
//                 c => (c.group && c.group === selectedGroup) || c.type === "switchGroup",
//               )
//           )
//             ?.filter(c => !(!sd.config?.data && c.type === "switchDataset"))
//             ?.map(c => c.id),
//         )
//         .filter(id => !!id),
//     [
//       selectedGroup,
//       dataset.components,
//       dataset.fieldGroups,
//       dataset.config?.data,
//       defaultTemplate?.components,
//     ],
//   );

//   // 1. Filter only active selected datasets
//   const activeSelectedDatasets = useMemo(
//     () =>
//       selectedDatasets
//         ?.map(sd => {
//           const psd = { ...sd };

//           const isActive = (
//             !sd.components?.find(c => c.type === "switchGroup") || !sd.fieldGroups
//               ? sd.components
//               : sd.components.filter(
//                   c => (c.group && c.group === selectedGroup) || c.type === "switchGroup",
//                 )
//           )
//             ?.filter(c => !(!sd.config?.data && c.type === "switchDataset"))
//             ?.map(c => c.id);

//           return isActive ? psd : undefined;
//         })
//         .filter(d => !!d) as DataCatalogItem[],
//     [selectedDatasets, selectedGroup],
//   );

//   // 2. collect all overrides
//   // 3. batch override
//   // 4. pass selected datasets (fieldComps) to UI

//   // Needs to flush anything that was active and now isn't, and vice versa. And update the scene when any changes occur
//   //      Needs to know what is active/inactive at all times
//   //      Needs to be available even when selection tab isn't selected
//   //      Needs to be single source of truth
//   //      MIGHT need to know format (aka if geojson, apply point, polyine, etc overrides only. If 3d tiles, 3d tile overrides only. Etc)

//   // Where could this be?
//   //      sidebar hooks
//   //      Just above DatasetCards (where the activeIds are calculated)
//   //      Jotai state

//   // Can I use lodash's merge again so we can create a

//   // TEMPLATES
//   // Get list of templates
//   // Check if defaultTemplate
//   //

//   // COMPONENTS
//   // Get list of dataset.components
//   // gather overrides

//   // Check activeIDs based on Switch group, TEMPLATES OR COMPONENTS
//   // update scene based on active dataset's components
//   return { processedSelectedDatasets, handleCurrentGroupChange };
// };
