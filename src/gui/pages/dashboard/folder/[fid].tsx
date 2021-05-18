import React, { useEffect, useMemo, useState } from "react";
import Breadcrumb from "../../../components/Breadcrumb";
import Layout from "../../../components/Layout";
import { isLoggedIn } from "../../../utils/isLoggedIn";
import axios from "../../../utils/axios.instance";
import cogoToast from "../../../utils/cogoToast.instance";
import {
  RiArrowLeftLine,
  RiFileUploadLine,
  RiFolderAddLine,
} from "react-icons/ri";
import NewFolderModal from "../../../components/Modals/NewFolderModal";
import FileUploadModal from "../../../components/Modals/FileUploadModal";
import DashboardFolders from "../../../components/Dashboard/DashboardFolders";
import DashboardFiles from "../../../components/Dashboard/DashboardFiles";
import { UserConsumer } from "../../../context/userContext";
import {
  ICurrentDirectory,
  IDirectory,
  IFile,
} from "../../../utils/interfaces";

interface ComponentProps {
  folderId: string;
}

export default function SpecificFolder({ folderId }: ComponentProps) {
  const [files, setFiles] = useState<IFile[]>([]);
  const [directories, setDirectories] = useState<IDirectory[]>([]);
  const [currentDirectory, setCurrentDirectory] = useState<ICurrentDirectory>({
    ID: -1,
    Name: "Loading...",
    Path: "",
    PathHash: "",
  });

  const [newFolderModalVisible, setNewFolderModalVisible] = useState<boolean>(
    false
  );
  const [newUploadModalVisible, setNewUploadModalVisible] = useState<boolean>(
    false
  );

  async function getDirectoryList() {
    try {
      const res = await axios.get(`/directory/${folderId}`);

      if (res.status === 200) {
        setCurrentDirectory(res.data.CurrentDir);
        setFiles(res.data.Files || []);
        setDirectories(res.data.SubDirectories || []);
      }
    } catch (err) {
      cogoToast.error("There has been an error, please try again");
    }
  }

  useEffect(() => {
    getDirectoryList();
  }, []);

  async function onNewFolder() {
    setNewFolderModalVisible(false);
    getDirectoryList();
  }

  async function onNewFile() {
    setNewUploadModalVisible(false);
    getDirectoryList();
  }

  return (
    <Layout title="My Dashboard">
      <section className="my-16 container">
        <Breadcrumb
          links={[
            {
              href: "/",
              title: "Synche",
            },
            {
              title: "Dashboard",
            },
          ]}
        />
        <UserConsumer>
          {(user) => (
            <h1 className="title">
              Welcome back{user.Name ? `, ${user.Name}` : null}
            </h1>
          )}
        </UserConsumer>

        <div>
          <NewFolderModal
            currentPathID={currentDirectory.ID}
            isOpen={newFolderModalVisible}
            onSubmit={onNewFolder}
            onGoBack={() => setNewFolderModalVisible(false)}
          />
          <FileUploadModal
            currentPathID={currentDirectory.ID}
            isOpen={newUploadModalVisible}
            onSubmit={onNewFile}
            onGoBack={() => setNewUploadModalVisible(false)}
          />
          <div className="my-8 w-full flex justify-end">
            <div className="w-full md:w-2/3 flex justify-end">
              {currentDirectory.hasOwnProperty("ParentDirectoryID") ? (
                <div className="w-1/3 px-2">
                  <a
                    href={`/dashboard/folder/${
                      currentDirectory.ParentDirectoryID || ""
                    }`}
                    className="secondary-button flex justify-center items-center"
                  >
                    <RiArrowLeftLine className="icon" />
                    <span>Go Back</span>
                  </a>
                </div>
              ) : null}

              <div className="w-1/3 px-2">
                <button
                  className="secondary-button flex justify-center items-center"
                  onClick={() => setNewFolderModalVisible(true)}
                >
                  <RiFolderAddLine className="icon" />
                  <span>New Folder</span>
                </button>
              </div>

              <div className="w-1/3 px-2">
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

          <h2 className="px-2 md:px-0 text-xl font-bold my-4">
            Current Directory: {currentDirectory.Name}
          </h2>
          <DashboardFolders directories={directories} />
          <DashboardFiles files={files} />
        </div>
      </section>
    </Layout>
  );
}

export const getServerSideProps = async ({ req, params }) => {
  if (!isLoggedIn(req.cookies.accessToken || "")) {
    return {
      redirect: {
        permanent: false,
        destination: "/login",
      },
    };
  }

  const folderId = params.fid || "";

  return {
    props: {
      folderId: folderId,
    },
  };
};
