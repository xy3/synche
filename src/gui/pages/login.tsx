import { Formik, Form, Field, ErrorMessage } from "formik";
import Layout from "../components/Layout";
import * as yup from "yup";
import Link from "next/link";
import axios from "../utils/axios.instance";
import cogoToast from "../utils/cogoToast.instance";
import Cookies from "js-cookie";
import { RiLockLine } from "react-icons/ri";
import { isLoggedIn } from "../utils/isLoggedIn";
import Breadcrumb from "../components/Breadcrumb";

export default function Login() {
  const validationSchema = yup.object().shape({
    email: yup.string().required("Email is required"),
    password: yup.string().required("Password is required"),
  });

  return (
    <Layout title="Log In">
      <section className="my-16 container">
        <Breadcrumb
          links={[
            {
              href: "/",
              title: "Synche",
            },
            {
              title: "Log In",
            },
          ]}
        />
        <h1 className="title">Log In</h1>
        <div className="my-8 p-4 bg-white shadow-sm">
          <Formik
            initialValues={{ email: "", password: "" }}
            onSubmit={async (values) => {
              try {
                const res = await axios.post(`/login`, {
                  email: values.email,
                  password: values.password,
                });

                if (res.status === 200) {
                  Cookies.set("token", res.data.accessToken || "test");
                  Cookies.set("refreshToken", res.data.refreshToken || "test");
                  window.location.href = "/";
                }
              } catch (err) {
                cogoToast.error("There has been an error, please try again");
              }
            }}
            validationSchema={validationSchema}
          >
            <Form>
              <div className="my-4">
                <label className="block text-gray-500 text-sm">Email *</label>
                <Field
                  type="text"
                  name="email"
                  className="input"
                  placeholder="Please enter your email..."
                />
                <div className="mb-2 text-red-500 text-sm">
                  <ErrorMessage name="email" />
                </div>
              </div>

              <div className="my-4">
                <label className="block text-gray-500 text-sm">
                  Password *
                </label>
                <Field
                  as="input"
                  type="password"
                  name="password"
                  className="input"
                  placeholder="Please enter your password..."
                />
                <div className="mb-2 text-red-500 text-sm">
                  <ErrorMessage name="password" />
                </div>
              </div>

              <div className="my-8">
                <button className="primary-button flex justify-center items-center">
                  <RiLockLine className="icon" />
                  <span>Log In</span>
                </button>
                <p className="my-2 text-gray-500 text-sm text-center">
                  If you don't have an account, you can{" "}
                  <Link href="/signup">
                    <a className="link">sign up</a>
                  </Link>
                </p>
              </div>
            </Form>
          </Formik>
        </div>
      </section>
    </Layout>
  );
}

export const getServerSideProps = async ({ req }) => {
  if (isLoggedIn(req.cookies.token || "")) {
    return {
      redirect: {
        permanent: false,
        destination: "/dashboard",
      },
    };
  }

  return {
    props: {},
  };
};
