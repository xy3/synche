import React, { useMemo } from "react";
import DashboardContextMenu from "./DashboardContextMenu";
import DashboardFiles from "./DashboardFiles";
import DashboardFolders from "./DashboardFolders";
import DashboardSavedItems from "./DashboardSavedItems";

interface IFile {
  id: string;
  name: string;
  lastDateModified: string;
  saved: boolean;
}

interface IFolder {
  id: string;
  name: string;
  lastDateModified: string;
  saved: boolean;
}

interface ComponentProps {
  path: string;
  folders: Array<IFolder>;
  files: Array<IFile>;
}

export default function DashboardTree({
  path,
  folders,
  files,
}: ComponentProps) {
  const savedItems = useMemo(() => {
    const foldersWithKey = folders.map((f) => {
      return {
        ...f,
        label: "folder",
      };
    });
    const filesWithKey = files.map((f) => {
      return {
        ...f,
        label: "file",
      };
    });

    const items = [...foldersWithKey, ...filesWithKey];

    return items.filter((item) => {
      return item.saved;
    });
  }, [folders, files]);

  return (
    <div>
      <DashboardContextMenu
        canCreateFolders={path === "/"}
        displayGoBack={path !== "/"}
      />
      {savedItems.length > 0 ? (
        <>
          <h2 className="px-2 md:px-0 text-xl font-bold my-4">
            My Saved Items
          </h2>
          <DashboardSavedItems items={savedItems} />
        </>
      ) : null}

      <h2 className="px-2 md:px-0 text-xl font-bold my-4">
        Current Directory: {path}
      </h2>
      {folders.length > 0 ? <DashboardFolders folders={folders} /> : null}
      {files.length > 0 ? <DashboardFiles files={files} /> : null}
    </div>
  );
}
