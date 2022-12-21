import mapVideo from "@web/extensions/sidebar/core/assets/mapVideo.png";
import { Button, Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

import useHooks from "./hooks";

const MapVideo: React.FC = () => {
  const { handleClose } = useHooks();
  return (
    <div>
      <CloseButton>
        <Icon size={32} icon="close" onClick={handleClose} />
      </CloseButton>
      <img src={mapVideo} />{" "}
    </div>
  );
};
export default MapVideo;

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
