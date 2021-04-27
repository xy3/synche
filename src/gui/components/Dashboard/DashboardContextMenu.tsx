import classNames from "classnames";
import Link from "next/link";
import { useState } from "react";
import {
  RiFolderAddLine,
  RiFileUploadLine,
  RiArrowLeftLine,
} from "react-icons/ri";
import FileUploadModal from "../Modals/FileUploadModal";
import NewFolderModal from "../Modals/NewFolderModal";

interface ComponentProps {
  canCreateFolders: boolean;
  displayGoBack: boolean;
}

export default function DashboardContextMenu({
  canCreateFolders,
  displayGoBack,
}: ComponentProps) {
  const [newFolderModalVisible, setNewFolderModalVisible] = useState<boolean>(
    false
  );
  const [newUploadModalVisible, setNewUploadModalVisible] = useState<boolean>(
    false
  );

  const newFolderClassname = classNames("w-1/2 px-2", {
    hidden: !canCreateFolders,
  });

  const displayGoBackClassname = classNames("w-1/2 px-2", {
    hidden: !displayGoBack,
  });

  return (
    <>
      <NewFolderModal
        isOpen={newFolderModalVisible}
        onSubmit={() => setNewFolderModalVisible(false)}
        onGoBack={() => setNewFolderModalVisible(false)}
      />
      <FileUploadModal
        isOpen={newUploadModalVisible}
        onSubmit={() => setNewUploadModalVisible(false)}
        onGoBack={() => setNewUploadModalVisible(false)}
      />
      <div className="my-8 w-full flex justify-end">
        <div className="w-full md:w-1/2 flex justify-end">
          <div className={displayGoBackClassname}>
            <Link href="/dashboard">
              <a className="secondary-button flex justify-center items-center">
                <RiArrowLeftLine className="icon" />
                <span>Go Back</span>
              </a>
            </Link>
          </div>

          <div className={newFolderClassname}>
            <button
              className="secondary-button flex justify-center items-center"
              onClick={() => setNewFolderModalVisible(true)}
            >
              <RiFolderAddLine className="icon" />
              <span>New Folder</span>
            </button>
          </div>

          <div className="w-1/2 px-2">
            <button
              className="primary-button"
              onClick={() => setNewUploadModalVisible(true)}
            >
              <RiFileUploadLine className="icon" />
              <span>Upload File</span>
            </button>
          </div>
        </div>
      </div>
    </>
  );
}
