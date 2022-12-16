import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
// import BetterScroll from "better-scroll";
import { useMemo } from "react";

const Editor: React.FC = () => {
  const stories = useMemo(() => [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], []);

  // useEffect(() => {
  //   const bs = new BetterScroll(".wrapper", {
  //     scrollX: true,
  //     probeType: 3,
  //   });

  //   setInterval(() => {
  //     bs.refresh();
  //   }, 5000);
  // }, []);

  return (
    <Wrapper className="wrapper">
      <Content>
        {stories.map((story, index) => (
          <Story key={index}>{story}</Story>
        ))}
        <CreateStory>
          <Icon icon="cornersOut" size={24} />
          <CreateText>Capture Scene</CreateText>
        </CreateStory>
      </Content>
    </Wrapper>
  );
};

const Wrapper = styled.div`
  height: 100%;
  flex: 1;
  white-space: nowrap;
`;

const Content = styled.div`
  height: 100%;
  padding: 12px;
  display: flex;
  width: 2000px;
  flex-wrap: nowrap;
  gap: 12px;
`;

const Story = styled.div`
  width: 170px;
  height: 100%;
  flex-shrink: 0;
  background: #f8f8f8;
  border-radius: 8px;
  border: 1px solid #c7c5c5;
`;

const CreateStory = styled.div`
  width: 170px;
  height: 100%;
  flex-shrink: 0;
  background: #fff;
  border-radius: 8px;
  border: 1px solid #00bebe;
  color: #00bebe;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
`;

const CreateText = styled.div`
  font-weight: 500;
  font-size: 14px;
  line-height: 21px;
`;

export default Editor;
