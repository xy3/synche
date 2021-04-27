import {
  RiArrowDownSFill,
  RiStarLine,
  RiArrowUpSFill,
  RiFileLine,
  RiFileDownloadLine,
  RiFileReduceLine,
  RiFileExcelLine,
  RiFileGifLine,
  RiFileMusicLine,
  RiFilePdfLine,
  RiFileZipLine,
} from "react-icons/ri";
import classNames from "classnames";
import { ReactChild, useState } from "react";
import dayjs from "dayjs";

import relativeTime from "dayjs/plugin/relativeTime";

dayjs.extend(relativeTime);

interface ComponentProps {
  id: string;
  name: string;
  lastDateModified: string;
  saved: boolean;
}

export default function File({
  name,
  id,
  lastDateModified,
  saved,
}: ComponentProps) {
  const [dropDownVisible, setDropDownVisible] = useState<boolean>(false);
  const container = classNames(
    "flex flex-wrap justify-start md:justify-end items-center border-b border-l border-r border-blue-100 bg-white p-4",
    {
      hidden: dropDownVisible === false,
      visible: dropDownVisible === true,
    }
  );

  function getIconFromFileName(extension: string): ReactChild {
    if (["xlsx", "xls"].includes(extension)) {
      return <RiFileExcelLine className="icon text-indigo-500" />;
    } else if (["gif"].includes(extension)) {
      return <RiFileGifLine className="icon text-indigo-500" />;
    } else if (
      ["mp3", "wav", "ogg", "flac", "mid", "aif"].includes(extension)
    ) {
      return <RiFileMusicLine className="icon text-indigo-500" />;
    } else if (["pdf"].includes(extension)) {
      return <RiFilePdfLine className="icon text-indigo-500" />;
    } else if (["zip", "7z", "rar"].includes(extension)) {
      return <RiFileZipLine className="icon text-indigo-500" />;
    } else {
      return <RiFileLine className="icon text-indigo-500" />;
    }
  }

  return (
    <div className="my-2 w-full shadow-sm bg-indigo-50 text-gray-600">
      <button
        className="py-6 px-4 w-full flex justify-between items-center"
        onClick={() => setDropDownVisible(!dropDownVisible)}
      >
        <div className="flex items-center">
          {getIconFromFileName(name.split(".").pop())}
          <div className="flex flex-col">
            <p className="text-left text-lg text-indigo-500 font-semibold">
              {name}
            </p>
            <p className="text-left text-sm text-indigo-400">
              Modified{" "}
              {lastDateModified ? dayjs(lastDateModified).fromNow() : "N/A"}
            </p>
          </div>
        </div>
        {dropDownVisible ? (
          <RiArrowUpSFill className="icon text-indigo-500" />
        ) : (
          <RiArrowDownSFill className="icon text-indigo-500" />
        )}
      </button>

      <div className={container}>
        <button className="m-2 text-sm context-button bg-indigo-50 text-indigo-500">
          <RiFileDownloadLine className="icon" />
          <span>Download</span>
        </button>

        <button className="m-2 text-sm context-button bg-yellow-100 text-yellow-600">
          <RiStarLine className="icon" />
          <span>{saved ? "Remove from saved" : "Add to saved"}</span>
        </button>

        <button className="m-2 text-sm context-button bg-red-50 text-red-500">
          <RiFileReduceLine className="icon" />
          <span>Delete file</span>
        </button>
      </div>
    </div>
  );
}
