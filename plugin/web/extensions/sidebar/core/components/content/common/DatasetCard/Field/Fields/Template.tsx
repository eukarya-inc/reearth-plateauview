import { SharedFieldProps } from "./types";

type Props = SharedFieldProps<"template"> & {};

const Template: React.FC<Props> = () => {
  return <p>I AM Template</p>;
};

export default Template;
