import { Button, Icon } from "@web/sharedComponents";
import Video from "@web/sharedComponents/Video";
import { styled } from "@web/theme";

import useHooks from "./hooks";

const ClipVideo: React.FC = () => {
  const { handleClose } = useHooks();
  return (
    <div>
      <CloseButton>
        <Icon size={32} icon="close" onClick={handleClose} />
      </CloseButton>
      <Video width="560" height="315" src="https://www.youtube.com/embed/HQ2lDxVnJ9A" />
    </div>
  );
};
export default ClipVideo;

const CloseButton = styled(Button)`
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  right: 0;
  height: 32px;
  width: 32px;
  border: none;
  background: #00bebe;
  color: white;
  cursor: pointer;
`;
