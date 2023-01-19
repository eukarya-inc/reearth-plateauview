import { styled } from "@web/theme";

type Props = {
  onClose: () => void;
};

const Backdrop: React.FC<Props> = ({ onClose }) => {
  return <BackdropWrapper onClick={onClose} />;
};
export default Backdrop;
const BackdropWrapper = styled.div`
  position: fixed;
  top: 0;
  right: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  z-index: 20;
  background-color: rgba(0, 0, 0, 0.75);
`;
