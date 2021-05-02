import {
  RiArrowDownSFill,
  RiArrowUpSFill,
  RiFolderLine,
  RiFolderTransferLine,
  RiFolderReduceLine,
} from "react-icons/ri";
import classNames from "classnames";
import { useState } from "react";
import { IDirectory } from "../../../utils/interfaces";
import FolderDeleteModal from "../../Modals/FolderDeleteModal";

interface ComponentProps {
  data: IDirectory;
}

export default function Folder({ data }: ComponentProps) {
  const [
    deleteFolderModalVisible,
    setDeleteFolderModalVisible,
  ] = useState<boolean>(false);

  const [dropDownVisible, setDropDownVisible] = useState<boolean>(false);
  const container = classNames(
    "flex flex-wrap justify-start md:justify-end items-center border-b border-l border-r border-blue-100 bg-white p-4",
    {
      hidden: dropDownVisible === false,
      visible: dropDownVisible === true,
    }
  );

  return (
    <>
      <FolderDeleteModal
        folderId={data.ID}
        isOpen={deleteFolderModalVisible}
        onSubmit={() => window.location.reload()}
        onGoBack={() => setDeleteFolderModalVisible(false)}
      />
      <div className="my-2 w-full shadow-sm bg-blue-50 text-gray-600">
        <button
          className="py-6 px-4 w-full flex justify-between items-center"
          onClick={() => setDropDownVisible(!dropDownVisible)}
        >
          <div className="flex items-center">
            <RiFolderLine className="icon text-blue-500" />
            <div className="flex flex-col">
              <p className="text-left text-lg text-blue-500 font-semibold">
                {data.Name}
              </p>
              <p className="text-left text-sm text-blue-400">
                {data.FileCount ? `${data.FileCount} files` : null}
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
          <a
            href={`/dashboard/folder/${data.ID}`}
            className="m-2 text-sm context-button bg-indigo-50 text-indigo-500"
          >
            <RiFolderTransferLine className="icon" />
            <span>View files</span>
          </a>

          <button
            className="m-2 text-sm context-button bg-red-50 text-red-500"
            onClick={() => setDeleteFolderModalVisible(true)}
          >
            <RiFolderReduceLine className="icon" />
            <span>Delete folder</span>
          </button>
        </div>
      </div>
    </>
  );
}
