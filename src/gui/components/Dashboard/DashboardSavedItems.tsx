import { useState } from "react";
import { RiAddLine, RiSubtractLine } from "react-icons/ri";
import Folder from "./Folder/Folder";
import File from "./File/File";
import classNames from "classnames";

interface ISavedItem {
  id: string;
  name: string;
  lastDateModified: string;
  saved: boolean;
  label: string;
}

interface ComponentProps {
  items: Array<ISavedItem>;
}

export default function DashboardSavedItems({ items }: ComponentProps) {
  const [expandedContent, setExpandedContent] = useState<boolean>(true);
  const contentClassName = classNames({
    hidden: !expandedContent,
    visible: expandedContent,
  });

  return (
    <div className="my-8">
      <div className="w-full flex justify-between items-center">
        <h2 className="subtitle">Saved Items ({items.length})</h2>
        <button onClick={() => setExpandedContent(!expandedContent)}>
          {expandedContent ? <RiSubtractLine /> : <RiAddLine />}
        </button>
      </div>

      <div className={contentClassName}>
        <div className="w-full flex flex-col">
          {items.map((item) => {
            if (item.label === "folder") {
              return (
                <Folder
                  key={item.id}
                  id={item.id}
                  name={item.name}
                  lastDateModified={item.lastDateModified}
                  saved={item.saved}
                />
              );
            } else {
              return (
                <File
                  key={item.id}
                  id={item.id}
                  name={item.name}
                  lastDateModified={item.lastDateModified}
                  saved={item.saved}
                />
              );
            }
          })}
        </div>
      </div>
    </div>
  );
}
