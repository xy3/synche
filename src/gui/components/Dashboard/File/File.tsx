import {
  RiArrowDownSFill,
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
import { IFile } from "../../../utils/interfaces";
import Cookies from "js-cookie";
import cogoToast from "../../../utils/cogoToast.instance";
import FileDeleteModal from "../../Modals/FileDeleteModal";
import { saveAs } from "file-saver";

interface ComponentProps {
  data: IFile;
}

function saveFile(blob, filename) {
  if (window.navigator.msSaveOrOpenBlob) {
    window.navigator.msSaveOrOpenBlob(blob, filename);
  } else {
    const a = document.createElement("a");
    document.body.appendChild(a);
    const url = window.URL.createObjectURL(blob);
    a.href = url;
    a.download = filename;
    a.click();
    setTimeout(() => {
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);
    }, 0);
  }
}

export default function File({ data }: ComponentProps) {
  const [deleteFileModalVisible, setDeleteFileModalVisible] = useState<boolean>(
    false
  );

  const [dropDownVisible, setDropDownVisible] = useState<boolean>(false);
  const container = classNames(
    "flex flex-wrap justify-start md:justify-end items-center border-b border-l border-r border-blue-100 bg-white p-4",
    {
      hidden: dropDownVisible === false,
      visible: dropDownVisible === true,
    }
  );

  async function initiateDownload(fileId: string) {
    try {
      fetch(`${process.env.NEXT_PUBLIC_BASE_URL}/download/${fileId}`, {
        method: "GET",
        headers: {
          "X-Token": Cookies.get("accessToken"),
        },
      })
        .then((res) => {
          return res.blob();
        })
        .then((blob) => {
          saveAs(blob, data.Name);
        })
        .catch((err) => {
          cogoToast.error("There has been an error downloading the file");
        });
    } catch (err) {
      cogoToast.error("There has been an error, please try again");
    }
  }

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
    <>
      <FileDeleteModal
        fileId={data.ID}
        isOpen={deleteFileModalVisible}
        onSubmit={() => window.location.reload()}
        onGoBack={() => setDeleteFileModalVisible(false)}
      />
      <div className="my-2 w-full shadow-sm bg-indigo-50 text-gray-600">
        <button
          className="py-6 px-4 w-full flex justify-between items-center"
          onClick={() => setDropDownVisible(!dropDownVisible)}
        >
          <div className="flex items-center">
            {getIconFromFileName(data.Name.split(".").pop())}
            <div className="flex flex-col">
              <p className="text-left text-lg text-indigo-500 font-semibold">
                {data.Name}
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
          <button
            className="m-2 text-sm context-button bg-indigo-50 text-indigo-500"
            onClick={() => initiateDownload(data.ID.toString())}
          >
            <RiFileDownloadLine className="icon" />
            <span>Download</span>
          </button>

          <button
            className="m-2 text-sm context-button bg-red-50 text-red-500"
            onClick={() => setDeleteFileModalVisible(true)}
          >
            <RiFileReduceLine className="icon" />
            <span>Delete file</span>
          </button>
        </div>
      </div>
    </>
  );
}
