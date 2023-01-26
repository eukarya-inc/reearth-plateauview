import { useEffect, useState } from "react";

import { postMsg } from "../utils";

import DesktopSidebar from "./components/Desktop";
import MobileSidebar from "./components/Mobile";

export type Props = {
  className?: string;
};

const Sidebar: React.FC<Props> = ({ className }) => {
  const [isMobile, setIsMobile] = useState<boolean>(false);

  useEffect(() => {
    postMsg({ action: "checkIfMobile" });
  }, []);

  useEffect(() => {
    const eventListenerCallback = (e: MessageEvent<any>) => {
      if (e.source !== parent) return;
      if (e.data.action === "checkIfMobile") {
        if (e.data.payload) {
          setIsMobile(e.data.payload);
        }
      }
    };
    addEventListener("message", e => eventListenerCallback(e));
    return () => {
      removeEventListener("message", eventListenerCallback);
    };
  }, []);

  return isMobile ? (
    <MobileSidebar className={className} />
  ) : (
    <DesktopSidebar className={className} />
  );
};

export default Sidebar;
