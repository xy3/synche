import { useState } from "react";
import { RiAddLine, RiSubtractLine } from "react-icons/ri";
import Folder from "./Folder/Folder";
import classNames from "classnames";
import { IDirectory } from "../../utils/interfaces";
import Skeleton from "../Skeleton";

interface ComponentProps {
  directories: Array<IDirectory>;
}

export default function DashboardFolders({ directories }: ComponentProps) {
  const [expandedContent, setExpandedContent] = useState<boolean>(true);
  const contentClassName = classNames({
    hidden: !expandedContent,
    visible: expandedContent,
  });

  return (
    <div className="my-8">
      <div className="w-full flex justify-between items-center">
        <h2 className="subtitle">Folders ({directories.length})</h2>
        <button onClick={() => setExpandedContent(!expandedContent)}>
          {expandedContent ? <RiSubtractLine /> : <RiAddLine />}
        </button>
      </div>

      <div className={contentClassName}>
        <div className="w-full flex flex-col">
          {directories.length > 0 ? (
            directories.map((folder) => {
              return <Folder key={folder.ID} data={folder} />;
            })
          ) : (
            <Skeleton message="There are no directories" />
          )}
        </div>
      </div>
    </div>
  );
}
