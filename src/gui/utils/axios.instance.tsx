import axios from "axios";
import cogoToast from "./cogoToast.instance";
import Cookies from "js-cookie";

const instance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_BASE_URL,
  validateStatus: function () {
    return true;
  },
});

instance.interceptors.response.use((res) => {
  if (res.status !== 200) {
    console.log(res.data);
    cogoToast.error("There has been an error, please try again");
  }

  return res;
});

instance.interceptors.request.use((req) => {
  if (Cookies.get("token")) {
    req.headers.Authorization = `Bearer ${Cookies.get("token")}`;
  }

  return req;
});

export default instance;
