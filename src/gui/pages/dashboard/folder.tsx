import React, { useEffect, useState } from "react";
import Breadcrumb from "../../components/Breadcrumb";
import Layout from "../../components/Layout";
import { isLoggedIn } from "../../utils/isLoggedIn";
import axios from "../../utils/axios.instance";
import cogoToast from "../../utils/cogoToast.instance";
import { IFile, IDirectory, ICurrentDirectory } from "../../utils/interfaces";
import DashboardFolders from "../../components/Dashboard/DashboardFolders";
import DashboardFiles from "../../components/Dashboard/DashboardFiles";
import NewFolderModal from "../../components/Modals/NewFolderModal";
import FileUploadModal from "../../components/Modals/FileUploadModal";
import { RiFileUploadLine, RiFolderAddLine } from "react-icons/ri";
import { UserConsumer } from "../../context/userContext";

export default function Dashboard() {
  const [files, setFiles] = useState<IFile[]>([]);
  const [directories, setDirectories] = useState<IDirectory[]>([]);
  const [currentDirectory, setCurrentDirectory] = useState<ICurrentDirectory>({
    ID: -1,
    Name: "Loading...",
    Path: "",
    PathHash: "",
    ParentDirectoryID: -1,
  });

  const [newFolderModalVisible, setNewFolderModalVisible] = useState<boolean>(
    false
  );
  const [newUploadModalVisible, setNewUploadModalVisible] = useState<boolean>(
    false
  );

  async function getDirectoryList() {
    try {
      const res = await axios.get(`/directory`);

      if (res.status === 200) {
        setCurrentDirectory(res.data.CurrentDir);
        setFiles(res.data.Files);
        setDirectories(res.data.SubDirectories);
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
            isOpen={newFolderModalVisible}
            onSubmit={onNewFolder}
            onGoBack={() => setNewFolderModalVisible(false)}
          />
          <FileUploadModal
            isOpen={newUploadModalVisible}
            onSubmit={onNewFile}
            onGoBack={() => setNewUploadModalVisible(false)}
          />
          <div className="my-8 w-full flex justify-end">
            <div className="w-full md:w-2/3 flex justify-end">
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

export const getServerSideProps = async ({ req }) => {
  if (!isLoggedIn(req.cookies.accessToken || "")) {
    return {
      redirect: {
        permanent: false,
        destination: "/login",
      },
    };
  }

  return {
    props: {},
  };
};
