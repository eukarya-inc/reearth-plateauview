import { SharedFieldProps } from "./types";

type Props = SharedFieldProps<"legend"> & {};

const Legend: React.FC<Props> = () => {
  return <p>I AM LEGEND</p>;
};

export default Legend;
