import { Camera } from "@web/extensions/sidebar/core/types";
import { styled } from "@web/theme";
import { useCallback, useState } from "react";

import { SharedFieldProps } from "./types";

type Props = SharedFieldProps<"idealZoom"> & {
  icon?: string;
  onCapture?: (camera: Partial<Camera>) => void;
};

const initialValues = {
  lat: 0,
  lng: 0,
  altitude: 0,
  heading: 0,
  pitch: 0,
  roll: 0,
};

const IdealZoom: React.FC<Props> = ({ inEditor, onCapture }) => {
  const [camera, setCamera] = useState<Camera>(initialValues);

  const handleLatitudeChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    console.log(e);
    // handleLatitudeChange
  }, []);

  const handleLongitudeChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    console.log(e);
    // handleLongitudeChange
  }, []);

  const handleAltitudeChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    console.log(e);
    // handleAltitudeChange
  }, []);

  const handleHeadingChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    console.log(e);
    // handleHeadingChange
  }, []);

  const handlePitchChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    console.log(e);
    // handlePitchChange
  }, []);

  const handleRollChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    console.log(e);
    // handleRollChange
  }, []);

  const handleCapture = useCallback(() => {
    if (!camera) return;
    onCapture?.(camera);
  }, [camera, onCapture]);

  const handleClean = useCallback(() => {
    setCamera(initialValues);
  }, []);

  return !inEditor ? (
    <div>
      <InnerWrapper>
        <Text>位置</Text>
        <InputWrapper>
          <Input
            type="number"
            placeholder="緯度" // Latitude
            value={camera.lat}
            onChange={handleLatitudeChange}
          />
          <Input
            type="number"
            placeholder="経度" // Longitude
            value={camera.lng}
            onChange={handleLongitudeChange}
          />
          <Input
            type="number"
            placeholder="高度" // Altitude
            value={camera.altitude}
            onChange={handleAltitudeChange}
          />
        </InputWrapper>
      </InnerWrapper>
      <InnerWrapper>
        <Text>ポーズ</Text>
        <InputWrapper>
          <Input
            type="number"
            placeholder="ヘッディング" // Heading
            value={camera.heading}
            onChange={handleHeadingChange}
          />
          <Input
            type="number"
            placeholder="ピッチ" // Pitch
            value={camera.pitch}
            onChange={handlePitchChange}
          />
          <Input
            type="number"
            placeholder="ロール" // Roll
            value={camera.roll}
            onChange={handleRollChange}
          />
        </InputWrapper>
      </InnerWrapper>
      <ButtonWrapper>
        <Button onClick={handleClean}>削除</Button>
        <Button onClick={handleCapture}>キャプチャー</Button>
      </ButtonWrapper>
    </div>
  ) : null;
};

export default IdealZoom;

const InnerWrapper = styled.div`
  display: flex;
  align-items: center;
`;

const Text = styled.p`
  margin: 0;
  width: 65px;
`;

const Input = styled.input`
  height: 32px;
  width: 64px;
  box-sizing: border-box;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  text-align: center;
`;

const InputWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 4px;
  margin-bottom: 8px;
`;

const ButtonWrapper = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 8px;
`;

const Button = styled.div`
  width: 100%;
  padding: 5px;
  border: 1px solid #d9d9d9;
  text-align: center;
  border-radius: 2px;
  user-select: none;
  cursor: pointer;

  :hover {
    background: #f4f4f4;
  }
`;
