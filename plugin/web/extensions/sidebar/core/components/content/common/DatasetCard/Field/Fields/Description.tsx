import { SharedFieldProps } from "./types";

type Props = SharedFieldProps<"description"> & {};

const Description: React.FC<Props> = () => {
  return <p>I AM Description</p>;
};

export default Description;
