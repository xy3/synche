import React, { useEffect } from "react";
import BareboneLayout from "../components/BareboneLayout";

export default function Dashboard() {
  useEffect(() => {
    setTimeout(() => {
      window.location.href = "/dashboard/folder";
    }, 5000);
  }, []);

  return (
    <BareboneLayout title="Redirecting...">
      <div className="w-full h-screen flex justify-center items-center">
        <div>
          <img src="/img/logo.png" className="w-32 h-auto mx-auto" />
          <h1 className="my-8 title text-center">Redirecting</h1>
          <p className="text-center">You will be redirected in 5 seconds</p>
        </div>
      </div>
    </BareboneLayout>
  );
}

export const getServerSideProps = async () => {
  return {
    redirect: {
      permanent: true,
      destination: "/dashboard/folder",
    },
  };
};
