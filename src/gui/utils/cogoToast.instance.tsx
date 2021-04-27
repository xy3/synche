import cogoToast from "cogo-toast";

const options = {
  hideAfter: 10,
};

const instance = {
  success: (message: string) => cogoToast.success(message, options),
  info: (message: string) => cogoToast.info(message, options),
  loading: (message: string) => cogoToast.loading(message, options),
  warn: (message: string) => cogoToast.warn(message, options),
  error: (message: string) => cogoToast.error(message, options),
};

export default instance;
