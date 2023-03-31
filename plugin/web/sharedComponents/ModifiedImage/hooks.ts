import { useState, useEffect } from "react";
import tinycolor from "tinycolor2";

type Props = {
  imageUrl: string;
  blendColor: string;
  width: number;
  height: number;
};

const useModifiedImage = ({ imageUrl, blendColor, width, height }: Props) => {
  const [modifiedImageUrl, setModifiedImageUrl] = useState("");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const image = new Image();
    image.crossOrigin = "anonymous";
    image.src = imageUrl;
    image.onload = () => {
      if (!image.complete) {
        console.error("Failed to load image due to cross-origin error");
        setLoading(false);
        return;
      }
      const canvas = document.createElement("canvas");
      canvas.width = width;
      canvas.height = height;
      const ctx = canvas.getContext("2d");
      if (!ctx) {
        setLoading(false);
        return;
      }
      ctx.drawImage(image, 0, 0, width, height);

      const mappedColor = tinycolor(blendColor).toRgb();
      const color = `rgba(${mappedColor.r}, ${mappedColor.g}, ${mappedColor.b}, ${mappedColor.a})`;

      // Convert the canvas to a data URL
      const canvas2 = document.createElement("canvas");
      canvas2.width = width;
      canvas2.height = height;
      const ctx2 = canvas2.getContext("2d");
      if (!ctx2) {
        setLoading(false);
        return;
      }
      ctx2.drawImage(image, 0, 0, width, height);
      ctx2.globalCompositeOperation = "multiply";
      ctx2.fillStyle = color;
      ctx2.fillRect(0, 0, canvas2.width, canvas2.height);
      ctx2.globalCompositeOperation = "destination-in";
      ctx2.drawImage(image, 0, 0, width, height);

      const modifiedImage = new Image();
      modifiedImage.src = canvas2.toDataURL();

      setModifiedImageUrl(modifiedImage.src);
      setLoading(false);
    };
    image.onerror = (event: string | Event) => {
      if (
        event instanceof Event &&
        event.target instanceof HTMLImageElement &&
        !event.target.complete
      ) {
        console.error("Failed to load image due to cross-origin error");
      }

      console.error(`Failed to load image: ${imageUrl}`);
      setLoading(false);
    };
  }, [imageUrl, blendColor, width, height]);

  return { modifiedImageUrl, loading };
};

export default useModifiedImage;
