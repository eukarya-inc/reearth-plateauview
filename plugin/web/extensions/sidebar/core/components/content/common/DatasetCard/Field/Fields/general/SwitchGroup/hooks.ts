import { array_move } from "@web/extensions/sidebar/utils";
import { useCallback, useEffect, useState } from "react";

import { groupItem, SwitchGroup } from "../../types";

export default (value: SwitchGroup) => {
  const [switchGroupObj, updateGroups] = useState<SwitchGroup>(value);
  const [groupsTitle, setGroupsTitle] = useState(value.title);
  const [currentGroup, setCurrentGroup] = useState<groupItem>(value.groups[0]);
  const [modifiedGroups, updateModifiedGroups] = useState<SwitchGroup | undefined>();

  //initialize the helper array (modifiedGroups)
  useEffect(() => {
    updateModifiedGroups({ type: "switchGroup", title: "Temp", groups: [] });
  }, []);

  //add empty item each time we press on add item
  const handAddItem = useCallback(() => {
    updateModifiedGroups(switchGroup => {
      if (!switchGroup) return;
      switchGroup.groups.length == 0
        ? switchGroup.groups.push({ id: 0, group: currentGroup.group, title: currentGroup.title })
        : switchGroup.groups.push({ id: 0, group: "", title: "" });
      return { ...switchGroup, groups: switchGroup.groups };
    });
  }, [currentGroup]);

  //modiy the group in the helper array
  const handleModifyGroup = (group: string, index: number) => {
    updateModifiedGroups(switchGroup => {
      let newArray: groupItem[] | undefined = undefined;
      if (!switchGroup) return;
      newArray = switchGroup.groups;
      newArray[index].group = group;
      return { ...switchGroup, groups: newArray };
    });
  };

  //modify the title in the helper array
  const handleModifyGroupTitle = (title: string, index: number) => {
    updateModifiedGroups(switchGroup => {
      let newArray: groupItem[] | undefined = undefined;
      if (!switchGroup) return;
      newArray = switchGroup.groups;
      newArray[index].title = title;
      return { ...switchGroup, groups: newArray };
    });
  };

  //modify the title of the switch group in config and reflect it on main switch group field component
  const handleTitleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setGroupsTitle(e.target.value);
  }, []);

  //handle on change group in main field component
  const handleChooseGroup = useCallback((item: groupItem) => {
    setCurrentGroup(item);
  }, []);

  const handleMoveUp = useCallback((idx: number) => {
    if (idx === 0) return;
    updateModifiedGroups(switchGroup => {
      let newItems: groupItem[] | undefined = undefined;
      if (!switchGroup) return;
      newItems = switchGroup.groups;
      array_move(newItems, idx, idx - 1);
      return { ...switchGroup, groups: newItems };
    });
  }, []);

  const handleMoveDown = useCallback(
    (idx: number) => {
      if (modifiedGroups?.groups && idx >= modifiedGroups.groups.length - 1) return;
      updateModifiedGroups(switchGroup => {
        let newItems: groupItem[] | undefined = undefined;
        if (!switchGroup) return;
        newItems = switchGroup.groups;
        array_move(newItems, idx, idx + 1);
        return { ...switchGroup, groups: newItems };
      });
    },
    [modifiedGroups?.groups],
  );

  const handleRemove = useCallback((idx: number) => {
    updateModifiedGroups(switchGroup => {
      let newItems: groupItem[] | undefined = undefined;
      if (!switchGroup) return;
      newItems = switchGroup.groups.filter((_, idx2) => idx2 != idx);
      return { ...switchGroup, groups: newItems };
    });
  }, []);

  useEffect(() => {
    updateGroups(switchGroup => {
      return {
        ...switchGroup,
        groupsTitle,
      };
    });
    updateGroups(switchGroup => {
      switchGroup.groups?.forEach(item1 => {
        const itemFromArr2 = modifiedGroups?.groups.find(item2 => item2.group == item1.group);
        if (itemFromArr2) {
          item1.title = itemFromArr2.title;
        }
      });
      return { ...switchGroup };
    });
  }, [groupsTitle, modifiedGroups]);

  return {
    switchGroupObj,
    groupsTitle,
    currentGroup,
    modifiedGroups,
    handleModifyGroupTitle,
    handleModifyGroup,
    handAddItem,
    handleRemove,
    handleMoveDown,
    handleMoveUp,
    handleTitleChange,
    handleChooseGroup,
  };
};
