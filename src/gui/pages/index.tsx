import Layout from "../components/Layout";
import Link from "next/link";
import { Formik, Form, Field, ErrorMessage } from "formik";
import { UserConsumer } from "../context/userContext";
import { SiGithub, SiGitlab, SiSwagger } from "react-icons/si";

export default function Index() {
  return (
    <Layout>
      <div className="bg-brand-dark-blue">
        <section className="p-4 py-24 container">
          <img src="/img/logo-lines.png" className="w-64 h-auto mx-auto" />
          <h1 className="text-center font-rubik text-4xl md:text-6xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-gray-50 to-gray-200 tracking-wide">
            Sync your files faster than ever.
          </h1>
          <p className="my-8 text-gray-300 text-center tracking-wider text-lg">
            Synche is a self-hosted storage solution that features multi-part concurrent file uploading.
            Put simply, Synche allows you to upload large files faster and more reliably than leading cloud-storage
            providers, all while maintaining complete control of all your files and the data storage location.
          </p>
          <UserConsumer>
            {(user) =>
              user.loggedIn ? (
                <div className="my-16 w-full flex justify-center">
                  <div className="w-2/3 md:w-1/3">
                    <Link href="/dashboard">
                      <a className="primary-button block">Dashboard</a>
                    </Link>
                  </div>
                </div>
              ) : (
                <div className="my-16 w-full flex flex-col items-center">
                  <div className="w-2/3 md:w-1/3">
                    <Link href="/signup">
                      <a className="primary-button block">Get Started</a>
                    </Link>
                  </div>
                  <Link href="/login">
                    <a className="mt-4 link">Already n user?</a>
                  </Link>
                </div>
              )
            }
          </UserConsumer>

          <div className="my-8 w-full flex justify-center">
            <a className="mx-4" href="#">
              <SiSwagger className="icon text-gray-300 w-8 h-8" />
            </a>
            <a className="mx-4" href="#">
              <SiGithub className="icon text-gray-300 w-8 h-8" />
            </a>
            <a className="mx-4" href="#">
              <SiGitlab className="icon text-gray-300 w-8 h-8" />
            </a>
          </div>
        </section>
      </div>

      <section className="py-16 md:py-32 bg-brand-blue text-blue-200">
        <div className="container p-4">
          <div className="my-16">
            <h2 className="font-rubik text-4xl text-blue-50">
              How do I use it?
            </h2>
            <p>
              Lorem ipsum dolor sit amet consectetur, adipisicing elit. Deleniti
              sunt nisi, magni labore repudiandae nobis animi repellat libero
              nesciunt cupiditate neque cum, cumque minima delectus
              exercitationem amet aspernatur beatae laborum.
            </p>
          </div>
          <div className="my-16">
            <h2 className="font-rubik text-4xl text-blue-50">
              What are the advantages?
            </h2>
            <p>
              Lorem ipsum dolor sit amet consectetur, adipisicing elit. Deleniti
              sunt nisi, magni labore repudiandae nobis animi repellat libero
              nesciunt cupiditate neque cum, cumque minima delectus
              exercitationem amet aspernatur beatae laborum.
            </p>
          </div>
          <div className="my-16">
            <h2 className="font-rubik text-4xl text-blue-50">Is it free?</h2>
            <p>
              Lorem ipsum dolor sit amet consectetur, adipisicing elit. Deleniti
              sunt nisi, magni labore repudiandae nobis animi repellat libero
              nesciunt cupiditate neque cum, cumque minima delectus
              exercitationem amet aspernatur beatae laborum.
            </p>
          </div>
        </div>
      </section>
      <section className="my-24 container p-4">
        <h2 className="font-rubik text-4xl text-brand-blue" id="contact">
          Questions? We got the answers!
        </h2>
        <Formik
          initialValues={{ message: "", email: "" }}
          onSubmit={async (values) => {
            console.log(values);
          }}
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
              <label className="block text-gray-500 text-sm">Message *</label>
              <Field
                as="textarea"
                rows={5}
                name="message"
                className="input resize-none"
                placeholder="Please enter your message..."
              />
              <div className="mb-2 text-red-500 text-sm">
                <ErrorMessage name="message" />
              </div>
            </div>

            <div className="my-8 w-full flex justify-end">
              <div className="w-2/3 md:w-1/2">
                <button className="primary-button flex justify-center items-center">
                  <span>Send</span>
                </button>
              </div>
            </div>
          </Form>
        </Formik>
      </section>
    </Layout>
  );
}
