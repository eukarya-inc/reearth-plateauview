export type groupItem = {
  title: string;
  group: string;
  id?: number;
};

export type SwitchGroupObj = {
  title: string;
  groups: groupItem[];
};
