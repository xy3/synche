import Link from "next/link";
import React from "react";
import Breadcrumb from "../components/Breadcrumb";
import Layout from "../components/Layout";

export default function Page404() {
  return (
    <Layout title="Not Found">
      <section className="my-16 container">
        <Breadcrumb
          links={[
            {
              href: "/",
              title: "Synche",
            },
            {
              title: "Page Not Found",
            },
          ]}
        />
        <h1 className="title">Not Found</h1>
        <p className="my-8 text-gray-600">
          This page does not exist, however you can go{" "}
          <Link href="/">
            <a className="text-blue-500 border-b border-blue-500">home</a>
          </Link>
        </p>
      </section>
    </Layout>
  );
}
