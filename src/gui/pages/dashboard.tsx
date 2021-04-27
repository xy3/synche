import { useEffect, useState } from "react";
import Breadcrumb from "../components/Breadcrumb";
import DashboardTree from "../components/Dashboard/DashboardTree";
import Layout from "../components/Layout";
import { isLoggedIn } from "../utils/isLoggedIn";

const folders = [
  {
    id: "1",
    name: "Folder #1",
    lastDateModified: "2021-04-20",
    saved: false,
  },
  {
    id: "2",
    name: "Folder #2",
    lastDateModified: "2021-04-20",
    saved: true,
  },
  {
    id: "3",
    name: "Folder #3",
    lastDateModified: "2021-04-15",
    saved: false,
  },
];

const files = [
  {
    id: "4",
    name: "feelgood.mp3",
    lastDateModified: "2021-04-24",
    saved: false,
  },
  {
    id: "5",
    name: "my_songs.zip",
    lastDateModified: "2021-02-20",
    saved: true,
  },
  {
    id: "6",
    name: "treasure.txt",
    lastDateModified: "2021-04-14",
    saved: true,
  },
];

export default function Dashboard() {
  const [name, setName] = useState<string>("");

  async function getName() {
    // Replace with API call
    setName("Stack");
  }

  useEffect(() => {
    getName();
  }, []);

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
        <h1 className="title">
          Welcome back,{" "}
          <span className="border-b border-brand-blue text-brand-blue">
            {name}
          </span>
        </h1>
        <DashboardTree path="/" folders={folders} files={files} />
      </section>
    </Layout>
  );
}

export const getServerSideProps = async ({ req }) => {
  if (!isLoggedIn(req.cookies.token || "")) {
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
