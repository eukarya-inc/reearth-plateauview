import { styled } from "@web/theme";

type Props = {
  id: string;
  type: "camera";
  title: string;
  icon?: string;
};

const Camera: React.FC<Props> = () => {
  return (
    <div>
      <InnerWrapper>
        <Text>Position</Text>
        <InputWrapper>
          <Input type="number" placeholder="Latitude" />
          <Input type="number" placeholder="Longitude" />
          <Input type="number" placeholder="Altitude" />
        </InputWrapper>
      </InnerWrapper>
      <InnerWrapper>
        <Text>Pose</Text>
        <InputWrapper>
          <Input type="number" placeholder="Heading" />
          <Input type="number" placeholder="Pitch" />
          <Input type="number" placeholder="Roll" />
        </InputWrapper>
      </InnerWrapper>
      <ButtonWrapper>
        <Button>Clean</Button>
        <Button>Capture</Button>
      </ButtonWrapper>
    </div>
  );
};

export default Camera;

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
