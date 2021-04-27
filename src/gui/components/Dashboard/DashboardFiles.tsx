import { useState } from "react";
import { RiAddLine, RiSubtractLine } from "react-icons/ri";
import classNames from "classnames";
import File from "./File/File";

interface IFile {
  id: string;
  name: string;
  lastDateModified: string;
  saved: boolean;
}

interface ComponentProps {
  files: Array<IFile>;
}

export default function DashboardFiles({ files }: ComponentProps) {
  const [expandedContent, setExpandedContent] = useState<boolean>(true);
  const contentClassName = classNames({
    hidden: !expandedContent,
    visible: expandedContent,
  });

  return (
    <div className="my-8">
      <div className="w-full flex justify-between items-center">
        <h2 className="subtitle">Files ({files.length})</h2>
        <button onClick={() => setExpandedContent(!expandedContent)}>
          {expandedContent ? <RiSubtractLine /> : <RiAddLine />}
        </button>
      </div>

      <div className={contentClassName}>
        <div className="w-full flex flex-col">
          {files.map((file) => {
            return (
              <File
                key={file.id}
                id={file.id}
                name={file.name}
                lastDateModified={file.lastDateModified}
                saved={file.saved}
              />
            );
          })}
        </div>
      </div>
    </div>
  );
}
