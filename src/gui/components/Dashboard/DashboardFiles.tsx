import { useState } from "react";
import { RiAddLine, RiSubtractLine } from "react-icons/ri";
import classNames from "classnames";
import File from "./File/File";
import { IFile } from "../../utils/interfaces";
import Skeleton from "../Skeleton";

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
          {files.length > 0 ? (
            files.map((file) => {
              return <File key={file.ID} data={file} />;
            })
          ) : (
            <Skeleton message="There are no files" />
          )}
        </div>
      </div>
    </div>
  );
}
