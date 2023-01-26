import { Tag } from "@web/extensions/sidebar/core/processCatalog";

export type UserDataItem = {
  type: string;
  id: string;
  name?: string;
  prefecture?: string;
  cityName?: string;
  description?: string;
  data?: string;
  dataUrl?: string;
  dataFormat?: string;
  tags?: Tag[];
};
