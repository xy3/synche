import axios from "axios";
import cogoToast from "./cogoToast.instance";
import Cookies from "js-cookie";
import { isString } from "lodash";

const instance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_BASE_URL,
  validateStatus: function () {
    return true;
  },
});

instance.interceptors.response.use((res) => {
  if (res.status !== 200) {
    if (res.data.message) {
      cogoToast.error(res.data.message);
    } else if (isString(res.data)) {
      cogoToast.error(res.data);
    } else {
      cogoToast.error("There has been an error, please try again");
    }
  }

  return res;
});

instance.interceptors.request.use((req) => {
  if (Cookies.get("accessToken") && Cookies.get("refreshToken")) {
    req.headers["X-Token"] = Cookies.get("accessToken");
    req.headers["X-Refresh-Token"] = Cookies.get("refreshToken");
  }

  return req;
});

export default instance;
