export type FieldGroup = "genera" | "point" | "polyline" | "polygon" | "3d-model" | "3d-tile";

export type FieldType = "idealZoom" | "description" | "legend";

export type SharedFieldProps<T extends FieldType> = {
  id: string;
  type: T;
  title: string;
  inEditor?: boolean;
};
