import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { useRef, useCallback, type WheelEvent } from "react";
import { DndProvider } from "react-dnd";
import { HTML5Backend } from "react-dnd-html5-backend";

import type { Camera, Story as StoryType } from "../../types";

import Story from "./story";

type Props = {
  stories: StoryType[];
  captureScene: () => void;
  viewStory: (camera: Camera) => void;
  recapture: (id: string) => void;
  deleteStory: (id: string) => void;
  editStory: (id: string) => void;
  moveStory: (dragIndex: number, hoverIndex: number) => void;
};

const Editor: React.FC<Props> = ({
  stories,
  captureScene,
  viewStory,
  recapture,
  deleteStory,
  editStory,
  moveStory,
}) => {
  const wrapperRef = useRef<HTMLDivElement>(null);

  const handleWheel = useCallback((e: WheelEvent) => {
    if (wrapperRef.current) {
      wrapperRef.current.scrollLeft += e.deltaY < 0 ? -30 : 30;
    }
  }, []);

  return (
    <Wrapper onWheel={handleWheel} ref={wrapperRef}>
      <DndProvider backend={HTML5Backend}>
        <Content>
          {stories?.map((story, index) => (
            <Story
              key={story.id}
              index={index}
              viewStory={viewStory}
              recapture={recapture}
              deleteStory={deleteStory}
              editStory={editStory}
              moveStory={moveStory}
              {...story}
            />
          ))}
          <CreateStory onClick={captureScene}>
            <Icon icon="cornersOut" size={24} />
            <CreateText>Capture Scene</CreateText>
          </CreateStory>
        </Content>
      </DndProvider>
    </Wrapper>
  );
};

const Wrapper = styled.div`
  height: 100%;
  flex: 1;
  overflow-x: auto;
  padding: 12px;
  scrollbar-width: thin;
`;

const Content = styled.div`
  display: flex;
  height: 100%;
  gap: 12px;
  float: left;
`;

const CreateStory = styled.div`
  width: 170px;
  height: 114px;
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
