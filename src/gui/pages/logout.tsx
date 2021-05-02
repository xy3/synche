import { useEffect } from "react";
import Cookies from "js-cookie";
import BareboneLayout from "../components/BareboneLayout";

export default function Logout() {
  useEffect(() => {
    Cookies.remove("accessToken");
    Cookies.remove("refreshToken");

    setTimeout(() => {
      window.location.href = "/";
    }, 5000);
  }, []);

  return (
    <BareboneLayout title="Logging Out...">
      <div className="w-full h-screen flex justify-center items-center">
        <div>
          <img src="/img/logo.png" className="w-32 h-auto mx-auto" />
          <h1 className="my-8 title text-center">Logging Out</h1>
          <p className="text-center">You will be redirected in 5 seconds</p>
        </div>
      </div>
    </BareboneLayout>
  );
}
