import {
  RiArrowDownSFill,
  RiStarLine,
  RiArrowUpSFill,
  RiFolderLine,
  RiFolderTransferLine,
  RiFolderReduceLine,
} from "react-icons/ri";
import classNames from "classnames";
import { useState } from "react";
import dayjs from "dayjs";
import { snakeCase } from "lodash";

import relativeTime from "dayjs/plugin/relativeTime";
import Link from "next/link";

dayjs.extend(relativeTime);

interface ComponentProps {
  id: string;
  name: string;
  lastDateModified: string;
  saved: boolean;
}

export default function Folder({
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

  return (
    <div className="my-2 w-full shadow-sm bg-blue-50 text-gray-600">
      <button
        className="py-6 px-4 w-full flex justify-between items-center"
        onClick={() => setDropDownVisible(!dropDownVisible)}
      >
        <div className="flex items-center">
          <RiFolderLine className="icon text-blue-500" />
          <div className="flex flex-col">
            <p className="text-left text-lg text-blue-500 font-semibold">
              {name}
            </p>
            <p className="text-left text-sm text-blue-400">
              Modified{" "}
              {lastDateModified ? dayjs(lastDateModified).fromNow() : "N/A"}
            </p>
          </div>
        </div>
        {dropDownVisible ? (
          <RiArrowUpSFill className="icon text-blue-500" />
        ) : (
          <RiArrowDownSFill className="icon text-blue-500" />
        )}
      </button>

      <div className={container}>
        <Link href={`/dashboard/folder/${snakeCase(name)}`}>
          <a className="m-2 text-sm context-button bg-indigo-50 text-indigo-500">
            <RiFolderTransferLine className="icon" />
            <span>View files</span>
          </a>
        </Link>

        <button className="m-2 text-sm context-button bg-yellow-100 text-yellow-600">
          <RiStarLine className="icon" />
          <span>{saved ? "Remove from saved" : "Add to saved"}</span>
        </button>

        <button className="m-2 text-sm context-button bg-red-50 text-red-500">
          <RiFolderReduceLine className="icon" />
          <span>Delete folder</span>
        </button>
      </div>
    </div>
  );
}
