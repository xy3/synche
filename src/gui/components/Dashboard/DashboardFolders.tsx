import { useState } from "react";
import { RiAddLine, RiSubtractLine } from "react-icons/ri";
import Folder from "./Folder/Folder";
import classNames from "classnames";

interface IFolder {
  id: string;
  name: string;
  lastDateModified: string;
  saved: boolean;
}

interface ComponentProps {
  folders: Array<IFolder>;
}

export default function DashboardFolders({ folders }: ComponentProps) {
  const [expandedContent, setExpandedContent] = useState<boolean>(true);
  const contentClassName = classNames({
    hidden: !expandedContent,
    visible: expandedContent,
  });

  return (
    <div className="my-8">
      <div className="w-full flex justify-between items-center">
        <h2 className="subtitle">Folders ({folders.length})</h2>
        <button onClick={() => setExpandedContent(!expandedContent)}>
          {expandedContent ? <RiSubtractLine /> : <RiAddLine />}
        </button>
      </div>

      <div className={contentClassName}>
        <div className="w-full flex flex-col">
          {folders.map((folder) => {
            return (
              <Folder
                key={folder.id}
                id={folder.id}
                name={folder.name}
                lastDateModified={folder.lastDateModified}
                saved={folder.saved}
              />
            );
          })}
        </div>
      </div>
    </div>
  );
}
