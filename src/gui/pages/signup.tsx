import { Formik, Form, Field, ErrorMessage } from "formik";
import Layout from "../components/Layout";
import * as yup from "yup";
import Link from "next/link";
import axios from "../utils/axios.instance";
import cogoToast from "../utils/cogoToast.instance";
import Router from "next/router";
import { useState } from "react";
import { isLoggedIn } from "../utils/isLoggedIn";
import Breadcrumb from "../components/Breadcrumb";

export default function Signup() {
  const [btnDisabled, setBtnDisabled] = useState<boolean>(false);

  const validationSchema = yup.object().shape({
    name: yup.string().optional(),
    email: yup
      .string()
      .email("Must be a valid email")
      .required("Email is required"),
    password: yup
      .string()
      .min(6, "Password must be at least 6 characters")
      .required("Password is required"),
    password2: yup
      .string()
      .oneOf([yup.ref("password")], "Password doesn't match")
      .required("Password is required"),
  });

  return (
    <Layout title="Sign Up">
      <section className="my-16 container">
        <Breadcrumb
          links={[
            {
              href: "/",
              title: "Synche",
            },
            {
              title: "Sign Up",
            },
          ]}
        />
        <h1 className="title">Sign Up</h1>
        <div className="my-8 p-4 bg-white shadow-sm">
          <Formik
            initialValues={{ email: "", password: "", password2: "", name: "" }}
            onSubmit={async (values) => {
              try {
                setBtnDisabled(true);

                const url = new URLSearchParams();

                url.append("email", values.email);
                url.append("password", values.password);
                if (values.name) {
                  url.append("name", values.name);
                }

                const res = await axios.post(`/register?${url.toString()}`);

                setBtnDisabled(false);

                if (res.status === 200) {
                  cogoToast.success(
                    "Thank you for registering, you can now log in!"
                  );
                  Router.push("/login");
                }
              } catch (err) {
                setBtnDisabled(false);
                cogoToast.error("There has been an error, please try again");
              }
            }}
            validationSchema={validationSchema}
          >
            <Form>
              <div className="my-4">
                <label className="block text-gray-500 text-sm">Name</label>
                <Field
                  as="input"
                  type="text"
                  name="name"
                  className="input"
                  placeholder="Please enter your name (optional)..."
                />
                <div className="mb-2 text-red-500 text-sm">
                  <ErrorMessage name="name" />
                </div>
              </div>

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

              <div className="my-4">
                <label className="block text-gray-500 text-sm">
                  Re-enter your Password *
                </label>
                <Field
                  as="input"
                  type="password"
                  name="password2"
                  className="input"
                  placeholder="Please re-enter your password..."
                />
                <div className="mb-2 text-red-500 text-sm">
                  <ErrorMessage name="password2" />
                </div>
              </div>

              <div className="my-8">
                <button
                  className="primary-button"
                  disabled={btnDisabled}
                  type="submit"
                >
                  Sign Up
                </button>
                <p className="my-2 text-gray-500 text-sm text-center">
                  If you already have an account, you can{" "}
                  <Link href="/login">
                    <a className="link">log in</a>
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
  if (isLoggedIn(req.cookies.accessToken || "")) {
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
