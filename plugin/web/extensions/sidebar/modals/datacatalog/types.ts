import { AdditionalData } from "../../core/types";

import { DataCatalogItem } from "./api/api";

export type UserDataItem = Partial<DataCatalogItem> & {
  description?: string;
  additionalData?: AdditionalData;
};
