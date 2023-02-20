export type InitData = {
  viewport: Viewport;
  data: RawDatasetData;
};

export type RawDatasetData = {
  title: string;
  dataID: string;
  searchIndex: {
    url: string;
  }[];
};

export type SearchIndex = {
  baseURL: string;
  indexRoot: {
    indexes: {
      [key: string]: {
        kind: string;
        values: {
          [key: string]: {
            count: number;
            url: string;
          };
        };
      };
    };
  };
  resultsData?: any[];
};

export type SearchResults = {
  threeDTilesId: string;
  results: Result[];
};

//
export type DatasetIndexes = {
  title: string;
  dataID: string;
  indexes: IndexData[];
};

export type IndexData = {
  field: string;
  values: string[];
};

export type Condition = {
  field: string;
  values: string[];
};

// from index
export type Result = {
  gml_id: string;
  Longitude: string;
  Latitude: string;
  Height: string;
};

// reearth types
export type Viewport = {
  width: number;
  height: number;
  isMobile: boolean;
};
