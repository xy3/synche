import React, { useEffect, useState } from "react";
import Breadcrumb from "../../../components/Breadcrumb";
import DashboardTree from "../../../components/Dashboard/DashboardTree";
import Layout from "../../../components/Layout";
import { isLoggedIn } from "../../../utils/isLoggedIn";

interface ComponentProps {
  folderId: string;
}

const files = [
  {
    id: "7",
    name: "example1.txt",
    lastDateModified: "2021-01-24",
    saved: false,
  },
  {
    id: "8",
    name: "some_example.xlsx",
    lastDateModified: "2021-01-20",
    saved: false,
  },
];

export default function SpecificFolder({ folderId }: ComponentProps) {
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
              href: "/dashboard",
              title: "Dashboard",
            },
            {
              title: folderId,
            },
          ]}
        />
        <h1 className="title">
          Welcome back,{" "}
          <span className="border-b border-brand-blue text-brand-blue">
            {name}
          </span>
        </h1>
        <DashboardTree path={`/${folderId}`} folders={[]} files={files} />
      </section>
    </Layout>
  );
}

export const getServerSideProps = async ({ req, params }) => {
  if (!isLoggedIn(req.cookies.token || "")) {
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
