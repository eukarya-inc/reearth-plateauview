import classes from "./Wrapper.module.css";

type Props = {
  className?: string;
};
const Wrapper: React.FC<Props> = () => {
  return (
    <div className={classes.Wrapper}>
      <h2>Hello World</h2>
    </div>
  );
};
export default Wrapper;
